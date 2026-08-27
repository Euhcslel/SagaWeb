package service

import (
	"errors"
	"strings"
	"github.com/Euhcslel/SagaWeb/internal/domain/dealer_registration_requests"
	"github.com/Euhcslel/SagaWeb/internal/domain/enums"
	errs "github.com/Euhcslel/SagaWeb/internal/errors"
	"github.com/Euhcslel/SagaWeb/internal/repository"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// SignIn проверяет учетные данные пользователя и возвращает его идентификатор, если они верны.
func SignIn(email, password string) (int64, error) {
	user, err := repository.GetUserByEmail(email)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return 0, errs.ErrInvalidCredentials
	} else if err != nil {
		return 0, err
	}

	userID := user.ID

	dbPassword := user.PasswordHash
	err = bcrypt.CompareHashAndPassword(dbPassword, []byte(password))
	if err != nil {
		return 0, errs.ErrInvalidCredentials
	}
	return userID, nil
}

// SignUp создает новую заявку на регистрацию дилера с указанными данными.
func SignUp(fullname, company, phone, email string) error {
	if strings.TrimSpace(fullname) == "" || strings.TrimSpace(email) == "" || strings.TrimSpace(phone) == "" {
		return errs.ErrBadRequest
	}

	request := dealer_registration_requests.DealerRegistrationRequest{
		Company:     company,
		Fullname:    fullname,
		PhoneNumber: phone,
		Email:       email,
		Status:      enums.RegistrationRequestStatusPending,
	}

	return repository.CreateDealerRegistrationRequest(request)
}

// SignOut удаляет сессию пользователя по указанному токену сессии.
func SignOut(sessionToken string) error {
	return repository.DeleteSession(sessionToken)
}

// GetDealersRegistrationRequests возвращает все заявки на регистрацию дилеров.
func GetDealersRegistrationRequests() ([]dealer_registration_requests.DealerRegistrationRequest, error) {
	return repository.GetAllDealerRegistrationRequests()
}
