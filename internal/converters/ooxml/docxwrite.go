package ooxmlconv

import (
	"archive/zip"
	"os"
	"strings"

	"github.com/AliFnieer/morph/internal/registry"
)

type docxWriter struct{}

func init() {
	registry.Register(&docxWriter{})
}

func (c *docxWriter) Name() string   { return "morph-ooxml-writer" }
func (c *docxWriter) From() []string { return []string{"txt", "md"} }
func (c *docxWriter) To() []string   { return []string{"docx"} }

func (c *docxWriter) Convert(in, out string) error {
	data, err := os.ReadFile(in)
	if err != nil {
		return err
	}
	text := string(data)
	zf, err := os.Create(out)
	if err != nil {
		return err
	}
	defer zf.Close()
	w := zip.NewWriter(zf)
	defer w.Close()

	// [Content_Types].xml
	ct := `[Content_Types].xml`
	f, _ := w.Create(ct)
	f.Write([]byte(`<?xml version="1.0" encoding="UTF-8"?><Types xmlns="http://schemas.openxmlformats.org/package/2006/content-types"><Default Extension="rels" ContentType="application/vnd.openxmlformats-package.relationships+xml"/><Default Extension="xml" ContentType="application/xml"/><Override PartName="/word/document.xml" ContentType="application/vnd.openxmlformats-officedocument.wordprocessingml.document.main+xml"/></Types>`))

	// _rels/.rels
	f, _ = w.Create("_rels/.rels")
	f.Write([]byte(`<?xml version="1.0" encoding="UTF-8"?><Relationships xmlns="http://schemas.openxmlformats.org/package/2006/relationships"><Relationship Id="rId1" Type="http://schemas.openxmlformats.org/officeDocument/2006/relationships/officeDocument" Target="word/document.xml"/></Relationships>`))

	// word/document.xml
	f, _ = w.Create("word/document.xml")
	doc := `<?xml version="1.0" encoding="UTF-8" standalone="yes"?><w:document xmlns:w="http://schemas.openxmlformats.org/wordprocessingml/2006/main"><w:body>`
	for _, line := range strings.Split(text, "\n") {
		doc += `<w:p><w:r><w:t>` + escapeXML(line) + `</w:t></w:r></w:p>`
	}
	doc += `<w:sectPr/></w:body></w:document>`
	f.Write([]byte(doc))

	return nil
}

func escapeXML(s string) string {
	s = strings.ReplaceAll(s, "&", "&amp;")
	s = strings.ReplaceAll(s, "<", "&lt;")
	s = strings.ReplaceAll(s, ">", "&gt;")
	return s
}
