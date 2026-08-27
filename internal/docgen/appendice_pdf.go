package docgen

import (
	"fmt"
	"time"

	"github.com/Euhcslel/SagaWeb/internal/types"
)

// GenerateAppendice генерирует PDF приложения к договору для конкретного заказа.
func GenerateAppendice(
	order *types.Order,
	customer CustomerInfo,
	appendiceNumber int,
	contractNumber string,
	contractDate time.Time,
) ([]byte, error) {
	d := newDoc()

	d.pdf.SetFont("DejaVu", "B", 11)
	d.pdf.CellFormat(0, 6, fmt.Sprintf("Приложение № %d", appendiceNumber), "", 1, "L", false, 0, "")
	d.pdf.SetFont("DejaVu", "", 9)
	d.pdf.CellFormat(0, 5, fmt.Sprintf("к Договору № %s от %s", contractNumber, formatDateRu(contractDate)), "", 1, "L", false, 0, "")
	d.pdf.Ln(3)

	d.pdf.SetFont("DejaVu", "", 9)
	d.pdf.CellFormat(0, 5, "Поставщик:", "", 1, "L", false, 0, "")
	d.pdf.SetFont("DejaVu", "B", 9)
	d.pdf.CellFormat(0, 5, Cfg.CompanyName, "", 1, "L", false, 0, "")
	d.pdf.SetFont("DejaVu", "", 9)
	if Cfg.CompanyINN != "" || Cfg.CompanyKPP != "" {
		d.pdf.CellFormat(0, 5, "ИНН/КПП "+Cfg.CompanyINN+"/"+Cfg.CompanyKPP, "", 1, "L", false, 0, "")
	}
	d.pdf.Ln(3)

	d.pdf.SetFont("DejaVu", "B", 12)
	d.pdf.CellFormat(0, 7, fmt.Sprintf("Приложение № %d от %s г.", appendiceNumber, time.Now().Format("02.01.2006")), "", 1, "C", false, 0, "")
	d.pdf.Ln(4)

	d.drawCustomerBlock(customer)
	d.pdf.Ln(5)

	for i, g := range order.Gates {
		d.drawGateSectionSingleWholesale(i+1, g)
		d.pdf.Ln(4)
	}

	if len(order.Products) > 0 {
		d.drawProductsSectionSingleWholesale(len(order.Gates)+1, order.Products)
		d.pdf.Ln(4)
	}

	d.drawGrandTotalSingleWholesale(order)
	d.pdf.Ln(10)
	d.drawAppendiceSignatures(customer)

	return renderPDF(d)
}

// drawCustomerBlock формирует блок с информацией о клиенте.
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

// drawAppendiceSignatures формирует места для подписей в документе и информацию о гарантии.
func (d *Doc) drawAppendiceSignatures(c CustomerInfo) {
	sigW := contentWidth / 2
	d.pdf.SetFont("DejaVu", "", 9)
	supplier := "Поставщик __________________"
	if Cfg.Representative != "" {
		supplier = "Поставщик ____________ /" + Cfg.Representative
	}
	d.pdf.CellFormat(sigW, 6, supplier, "", 0, "L", false, 0, "")
	buyer := "Покупатель __________________"
	if c.ContactPerson != "" {
		buyer = "Покупатель ____________ /" + c.ContactPerson
	}
	d.pdf.CellFormat(sigW, 6, buyer, "", 1, "L", false, 0, "")
	d.pdf.Ln(8)
	if Cfg.WarrantyEquip != "" || Cfg.WarrantyAuto != "" {
		d.addLine(fmt.Sprintf("Гарантия на оборудование — %s, на автоматику — %s.", Cfg.WarrantyEquip, Cfg.WarrantyAuto))
	}
}
