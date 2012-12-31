package gats

import (
	"exp/html"
	//"fmt"
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
	cont := makeContext(data, nil)
	template := goquery.NewDocumentFromNode(rootNode).Find("html") // wrap goquery around the DOM
	handleGatsRemoves(template)
	err = fillInTemplate(template, cont) // use goquery to process the template
	if err == nil {
		html.Render(out, rootNode) // render the DOM back to html and send it off
	}
	return err
}

func fillInTemplate(scope *goquery.Selection, cont *context) error {
	// filling in the template happens children first so that the closure of
	// available data is readily available (if the most senior elements were
	// done first there'd be no way to know the reflection path within data)
	scope.Find("[gatsrepeatover]").Each(func(_ int, sel *goquery.Selection) {
		fieldName, found := sel.Attr("gatsrepeatover")
		if !found { // this can happen when the Find finds nested repeatovers
			return
		}
		length, err := getLength(fieldName, cont)
		insertionPoint := sel.Next()
		template := sel.Remove()
		if err == nil {
			for i := 0; i < length; i++ {
				c, e := getItem(fieldName, i, cont)
				if e == nil {
					sel = template.Clone().InsertBefore(insertionPoint)
					fillInTemplate(sel, c) // TODO: stop ignoring the returned errors here
				}
			}
		}
	})

	// do the actual work of filling in the template
	e := handleGatsIf(scope, cont)
	if e != nil {
		return e
	}
	return nil
}

func handleGatsRemoves(t *goquery.Selection) {
	t.Find("[gatsremove]").Remove()
}

func handleGatsIf(t *goquery.Selection, cont *context) error {
	var result error = nil
	t.Find("[gatsif]").Each(func(_ int, sel *goquery.Selection) {
		fieldName, _ := sel.Attr("gatsif")
		show, _ := getBool(fieldName, cont)
		if !show {
			sel.Remove()
		}
		sel.RemoveAttr("gatsif")
	})
	return result
}
