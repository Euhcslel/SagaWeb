// Package order_gate_manual_drives предоставляет модель для работы с ручным приводом у ворот в заказе.
// Данные хранятся в таблице order_gate_manual_drives.

package order_gate_manual_drives

const TableNameOrderGateManualDrive = "order_gate_manual_drives"

type OrderGateManualDrive struct {
	OrderID     int64 `gorm:"column:order_id;primaryKey"`
	RowNumber   int64 `gorm:"column:row_number;primaryKey"`
	ChainLength int32 `gorm:"column:chain_length"`
}

func (*OrderGateManualDrive) TableName() string {
	return TableNameOrderGateManualDrive
}
