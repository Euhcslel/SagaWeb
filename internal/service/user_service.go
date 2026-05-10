package service

import (
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"net/smtp"
	"os"
	"time"

	"github.com/Euhcslel/SagaWeb/internal/database"
	"github.com/Euhcslel/SagaWeb/internal/domain/dealer_manager_assignments"
	"github.com/Euhcslel/SagaWeb/internal/domain/dealers"
	"github.com/Euhcslel/SagaWeb/internal/domain/enums"
	"github.com/Euhcslel/SagaWeb/internal/domain/users"
	errs "github.com/Euhcslel/SagaWeb/internal/errors"
	"github.com/Euhcslel/SagaWeb/internal/repository"
	"github.com/Euhcslel/SagaWeb/internal/types"

	"golang.org/x/crypto/bcrypt"
)

func GetUserDealers(user *users.User) ([]dealer_manager_assignments.DealerManagerAssignment, error) {
	role := user.Role

	if role == enums.ManagerRole {
		dealers, err := repository.GetUserDealers(database.DB, user.ID)
		if err != nil {
			return nil, err
		}
		return dealers, nil
	} else {
		return nil, errs.ErrForbidden
	}
}

type UserInfo struct {
	ID       int64
	Fullname string
	Email    string
	Phone    string
	Role     enums.Role
	IsDealer bool

	Dealer *DealerInfo
}

type DealerInfo struct {
	CompanyName string
	Address     string
}

func GetUserInfo(user *users.User) (*UserInfo, error) {
	userInfo := &UserInfo{
		ID:       user.ID,
		Fullname: user.Fullname,
		Email:    user.Email,
		Phone:    user.PhoneNumber,
		Role:     user.Role,
		IsDealer: user.Role == enums.DealerRole,
	}

	if userInfo.IsDealer {
		dealer, err := repository.GetDealerInfo(database.DB, user.ID)
		if err != nil {
			return nil, err
		}

		userInfo.Dealer = &DealerInfo{
			CompanyName: dealer.Company.Name,
			Address:     *dealer.Address,
		}
	}

	return userInfo, nil
}

func UpdateUserInfo(user *users.User, userInfo types.UpdatedUserInfo) error {
	tx := database.DB.Begin()
	if tx.Error != nil {
		return tx.Error
	}
	defer tx.Rollback()

	if err := repository.UpdateUserInfo(tx, user.ID, userInfo); err != nil {
		return err
	}

	if user.Role == enums.DealerRole {
		dealer, err := repository.GetDealerByID(tx, user.ID)
		if err != nil {
			return err
		}

		if err := repository.UpdateDealerInfo(tx, dealer.UserID, userInfo); err != nil {
			return err
		}

		if err := repository.UpdateCompanyInfo(tx, dealer.CompanyID, userInfo); err != nil {
			return err
		}
	}

	if err := tx.Commit().Error; err != nil {
		return err
	}

	return nil
}

func ConfirmDealerRegistrationRequest(manager *users.User, requestID int64) error {
	password, err := generatePassword(18)
	if err != nil {
		return err
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	tx := database.DB.Begin()

	request, err := repository.GetRegistrationRequestByID(tx, requestID)
	if err != nil {
		return err
	}

	newUser := users.User{
		Fullname:     request.Fullname,
		PhoneNumber:  request.PhoneNumber,
		Email:        request.Email,
		PasswordHash: hash,
		Role:         enums.DealerRole,
	}
	if err := repository.CreateUser(tx, &newUser); err != nil {
		return err
	}

	company, err := repository.GetOrCreateCompanyByName(tx, request.Company)
	if err != nil {
		return err
	}

	newDealer := dealers.Dealer{
		UserID:    newUser.ID,
		CompanyID: company.ID,
	}
	if err := repository.CreateDealer(tx, &newDealer); err != nil {
		return err
	}

	if err := repository.AttachDealerToManager(tx, newUser.ID, manager.ID); err != nil {
		return err
	}

	if err := repository.DeleteRegistrationRequest(tx, requestID); err != nil {
		return err
	}

	if err := sendPasswordEmail(newUser.Email, newUser.Fullname, password); err != nil {
		return err
	}

	if err := tx.Commit().Error; err != nil {
		return err
	}

	return nil
}

func generatePassword(byteLength int) (string, error) {
	b := make([]byte, byteLength)

	if _, err := rand.Read(b); err != nil {
		return "", err
	}

	return base64.RawURLEncoding.EncodeToString(b), nil
}

func sendPasswordEmail(toEmail, name, password string) error {
	smtpHost := os.Getenv("SMTP_HOST")
	smtpPort := os.Getenv("SMTP_PORT")
	smtpUser := os.Getenv("SMTP_USER")
	smtpPass := os.Getenv("SMTP_PASSWORD")
	fromEmail := os.Getenv("SMTP_FROM_EMAIL")

	if smtpHost == "" || smtpPort == "" || smtpUser == "" || smtpPass == "" || fromEmail == "" {
		return errors.New("smtp config is incomplete")
	}

	auth := smtp.PlainAuth("", smtpUser, smtpPass, smtpHost)

	subject := "Ваш доступ к системе"

	body := fmt.Sprintf(`Здравствуйте, %s!

	Для вас создан аккаунт.


	Логин: %s
	Временный пароль: %s

	С уважением, Saga Doors`,
		name, toEmail, password)

	message := buildEmailMessage(fromEmail, toEmail, subject, body)

	addr := smtpHost + ":" + smtpPort

	return smtp.SendMail(
		addr,
		auth,
		fromEmail,
		[]string{toEmail},
		[]byte(message),
	)
}

func buildEmailMessage(from, to, subject, body string) string {
	return fmt.Sprintf(
		"From: %s\r\n"+
			"To: %s\r\n"+
			"Subject: %s\r\n"+
			"MIME-Version: 1.0\r\n"+
			"Content-Type: text/plain; charset=\"UTF-8\"\r\n"+
			"Date: %s\r\n"+
			"\r\n"+
			"%s\r\n",
		from,
		to,
		subject,
		time.Now().Format(time.RFC1123Z),
		body,
	)
}

func RejectDealerRegistrationRequest(requestID int64) error {
	return repository.DeleteRegistrationRequest(database.DB, requestID)
}
