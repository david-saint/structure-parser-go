package parser

import (
	"encoding/json"
	"fmt"
	"html"
	"os"
	"strings"
	"sync"
)

type DomNode struct {
	Element    string
	Value      string
	Attributes []AttrNode
	Css        []CssNode
	Children   []DomNode
}

type CssNode struct {
	Property string
	Value    string
}

type AttrNode struct {
	Key   string
	Value string
}

type Stringer interface {
	String() string
}

func (c CssNode) String() string {
	return fmt.Sprintf("%s: %s;", c.Property, c.Value)
}

func (a AttrNode) String() string {
	return fmt.Sprintf("%s=\"%s\"", a.Key, a.Value)
}

func Parse(f *os.File) (r string, err error) {
	var domTree []DomNode
	var wg sync.WaitGroup

	decoder := json.NewDecoder(f)
	err = decoder.Decode(&domTree)

	if err != nil {
		return
	}

	for _, node := range domTree {
		wg.Add(1)
		go func() {
			r += node.parse()
			wg.Done()
		}()
	}
	wg.Wait()
	return
}

func (t *DomNode) parse() string {
	if t.Element == ElementTextNode {
		return t.Value
	}

	var style string
	var attributes string
	var wg sync.WaitGroup
	tag := html.EscapeString(t.Element)
	html := "<" + tag

	if len(t.Css) > 0 {
		style = GenerateStyles(" ", t.Css...)
		html += " style=\"" + style + "\""
	}

	if len(t.Attributes) > 0 {
		attributes = GenerateAttributes(" ", t.Attributes...)
		html += " " + attributes
	}

	if len(t.Children) == 0 {
		html += "/>"
		return html
	}

	html += ">"

	for _, node := range t.Children {
		wg.Add(1)
		go func() {
			html += node.parse()
			wg.Done()
		}()
		wg.Wait()
	}

	html += "</" + tag + ">"

	return html
}

func GenerateStyles(separator string, s ...CssNode) string {
	var list []string
	for _, kv := range s {
		list = append(list, kv.String())
	}
	return strings.Join(list[:], separator)
}

func GenerateAttributes(separator string, s ...AttrNode) string {
	var list []string
	for _, kv := range s {
		list = append(list, kv.String())
	}
	return strings.Join(list[:], separator)
}
