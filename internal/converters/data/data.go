package dataconv

import (
	"encoding/json"
	"os"
	"strings"

	"github.com/AliFnieer/morph/internal/registry"
	"github.com/pelletier/go-toml/v2"
	"gopkg.in/yaml.v3"
)

type conv struct{}

func init() {
	registry.Register(&conv{})
}

func (c *conv) Name() string { return "morph-data" }
func (c *conv) From() []string {
	return []string{"yaml", "yml", "toml", "json", "jsonl", "ndjson", "xml", "ini", "env"}
}
func (c *conv) To() []string {
	return []string{"yaml", "yml", "toml", "json", "jsonl", "ndjson", "xml", "ini", "txt"}
}

func (c *conv) Convert(in, out string) error {
	data, err := os.ReadFile(in)
	if err != nil {
		return err
	}
	inExt := getExt(in)
	outExt := getExt(out)

	var v interface{}

	if inExt == "yaml" || inExt == "yml" {
		if err := yaml.Unmarshal(data, &v); err != nil {
			return err
		}
	} else if inExt == "toml" {
		if err := toml.Unmarshal(data, &v); err != nil {
			return err
		}
	} else if inExt == "json" {
		if err := json.Unmarshal(data, &v); err != nil {
			return err
		}
	} else if inExt == "jsonl" || inExt == "ndjson" {
		lines := strings.Split(string(data), "\n")
		var arr []interface{}
		for _, l := range lines {
			l = strings.TrimSpace(l)
			if l == "" {
				continue
			}
			var item interface{}
			if err := json.Unmarshal([]byte(l), &item); err != nil {
				return err
			}
			arr = append(arr, item)
		}
		v = arr
		if outExt == "jsonl" || outExt == "ndjson" {
			var b strings.Builder
			for _, item := range arr {
				jb, _ := json.Marshal(item)
				b.Write(jb)
				b.WriteString("\n")
			}
			return os.WriteFile(out, []byte(b.String()), 0644)
		}
	} else if inExt == "ini" || inExt == "env" {
		m := map[string]string{}
		for _, line := range strings.Split(string(data), "\n") {
			line = strings.TrimSpace(line)
			if line == "" || strings.HasPrefix(line, "#") || strings.HasPrefix(line, ";") {
				continue
			}
			parts := strings.SplitN(line, "=", 2)
			if len(parts) == 2 {
				m[strings.TrimSpace(parts[0])] = strings.TrimSpace(parts[1])
			}
		}
		v = m
		if outExt == "json" {
			b, _ := json.MarshalIndent(v, "", "  ")
			return os.WriteFile(out, b, 0644)
		}
	} else if inExt == "xml" {
		if outExt == "json" {
			m := map[string]string{"raw_xml": string(data)}
			b, _ := json.MarshalIndent(m, "", "  ")
			return os.WriteFile(out, b, 0644)
		}
	} else {
		// try as text
		v = string(data)
	}

	if outExt == "json" {
		b, err := json.MarshalIndent(v, "", "  ")
		if err != nil {
			b = data
		}
		return os.WriteFile(out, b, 0644)
	}
	if outExt == "yaml" || outExt == "yml" {
		b, err := yaml.Marshal(v)
		if err != nil {
			return err
		}
		return os.WriteFile(out, b, 0644)
	}
	if outExt == "toml" {
		b, err := toml.Marshal(v)
		if err != nil {
			return err
		}
		return os.WriteFile(out, b, 0644)
	}
	if outExt == "jsonl" || outExt == "ndjson" {
		if arr, ok := v.([]interface{}); ok {
			var b strings.Builder
			for _, item := range arr {
				jb, _ := json.Marshal(item)
				b.Write(jb)
				b.WriteString("\n")
			}
			return os.WriteFile(out, []byte(b.String()), 0644)
		}
		jb, _ := json.Marshal(v)
		return os.WriteFile(out, append(jb, '\n'), 0644)
	}
	if outExt == "xml" {
		s := "<root>" + string(data) + "</root>"
		return os.WriteFile(out, []byte(s), 0644)
	}
	if outExt == "ini" || outExt == "env" {
		if m, ok := v.(map[string]interface{}); ok {
			var b strings.Builder
			for k, val := range m {
				b.WriteString(k)
				b.WriteString("=")
				b.WriteString(toStr(val))
				b.WriteString("\n")
			}
			return os.WriteFile(out, []byte(b.String()), 0644)
		}
		if m, ok := v.(map[string]string); ok {
			var b strings.Builder
			for k, val := range m {
				b.WriteString(k)
				b.WriteString("=")
				b.WriteString(val)
				b.WriteString("\n")
			}
			return os.WriteFile(out, []byte(b.String()), 0644)
		}
	}
	if outExt == "txt" {
		if s, ok := v.(string); ok {
			return os.WriteFile(out, []byte(s), 0644)
		}
		jb, _ := json.MarshalIndent(v, "", "  ")
		return os.WriteFile(out, jb, 0644)
	}

	return os.WriteFile(out, data, 0644)
}

func toStr(v interface{}) string {
	if s, ok := v.(string); ok {
		return s
	}
	b, _ := json.Marshal(v)
	return string(b)
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

type Conv = conv
