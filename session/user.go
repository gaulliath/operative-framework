package session

import (
	"errors"
	"golang.org/x/crypto/bcrypt"
	"strings"
)

type User struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Key      string `json:"api_key"`
}

func (s *Session) NewUser(username string, password string) (*User, error) {
	_, err := s.GetUser(username)
	if err == nil {
		return nil, errors.New("User already found")
	}

	cPassword, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	if err != nil {
		return nil, err
	}

	return &User{
		Username: username,
		Password: string(cPassword),
		Key:      "",
	}, nil
}

func (s *Session) GetUser(username string) (*User, error) {
	for _, user := range s.Users {
		if strings.ToLower(user.Username) == strings.ToLower(username) {
			return user, nil
		}
	}

	return nil, errors.New("User not found in current session")
}
