package service

import (
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"net/smtp"
	"os"
	"path/filepath"
	"time"

	"github.com/Euhcslel/SagaWeb/internal/database"
	"github.com/Euhcslel/SagaWeb/internal/domain/dealer_contract"
	"github.com/Euhcslel/SagaWeb/internal/domain/dealer_manager_assignments"
	"github.com/Euhcslel/SagaWeb/internal/domain/dealers"
	"github.com/Euhcslel/SagaWeb/internal/domain/enums"
	"github.com/Euhcslel/SagaWeb/internal/domain/users"
	errs "github.com/Euhcslel/SagaWeb/internal/errors"
	"github.com/Euhcslel/SagaWeb/internal/repository"
	"github.com/Euhcslel/SagaWeb/internal/types"

	"golang.org/x/crypto/bcrypt"
)

type DealerWithContract struct {
	dealer_manager_assignments.DealerManagerAssignment
	Contract *dealer_contract.DealerContract
}

func GetUserDealers(user *users.User) ([]DealerWithContract, error) {
	if user.Role != enums.ManagerRole && user.Role != enums.AdminRole {
		return nil, errs.ErrForbidden
	}

	assignments, err := repository.GetUserDealers(database.DB, user.ID)
	if err != nil {
		return nil, err
	}

	dealerIDs := make([]int64, len(assignments))
	for i, a := range assignments {
		dealerIDs[i] = a.DealerID
	}

	contractMap, err := repository.GetDealerContractsByDealerIDs(database.DB, dealerIDs)
	if err != nil {
		return nil, err
	}

	result := make([]DealerWithContract, len(assignments))
	for i, a := range assignments {
		result[i] = DealerWithContract{
			DealerManagerAssignment: a,
			Contract:                contractMap[a.DealerID],
		}
	}
	return result, nil
}

func AttachContractToDealer(user *users.User, dealerID int64, contractNumber string, signedAt time.Time, file multipart.File, handler *multipart.FileHeader) error {
	if user.Role != enums.ManagerRole && user.Role != enums.AdminRole {
		return errs.ErrForbidden
	}

	dir := os.Getenv("CONTRACTS_DIRECTORY")

	existing, err := repository.GetDealerContract(database.DB, dealerID)
	if err != nil {
		return err
	}
	if existing != nil && existing.Path != nil && dir != "" {
		_ = os.Remove(filepath.Join(dir, *existing.Path))
	}

	contract := dealer_contract.DealerContract{
		DealerID:       dealerID,
		ContractNumber: contractNumber,
		SignedAt:       signedAt,
	}

	if file != nil && handler != nil {
		if dir == "" {
			return fmt.Errorf("переменная CONTRACTS_DIRECTORY не задана")
		}
		dst, err := os.OpenFile(filepath.Join(dir, handler.Filename), os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
		if err != nil {
			return err
		}
		defer dst.Close()
		if _, err = io.Copy(dst, file); err != nil {
			return err
		}
		contract.Path = &handler.Filename
	}

	return repository.UpsertDealerContract(database.DB, contract)
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

		address := ""
		if dealer.Address != nil {
			address = *dealer.Address
		}

		userInfo.Dealer = &DealerInfo{
			CompanyName: dealer.Company.Name,
			Address:     address,
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
