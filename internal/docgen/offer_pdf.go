package docgen

import (
	"bytes"
	"fmt"
	"time"

	"github.com/Euhcslel/SagaWeb/internal/types"
)

// GenerateClientOffer генерирует КП для неавторизованного клиента (без поставщика и покупателя).
func GenerateClientOffer(order *types.Order, offerValidityDays int) ([]byte, error) {
	d := newDoc()

	d.addTitle("КОММЕРЧЕСКОЕ ПРЕДЛОЖЕНИЕ", 14)
	d.addLine("от " + time.Now().Format("02.01.2006") + " г.")
	d.pdf.Ln(6)

	for i, g := range order.Gates {
		d.drawGateSectionSingle(i+1, g)
		d.pdf.Ln(4)
	}
	if len(order.Products) > 0 {
		d.drawProductsSectionSingle(len(order.Gates)+1, order.Products)
		d.pdf.Ln(4)
	}

	d.drawGrandTotalSingle(order)
	d.pdf.Ln(6)
	d.addLine(fmt.Sprintf("КП действительно %d дней", offerValidityDays))

	return renderPDF(d)
}

// GenerateDealerOfferForClient генерирует КП дилера для покупателя (только розничные цены).
// orderID — номер заказа; если nil, номер в заголовок не включается.
func GenerateDealerOfferForClient(
	order *types.Order,
	orderID *int64,
	supplier SupplierInfo,
	clientName string,
) ([]byte, error) {
	d := newDoc()

	drawOfferHeader(d, orderID)
	drawSupplierBuyerBlock(d, supplier, clientName)
	d.pdf.Ln(4)

	for i, g := range order.Gates {
		d.drawGateSectionSingle(i+1, g)
		d.pdf.Ln(4)
	}
	if len(order.Products) > 0 {
		d.drawProductsSectionSingle(len(order.Gates)+1, order.Products)
		d.pdf.Ln(4)
	}

	d.drawGrandTotalSingle(order)
	d.pdf.Ln(6)
	d.addLine(fmt.Sprintf("КП действительно %d дней", supplier.OfferValidityDays))

	return renderPDF(d)
}

// GenerateDealerOfferForSelf генерирует КП дилера для себя (розничные + дилерские цены).
// orderID — номер заказа; если nil, номер в заголовок не включается.
func GenerateDealerOfferForSelf(
	order *types.Order,
	orderID *int64,
	supplier SupplierInfo,
	clientName string,
) ([]byte, error) {
	d := newDoc()

	drawOfferHeader(d, orderID)
	drawSupplierBuyerBlock(d, supplier, clientName)
	d.pdf.Ln(4)

	for i, g := range order.Gates {
		d.drawGateSectionDual(i+1, g)
		d.pdf.Ln(4)
	}
	if len(order.Products) > 0 {
		d.drawProductsSectionDual(len(order.Gates)+1, order.Products)
		d.pdf.Ln(4)
	}

	d.drawGrandTotalDual(order)
	d.pdf.Ln(6)
	d.addLine(fmt.Sprintf("КП действительно %d дней", supplier.OfferValidityDays))

	return renderPDF(d)
}

// drawOfferHeader рисует заголовок КП с датой и опциональным номером заказа.
func drawOfferHeader(d *Doc, orderID *int64) {
	title := "КОММЕРЧЕСКОЕ ПРЕДЛОЖЕНИЕ"
	if orderID != nil {
		title = fmt.Sprintf("КОММЕРЧЕСКОЕ ПРЕДЛОЖЕНИЕ  Заказ №%d", *orderID)
	}
	d.addTitle(title, 13)
	d.pdf.SetFont("DejaVu", "", 9)
	d.pdf.CellFormat(0, 6, "от "+time.Now().Format("02.01.2006")+" г.", "", 1, "C", false, 0, "")
	d.pdf.Ln(4)
}

// drawSupplierBuyerBlock рисует блок «Поставщик | Покупатель».
func drawSupplierBuyerBlock(d *Doc, supplier SupplierInfo, clientName string) {
	d.pdf.SetDrawColor(180, 180, 180)
	d.pdf.Line(marginX, d.pdf.GetY(), marginX+contentWidth, d.pdf.GetY())
	d.pdf.Ln(3)

	colW := contentWidth / 2

	d.pdf.SetFont("DejaVu", "B", 8)
	d.pdf.CellFormat(colW, 5, "ПОСТАВЩИК", "", 0, "L", false, 0, "")
	d.pdf.CellFormat(colW, 5, "ПОКУПАТЕЛЬ", "", 1, "L", false, 0, "")

	d.pdf.SetFont("DejaVu", "B", 9)
	d.pdf.CellFormat(colW, 5, supplier.Name, "", 0, "L", false, 0, "")
	d.pdf.SetFont("DejaVu", "", 9)
	d.pdf.CellFormat(colW, 5, clientName, "", 1, "L", false, 0, "")

	d.pdf.SetFont("DejaVu", "", 8)
	if supplier.Phone != "" {
		d.pdf.CellFormat(colW, 5, "Тел.: "+supplier.Phone, "", 0, "L", false, 0, "")
		d.pdf.CellFormat(colW, 5, "", "", 1, "L", false, 0, "")
	}
	if supplier.Email != "" {
		d.pdf.CellFormat(colW, 5, "Email: "+supplier.Email, "", 0, "L", false, 0, "")
		d.pdf.CellFormat(colW, 5, "", "", 1, "L", false, 0, "")
	}

	d.pdf.Ln(2)
	d.pdf.SetDrawColor(180, 180, 180)
	d.pdf.Line(marginX, d.pdf.GetY(), marginX+contentWidth, d.pdf.GetY())
	d.pdf.Ln(4)
}

// renderPDF проверяет ошибки fpdf и возвращает байты PDF.
func renderPDF(d *Doc) ([]byte, error) {
	if err := d.pdf.Error(); err != nil {
		return nil, fmt.Errorf("ошибка построения PDF: %w", err)
	}
	var buf bytes.Buffer
	if err := d.pdf.Output(&buf); err != nil {
		return nil, fmt.Errorf("ошибка вывода PDF: %w", err)
	}
	return buf.Bytes(), nil
}
