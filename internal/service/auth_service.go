package service

import (
	"errors"
	"github.com/Euhcslel/SagaWeb/internal/domain/dealer_registration_requests"
	"github.com/Euhcslel/SagaWeb/internal/domain/enums"
	errs "github.com/Euhcslel/SagaWeb/internal/errors"
	"github.com/Euhcslel/SagaWeb/internal/repository"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func SignIn(email, password string) (int64, error) {
	user, err := repository.GetUserByEmail(email)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return 0, errs.ErrInvalidCredentials
	} else if err != nil {
		return 0, err
	}

	userID := user.ID

	dbPassword := user.PasswordHash
	err = bcrypt.CompareHashAndPassword([]byte(dbPassword), []byte(password))
	if err != nil {
		return 0, errs.ErrInvalidCredentials
	}
	return userID, nil
}

func SignUp(fullname, company, phone, email string) error {
	request := dealer_registration_requests.DealerRegistrationRequest{
		Company:     company,
		Fullname:    fullname,
		PhoneNumber: phone,
		Email:       email,
		Status:      enums.RegistrationRequestStatusPending,
	}

	return repository.CreateDealerRegistrationRequest(request)
}

func SignOut(sessionToken string) error {
	return repository.DeleteSession(sessionToken)
}

func GetDealersRegistrationRequests() ([]dealer_registration_requests.DealerRegistrationRequest, error) {
	return repository.GetAllDealerRegistrationRequests()
}
