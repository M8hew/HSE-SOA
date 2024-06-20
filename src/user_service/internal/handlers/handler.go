package handlers

import (
	"crypto/md5"
	"fmt"
	"log"
	"net/http"
	"slices"
	"strings"
	"time"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"

	"user_service/api"
	pb "user_service/api/proto"
)

var ErrAuth = fmt.Errorf("authorization error")

type ServerHandler struct {
	db                  dbWrapper
	keys                rsaKeys
	userSessionLifeTime time.Duration
	contentService      pb.ContentServiceClient
	statService         pb.StatServiceClient
	kafkaProducer       *kafka.Producer
}

func hashPassword(password, salt string) (passwordHash [16]byte) {
	hash := md5.New()
	hash.Write([]byte(password + salt))
	copy(passwordHash[:], hash.Sum(nil)[:16])
	return
}

func checkAuth(ctx echo.Context, keys rsaKeys) (id int, err error) {
	authHeader := ctx.Request().Header.Get("Authorization")
	parts := strings.SplitN(authHeader, " ", 2)
	if len(parts) != 2 || parts[0] != "Bearer" {
		err = ctx.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid authorization header format"})
		if err == nil {
			err = ErrAuth
		}
		return
	}
	tokenString := parts[1]

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (any, error) { return keys.jwtPublic, nil })
	if err != nil {
		log.Println(err.Error())
		err = ctx.JSON(http.StatusBadRequest, map[string]string{"error": "Error parsing token"})
		if err == nil {
			err = ErrAuth
		}
		return
	}

	if !token.Valid {
		err = ctx.JSON(http.StatusUnauthorized, map[string]string{"error": "Unauthorized, invalid or expired token"})
		if err == nil {
			err = ErrAuth
		}
		return
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		err = ctx.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid claims in token"})
		if err == nil {
			err = ErrAuth
		}
		return
	}

	log.Println(claims)

	expTime, err := claims.GetExpirationTime()
	if err != nil || expTime == nil {
		err = ctx.JSON(http.StatusBadRequest, map[string]string{"error": "Error parsing exp field"})
		if err == nil {
			err = ErrAuth
		}
		return
	}

	if !ok || time.Now().After(expTime.Time) {
		err = ctx.JSON(http.StatusUnauthorized, map[string]string{"error": "Unauthorized, invalid or expired token"})
		if err == nil {
			err = ErrAuth
		}
		return
	}

	userId, ok := claims["user_id"].(float64)
	if !ok {
		err = ctx.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid user claims in token"})
		if err == nil {
			err = ErrAuth
		}
		return
	}
	return int(userId), nil
}

// User login
// (POST /login)
func (s *ServerHandler) PostLogin(ctx echo.Context) error {
	log.Println("Login request")
	loginRequest := api.PostLoginJSONRequestBody{}
	if err := ctx.Bind(&loginRequest); err != nil {
		log.Println(err.Error())
		return ctx.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request payload"})
	}

	userPassword, id, err := s.db.getUserPasswordId(*loginRequest.Username)
	if err != nil {
		log.Println(err.Error())
		return ctx.JSON(http.StatusInternalServerError, map[string]string{"error": "User not found"})
	}

	if hPassword := hashPassword(*loginRequest.Password, *loginRequest.Username); slices.Compare(userPassword, hPassword[:]) != 0 {
		return ctx.JSON(http.StatusUnauthorized, map[string]string{"error": "Incorrect login/password"})
	}

	expirationDate := time.Now().Add(s.userSessionLifeTime)
	token := jwt.NewWithClaims(jwt.SigningMethodRS256, jwt.MapClaims{
		"user_id": id,
		"nbf":     time.Now().Unix(),
		"exp":     expirationDate.Unix(),
	})

	tokenString, err := token.SignedString(s.keys.jwtPrivate)
	if err != nil {
		log.Println(err.Error())
		return ctx.JSON(http.StatusInternalServerError, map[string]string{"error": "Error while sighning token"})
	}

	response := api.AuthSuccess{
		Token:          &tokenString,
		ExpirationDate: &expirationDate,
	}
	return ctx.JSON(http.StatusOK, response)
}

// Register a new user
// (POST /register)
func (s *ServerHandler) PostRegister(ctx echo.Context) error {
	log.Println("Register request")
	registerRequest := api.PostRegisterJSONRequestBody{}
	if err := ctx.Bind(&registerRequest); err != nil {
		log.Println(err.Error())
		return ctx.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request payload"})
	}

	hPassword := hashPassword(*registerRequest.Password, *registerRequest.Username)
	id, err := s.db.addNewUser(*registerRequest.Username, hPassword)
	if err != nil {
		log.Println(err.Error())
		return ctx.JSON(http.StatusConflict, map[string]string{"error": "User already exists"})
	}

	expirationDate := time.Now().Add(s.userSessionLifeTime)
	token := jwt.NewWithClaims(jwt.SigningMethodRS256, jwt.MapClaims{
		"user_id": id,
		"nbf":     time.Now().Unix(),
		"exp":     expirationDate.Unix(),
	})

	tokenString, err := token.SignedString(s.keys.jwtPrivate)
	if err != nil {
		log.Println(err.Error())
		return ctx.JSON(http.StatusInternalServerError, map[string]string{"error": "Error while sighning token"})
	}

	response := api.AuthSuccess{
		Token:          &tokenString,
		ExpirationDate: &expirationDate,
	}
	return ctx.JSON(http.StatusOK, response)
}

// Update user data
// (PUT /update)
func (s *ServerHandler) PutUpdate(ctx echo.Context) error {
	log.Println("Put request")

	userId, err := checkAuth(ctx, s.keys)
	if err != nil {
		return err
	}

	updateRequest := api.PutUpdateJSONRequestBody{}
	if err = ctx.Bind(&updateRequest); err != nil {
		log.Println(err.Error())
		return ctx.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request payload"})
	}

	if err = s.db.updateUser(int(userId), updateRequest); err != nil {
		log.Println(err.Error())
		return ctx.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request payload"})
	}

	return ctx.JSON(http.StatusOK, map[string]string{"message": "User data updated successfully"})
}
