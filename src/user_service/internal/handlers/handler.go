package handlers

import (
	"crypto/md5"
	"log"
	"net/http"
	"slices"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"

	"user_service/api"
)

type ServerHandler struct {
	db                  dbWrapper
	keys                rsaKeys
	userSessionLifeTime time.Duration
}

func hashPassword(password, salt string) (passwordHash [16]byte) {
	hash := md5.New()
	hash.Write([]byte(password + salt))
	copy(passwordHash[:], hash.Sum(nil)[:16])
	return
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

	userPassword, err := s.db.getUserPassword(*loginRequest.Username)
	if err != nil {
		log.Println(err.Error())
		return ctx.JSON(http.StatusInternalServerError, map[string]string{"error": "User not found"})
	}

	if hPassword := hashPassword(*loginRequest.Password, *loginRequest.Username); slices.Compare(userPassword, hPassword[:]) != 0 {
		return ctx.JSON(http.StatusUnauthorized, map[string]string{"error": "Incorrect login/password"})
	}

	expirationDate := time.Now().Add(s.userSessionLifeTime)
	token := jwt.NewWithClaims(jwt.SigningMethodRS256, jwt.MapClaims{
		"username": loginRequest.Username,
		"nbf":      time.Now().Unix(),
		"exp":      expirationDate.Unix(),
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
	if err := s.db.addNewUser(*registerRequest.Username, hPassword); err != nil {
		log.Println(err.Error())
		return ctx.JSON(http.StatusConflict, map[string]string{"error": "User already exists"})
	}

	expirationDate := time.Now().Add(s.userSessionLifeTime)
	token := jwt.NewWithClaims(jwt.SigningMethodRS256, jwt.MapClaims{
		"username": registerRequest.Username,
		"nbf":      time.Now().Unix(),
		"exp":      expirationDate.Unix(),
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

	authHeader := ctx.Request().Header.Get("Authorization")
	parts := strings.SplitN(authHeader, " ", 2)
	if len(parts) != 2 || parts[0] != "Bearer" {
		return ctx.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid authorization header format"})
	}
	tokenString := parts[1]

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (any, error) { return s.keys.jwtPublic, nil })
	if err != nil {
		log.Println(err.Error())
		return ctx.JSON(http.StatusBadRequest, map[string]string{"error": "Error parsing token"})
	}

	if !token.Valid {
		return ctx.JSON(http.StatusUnauthorized, map[string]string{"error": "Unauthorized, invalid or expired token"})
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return ctx.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid claims in token"})
	}

	log.Println(claims)

	expTime, err := claims.GetExpirationTime()
	if err != nil || expTime == nil {
		return ctx.JSON(http.StatusBadRequest, map[string]string{"error": "Error parsing exp field"})
	}

	if !ok || time.Now().After(expTime.Time) {
		return ctx.JSON(http.StatusUnauthorized, map[string]string{"error": "Unauthorized, invalid or expired token"})
	}

	username, ok := claims["username"].(string)
	if !ok {
		return ctx.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid claims in token"})
	}

	updateRequest := api.PutUpdateJSONRequestBody{}
	if err = ctx.Bind(&updateRequest); err != nil {
		log.Println(err.Error())
		return ctx.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request payload"})
	}

	if err = s.db.updateUser(username, updateRequest); err != nil {
		log.Println(err.Error())
		return ctx.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request payload"})
	}

	return ctx.JSON(http.StatusOK, map[string]string{"message": "User data updated successfully"})
}
