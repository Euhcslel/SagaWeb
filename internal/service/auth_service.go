package service

import (
	"project/internal/domain/dealers_reg_requests"
	"project/internal/domain/enums"
	"project/internal/repository"

	"golang.org/x/crypto/bcrypt"
)

func SignIn(username, password string) (int64, error) {
	user, err := repository.GetUserByUsername(username)
	if err != nil {
		return 0, err
	}

	userId := user.ID

	dbPassword := user.PasswordHash
	err = bcrypt.CompareHashAndPassword([]byte(dbPassword), []byte(password))
	if err != nil {
		return 0, err
	}
	return userId, nil
}

func SignUp(fullname, company, phone, email, password string) error {
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	request := dealers_reg_requests.DealerRegRequest{
		Company:      company,
		Fullname:     fullname,
		PhoneNumber:  phone,
		Email:        email,
		PasswordHash: passwordHash,
		Status:       enums.RegRequestStatusPending,
	}

	return repository.CreateDealerRegRequest(request)
}

func SignOut(sessionToken string) error {
	return repository.DeleteSession(sessionToken)
}
