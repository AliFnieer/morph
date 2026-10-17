package ooxmlconv

import (
	"archive/zip"
	"encoding/xml"
	"os"
	"strings"

	"github.com/AliFnieer/morph/internal/registry"
)

type conv struct{}

func init() {
	registry.Register(&conv{})
}

func (c *conv) Name() string   { return "morph-ooxml" }
func (c *conv) From() []string { return []string{"docx", "pptx", "xlsx"} }
func (c *conv) To() []string   { return []string{"txt", "md", "html", "csv", "json", "pdf"} }

func (c *conv) Convert(in, out string) error {
	inExt := getExt(in)
	outExt := getExt(out)
	r, err := zip.OpenReader(in)
	if err != nil {
		return err
	}
	defer r.Close()

	var all string
	if inExt == "docx" {
		all = extractDocx(r)
	} else if inExt == "xlsx" {
		all = extractXlsx(r)
	} else if inExt == "pptx" {
		all = extractPptx(r)
	}

	if outExt == "pdf" {
		content := all
		content = strings.ReplaceAll(content, `\`, `\\`)
		content = strings.ReplaceAll(content, `(`, `\(`)
		content = strings.ReplaceAll(content, `)`, `\)`)
		var b strings.Builder
		b.WriteString("%PDF-1.4\n")
		b.WriteString("1 0 obj<</Type/Catalog/Pages 2 0 R>>endobj\n")
		b.WriteString("2 0 obj<</Type/Pages/Kids[3 0 R]/Count 1>>endobj\n")
		b.WriteString("3 0 obj<</Type/Page/Parent 2 0 R/MediaBox[0 0 612 792]/Contents 4 0 R/Resources<</Font<</F1<</Type/Font/Subtype/Type1/BaseFont/Helvetica>>>>>>endobj\n")
		stream := "BT /F1 10 Tf 50 750 Td (" + content + ") Tj ET\n"
		b.WriteString("4 0 obj<</Length " + itoa(len(stream)) + ">>stream\n")
		b.WriteString(stream)
		b.WriteString("endstream\nendobj\n")
		xref := b.Len()
		b.WriteString("xref\n0 5\n0000000000 65535 f \n")
		b.WriteString("trailer<</Size 5/Root 1 0 R>>\nstartxref\n")
		b.WriteString(itoa(xref))
		b.WriteString("\n%%EOF\n")
		return os.WriteFile(out, []byte(b.String()), 0644)
	}
	if outExt == "html" {
		s := "<html><body><pre>" + htmlEscape(all) + "</pre></body></html>"
		return os.WriteFile(out, []byte(s), 0644)
	}
	if outExt == "csv" && inExt == "xlsx" {
		return os.WriteFile(out, []byte(all), 0644)
	}
	if outExt == "json" {
		enc := []byte("{\"text\":\"" + strings.ReplaceAll(all, "\"", "\\\"") + "\"}")
		return os.WriteFile(out, enc, 0644)
	}
	return os.WriteFile(out, []byte(all), 0644)
}

func extractDocx(r *zip.ReadCloser) string {
	for _, f := range r.File {
		if f.Name == "word/document.xml" {
			rc, err := f.Open()
			if err == nil {
				buf := make([]byte, 10*1024*1024)
				n, _ := rc.Read(buf)
				rc.Close()
				return extractTextXML(string(buf[:n]))
			}
		}
	}
	return ""
}

func extractPptx(r *zip.ReadCloser) string {
	var parts []string
	for _, f := range r.File {
		if strings.HasPrefix(f.Name, "ppt/slides/slide") && strings.HasSuffix(f.Name, ".xml") {
			rc, err := f.Open()
			if err == nil {
				buf := make([]byte, 5*1024*1024)
				n, _ := rc.Read(buf)
				rc.Close()
				parts = append(parts, extractTextXML(string(buf[:n])))
			}
		}
	}
	return strings.Join(parts, "\n\n")
}

func extractXlsx(r *zip.ReadCloser) string {
	// try to extract shared strings
	shared := map[string]string{}
	for _, f := range r.File {
		if f.Name == "xl/sharedStrings.xml" {
			rc, err := f.Open()
			if err == nil {
				buf := make([]byte, 5*1024*1024)
				n, _ := rc.Read(buf)
				rc.Close()
				shared = parseSharedStrings(string(buf[:n]))
			}
		}
	}
	// extract first sheet
	for _, f := range r.File {
		if strings.HasPrefix(f.Name, "xl/worksheets/sheet") && strings.HasSuffix(f.Name, ".xml") {
			rc, err := f.Open()
			if err == nil {
				buf := make([]byte, 10*1024*1024)
				n, _ := rc.Read(buf)
				rc.Close()
				return parseSheet(string(buf[:n]), shared)
			}
		}
	}
	return ""
}

func parseSharedStrings(s string) map[string]string {
	res := map[string]string{}
	// crude: split by <t> and </t>
	parts := strings.Split(s, "<t")
	for i := 1; i < len(parts); i++ {
		end := strings.Index(parts[i], "</t>")
		if end > 0 {
			// skip tag attrs
			if idx := strings.Index(parts[i], ">"); idx >= 0 && idx < end {
				res[string(rune('0'+len(res)))] = parts[i][idx+1 : end]
			}
		}
	}
	// better keyed by index - rebuild properly
	res = map[string]string{}
	idx := 0
	parts = strings.Split(s, "<si")
	for i := 1; i < len(parts); i++ {
		// look for <t>
		t := parts[i]
		if strings.Contains(t, "<t>") {
			start := strings.Index(t, "<t>") + 3
			end := strings.Index(t[start:], "</t>")
			if end > 0 {
				res[itoa(idx)] = t[start : start+end]
				idx++
				continue
			}
		}
		if strings.Contains(t, "<t ") {
			if start := strings.Index(t, ">"); start >= 0 {
				start++
				if end := strings.Index(t[start:], "</t>"); end > 0 {
					res[itoa(idx)] = t[start : start+end]
					idx++
					continue
				}
			}
		}
	}
	return res
}

func parseSheet(s string, shared map[string]string) string {
	var out strings.Builder
	// split rows
	rows := strings.Split(s, "<row")
	for i := 1; i < len(rows); i++ {
		var cells []string
		// split cells
		cellsParts := strings.Split(rows[i], "<c")
		for j := 1; j < len(cellsParts); j++ {
			cell := cellsParts[j]
			// check type
			if strings.Contains(cell, `t="s"`) {
				// shared string index
				if idx := strings.Index(cell, "<v>"); idx >= 0 {
					vend := strings.Index(cell[idx:], "</v>")
					if vend > 0 {
						num := cell[idx+3 : idx+vend]
						if val, ok := shared[num]; ok {
							cells = append(cells, val)
							continue
						}
						cells = append(cells, num)
						continue
					}
				}
			} else if strings.Contains(cell, `t="inlineStr"`) {
				if start := strings.Index(cell, "<t"); start >= 0 {
					if gt := strings.Index(cell[start:], ">"); gt >= 0 {
						start2 := start + gt + 1
						if end := strings.Index(cell[start2:], "</t>"); end > 0 {
							cells = append(cells, cell[start2:start2+end])
							continue
						}
					}
				}
			}
			if idx := strings.Index(cell, "<v>"); idx >= 0 {
				vend := strings.Index(cell[idx:], "</v>")
				if vend > 0 {
					cells = append(cells, cell[idx+3:idx+vend])
					continue
				}
			}
			cells = append(cells, "")
		}
		if len(cells) > 0 {
			out.WriteString(strings.Join(cells, ","))
			out.WriteString("\n")
		}
	}
	if out.Len() == 0 {
		return extractTextXML(s)
	}
	return out.String()
}

func extractTextXML(s string) string {
	var res strings.Builder
	inTag := false
	for i := 0; i < len(s); i++ {
		if s[i] == '<' {
			inTag = true
			continue
		}
		if s[i] == '>' {
			inTag = false
			continue
		}
		if !inTag {
			res.WriteByte(s[i])
		}
	}
	return strings.ReplaceAll(res.String(), "  ", " ")
}

func htmlEscape(s string) string {
	s = strings.ReplaceAll(s, "&", "&amp;")
	s = strings.ReplaceAll(s, "<", "&lt;")
	s = strings.ReplaceAll(s, ">", "&gt;")
	return s
}
func getExt(path string) string {
	for i := len(path) - 1; i >= 0; i-- {
		if path[i] == '.' {
			return strings.ToLower(path[i+1:])
		}
		if path[i] == '/' || path[i] == '\\' {
			break
		}
	}
	return ""
}
func itoa(i int) string {
	if i == 0 {
		return "0"
	}
	neg := false
	if i < 0 {
		neg = true
		i = -i
	}
	var b [32]byte
	pos := len(b)
	for i > 0 {
		pos--
		b[pos] = byte('0' + i%10)
		i /= 10
	}
	if neg {
		pos--
		b[pos] = '-'
	}
	return string(b[pos:])
}

type Conv = conv

var _ xml.Name
