package handlers

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"os"
	"user_service/api"

	_ "github.com/lib/pq"
)

type dbWrapper struct {
	*sql.DB
}

func NewDBWrapper() (dbWrapper, error) {
	username := os.Getenv("USER_DB_USER")
	password := os.Getenv("USER_DB_PASSWORD")
	host := "postgres"
	port := os.Getenv("USER_DB_PORT")
	dbname := os.Getenv("USER_DB")

	dbAddress := fmt.Sprintf("postgres://%s:%s@%s:%v/%s?sslmode=disable", username, password, host, port, dbname)
	db, err := sql.Open("postgres", dbAddress)
	if err != nil {
		return dbWrapper{}, fmt.Errorf("error opening db connection: %v", err)
	}

	err = db.Ping()
	if err != nil {
		return dbWrapper{}, err
	}

	return dbWrapper{db}, nil
}

func (w dbWrapper) getUserLogin(userID int) (login string, err error) {
	log.Printf("get user login by id: %d\n", userID)

	err = w.QueryRow(
		`SELECT userlogin 
		FROM UserCredentials 
		WHERE id = $1`,
		userID).Scan(&login)
	return
}

func (w dbWrapper) addNewUser(userLogin string, hashPassword [16]byte) (id int, err error) {
	log.Println("db adding new User")

	var count int
	err = w.QueryRow(
		`SELECT COUNT(*) 
		FROM UserCredentials 
		WHERE userlogin = $1`,
		userLogin).Scan(&count)
	if err != nil {
		return
	}

	if count > 0 {
		err = errors.New("user already exist")
		return
	}

	err = w.QueryRow(
		`INSERT INTO UserCredentials 
		(userlogin, userpassword) VALUES ($1, $2) 
		RETURNING id`,
		userLogin, hashPassword[:]).Scan(&id)
	log.Printf("New user id: %d\n", id)
	return
}

func (w dbWrapper) getUserPasswordId(userLogin string) ([]byte, int, error) {
	log.Println("db getting user password")
	var id int
	password := make([]byte, 16)
	err := w.QueryRow(
		`SELECT userpassword, id 
		FROM UserCredentials 
		WHERE userlogin = $1`,
		userLogin).Scan(&password, &id)

	if err != nil {
		return []byte{}, id, err
	}
	return password, id, nil
}

func (w dbWrapper) updateUser(userId int, userData api.PutUpdateJSONRequestBody) error {
	log.Println("db updateUser")

	query, err := w.Prepare(`
	INSERT INTO UserProfile (id, birthdate, email, first_name, second_name, phone_number)
	VALUES ($1, $2, $3, $4, $5, $6)
	ON CONFLICT (id)
	DO UPDATE SET
		birthdate = COALESCE(EXCLUDED.birthdate, UserProfile.birthdate),
		email = COALESCE(EXCLUDED.email, UserProfile.email),
		first_name = COALESCE(EXCLUDED.first_name, UserProfile.first_name),
		second_name = COALESCE(EXCLUDED.second_name, UserProfile.second_name),
		phone_number = COALESCE(EXCLUDED.phone_number, UserProfile.phone_number);`)
	if err != nil {
		return err
	}
	defer query.Close()

	_, err = query.Exec(userId, userData.DateOfBirth.Time, userData.Email, userData.Name,
		userData.Surname, userData.PhoneNumber)
	return err
}
