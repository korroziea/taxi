package driver

import "github.com/korroziea/taxi/driver-service/internal/domain"

type signUpReq struct {
	FirstName string `json:"first_name"`
	Phone     string `json:"phone"`
	Email     string `json:"email"`
	Password  string `json:"password"`
}

func (r signUpReq) toDomain() domain.SignUpDriver {
	driver := domain.SignUpDriver{
		FirstName: r.FirstName,
		Phone:     r.Phone,
		Email:     r.Email,
		Password:  r.Password,
	}

	return driver
}

type signInReq struct {
	Phone    string `json:"phone"`
	Password string `json:"password"`
}

func (r signInReq) toDomain() domain.SignInDriver {
	driver := domain.SignInDriver{
		Phone:    r.Phone,
		Password: r.Password,
	}

	return driver
}
