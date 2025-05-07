package domain

import (
	"time"

	"github.com/korroziea/taxi/user-service/pkg/utils"
)

type User struct {
	ID        string
	FirstName string
	Email     string
	Phone     string
	Password  string
	CreatedAt time.Time
	UpdatedAt time.Time
}

func GenUserID() (string, error) {
	id, err := utils.GenID()
	if err != nil {
		return "", err
	}

	return "user_" + id, nil
}

type SignUpUser struct {
	ID        string
	FirstName string
	Email     string
	Phone     string
	Password  string
}

type SignInUser struct {
	Phone    string
	Password string
}

type ProfileUser struct {
	ID       string
	Email    string
	Password string
}
