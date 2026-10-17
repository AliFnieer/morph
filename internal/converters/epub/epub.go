package epubconv

import (
	"archive/zip"
	"os"
	"strings"

	"github.com/AliFnieer/morph/internal/registry"
)

type conv struct{}

func init() {
	registry.Register(&conv{})
}

func (c *conv) Name() string   { return "morph-epub" }
func (c *conv) From() []string { return []string{"epub"} }
func (c *conv) To() []string   { return []string{"txt", "html", "md", "pdf"} }

func (c *conv) Convert(in, out string) error {
	r, err := zip.OpenReader(in)
	if err != nil {
		return err
	}
	defer r.Close()

	var parts []string
	for _, f := range r.File {
		if strings.HasSuffix(f.Name, ".html") || strings.HasSuffix(f.Name, ".xhtml") || strings.HasSuffix(f.Name, ".htm") {
			rc, err := f.Open()
			if err == nil {
				buf := make([]byte, 5*1024*1024)
				n, _ := rc.Read(buf)
				rc.Close()
				parts = append(parts, stripHTML(string(buf[:n])))
			}
		}
	}
	all := strings.Join(parts, "\n\n")
	outExt := getExt(out)
	if outExt == "html" {
		s := "<html><body><pre>" + htmlEscape(all) + "</pre></body></html>"
		return os.WriteFile(out, []byte(s), 0644)
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
	return os.WriteFile(out, []byte(all), 0644)
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
		if !inTag {
			res.WriteByte(s[i])
		}
	}
	s = strings.ReplaceAll(res.String(), "  ", " ")
	s = strings.ReplaceAll(s, "&nbsp;", " ")
	s = strings.ReplaceAll(s, "&amp;", "&")
	return s
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
