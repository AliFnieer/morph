package pdfconv

import (
	"os"
	"regexp"
	"strings"

	"github.com/AliFnieer/morph/internal/registry"
)

type pdfReader struct{}

func init() {
	registry.Register(&pdfReader{})
}

func (c *pdfReader) Name() string   { return "morph-pdf-reader" }
func (c *pdfReader) From() []string { return []string{"pdf"} }
func (c *pdfReader) To() []string   { return []string{"txt", "html", "md"} }

func (c *pdfReader) Convert(in, out string) error {
	data, err := os.ReadFile(in)
	if err != nil {
		return err
	}
	s := string(data)
	// crude extraction of text in parentheses
	re := regexp.MustCompile(`\(([^()]*)\)`)
	matches := re.FindAllStringSubmatch(s, -1)
	var parts []string
	for _, m := range matches {
		if len(m) > 1 {
			txt := m[1]
			txt = strings.ReplaceAll(txt, `\(`, "(")
			txt = strings.ReplaceAll(txt, `\)`, ")")
			txt = strings.ReplaceAll(txt, `\\`, `\`)
			parts = append(parts, txt)
		}
	}
	all := strings.Join(parts, " ")
	outExt := getExt(out)
	if outExt == "html" {
		all = "<html><body><pre>" + htmlEscape(all) + "</pre></body></html>"
	}
	return os.WriteFile(out, []byte(all), 0644)
}

func htmlEscape(s string) string {
	s = strings.ReplaceAll(s, "&", "&amp;")
	s = strings.ReplaceAll(s, "<", "&lt;")
	s = strings.ReplaceAll(s, ">", "&gt;")
	return s
}
