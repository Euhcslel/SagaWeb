package docgen

import (
	"bytes"
	"fmt"
	"time"

	"github.com/Euhcslel/SagaWeb/internal/types"
)

// Собирает и генерирует PDF коммерческого предложения для клиента на основе данных заказа
func GenerateOfferForClient(order *types.Order) ([]byte, error) {
	d := newDoc()

	d.addTitle("КОММЕРЧЕСКОЕ ПРЕДЛОЖЕНИЕ", 15)
	d.addLine("от " + time.Now().Format("02.01.2006") + " г.")
	d.pdf.Ln(4)

	for i, g := range order.Gates {
		d.drawGateSection(i+1, g)
		d.pdf.Ln(4)
	}

	if len(order.Products) > 0 {
		startNo := len(order.Gates) + 1
		d.drawProductsSection(startNo, order.Products)
		d.pdf.Ln(4)
	}

	d.drawGrandTotal(order)

	d.pdf.Ln(4)
	d.addLine("КП действительно 14 дней")

	if err := d.pdf.Error(); err != nil {
		return nil, fmt.Errorf("ошибка построения PDF: %w", err)
	}
	var buf bytes.Buffer
	if err := d.pdf.Output(&buf); err != nil {
		return nil, fmt.Errorf("ошибка вывода PDF: %w", err)
	}
	return buf.Bytes(), nil
}

// Собирает и генерирует PDF коммерческого предложения для дилера.
// В шапку добавляется информация о поставщике (компания и адрес дилера).
func GenerateOfferForDealer(order *types.Order, supplierName, supplierAddress string) ([]byte, error) {
	d := newDoc()

	d.addTitle("КОММЕРЧЕСКОЕ ПРЕДЛОЖЕНИЕ", 15)
	d.addLine("от " + time.Now().Format("02.01.2006") + " г.")
	d.pdf.Ln(2)
	d.addLine("Поставщик: " + supplierName)
	if supplierAddress != "" {
		d.addLine("Адрес: " + supplierAddress)
	}
	d.pdf.Ln(4)

	for i, g := range order.Gates {
		d.drawGateSection(i+1, g)
		d.pdf.Ln(4)
	}

	if len(order.Products) > 0 {
		startNo := len(order.Gates) + 1
		d.drawProductsSection(startNo, order.Products)
		d.pdf.Ln(4)
	}

	d.drawGrandTotal(order)

	d.pdf.Ln(4)
	d.addLine("КП действительно 14 дней")

	if err := d.pdf.Error(); err != nil {
		return nil, fmt.Errorf("ошибка построения PDF: %w", err)
	}
	var buf bytes.Buffer
	if err := d.pdf.Output(&buf); err != nil {
		return nil, fmt.Errorf("ошибка вывода PDF: %w", err)
	}
	return buf.Bytes(), nil
}
