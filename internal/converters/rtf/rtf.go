package rtfconv

import (
	"os"
	"regexp"
	"strings"

	"github.com/AliFnieer/morph/internal/registry"
)

type conv struct{}

func init() {
	registry.Register(&conv{})
}

func (c *conv) Name() string   { return "morph-rtf" }
func (c *conv) From() []string { return []string{"rtf"} }
func (c *conv) To() []string   { return []string{"txt", "md", "html", "pdf"} }

func (c *conv) Convert(in, out string) error {
	data, err := os.ReadFile(in)
	if err != nil {
		return err
	}
	s := string(data)
	re := regexp.MustCompile(`\\[a-z]+-?\d* ?`)
	s = re.ReplaceAllString(s, " ")
	re2 := regexp.MustCompile(`\\[^a-zA-Z]`)
	s = re2.ReplaceAllString(s, "")
	s = strings.ReplaceAll(s, "{", "")
	s = strings.ReplaceAll(s, "}", "")
	s = strings.ReplaceAll(s, "\r", "\n")
	outExt := getExt(out)
	if outExt == "pdf" {
		content := s
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
		s2 := "<html><body><pre>" + htmlEscape(s) + "</pre></body></html>"
		return os.WriteFile(out, []byte(s2), 0644)
	}
	return os.WriteFile(out, []byte(s), 0644)
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
