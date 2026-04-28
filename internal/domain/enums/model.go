package enums

import (
	"errors"
	errs "project/internal/errors"
	"project/internal/generated"
)

// Тип Статус заказа
type OrderStatus string

const (
	OrderStatusPending   OrderStatus = "pending"
	OrderStatusPaid      OrderStatus = "paid"
	OrderStatusCancelled OrderStatus = "cancelled"
	OrderStatusConfirmed OrderStatus = "confirmed"
	OrderStatusDone      OrderStatus = "done"
)

// Функция для получения статуса заказа из proto-типа
func GetOrderStatusFromProto(status *generated.OrderStatus) (OrderStatus, error) {
	switch *status {
	case generated.OrderStatus_ORDER_STATUS_CANCELLED:
		return OrderStatusCancelled, nil
	case generated.OrderStatus_ORDER_STATUS_DONE:
		return OrderStatusDone, nil
	case generated.OrderStatus_ORDER_STATUS_PENDING:
		return OrderStatusPending, nil
	case generated.OrderStatus_ORDER_STATUS_PAID:
		return OrderStatusPaid, nil
	case generated.OrderStatus_ORDER_STATUS_CONFIRMED:
		return OrderStatusConfirmed, nil
	}

	return "", errs.ErrInvalidOrderStatus
}

// Функция для получения всех статусов заказа
func GetAllOrderStatuses() []OrderStatus {
	return []OrderStatus{
		OrderStatusPending,
		OrderStatusPaid,
		OrderStatusCancelled,
		OrderStatusConfirmed,
		OrderStatusDone,
	}
}

// Функция для получения наименования статуса заказа
func (s OrderStatus) Label() string {
	switch s {
	case OrderStatusPending:
		return "Ожидает подтверждения"
	case OrderStatusPaid:
		return "Оплачен"
	case OrderStatusCancelled:
		return "Отменён"
	case OrderStatusDone:
		return "Завершен"
	case OrderStatusConfirmed:
		return "Подтвержден"
	default:
		return "Неизвестно"
	}
}

// Тип Статус заявки на регистрацию
type RegRequestStatus string

const (
	RegRequestStatusPending  RegRequestStatus = "pending"
	RegRequestStatusApproved RegRequestStatus = "approved"
	RegRequestStatusRejected RegRequestStatus = "rejected"
)

// Функция для получения всех статусов заявки на регистрацию
func GetAllRegRequestStatuses() []RegRequestStatus {
	return []RegRequestStatus{
		RegRequestStatusPending,
		RegRequestStatusApproved,
		RegRequestStatusRejected,
	}
}

// Функция для получения наименования статуса заявки на регистрацию
func (s RegRequestStatus) Label() string {
	switch s {
	case RegRequestStatusPending:
		return "Ожидает подтверждения"
	case RegRequestStatusApproved:
		return "Подтверждён"
	case RegRequestStatusRejected:
		return "Отклонён"
	default:
		return "Неизвестно"
	}
}

// Тип Тип ворот
type GateType string

const (
	GateTypeInd GateType = "ind"
	GateTypeRes GateType = "res"
)

// Функция для получения всех типов ворот
func GetAllGateTypes() []GateType {
	return []GateType{
		GateTypeInd,
		GateTypeRes,
	}
}

// Функция для преобразования строки в тип ворот
func DetermineGateType(name string) (GateType, error) {
	t := GateType(name)

	switch t {
	case GateTypeRes, GateTypeInd:
		return t, nil
	default:
		return "", errors.New("Такой тип ворот не существует")
	}
}

// Функция для получения наименования типа ворот
func (t GateType) Label() string {
	switch t {
	case GateTypeInd:
		return "Промышленные ворота"
	case GateTypeRes:
		return "Бытовые ворота"
	default:
		return "Неопределено"
	}
}

// Функция для получения типа ворот из proto-типа
func GateTypeFromProto(t generated.GateType) (GateType, error) {
	switch t {
	case generated.GateType_GATE_TYPE_IND:
		return GateTypeInd, nil
	case generated.GateType_GATE_TYPE_RES:
		return GateTypeRes, nil
	default:
		return "", errors.New("Такой тип ворот не существует")
	}
}

// Тип Тип привода
type DriveType string

const (
	IndDriveType    DriveType = "industrial"
	ResDriveType    DriveType = "residential"
	ManualDriveType DriveType = "manual"
)

// Функция для получения всех типов привода
func GetAllDriveTypes() []DriveType {
	return []DriveType{
		IndDriveType,
		ResDriveType,
		ManualDriveType,
	}
}

// Функция для преобразования строки в тип привода
func DetermineDriveType(name string) (DriveType, error) {
	t := DriveType(name)

	switch t {
	case IndDriveType, ResDriveType, ManualDriveType:
		return t, nil
	default:
		return "", errors.New("Такой тип привода не существует")
	}
}

// Функция для получения наименования типа привода
func (t DriveType) Label() string {
	switch t {
	case IndDriveType:
		return "Промышленный привод"
	case ResDriveType:
		return "Бытовой привод"
	case ManualDriveType:
		return "Ручной привод"
	default:
		return "Неопределено"
	}
}

// Функция для преобразования типа привода из proto-типа в enum
func GetDriveTypeFromProto(drive *generated.Drive) (DriveType, error) {
	driveType := drive.DriveType
	switch driveType.(type) {
	case *generated.Drive_Industrial:
		return IndDriveType, nil
	case *generated.Drive_Manual:
		return ManualDriveType, nil
	case *generated.Drive_Residential:
		return ResDriveType, nil
	}
	return "", errors.New("Такой тип привода не существует")
}

// Тип Роль
type Role string

const (
	ClientRole      Role = "client"
	ManagerRole     Role = "manager"
	DealerRole      Role = "dealer"
	AdminRole       Role = "admin"
	LogisticianRole Role = "logistician"
)

// Тип Тип документа (не используется в бд)
type DocumentType string

const (
	ContractDocumentType = "contract"
	OfferDocumentType    = "offer"
	BillDocumentType     = "bill"
	OtherDocumentType    = "other"
)

func GetDocumentTypeFromString(docType string) (DocumentType, error) {
	switch docType {
	case ContractDocumentType:
		return ContractDocumentType, nil
	case OfferDocumentType:
		return OfferDocumentType, nil
	case BillDocumentType:
		return BillDocumentType, nil
	case OtherDocumentType:
		return OtherDocumentType, nil
	default:
		return "", errs.ErrInvalidDocumentType
	}
}
