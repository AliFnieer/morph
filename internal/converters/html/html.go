package htmlconv

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

func (c *conv) Name() string   { return "morph-html" }
func (c *conv) From() []string { return []string{"html", "htm", "txt", "md"} }
func (c *conv) To() []string   { return []string{"html", "htm", "txt"} }

func (c *conv) Convert(in, out string) error {
	data, err := os.ReadFile(in)
	if err != nil {
		return err
	}
	inExt := getExt(in)
	outExt := getExt(out)
	if outExt == "txt" || outExt == "" {
		s := string(data)
		re := regexp.MustCompile(`<[^>]*>`)
		s = re.ReplaceAllString(s, "")
		s = strings.ReplaceAll(s, "&nbsp;", " ")
		s = strings.ReplaceAll(s, "&amp;", "&")
		s = strings.ReplaceAll(s, "&lt;", "<")
		s = strings.ReplaceAll(s, "&gt;", ">")
		return os.WriteFile(out, []byte(s), 0644)
	}
	if inExt == "txt" || inExt == "md" {
		s := "<html><body><pre>" + htmlEscape(string(data)) + "</pre></body></html>"
		return os.WriteFile(out, []byte(s), 0644)
	}
	return os.WriteFile(out, data, 0644)
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
			return path[i+1:]
		}
		if path[i] == '/' || path[i] == '\\' {
			break
		}
	}
	return ""
}

type Conv = conv
