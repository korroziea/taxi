package user

import "github.com/korroziea/taxi/user-service/internal/domain"

type signUpReq struct {
	FirstName string `json:"first_name" binding:"required"`
	Email     string `json:"email" binding:"required"`
	Phone     string `json:"phone", binding:"required"`
	Password  string `json:"password" binding:"required"`
}

func (r signUpReq) toDomain() domain.SignUpUser {
	user := domain.SignUpUser{
		FirstName: r.FirstName,
		Email:     r.Email,
		Phone:     r.Phone,
		Password:  r.Password,
	}

	return user
}

type signInReq struct {
	Phone    string `json:"phone", binding:"required"`
	Password string `json:"password" binding:"required"`
}

func (r signInReq) toDomain() domain.SignInUser {
	user := domain.SignInUser{
		Phone:    r.Phone,
		Password: r.Password,
	}

	return user
}

type signInResp struct {
	Token string `json:"token"`
}

type profileReq struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (r profileReq) toDomain() domain.ProfileUser {
	user := domain.ProfileUser{
		Email:    r.Email,
		Password: r.Password,
	}

	return user
}
