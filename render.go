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
	// NOTE: handleGatsRepeatOvers is mutually recursive with fillInTemplate
	for sel := scope.Find("[gatsrepeatover]"); sel.Length() > 0; sel = scope.Find("[gatsrepeatover]") {
		handleGatsRepeatOvers(sel.First(), cont)
	}

	// do the actual work of filling in the template
	e := handleGatsIf(scope, cont)
	if e != nil {
		return e
	}
	e = handleGatsAttributes(scope, cont)
	if e != nil {
		return e
	}
	e = handleGatsText(scope, cont)
	if e != nil {
		return e
	}
	return nil
}

func handleGatsRepeatOvers(sel *goquery.Selection, cont *context) {
	fieldName, _ := sel.Attr("gatsrepeatover")
	length, err := getLength(fieldName, cont)
	if err == nil {
		for i := 0; i < length; i++ {
			c, e := getItem(fieldName, i, cont)
			if e == nil {
				instance := sel.Clone().InsertBefore(sel)
				fillInTemplate(instance, c) // TODO: stop ignoring the returned errors here
				instance.RemoveAttr("gatsrepeatover")
			}
		}
	}
	sel.Remove()
}

func handleGatsRemoves(t *goquery.Selection) {
	t.Find("[gatsremove]").Remove()
}

func handleGatsIf(t *goquery.Selection, cont *context) error {
	var result error // = nil
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

func handleGatsAttributes(t *goquery.Selection, cont *context) error {
	var result error
	t.Find("[gatsattributes]").Each(func(_ int, sel *goquery.Selection) {
		fieldName, _ := sel.Attr("gatsattributes")
		attribs, err := getStringMap(fieldName, cont)
		if err != nil {
			result = err
			return
		}
		for k, v := range attribs {
			sel.SetAttr(k, v)
		}
		sel.RemoveAttr("gatsattributes")
	})
	return result
}

func handleGatsText(t *goquery.Selection, cont *context) error {
	var result error
	t.Find("[gatstext]").Each(func(_ int, sel *goquery.Selection) {
		fieldName, _ := sel.Attr("gatstext")
		text, err := getString(fieldName, cont)
		if err != nil {
			result = err
			return
		}
		sel.SetText(text)
		sel.RemoveAttr("gatstext")
	})
	return result
}
