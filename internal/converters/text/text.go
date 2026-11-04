package textconv

import (
	"encoding/csv"
	"encoding/json"
	"os"
	"strings"

	"github.com/AliFnieer/morph/internal/registry"
)

type conv struct{}

func init() {
	registry.Register(&conv{})
}

func (c *conv) Name() string { return "stdlib-text" }
func (c *conv) From() []string {
	return []string{"txt", "json", "csv", "md"}
}
func (c *conv) To() []string {
	return []string{"txt", "json", "csv", "md"}
}

func (c *conv) Convert(in, out string) error {
	data, err := os.ReadFile(in)
	if err != nil {
		return err
	}

	inExt := getExt(in)
	outExt := getExt(out)

	// Handle CSV <-> JSON as basic conversions
	if inExt == "csv" && outExt == "json" {
		return csvToJSON(data, out)
	}
	if inExt == "json" && outExt == "csv" {
		return jsonToCSV(data, out)
	}

	// Otherwise, just copy/normalize as plain text
	return os.WriteFile(out, data, 0644)
}

func csvToJSON(data []byte, out string) error {
	r := csv.NewReader(strings.NewReader(string(data)))
	records, err := r.ReadAll()
	if err != nil {
		return err
	}
	if len(records) == 0 {
		return os.WriteFile(out, []byte("[]"), 0644)
	}
	headers := records[0]
	var result []map[string]string
	for i := 1; i < len(records); i++ {
		row := map[string]string{}
		for j, h := range headers {
			if j < len(records[i]) {
				row[h] = records[i][j]
			} else {
				row[h] = ""
			}
		}
		result = append(result, row)
	}
	b, err := json.MarshalIndent(result, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(out, b, 0644)
}

func jsonToCSV(data []byte, out string) error {
	var arr []map[string]interface{}
	if err := json.Unmarshal(data, &arr); err != nil {
		return err
	}
	if len(arr) == 0 {
		return os.WriteFile(out, []byte(""), 0644)
	}
	// collect headers
	headers := map[string]bool{}
	for _, m := range arr {
		for k := range m {
			headers[k] = true
		}
	}
	var headerList []string
	for k := range headers {
		headerList = append(headerList, k)
	}
	// write
	outf, err := os.Create(out)
	if err != nil {
		return err
	}
	defer outf.Close()
	w := csv.NewWriter(outf)
	defer w.Flush()
	if err := w.Write(headerList); err != nil {
		return err
	}
	for _, m := range arr {
		var row []string
		for _, h := range headerList {
			if v, ok := m[h]; ok {
				row = append(row, toString(v))
			} else {
				row = append(row, "")
			}
		}
		if err := w.Write(row); err != nil {
			return err
		}
	}
	return nil
}

func toString(v interface{}) string {
	switch t := v.(type) {
	case string:
		return t
	default:
		b, _ := json.Marshal(t)
		return string(b)
	}
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

// Make sure init runs
type Conv = conv
