package service

import (
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

type User struct {
	User     string `bson:"user"`
	Password string `bson:"password"`
	Role     string `bson:"role"`
	DB       string `bson:"db"`
}

//returns a new User
func NewUser(username string, password string, role string) (*User, error) {
	hashedPassword, err := HashPassword(password)
	if err != nil {
		return nil, err
	}
	user := &User{
		User:     username,
		Password: hashedPassword,
		Role:     role,
	}
	return user, nil

}
func HashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("cannot hash password: %w", err)
	}
	return string(hashedPassword), nil
}
func IsCorrectPassword(account *User, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(account.Password), []byte(password))
	return err == nil
}
func (user *User) Clone() *User {

	return &User{
		User:     user.User,
		Role:     user.Role,
		Password: user.Password,
	}
}
