package gats

import (
	"exp/html"
	"fmt"
	"github.com/dunmatt/goquery"
	"io"
	"os"
)

func RenderTemplateFile(filename string, data interface{}, out io.Writer) error {
	f, err := os.Open(filename) // read the file in
	if err != nil {
		return err
	}
	defer f.Close()
	rootNode, err := html.Parse(f) // parse it
	if err != nil {
		return err
	}
	template := goquery.NewDocumentFromNode(rootNode) // wrap goquery around the DOM
	err = fillInTemplate(template, data)              // use goquery to process the template
	if err == nil {
		html.Render(out, rootNode) // render the DOM back to html and send it off
	}
	return err
}

func fillInTemplate(t *goquery.Document, data interface{}) error {
	handleGatsRemoves(t)
	e := handleGatsIf(t, data)
	if e != nil {
		return e
	}
	return nil
}

func handleGatsRemoves(t *goquery.Document) {
	t.Find("[gatsremove]").Remove()
}

func handleGatsIf(t *goquery.Document, data interface{}) error {
	var result error = nil
	t.Find("[gatsif]").Each(func(_ int, sel *goquery.Selection) {
		fieldName, _ := sel.Attr("gatsif")
		show, found := getBool(fieldName, data)
		if !found {
			result = fmt.Errorf("%v not found in the data!", fieldName)
			return
		}
		if !show {
			sel.Remove()
		}
		sel.RemoveAttr("gatsif")
	})
	return result
}
