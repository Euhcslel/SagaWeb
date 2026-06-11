package docgen

import (
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

// CustomerInfo содержит данные о заказчике для приложения к договору.
type CustomerInfo struct {
	OrganizationName string // Наименование организации или ФИО
	ContactPerson    string // Контактное лицо
	Phone            string // Телефон
}

// GenerateAppendice генерирует PDF приложения к договору для конкретного заказа.
func GenerateAppendice(
	order *types.Order,
	company CompanyInfo,
	customer CustomerInfo,
	appendiceNumber int,
	contractNumber string,
	contractDate time.Time,
) ([]byte, error) {
	d := newDoc()

	// Шапка
	d.pdf.SetFont("DejaVu", "B", 11)
	d.pdf.CellFormat(0, 6, fmt.Sprintf("Приложение № %d", appendiceNumber), "", 1, "L", false, 0, "")
	d.pdf.SetFont("DejaVu", "", 9)
	d.pdf.CellFormat(0, 5, fmt.Sprintf("к Договору № %s от %s", contractNumber, formatDateRu(contractDate)), "", 1, "L", false, 0, "")
	d.pdf.Ln(3)

	// Поставщик
	d.pdf.SetFont("DejaVu", "", 9)
	d.pdf.CellFormat(0, 5, "Поставщик:", "", 1, "L", false, 0, "")
	d.pdf.SetFont("DejaVu", "B", 9)
	d.pdf.CellFormat(0, 5, company.Name, "", 1, "L", false, 0, "")
	d.pdf.SetFont("DejaVu", "", 9)
	if company.INN != "" || company.KPP != "" {
		d.pdf.CellFormat(0, 5, "ИНН/КПП "+company.INN+"/"+company.KPP, "", 1, "L", false, 0, "")
	}
	d.pdf.Ln(5)

	// Центрированный заголовок приложения
	d.pdf.SetFont("DejaVu", "B", 12)
	d.pdf.CellFormat(0, 7, fmt.Sprintf("Приложение № %d от %sг.", appendiceNumber, time.Now().Format("02.01.2006")), "", 1, "C", false, 0, "")
	d.pdf.Ln(4)

	// Данные заказчика
	d.drawCustomerBlock(customer)
	d.pdf.Ln(5)

	// Ворота
	for i, g := range order.Gates {
		d.drawGateSectionSingle(i+1, g)
		d.pdf.Ln(4)
	}

	// Дополнительные товары
	if len(order.Products) > 0 {
		d.drawProductsSectionSingle(len(order.Gates)+1, order.Products)
		d.pdf.Ln(4)
	}

	// Итог
	d.drawGrandTotalSingle(order)
	d.pdf.Ln(10)
	d.drawAppendiceSignatures(company, customer)

	return renderPDF(d)
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

func (d *Doc) drawAppendiceSignatures(company CompanyInfo, c CustomerInfo) {
	sigW := contentWidth / 2
	d.pdf.SetFont("DejaVu", "", 9)
	d.pdf.CellFormat(sigW, 6, "Поставщик __________________", "", 0, "L", false, 0, "")
	buyer := "Покупатель __________________"
	if c.ContactPerson != "" {
		buyer = "Покупатель ____________ /" + c.ContactPerson
	}
	d.pdf.CellFormat(sigW, 6, buyer, "", 1, "L", false, 0, "")
	d.pdf.Ln(8)
	if company.WarrantyEquipment != "" || company.WarrantyAutomation != "" {
		d.addLine(fmt.Sprintf("Гарантия на оборудование — %s, на автоматику %s.", company.WarrantyEquipment, company.WarrantyAutomation))
	}
}
