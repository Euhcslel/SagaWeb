package docgen

import (
	_ "embed"
	"fmt"
	"strings"
	"time"

	"github.com/Euhcslel/SagaWeb/internal/domain/enums"
	"github.com/Euhcslel/SagaWeb/internal/domain/industrial_gate_drives"
	"github.com/Euhcslel/SagaWeb/internal/domain/order_gates"
	"github.com/Euhcslel/SagaWeb/internal/domain/order_products"
	"github.com/Euhcslel/SagaWeb/internal/types"
	"github.com/go-pdf/fpdf"
	"github.com/shopspring/decimal"
)

//go:embed fonts/DejaVuSans.ttf
var fontRegular []byte

//go:embed fonts/DejaVuSans-Bold.ttf
var fontBold []byte

const (
	marginX      = 18.0
	contentWidth = 210.0 - 2*marginX
	lineHeight   = 4.5
	minRowHeight = 7.0
)

// Ширины колонок для таблицы с одной ценой (КП для покупателя / клиента).
var colSingle = struct {
	No, Name, Unit, Price, Amount, Sum float64
}{
	No: 10, Name: 86, Unit: 14, Price: 26, Amount: 14, Sum: 24,
}

// Ширины колонок для таблицы с двумя ценами (КП для дилера).
var colDual = struct {
	No, Name, Unit, Retail, Dealer, Amount, SumRetail, SumDealer float64
}{
	No: 8, Name: 58, Unit: 10, Retail: 22, Dealer: 22, Amount: 12, SumRetail: 21, SumDealer: 21,
}

// Doc — обёртка над fpdf для генерации документа
type Doc struct {
	pdf *fpdf.Fpdf
}

// Функция, создающая новый PDF-документ. Настраивает шрифты, поля и страницу
func newDoc() *Doc {
	pdf := fpdf.New("P", "mm", "A4", "")
	pdf.SetMargins(marginX, 14, marginX)
	pdf.SetAutoPageBreak(true, 15)
	pdf.AddUTF8FontFromBytes("DejaVu", "", fontRegular)
	pdf.AddUTF8FontFromBytes("DejaVu", "B", fontBold)
	pdf.AddPage()
	return &Doc{pdf: pdf}
}

func (d *Doc) addLine(s string) {
	d.pdf.SetFont("DejaVu", "", 9)
	d.pdf.CellFormat(0, 6, s, "", 1, "L", false, 0, "")
}

func (d *Doc) addTitle(s string, size float64) {
	d.pdf.SetFont("DejaVu", "B", size)
	d.pdf.CellFormat(0, 8, s, "", 1, "C", false, 0, "")
}

// Функция, форматирующая decimal в строку вида «12 345,67»
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

// gateTitle формирует название ворот для документа.
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

// driveDescription возвращает строку описания привода.
func driveDescription(drive any) string {
	switch d := drive.(type) {
	case industrial_gate_drives.IndustrialGateDrive:
		return "Управление: " + d.Name
	case types.ResidentialDriveRail:
		return fmt.Sprintf("Управление: %s + %s", d.Drive.Name, d.Rail.Name)
	case types.ManualDrive:
		return fmt.Sprintf("Управление: ручное (цепь %d м)", d.ChainLength)
	}
	return "Управление"
}

// driveRetailPrice возвращает розничную цену привода.
func driveRetailPrice(drive any) decimal.Decimal {
	switch d := drive.(type) {
	case industrial_gate_drives.IndustrialGateDrive:
		return d.RetailPrice
	case types.ResidentialDriveRail:
		return d.Drive.RetailPrice.Add(d.Rail.RetailPrice)
	case types.ManualDrive:
		chain := decimal.NewFromInt32(d.ChainLength)
		return d.PriceInfo.ChainMeterRetailPrice.Mul(chain).Add(d.PriceInfo.RcpRetailPrice)
	}
	return decimal.Zero
}

// driveWholesalePrice возвращает дилерскую (оптовую) цену привода.
func driveWholesalePrice(drive any) decimal.Decimal {
	switch d := drive.(type) {
	case industrial_gate_drives.IndustrialGateDrive:
		return d.WholesalePrice
	case types.ResidentialDriveRail:
		return d.Drive.WholesalePrice.Add(d.Rail.WholesalePrice)
	case types.ManualDrive:
		chain := decimal.NewFromInt32(d.ChainLength)
		return d.PriceInfo.ChainMeterWholesalePrice.Mul(chain).Add(d.PriceInfo.RcpWholesalePrice)
	}
	return decimal.Zero
}

// ─── Таблица с одной ценой ────────────────────────────────────────────────────

func (d *Doc) drawTableHeaderSingle() {
	d.pdf.SetFont("DejaVu", "B", 8)
	d.pdf.SetFillColor(245, 245, 245)
	d.pdf.CellFormat(colSingle.No, 8, "№", "1", 0, "C", true, 0, "")
	d.pdf.CellFormat(colSingle.Name, 8, "Наименование и характеристика", "1", 0, "L", true, 0, "")
	d.pdf.CellFormat(colSingle.Unit, 8, "Ед.", "1", 0, "C", true, 0, "")
	d.pdf.CellFormat(colSingle.Price, 8, "Цена, ₽", "1", 0, "R", true, 0, "")
	d.pdf.CellFormat(colSingle.Amount, 8, "Кол.", "1", 0, "C", true, 0, "")
	d.pdf.CellFormat(colSingle.Sum, 8, "Сумма, ₽", "1", 1, "R", true, 0, "")
}

func (d *Doc) drawRowSingle(no, name, unit, price, qty, sum string, bold bool) {
	style := ""
	if bold {
		style = "B"
	}
	d.pdf.SetFont("DejaVu", style, 8)

	lines := d.pdf.SplitText(name, colSingle.Name-2)
	rows := len(lines)
	if rows == 0 {
		rows = 1
	}
	rowH := float64(rows) * lineHeight
	if rowH < minRowHeight {
		rowH = minRowHeight
	}
	yStart := d.pdf.GetY()

	d.pdf.CellFormat(colSingle.No, rowH, no, "1", 0, "C", false, 0, "")
	xName := d.pdf.GetX()
	d.pdf.MultiCell(colSingle.Name, rowH/float64(rows), name, "1", "L", false)
	d.pdf.SetXY(xName+colSingle.Name, yStart)
	d.pdf.CellFormat(colSingle.Unit, rowH, unit, "1", 0, "C", false, 0, "")
	d.pdf.CellFormat(colSingle.Price, rowH, price, "1", 0, "R", false, 0, "")
	d.pdf.CellFormat(colSingle.Amount, rowH, qty, "1", 0, "C", false, 0, "")
	d.pdf.CellFormat(colSingle.Sum, rowH, sum, "1", 1, "R", false, 0, "")
	d.pdf.SetY(yStart + rowH)
}

func (d *Doc) drawTextRowSingle(no, name string) {
	d.drawRowSingle(no, name, "", "", "", "", false)
}

// ─── Таблица с двумя ценами ───────────────────────────────────────────────────

func (d *Doc) drawTableHeaderDual() {
	d.pdf.SetFont("DejaVu", "B", 7)
	d.pdf.SetFillColor(245, 245, 245)
	d.pdf.CellFormat(colDual.No, 8, "№", "1", 0, "C", true, 0, "")
	d.pdf.CellFormat(colDual.Name, 8, "Наименование и характеристика", "1", 0, "L", true, 0, "")
	d.pdf.CellFormat(colDual.Unit, 8, "Ед.", "1", 0, "C", true, 0, "")
	d.pdf.CellFormat(colDual.Retail, 8, "Цена розн.", "1", 0, "R", true, 0, "")
	d.pdf.CellFormat(colDual.Dealer, 8, "Цена дил.", "1", 0, "R", true, 0, "")
	d.pdf.CellFormat(colDual.Amount, 8, "Кол.", "1", 0, "C", true, 0, "")
	d.pdf.CellFormat(colDual.SumRetail, 8, "Сумма розн.", "1", 0, "R", true, 0, "")
	d.pdf.CellFormat(colDual.SumDealer, 8, "Сумма дил.", "1", 1, "R", true, 0, "")
}

func (d *Doc) drawRowDual(no, name, unit, retail, dealer, qty, sumRetail, sumDealer string, bold bool) {
	style := ""
	if bold {
		style = "B"
	}
	d.pdf.SetFont("DejaVu", style, 7)

	lines := d.pdf.SplitText(name, colDual.Name-2)
	rows := len(lines)
	if rows == 0 {
		rows = 1
	}
	rowH := float64(rows) * lineHeight
	if rowH < minRowHeight {
		rowH = minRowHeight
	}
	yStart := d.pdf.GetY()

	d.pdf.CellFormat(colDual.No, rowH, no, "1", 0, "C", false, 0, "")
	xName := d.pdf.GetX()
	d.pdf.MultiCell(colDual.Name, rowH/float64(rows), name, "1", "L", false)
	d.pdf.SetXY(xName+colDual.Name, yStart)
	d.pdf.CellFormat(colDual.Unit, rowH, unit, "1", 0, "C", false, 0, "")
	d.pdf.CellFormat(colDual.Retail, rowH, retail, "1", 0, "R", false, 0, "")
	d.pdf.CellFormat(colDual.Dealer, rowH, dealer, "1", 0, "R", false, 0, "")
	d.pdf.CellFormat(colDual.Amount, rowH, qty, "1", 0, "C", false, 0, "")
	d.pdf.CellFormat(colDual.SumRetail, rowH, sumRetail, "1", 0, "R", false, 0, "")
	d.pdf.CellFormat(colDual.SumDealer, rowH, sumDealer, "1", 1, "R", false, 0, "")
	d.pdf.SetY(yStart + rowH)
}

func (d *Doc) drawTextRowDual(no, name string) {
	d.drawRowDual(no, name, "", "", "", "", "", "", false)
}

// ─── Секция ворот ─────────────────────────────────────────────────────────────

// Функция, рисующая секцию ворот для КП с одной колонкой цен
func (d *Doc) drawGateSectionSingle(itemNo int, g types.Gate) {
	d.pdf.SetFont("DejaVu", "B", 10)
	d.pdf.SetFillColor(235, 235, 235)
	d.pdf.CellFormat(contentWidth, 7, gateTitle(g.Gate), "1", 1, "L", true, 0, "")
	d.drawTableHeaderSingle()

	drivePrice := driveRetailPrice(g.Drive)
	optionsTotal := decimal.Zero
	for _, opt := range g.Options {
		optionsTotal = optionsTotal.Add(opt.RetailPrice)
	}
	gateBasePrice := g.Gate.GateRetailPrice.Sub(drivePrice).Sub(optionsTotal)
	amount := decimal.NewFromInt32(g.Gate.Amount)

	d.drawRowSingle(
		fmt.Sprintf("%d", itemNo),
		gateTitle(g.Gate),
		"шт.",
		money(gateBasePrice),
		fmt.Sprintf("%d", g.Gate.Amount),
		money(gateBasePrice.Mul(amount)),
		false,
	)
	d.drawTextRowSingle("", fmt.Sprintf("Размер проёма (ширина х высота): %d х %d мм", g.Gate.Width, g.Gate.Height))
	if g.Gate.Headroom > 0 {
		d.drawTextRowSingle("", fmt.Sprintf("Притолока: %d мм", g.Gate.Headroom))
	}

	sub := 0
	subNo := func() string { sub++; return fmt.Sprintf("%d.%d", itemNo, sub) }

	d.drawTextRowSingle(subNo(), "Цвет панели (наружный): "+g.Gate.ColorOut.Code)
	d.drawTextRowSingle(subNo(), "Тип подъёма: "+g.Gate.LiftType.Name)
	d.drawTextRowSingle(subNo(), fmt.Sprintf("Количество циклов: %v", g.Gate.CycleAmount.Amount))

	if len(g.Options) > 0 {
		optHeader := subNo()
		d.drawTextRowSingle(optHeader, "Опции:")
		for i, opt := range g.Options {
			optSum := opt.RetailPrice.Mul(amount)
			d.drawRowSingle(
				fmt.Sprintf("%s.%d", optHeader, i+1),
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
		d.drawRowSingle(
			subNo(),
			driveDescription(g.Drive),
			"шт.",
			money(drivePrice),
			fmt.Sprintf("%d", g.Gate.Amount),
			money(drivePrice.Mul(amount)),
			false,
		)
	}

	labelW := colSingle.No + colSingle.Name + colSingle.Unit + colSingle.Price + colSingle.Amount
	d.pdf.SetFont("DejaVu", "B", 8)
	d.pdf.CellFormat(labelW, minRowHeight, "Итого ворота:", "1", 0, "R", false, 0, "")
	d.pdf.CellFormat(colSingle.Sum, minRowHeight, money(g.Gate.GateRetailPrice.Mul(amount)), "1", 1, "R", false, 0, "")
}

// drawGateSectionDual рисует секцию ворот для КП с двумя колонками цен.
func (d *Doc) drawGateSectionDual(itemNo int, g types.Gate) {
	d.pdf.SetFont("DejaVu", "B", 10)
	d.pdf.SetFillColor(235, 235, 235)
	d.pdf.CellFormat(contentWidth, 7, gateTitle(g.Gate), "1", 1, "L", true, 0, "")
	d.drawTableHeaderDual()

	retailDrive := driveRetailPrice(g.Drive)
	dealerDrive := driveWholesalePrice(g.Drive)

	retailOptions := decimal.Zero
	dealerOptions := decimal.Zero
	for _, opt := range g.Options {
		retailOptions = retailOptions.Add(opt.RetailPrice)
		dealerOptions = dealerOptions.Add(opt.WholesalePrice)
	}

	gateBaseRetail := g.Gate.GateRetailPrice.Sub(retailDrive).Sub(retailOptions)
	gateBaseDealer := g.Gate.GateWholesalePrice.Sub(dealerDrive).Sub(dealerOptions)
	amount := decimal.NewFromInt32(g.Gate.Amount)

	d.drawRowDual(
		fmt.Sprintf("%d", itemNo),
		gateTitle(g.Gate),
		"шт.",
		money(gateBaseRetail),
		money(gateBaseDealer),
		fmt.Sprintf("%d", g.Gate.Amount),
		money(gateBaseRetail.Mul(amount)),
		money(gateBaseDealer.Mul(amount)),
		false,
	)
	d.drawTextRowDual("", fmt.Sprintf("Размер проёма (ширина х высота): %d х %d мм", g.Gate.Width, g.Gate.Height))
	if g.Gate.Headroom > 0 {
		d.drawTextRowDual("", fmt.Sprintf("Притолока: %d мм", g.Gate.Headroom))
	}

	sub := 0
	subNo := func() string { sub++; return fmt.Sprintf("%d.%d", itemNo, sub) }

	d.drawTextRowDual(subNo(), "Цвет панели (наружный): "+g.Gate.ColorOut.Code)
	d.drawTextRowDual(subNo(), "Тип подъёма: "+g.Gate.LiftType.Name)
	d.drawTextRowDual(subNo(), fmt.Sprintf("Количество циклов: %v", g.Gate.CycleAmount.Amount))

	if len(g.Options) > 0 {
		optHeader := subNo()
		d.drawTextRowDual(optHeader, "Опции:")
		for i, opt := range g.Options {
			d.drawRowDual(
				fmt.Sprintf("%s.%d", optHeader, i+1),
				"   "+opt.Name,
				"шт.",
				money(opt.RetailPrice),
				money(opt.WholesalePrice),
				fmt.Sprintf("%d", g.Gate.Amount),
				money(opt.RetailPrice.Mul(amount)),
				money(opt.WholesalePrice.Mul(amount)),
				false,
			)
		}
	}

	if g.Drive != nil && (!retailDrive.IsZero() || !dealerDrive.IsZero()) {
		d.drawRowDual(
			subNo(),
			driveDescription(g.Drive),
			"шт.",
			money(retailDrive),
			money(dealerDrive),
			fmt.Sprintf("%d", g.Gate.Amount),
			money(retailDrive.Mul(amount)),
			money(dealerDrive.Mul(amount)),
			false,
		)
	}

	labelW := colDual.No + colDual.Name + colDual.Unit + colDual.Retail + colDual.Dealer + colDual.Amount
	d.pdf.SetFont("DejaVu", "B", 8)
	d.pdf.CellFormat(labelW, minRowHeight, "Итого ворота:", "1", 0, "R", false, 0, "")
	d.pdf.CellFormat(colDual.SumRetail, minRowHeight, money(g.Gate.GateRetailPrice.Mul(amount)), "1", 0, "R", false, 0, "")
	d.pdf.CellFormat(colDual.SumDealer, minRowHeight, money(g.Gate.GateWholesalePrice.Mul(amount)), "1", 1, "R", false, 0, "")
}

// ─── Секция товаров ───────────────────────────────────────────────────────────

func (d *Doc) drawProductsSectionSingle(startNo int, products []order_products.OrderProduct) {
	d.pdf.SetFont("DejaVu", "B", 10)
	d.pdf.SetFillColor(235, 235, 235)
	d.pdf.CellFormat(contentWidth, 7, "Дополнительные товары", "1", 1, "L", true, 0, "")
	d.drawTableHeaderSingle()
	for i, p := range products {
		amount := decimal.NewFromInt32(p.Amount)
		d.drawRowSingle(
			fmt.Sprintf("%d", startNo+i),
			p.Product.Name,
			"шт.",
			money(p.Product.RetailPrice),
			fmt.Sprintf("%d", p.Amount),
			money(p.Product.RetailPrice.Mul(amount)),
			false,
		)
	}
}

func (d *Doc) drawProductsSectionDual(startNo int, products []order_products.OrderProduct) {
	d.pdf.SetFont("DejaVu", "B", 10)
	d.pdf.SetFillColor(235, 235, 235)
	d.pdf.CellFormat(contentWidth, 7, "Дополнительные товары", "1", 1, "L", true, 0, "")
	d.drawTableHeaderDual()
	for i, p := range products {
		amount := decimal.NewFromInt32(p.Amount)
		d.drawRowDual(
			fmt.Sprintf("%d", startNo+i),
			p.Product.Name,
			"шт.",
			money(p.Product.RetailPrice),
			money(p.Product.WholesalePrice),
			fmt.Sprintf("%d", p.Amount),
			money(p.Product.RetailPrice.Mul(amount)),
			money(p.Product.WholesalePrice.Mul(amount)),
			false,
		)
	}
}

func (d *Doc) drawGateSectionSingleWholesale(itemNo int, g types.Gate) {
	d.pdf.SetFont("DejaVu", "B", 10)
	d.pdf.SetFillColor(235, 235, 235)
	d.pdf.CellFormat(contentWidth, 7, gateTitle(g.Gate), "1", 1, "L", true, 0, "")
	d.drawTableHeaderSingle()

	drivePrice := driveWholesalePrice(g.Drive)
	optionsTotal := decimal.Zero
	for _, opt := range g.Options {
		optionsTotal = optionsTotal.Add(opt.WholesalePrice)
	}
	gateBasePrice := g.Gate.GateWholesalePrice.Sub(drivePrice).Sub(optionsTotal)
	amount := decimal.NewFromInt32(g.Gate.Amount)

	d.drawRowSingle(
		fmt.Sprintf("%d", itemNo),
		gateTitle(g.Gate),
		"шт.",
		money(gateBasePrice),
		fmt.Sprintf("%d", g.Gate.Amount),
		money(gateBasePrice.Mul(amount)),
		false,
	)
	d.drawTextRowSingle("", fmt.Sprintf("Размер проёма (ширина х высота): %d х %d мм", g.Gate.Width, g.Gate.Height))
	if g.Gate.Headroom > 0 {
		d.drawTextRowSingle("", fmt.Sprintf("Притолока: %d мм", g.Gate.Headroom))
	}

	sub := 0
	subNo := func() string { sub++; return fmt.Sprintf("%d.%d", itemNo, sub) }

	d.drawTextRowSingle(subNo(), "Цвет панели (наружный): "+g.Gate.ColorOut.Code)
	d.drawTextRowSingle(subNo(), "Тип подъёма: "+g.Gate.LiftType.Name)
	d.drawTextRowSingle(subNo(), fmt.Sprintf("Количество циклов: %v", g.Gate.CycleAmount.Amount))

	if len(g.Options) > 0 {
		optHeader := subNo()
		d.drawTextRowSingle(optHeader, "Опции:")
		for i, opt := range g.Options {
			optSum := opt.WholesalePrice.Mul(amount)
			d.drawRowSingle(
				fmt.Sprintf("%s.%d", optHeader, i+1),
				"   "+opt.Name,
				"шт.",
				money(opt.WholesalePrice),
				fmt.Sprintf("%d", g.Gate.Amount),
				money(optSum),
				false,
			)
		}
	}

	if g.Drive != nil && !drivePrice.IsZero() {
		d.drawRowSingle(
			subNo(),
			driveDescription(g.Drive),
			"шт.",
			money(drivePrice),
			fmt.Sprintf("%d", g.Gate.Amount),
			money(drivePrice.Mul(amount)),
			false,
		)
	}

	labelW := colSingle.No + colSingle.Name + colSingle.Unit + colSingle.Price + colSingle.Amount
	d.pdf.SetFont("DejaVu", "B", 8)
	d.pdf.CellFormat(labelW, minRowHeight, "Итого ворота:", "1", 0, "R", false, 0, "")
	d.pdf.CellFormat(colSingle.Sum, minRowHeight, money(g.Gate.GateWholesalePrice.Mul(amount)), "1", 1, "R", false, 0, "")
}

func (d *Doc) drawProductsSectionSingleWholesale(startNo int, products []order_products.OrderProduct) {
	d.pdf.SetFont("DejaVu", "B", 10)
	d.pdf.SetFillColor(235, 235, 235)
	d.pdf.CellFormat(contentWidth, 7, "Дополнительные товары", "1", 1, "L", true, 0, "")
	d.drawTableHeaderSingle()
	for i, p := range products {
		amount := decimal.NewFromInt32(p.Amount)
		d.drawRowSingle(
			fmt.Sprintf("%d", startNo+i),
			p.Product.Name,
			"шт.",
			money(p.Product.WholesalePrice),
			fmt.Sprintf("%d", p.Amount),
			money(p.Product.WholesalePrice.Mul(amount)),
			false,
		)
	}
}

// ─── Итого ────────────────────────────────────────────────────────────────────

func (d *Doc) drawGrandTotalSingle(order *types.Order) {
	total := decimal.Zero
	for _, g := range order.Gates {
		total = total.Add(g.Gate.GateRetailPrice.Mul(decimal.NewFromInt32(g.Gate.Amount)))
	}
	for _, p := range order.Products {
		total = total.Add(p.Product.RetailPrice.Mul(decimal.NewFromInt32(p.Amount)))
	}
	const totalSumW = 34.0
	labelW := contentWidth - totalSumW
	d.pdf.SetFont("DejaVu", "B", 10)
	d.pdf.CellFormat(labelW, 8, "ИТОГО К ОПЛАТЕ:", "1", 0, "R", false, 0, "")
	d.pdf.CellFormat(totalSumW, 8, money(total), "1", 1, "R", false, 0, "")
}

func (d *Doc) drawGrandTotalSingleWholesale(order *types.Order) {
	total := decimal.Zero
	for _, g := range order.Gates {
		total = total.Add(g.Gate.GateWholesalePrice.Mul(decimal.NewFromInt32(g.Gate.Amount)))
	}
	for _, p := range order.Products {
		total = total.Add(p.Product.WholesalePrice.Mul(decimal.NewFromInt32(p.Amount)))
	}
	const totalSumW = 34.0
	labelW := contentWidth - totalSumW
	d.pdf.SetFont("DejaVu", "B", 10)
	d.pdf.CellFormat(labelW, 8, "ИТОГО К ОПЛАТЕ:", "1", 0, "R", false, 0, "")
	d.pdf.CellFormat(totalSumW, 8, money(total), "1", 1, "R", false, 0, "")
}

func (d *Doc) drawGrandTotalDual(order *types.Order) {
	totalRetail := decimal.Zero
	totalDealer := decimal.Zero
	for _, g := range order.Gates {
		amt := decimal.NewFromInt32(g.Gate.Amount)
		totalRetail = totalRetail.Add(g.Gate.GateRetailPrice.Mul(amt))
		totalDealer = totalDealer.Add(g.Gate.GateWholesalePrice.Mul(amt))
	}
	for _, p := range order.Products {
		amt := decimal.NewFromInt32(p.Amount)
		totalRetail = totalRetail.Add(p.Product.RetailPrice.Mul(amt))
		totalDealer = totalDealer.Add(p.Product.WholesalePrice.Mul(amt))
	}
	const totalDualSumW = 26.0
	labelW := contentWidth - totalDualSumW*2
	d.pdf.SetFont("DejaVu", "B", 10)
	d.pdf.CellFormat(labelW, 8, "ИТОГО К ОПЛАТЕ:", "1", 0, "R", false, 0, "")
	d.pdf.CellFormat(totalDualSumW, 8, money(totalRetail), "1", 0, "R", false, 0, "")
	d.pdf.CellFormat(totalDualSumW, 8, money(totalDealer), "1", 1, "R", false, 0, "")
}

// ─── Общие утилиты ────────────────────────────────────────────────────────────

var ruMonths = [...]string{
	"января", "февраля", "марта", "апреля", "мая", "июня",
	"июля", "августа", "сентября", "октября", "ноября", "декабря",
}

func formatDateRu(t time.Time) string {
	return fmt.Sprintf("«%d» %s %d г.", t.Day(), ruMonths[t.Month()-1], t.Year())
}

func ptrOrEmpty(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}
