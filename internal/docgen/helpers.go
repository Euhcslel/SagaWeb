package docgen

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/Euhcslel/SagaWeb/internal/domain/enums"
	"github.com/Euhcslel/SagaWeb/internal/domain/order_gates"
	"github.com/Euhcslel/SagaWeb/internal/domain/order_products"
	"github.com/Euhcslel/SagaWeb/internal/types"
	"github.com/go-pdf/fpdf"
	"github.com/shopspring/decimal"
)

const (
	marginX      = 18.0
	contentWidth = 210.0 - 2*marginX
	lineHeight   = 4.5 // высота одной строки текста внутри ячейки
	minRowHeight = 7.0 // минимальная высота строки таблицы
)

var columnWidth = struct {
	No, Name, Unit, Price, Amount, Sum float64
}{
	No: 10, Name: 86, Unit: 14, Price: 26, Amount: 14, Sum: 24,
}

// Doc — обертка над fpdf для генерации документа.
type Doc struct {
	pdf *fpdf.Fpdf
}

// Cоздает новый PDF-документ и настраивает шрифты, поля и страницу.
func newDoc() *Doc {
	pdf := fpdf.New("P", "mm", "A4", "")
	pdf.SetMargins(marginX, 14, marginX)
	pdf.SetAutoPageBreak(true, 15)

	const fontDir = "C:/Windows/Fonts"
	pdf.AddUTF8Font("DejaVu", "", fontDir+"/arial.ttf")
	pdf.AddUTF8Font("DejaVu", "B", fontDir+"/arialbd.ttf")

	pdf.AddPage()
	return &Doc{pdf: pdf}
}

// Добавляет обычную строку текста
func (d *Doc) addLine(s string) {
	d.pdf.SetFont("DejaVu", "", 9)
	d.pdf.CellFormat(0, 6, s, "", 1, "L", false, 0, "")
}


// Добавляет заголовок указанного размера
func (d *Doc) addTitle(s string, size float64) {
	d.pdf.SetFont("DejaVu", "B", size)
	d.pdf.CellFormat(0, 8, s, "", 1, "L", false, 0, "")
}

// Рисует строку таблицы с автоматическим переносом текста
// в колонке "Наименование"
func (d *Doc) drawRow(no, name, unit, price, qty, sum string, bold bool) {
	style := ""
	if bold {
		style = "B"
	}
	d.pdf.SetFont("DejaVu", style, 8)

	lines := d.pdf.SplitText(name, columnWidth.Name-2)
	rows := len(lines)
	if rows == 0 {
		rows = 1
	}
	rowH := float64(rows) * lineHeight
	if rowH < minRowHeight {
		rowH = minRowHeight
	}

	yStart := d.pdf.GetY()

	d.pdf.CellFormat(columnWidth.No, rowH, no, "1", 0, "C", false, 0, "")

	xName := d.pdf.GetX()
	cellLineH := rowH / float64(rows)
	d.pdf.MultiCell(columnWidth.Name, cellLineH, name, "1", "L", false)
	d.pdf.SetXY(xName+columnWidth.Name, yStart)

	d.pdf.CellFormat(columnWidth.Unit, rowH, unit, "1", 0, "C", false, 0, "")
	d.pdf.CellFormat(columnWidth.Price, rowH, price, "1", 0, "R", false, 0, "")
	d.pdf.CellFormat(columnWidth.Amount, rowH, qty, "1", 0, "C", false, 0, "")
	d.pdf.CellFormat(columnWidth.Sum, rowH, sum, "1", 1, "R", false, 0, "")

	d.pdf.SetY(yStart + rowH)

	d.pdf.SetY(yStart + rowH)
}

// Добавляет строку-пояснение без цен и количества
func (d *Doc) drawTextRow(no, name string) {
	d.drawRow(no, name, "", "", "", "", false)
}

// Форматирует сумму в денежный вид (например: 12 345,67)
func money(v decimal.Decimal) string {
	s := v.StringFixed(2)
	parts := strings.SplitN(s, ".", 2)
	intPart, decPart := parts[0], parts[1]
	var b strings.Builder
	n := len(intPart)
	for i, ch := range intPart {
		if i > 0 && (n-i)%3 == 0 {
			b.WriteByte(' ')
		}
		b.WriteRune(ch)
	}
	return b.String() + "," + decPart
}

// Получает строковое поле структуры по имени через reflect
func getStringField(v any, fieldName string) string {
	rv := reflect.ValueOf(v)
	if rv.Kind() == reflect.Ptr {
		rv = rv.Elem()
	}
	if rv.Kind() != reflect.Struct {
		return ""
	}
	f := rv.FieldByName(fieldName)
	if !f.IsValid() || f.Kind() != reflect.String {
		return ""
	}
	return f.String()
}

// Получает поле типа decimal.Decimal по имени через reflect
func getDecimalField(v any, fieldName string) decimal.Decimal {
	rv := reflect.ValueOf(v)
	if rv.Kind() == reflect.Ptr {
		rv = rv.Elem()
	}
	if rv.Kind() != reflect.Struct {
		return decimal.Zero
	}
	f := rv.FieldByName(fieldName)
	if !f.IsValid() {
		return decimal.Zero
	}
	if d, ok := f.Interface().(decimal.Decimal); ok {
		return d
	}
	return decimal.Zero
}

// Формирует название ворот для документа
func gateTitle(g order_gates.OrderGate) string {
	base := "Секционные ворота"
	switch g.GateType {
	case enums.GateTypeInd:
		base = "Промышленные секционные ворота"
	case enums.GateTypeRes:
		base = "Гаражные секционные ворота"
	}
	return fmt.Sprintf("%s, проём %d×%d мм", base, g.Width, g.Height)
}

// Формирует описание привода для вывода в документ
func driveDescription(drive any) string {
	switch d := drive.(type) {
	case types.ResidentialDriveRail:
		return fmt.Sprintf("Управление: %s + %s", d.Drive.Name, d.Rail.Name)
	case types.ManualDrive:
		return fmt.Sprintf("Управление: ручное (цепь %d м)", d.ChainLength)
	}

	if name := getStringField(drive, "Name"); name != "" {
		return "Управление: " + name
	}
	return "Управление"
}

// Возвращает розничную стоимость привода
func driveRetailPrice(drive any) decimal.Decimal {
	switch d := drive.(type) {
	case types.ResidentialDriveRail:
		return d.Drive.RetailPrice.Add(d.Rail.RetailPrice)
	case types.ManualDrive:
		chain := decimal.NewFromInt32(d.ChainLength)
		return d.PriceInfo.ChainMeterRetailPrice.Mul(chain).Add(d.PriceInfo.RcpRetailPrice)
	}

	return getDecimalField(drive, "RetailPrice")
}

// Выводит в документ информацию по одним воротам:
// характеристики, опции, привод и итоговую стоимость
func (d *Doc) drawGateSection(itemNo int, g types.Gate) {
	d.pdf.SetFont("DejaVu", "B", 10)
	d.pdf.SetFillColor(235, 235, 235)
	d.pdf.CellFormat(contentWidth, 7, gateTitle(g.Gate), "1", 1, "L", true, 0, "")

	d.drawTableHeader()

	drivePrice := driveRetailPrice(g.Drive)

	optionsTotal := decimal.Zero
	for _, opt := range g.Options {
		optionsTotal = optionsTotal.Add(opt.RetailPrice)
	}

	gateOnlyPrice := g.Gate.GateRetailPrice.Sub(drivePrice).Sub(optionsTotal)
	amount := decimal.NewFromInt32(g.Gate.Amount)

	d.drawRow(
		fmt.Sprintf("%d", itemNo),
		gateTitle(g.Gate),
		"шт.",
		money(gateOnlyPrice),
		fmt.Sprintf("%d", g.Gate.Amount),
		money(gateOnlyPrice.Mul(amount)),
		false,
	)

	d.drawTextRow("", fmt.Sprintf("Размер проёма (ширина х высота): %d х %d мм", g.Gate.Width, g.Gate.Height))
	if g.Gate.Headroom > 0 {
		d.drawTextRow("", fmt.Sprintf("Притолока: %d мм", g.Gate.Headroom))
	}

	sub := 0
	subNo := func() string { sub++; return fmt.Sprintf("%d.%d", itemNo, sub) }

	d.drawTextRow(subNo(), "Цвет панели (наружный): "+g.Gate.ColorOut.Code)
	d.drawTextRow(subNo(), "Тип подъёма: "+g.Gate.LiftType.Name)
	d.drawTextRow(subNo(), fmt.Sprintf("Количество циклов: %v", g.Gate.CycleAmount.Amount))

	if len(g.Options) > 0 {
		optionsHeaderNo := subNo()
		d.drawTextRow(optionsHeaderNo, "Опции:")

		for i, opt := range g.Options {
			optNo := fmt.Sprintf("%s.%d", optionsHeaderNo, i+1)
			optSum := opt.RetailPrice.Mul(amount)
			d.drawRow(
				optNo,
				"   "+opt.Name,
				"шт.",
				money(opt.RetailPrice),
				fmt.Sprintf("%d", g.Gate.Amount),
				money(optSum),
				false,
			)
		}
	}

	if g.Drive != nil && !drivePrice.IsZero() {
		driveSum := drivePrice.Mul(amount)
		d.drawRow(
			subNo(),
			driveDescription(g.Drive),
			"шт.",
			money(drivePrice),
			fmt.Sprintf("%d", g.Gate.Amount),
			money(driveSum),
			false,
		)
	}

	d.pdf.SetFont("DejaVu", "B", 8)
	labelW := columnWidth.No + columnWidth.Name + columnWidth.Unit + columnWidth.Price + columnWidth.Amount
	d.pdf.CellFormat(labelW, minRowHeight, "Итого ворота:", "1", 0, "R", false, 0, "")
	d.pdf.CellFormat(columnWidth.Sum, minRowHeight, money(g.Gate.GateRetailPrice.Mul(amount)), "1", 1, "R", false, 0, "")
}

// Рисует шапку таблицы
func (d *Doc) drawTableHeader() {
	d.pdf.SetFont("DejaVu", "B", 8)
	d.pdf.SetFillColor(245, 245, 245)
	d.pdf.CellFormat(columnWidth.No, 8, "№", "1", 0, "C", true, 0, "")
	d.pdf.CellFormat(columnWidth.Name, 8, "Наименование и характеристика", "1", 0, "L", true, 0, "")
	d.pdf.CellFormat(columnWidth.Unit, 8, "Ед.", "1", 0, "C", true, 0, "")
	d.pdf.CellFormat(columnWidth.Price, 8, "Цена, ₽", "1", 0, "R", true, 0, "")
	d.pdf.CellFormat(columnWidth.Amount, 8, "Кол.", "1", 0, "C", true, 0, "")
	d.pdf.CellFormat(columnWidth.Sum, 8, "Сумма, ₽", "1", 1, "R", true, 0, "")
}

// Выводит таблицу дополнительных товаров
func (d *Doc) drawProductsSection(startNo int, products []order_products.OrderProduct) {
	d.pdf.SetFont("DejaVu", "B", 10)
	d.pdf.SetFillColor(235, 235, 235)
	d.pdf.CellFormat(contentWidth, 7, "Дополнительные товары", "1", 1, "L", true, 0, "")

	d.drawTableHeader()

	for i, p := range products {
		amount := decimal.NewFromInt32(p.Amount)
		sum := p.Product.RetailPrice.Mul(amount)

		d.drawRow(
			fmt.Sprintf("%d", startNo+i),
			p.Product.Name,
			"шт.",
			money(p.Product.RetailPrice),
			fmt.Sprintf("%d", p.Amount),
			money(sum),
			false,
		)
	}
}

// Рассчитывает и выводит общую сумму заказа
func (d *Doc) drawGrandTotal(order *types.Order) {
	total := decimal.Zero
	for _, g := range order.Gates {
		amount := decimal.NewFromInt32(g.Gate.Amount)
		total = total.Add(g.Gate.GateRetailPrice.Mul(amount))
	}
	for _, p := range order.Products {
		amount := decimal.NewFromInt32(p.Amount)
		total = total.Add(p.Product.RetailPrice.Mul(amount))
	}

	d.pdf.SetFont("DejaVu", "B", 10)
	labelW := columnWidth.No + columnWidth.Name + columnWidth.Unit + columnWidth.Price + columnWidth.Amount
	d.pdf.CellFormat(labelW, 8, "ИТОГО К ОПЛАТЕ:", "1", 0, "R", false, 0, "")
	d.pdf.CellFormat(columnWidth.Sum, 8, money(total), "1", 1, "R", false, 0, "")
}
