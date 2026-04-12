package enums

import (
	"errors"
	"project/internal/generated"
)

// Тип Статус заказа
type OrderStatus string

const (
	OrderStatusNew       OrderStatus = "new"
	OrderStatusPaid      OrderStatus = "paid"
	OrderStatusCancelled OrderStatus = "cancelled"
)

// Функция для получения всех статусов заказа
func GetAllOrderStatuses() []OrderStatus {
	return []OrderStatus{
		OrderStatusNew,
		OrderStatusPaid,
		OrderStatusCancelled,
	}
}

// Функция для получения наименования статуса заказа
func (s OrderStatus) Label() string {
	switch s {
	case OrderStatusNew:
		return "Новый"
	case OrderStatusPaid:
		return "Оплачен"
	case OrderStatusCancelled:
		return "Отменён"
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
func GetDriveTypeFromProto(drive *generated.Drive) DriveType {
	driveType := drive.DriveType
	switch driveType.(type) {
	case *generated.Drive_Industrial:
		return IndDriveType
	case *generated.Drive_Manual:
		return ManualDriveType
	case *generated.Drive_Residential:
		return ResDriveType
	}
	return "unknown"
}
