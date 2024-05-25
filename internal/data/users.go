package data

import (
	"context"
	"log"

	"github.com/vaidik-bajpai/hackernews/prisma/db"

	"github.com/vaidik-bajpai/hackernews/database"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
	Password string `json:"password"`
}

func (user *User) Create() {
	hash, err := HashPassword(user.Password)
	if err != nil {
		log.Fatal(err)
	}
	_, err = database.Db.User.CreateOne(
		db.User.Name.Set(user.Username),
		db.User.Password.Set(hash),
	).Exec(context.Background())
	if err != nil {
		log.Fatal(err)
	}
}

// HashPassword hashes given password
func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

// CheckPassword hash compares raw password with it's hashed values
func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func GetUserIdByUsername(username string) (int, error) {
	userID, err := database.Db.User.FindUnique(
		db.User.Name.Equals(username),
	).Exec(context.Background())

	if err != nil {
		log.Fatal(err)
	}

	return userID.ID, nil
}

func (user *User) Authenticate() bool {
	newUser, err := database.Db.User.FindUnique(
		db.User.ID.Equals(user.ID),
	).Exec(context.Background())

	if err != nil {
		log.Fatal(err)
	}

	return CheckPasswordHash(user.Password, newUser.Password)
}
