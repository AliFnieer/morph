package pdfconv

import (
	"bytes"
	"os"
	"strings"
	"unicode"

	"github.com/AliFnieer/morph/internal/registry"
)

type conv struct{}

func init() {
	registry.Register(&conv{})
}

func (c *conv) Name() string { return "morph-pdf" }
func (c *conv) From() []string {
	return []string{"txt", "md", "text", "json", "csv", "log", "html", "htm"}
}
func (c *conv) To() []string {
	return []string{"pdf"}
}

func (c *conv) Convert(in, out string) error {
	data, err := os.ReadFile(in)
	if err != nil {
		return err
	}
	content := string(data)
	// For HTML, do basic strip
	ext := getExt(in)
	if ext == "html" || ext == "htm" {
		content = stripHTML(content)
	}
	// Wrap text
	lines := wrapText(content, 80)
	var buf bytes.Buffer
	buf.WriteString("%PDF-1.4\n")
	buf.WriteString("1 0 obj<</Type/Catalog/Pages 2 0 R>>endobj\n")
	buf.WriteString("2 0 obj<</Type/Pages/Kids[3 0 R]/Count 1>>endobj\n")
	buf.WriteString("3 0 obj<</Type/Page/Parent 2 0 R/MediaBox[0 0 612 792]/Contents 4 0 R/Resources<</Font<</F1<</Type/Font/Subtype/Type1/BaseFont/Helvetica>>>>>>endobj\n")
	// Build content stream with line breaks
	var stream bytes.Buffer
	stream.WriteString("BT /F1 10 Tf 50 750 Td 14 TL\n")
	for i, line := range lines {
		if i == 0 {
			stream.WriteString("(")
			stream.WriteString(escapePDF(line))
			stream.WriteString(") Tj T*\n")
		} else {
			stream.WriteString("(")
			stream.WriteString(escapePDF(line))
			stream.WriteString(") Tj T*\n")
		}
		if i > 200 { // limit very long files
			break
		}
	}
	stream.WriteString("ET\n")
	streamStr := stream.String()
	buf.WriteString("4 0 obj<</Length " + itoa(len(streamStr)) + ">>stream\n")
	buf.WriteString(streamStr)
	buf.WriteString("endstream\nendobj\n")
	xref := buf.Len()
	buf.WriteString("xref\n0 5\n0000000000 65535 f \n")
	buf.WriteString("trailer<</Size 5/Root 1 0 R>>\nstartxref\n")
	buf.WriteString(itoa(xref))
	buf.WriteString("\n%%EOF\n")
	return os.WriteFile(out, buf.Bytes(), 0644)
}

func escapePDF(s string) string {
	s = strings.ReplaceAll(s, `\`, `\\`)
	s = strings.ReplaceAll(s, `(`, `\(`)
	s = strings.ReplaceAll(s, `)`, `\)`)
	return s
}

func wrapText(s string, width int) []string {
	// simple word wrap
	var result []string
	for _, line := range strings.Split(s, "\n") {
		if len(line) <= width {
			result = append(result, line)
			continue
		}
		var current strings.Builder
		for _, word := range strings.Fields(line) {
			if current.Len()+len(word)+1 > width && current.Len() > 0 {
				result = append(result, current.String())
				current.Reset()
			}
			if current.Len() > 0 {
				current.WriteString(" ")
			}
			current.WriteString(word)
		}
		if current.Len() > 0 {
			result = append(result, current.String())
		}
	}
	if len(result) == 0 {
		result = append(result, "")
	}
	return result
}

func stripHTML(s string) string {
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
		if !inTag && unicode.IsPrint(rune(s[i])) {
			res.WriteByte(s[i])
		}
	}
	return res.String()
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
