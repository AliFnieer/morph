package registry

import (
	"encoding/json"
	"fmt"
	"sort"
)

type Converter interface {
	Convert(in, out string) error
	Name() string
	From() []string
	To() []string
}

var converters []Converter

func Register(c Converter) {
	converters = append(converters, c)
}

func Get(from, to string) Converter {
	for _, c := range converters {
		if contains(c.From(), from) && contains(c.To(), to) {
			return c
		}
	}
	return nil
}

func List() {
	fmt.Print("Supported conversions (registered):\n\n")
	groups := map[string]map[string]bool{}
	for _, c := range converters {
		cat := category(c.Name())
		if groups[cat] == nil {
			groups[cat] = map[string]bool{}
		}
		for _, f := range c.From() {
			for _, t := range c.To() {
				groups[cat][fmt.Sprintf("%s -> %s", f, t)] = true
			}
		}
	}
	var cats []string
	for k := range groups {
		cats = append(cats, k)
	}
	sort.Strings(cats)
	for _, cat := range cats {
		fmt.Printf("[%s]\n", cat)
		var items []string
		for k := range groups[cat] {
			items = append(items, k)
		}
		sort.Strings(items)
		for _, it := range items {
			fmt.Printf("  %s\n", it)
		}
		fmt.Println()
	}
}

func ListJSON() {
	type pair struct {
		From     string `json:"from"`
		To       string `json:"to"`
		Via      string `json:"via"`
		Category string `json:"category"`
	}
	var out []pair
	type p struct{ f, t string }
	seen := map[p]bool{}
	for _, c := range converters {
		cat := category(c.Name())
		for _, f := range c.From() {
			for _, t := range c.To() {
				pr := p{f, t}
				if seen[pr] {
					continue
				}
				seen[pr] = true
				out = append(out, pair{From: f, To: t, Via: c.Name(), Category: cat})
			}
		}
	}
	b, _ := json.MarshalIndent(out, "", "  ")
	fmt.Println(string(b))
}

func contains(list []string, s string) bool {
	for _, v := range list {
		if v == s {
			return true
		}
	}
	return false
}

func category(name string) string {
	switch {
	case name == "stdlib-image" || name == "x/image" || name == "x/image-tiff":
		return "Images"
	case name == "oksvg":
		return "Vector"
	case name == "morph-data":
		return "Data/Config"
	case name == "stdlib-text" || name == "morph-html":
		return "Text/Markup"
	case name == "morph-pdf" || name == "morph-pdf-minimal":
		return "Documents"
	case name == "morph-office-minimal":
		return "Office"
	case name == "morph-ooxml":
		return "Office (OOXML)"
	case name == "morph-odf":
		return "Office (ODF)"
	case name == "morph-epub":
		return "Ebooks"
	case name == "morph-rtf":
		return "Documents"
	case name == "morph-archive":
		return "Archives"
	case name == "morph-audio":
		return "Audio"
	case name == "morph-video-stub":
		return "Video"
	case name == "morph-raw-copy":
		return "Raw/Other"
	default:
		return "Other"
	}
}

type Pair struct {
	From string
	To   string
}

func AllPairs() []Pair {
	type p struct{ f, t string }
	seen := map[p]bool{}
	var out []Pair
	for _, c := range converters {
		for _, f := range c.From() {
			for _, t := range c.To() {
				pr := p{f, t}
				if seen[pr] {
					continue
				}
				seen[pr] = true
				out = append(out, Pair{From: f, To: t})
			}
		}
	}
	return out
}
