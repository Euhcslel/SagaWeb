package models

import (
	"errors"
	"project/pkg/proto_files"
)

type OrderStatus string

const (
	OrderStatusNew       OrderStatus = "new"
	OrderStatusPaid      OrderStatus = "paid"
	OrderStatusCancelled OrderStatus = "cancelled"
)

func GetAllOrderStatuses() []OrderStatus {
	return []OrderStatus{
		OrderStatusNew,
		OrderStatusPaid,
		OrderStatusCancelled,
	}
}

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

type RegRequestStatus string

const (
	RegRequestStatusPending  RegRequestStatus = "pending"
	RegRequestStatusApproved RegRequestStatus = "approved"
	RegRequestStatusRejected RegRequestStatus = "rejected"
)

func GetAllRegRequestStatuses() []RegRequestStatus {
	return []RegRequestStatus{
		RegRequestStatusPending,
		RegRequestStatusApproved,
		RegRequestStatusRejected,
	}
}

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

type GateType string

const (
	GateTypeInd GateType = "ind"
	GateTypeRes GateType = "res"
)

func GetAllGateTypes() []GateType {
	return []GateType{
		GateTypeInd,
		GateTypeRes,
	}
}

func DetermineGateType(name string) (GateType, error) {
	t := GateType(name)

	switch t {
	case GateTypeRes, GateTypeInd:
		return t, nil
	default:
		return "", errors.New("Такой тип ворот не существует")
	}
}

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

func GateTypeFromProto(t proto_files.GateType) (GateType, error) {
	switch t {
	case proto_files.GateType_GATE_TYPE_IND:
		return GateTypeInd, nil
	case proto_files.GateType_GATE_TYPE_RES:
		return GateTypeRes, nil
	default:
		return "", errors.New("Такой тип ворот не существует")
	}
}
