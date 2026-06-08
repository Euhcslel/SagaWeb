package docgen

import (
	"bytes"
	"fmt"
	"time"

	"github.com/Euhcslel/SagaWeb/internal/types"
)

var ruMonths = [...]string{
	"января", "февраля", "марта", "апреля", "мая", "июня",
	"июля", "августа", "сентября", "октября", "ноября", "декабря",
}

func formatDateRu(t time.Time) string {
	return fmt.Sprintf("«%d» %s %d г.", t.Day(), ruMonths[t.Month()-1], t.Year())
}

// CustomerInfo содержит данные о заказчике для приложения к договору
type CustomerInfo struct {
	OrganizationName string // Наименование организации или ФИО
	ContactPerson    string // Контактное лицо
	Phone            string // Телефон
}

// GenerateAppendix генерирует PDF приложения к договору для конкретного заказа
func GenerateAppendix(
	order *types.Order,
	customer CustomerInfo,
	appendixNumber int,
	contractNumber string,
	contractDate time.Time,
) ([]byte, error) {
	d := newDoc()

	// Шапка
	d.pdf.SetFont("DejaVu", "B", 11)
	d.pdf.CellFormat(0, 6, fmt.Sprintf("Приложение № %d", appendixNumber), "", 1, "L", false, 0, "")
	d.pdf.SetFont("DejaVu", "", 9)
	d.pdf.CellFormat(0, 5, fmt.Sprintf("к Договору № %s от %s", contractNumber, formatDateRu(contractDate)), "", 1, "L", false, 0, "")
	d.pdf.Ln(3)

	// Поставщик
	d.pdf.SetFont("DejaVu", "", 9)
	d.pdf.CellFormat(0, 5, "Поставщик:", "", 1, "L", false, 0, "")
	d.pdf.SetFont("DejaVu", "B", 9)
	d.pdf.CellFormat(0, 5, "ООО «САГА»", "", 1, "L", false, 0, "")
	d.pdf.SetFont("DejaVu", "", 9)
	d.pdf.CellFormat(0, 5, "ИНН/КПП 6671287803/667801001", "", 1, "L", false, 0, "")
	d.pdf.Ln(5)

	// Центрированный заголовок приложения
	d.pdf.SetFont("DejaVu", "B", 12)
	d.pdf.CellFormat(0, 7, fmt.Sprintf("Приложение № %d от %sг.", appendixNumber, time.Now().Format("02.01.2006")), "", 1, "C", false, 0, "")
	d.pdf.Ln(4)

	// Данные заказчика
	d.drawCustomerBlock(customer)
	d.pdf.Ln(5)

	// Ворота
	for i, g := range order.Gates {
		d.drawGateSection(i+1, g)
		d.pdf.Ln(4)
	}

	// Дополнительные товары
	if len(order.Products) > 0 {
		startNo := len(order.Gates) + 1
		d.drawProductsSection(startNo, order.Products)
		d.pdf.Ln(4)
	}

	// Итог
	d.drawGrandTotal(order)

	d.pdf.Ln(10)
	d.drawAppendixSignatures(customer)

	if err := d.pdf.Error(); err != nil {
		return nil, fmt.Errorf("ошибка построения PDF: %w", err)
	}
	var buf bytes.Buffer
	if err := d.pdf.Output(&buf); err != nil {
		return nil, fmt.Errorf("ошибка вывода PDF: %w", err)
	}
	return buf.Bytes(), nil
}

func (d *Doc) drawCustomerBlock(c CustomerInfo) {
	const labelW = 72.0
	valueW := contentWidth - labelW

	row := func(label, value string) {
		d.pdf.SetFont("DejaVu", "B", 9)
		d.pdf.CellFormat(labelW, minRowHeight, label, "1", 0, "L", false, 0, "")
		d.pdf.SetFont("DejaVu", "", 9)
		d.pdf.CellFormat(valueW, minRowHeight, value, "1", 1, "L", false, 0, "")
	}

	row("Наименование организации/ФИО:", c.OrganizationName)
	row("Контактное лицо:", c.ContactPerson)
	row("Тел.:", c.Phone)
}

func (d *Doc) drawAppendixSignatures(c CustomerInfo) {
	sigW := contentWidth / 2
	d.pdf.SetFont("DejaVu", "", 9)
	d.pdf.CellFormat(sigW, 6, "Поставщик __________________", "", 0, "L", false, 0, "")
	buyer := "Покупатель __________________"
	if c.ContactPerson != "" {
		buyer = "Покупатель ____________ /" + c.ContactPerson
	}
	d.pdf.CellFormat(sigW, 6, buyer, "", 1, "L", false, 0, "")
	d.pdf.Ln(8)
	d.addLine("Гарантия на оборудование - 12 месяцев, на автоматику 24 месяца.")
}
