package docgen

import (
	"archive/zip"
	"bytes"
	"fmt"
	"io"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/Euhcslel/SagaWeb/internal/domain/dealers"
	"github.com/Euhcslel/SagaWeb/internal/domain/users"
)

// GenerateContract генерирует договор на основе DOCX шаблона и данных дилера.
func GenerateContract(dealer *dealers.Dealer, dealerUser *users.User, contractNumber int, date time.Time) ([]byte, error) {
	templateData, err := os.ReadFile(Cfg.ContractTemplatePath)
	if err != nil {
		return nil, fmt.Errorf("открытие шаблона договора: %w", err)
	}

	replacements := map[string]string{
		"{{CONTRACT_NUMBER}}": strconv.Itoa(contractNumber),
		"{{DAY}}":             fmt.Sprintf("%02d", date.Day()),
		"{{MONTH}}":           ruMonths[date.Month()-1],
		"{{YEAR}}":            fmt.Sprintf("%d", date.Year()),
		"{{DEALER_COMPANY}}":  dealer.Company.Name,
		"{{DEALER_ADDRESS}}":  ptrOrEmpty(dealer.Address),
		"{{DEALER_FULLNAME}}": dealerUser.Fullname,
		"{{DEALER_EMAIL}}":    dealerUser.Email,
		"{{DEALER_TIN}}":      ptrOrEmpty(dealer.Company.INN),
		"{{DEALER_KPP}}":      ptrOrEmpty(dealer.Company.KPP),
		"{{DEALER_OGRN}}":     ptrOrEmpty(dealer.Company.OGRN),
	}

	return fillDocxTemplate(templateData, replacements)
}

// fillDocxTemplate открывает DOCX как ZIP и заменяет плейсхолдеры в XML.
func fillDocxTemplate(docxData []byte, replacements map[string]string) ([]byte, error) {
	r, err := zip.NewReader(bytes.NewReader(docxData), int64(len(docxData)))
	if err != nil {
		return nil, fmt.Errorf("чтение DOCX: %w", err)
	}

	var buf bytes.Buffer
	w := zip.NewWriter(&buf)

	for _, f := range r.File {
		content, err := readZipEntry(f)
		if err != nil {
			return nil, fmt.Errorf("чтение %s: %w", f.Name, err)
		}

		if isWordXML(f.Name) {
			xmlStr := string(content)
			for placeholder, value := range replacements {
				xmlStr = replaceSplitPlaceholder(xmlStr, placeholder, xmlEscape(value))
			}
			content = []byte(xmlStr)
		}

		fw, err := w.Create(f.Name)
		if err != nil {
			return nil, err
		}
		if _, err := fw.Write(content); err != nil {
			return nil, err
		}
	}

	if err := w.Close(); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

// replaceSplitPlaceholder заменяет {{KEY}} даже когда Word разбил его на
// несколько <w:r> runs из-за автоопределения языка или проверки орфографии.
func replaceSplitPlaceholder(xmlContent, placeholder, replacement string) string {
	const runBoundary = `(?:</w:t></w:r>(?:<w:r[^>]*>)(?s:(?:<w:rPr>.*?</w:rPr>)?)<w:t[^>]*>)?`

	var sb strings.Builder
	for i, ch := range []rune(placeholder) {
		if i > 0 {
			sb.WriteString(runBoundary)
		}
		sb.WriteString(regexp.QuoteMeta(string(ch)))
	}

	re, err := regexp.Compile(sb.String())
	if err != nil {
		return strings.ReplaceAll(xmlContent, placeholder, replacement)
	}
	return re.ReplaceAllString(xmlContent, replacement)
}

// readZipEntry открывает файл zip-архива и читает его содержимое в байтовый слайс.
func readZipEntry(f *zip.File) ([]byte, error) {
	rc, err := f.Open()
	if err != nil {
		return nil, err
	}
	defer rc.Close()
	return io.ReadAll(rc)
}

// isWordXML проверяет является ли файл xml-документом из папки word.
func isWordXML(name string) bool {
	return strings.HasPrefix(name, "word/") && strings.HasSuffix(name, ".xml")
}

// xmlEscape заменяет специальные символы.
func xmlEscape(s string) string {
	s = strings.ReplaceAll(s, "&", "&amp;")
	s = strings.ReplaceAll(s, "<", "&lt;")
	s = strings.ReplaceAll(s, ">", "&gt;")
	s = strings.ReplaceAll(s, `"`, "&quot;")
	return s
}

