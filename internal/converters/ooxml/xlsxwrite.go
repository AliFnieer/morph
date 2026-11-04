package ooxmlconv

import (
	"archive/zip"
	"encoding/csv"
	"os"
	"strings"

	"github.com/AliFnieer/morph/internal/registry"
)

type xlsxWriter struct{}

func init() {
	registry.Register(&xlsxWriter{})
}

func (c *xlsxWriter) Name() string   { return "morph-ooxml-xlsx-writer" }
func (c *xlsxWriter) From() []string { return []string{"csv", "txt"} }
func (c *xlsxWriter) To() []string   { return []string{"xlsx"} }

func (c *xlsxWriter) Convert(in, out string) error {
	data, err := os.ReadFile(in)
	if err != nil {
		return err
	}
	var rows [][]string
	if getExt(in) == "csv" {
		r := csv.NewReader(strings.NewReader(string(data)))
		rows, err = r.ReadAll()
		if err != nil {
			return err
		}
	} else {
		for _, line := range strings.Split(string(data), "\n") {
			rows = append(rows, []string{line})
		}
	}
	zf, err := os.Create(out)
	if err != nil {
		return err
	}
	defer zf.Close()
	w := zip.NewWriter(zf)
	defer w.Close()

	if err := writeZipPart(w, "[Content_Types].xml", `<?xml version="1.0" encoding="UTF-8"?><Types xmlns="http://schemas.openxmlformats.org/package/2006/content-types"><Default Extension="rels" ContentType="application/vnd.openxmlformats-package.relationships+xml"/><Default Extension="xml" ContentType="application/xml"/><Override PartName="/xl/workbook.xml" ContentType="application/vnd.openxmlformats-officedocument.spreadsheetml.sheet.main+xml"/><Override PartName="/xl/worksheets/sheet1.xml" ContentType="application/vnd.openxmlformats-officedocument.spreadsheetml.worksheet+xml"/></Types>`); err != nil {
		return err
	}
	if err := writeZipPart(w, "_rels/.rels", `<?xml version="1.0" encoding="UTF-8"?><Relationships xmlns="http://schemas.openxmlformats.org/package/2006/relationships"><Relationship Id="rId1" Type="http://schemas.openxmlformats.org/officeDocument/2006/relationships/officeDocument" Target="xl/workbook.xml"/></Relationships>`); err != nil {
		return err
	}
	if err := writeZipPart(w, "xl/workbook.xml", `<?xml version="1.0" encoding="UTF-8"?><workbook xmlns="http://schemas.openxmlformats.org/spreadsheetml/2006/main" xmlns:r="http://schemas.openxmlformats.org/officeDocument/2006/relationships"><sheets><sheet name="Sheet1" sheetId="1" r:id="rId1"/></sheets></workbook>`); err != nil {
		return err
	}
	if err := writeZipPart(w, "xl/_rels/workbook.xml.rels", `<?xml version="1.0" encoding="UTF-8"?><Relationships xmlns="http://schemas.openxmlformats.org/package/2006/relationships"><Relationship Id="rId1" Type="http://schemas.openxmlformats.org/officeDocument/2006/relationships/worksheet" Target="worksheets/sheet1.xml"/></Relationships>`); err != nil {
		return err
	}
	sheet := `<?xml version="1.0" encoding="UTF-8"?><worksheet xmlns="http://schemas.openxmlformats.org/spreadsheetml/2006/main"><sheetData>`
	for r := 0; r < len(rows); r++ {
		sheet += `<row r="` + itoa(r+1) + `">`
		for c := 0; c < len(rows[r]); c++ {
			col := string(rune('A' + c))
			sheet += `<c r="` + col + itoa(r+1) + `" t="inlineStr"><is><t>` + escapeXML(rows[r][c]) + `</t></is></c>`
		}
		sheet += `</row>`
	}
	sheet += `</sheetData></worksheet>`
	if err := writeZipPart(w, "xl/worksheets/sheet1.xml", sheet); err != nil {
		return err
	}
	return nil
}
