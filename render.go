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

	handleGatsRemoves(template)          // take out everything that definitely won't be shown
	err = fillInTemplate(template, cont) // process the template
	handleGatsOmitTag(template)          // make sure this one comes last, it can interfere with other attributes
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
	e = handleGatsContent(scope, cont)
	if e != nil {
		return e
	}
	e = handleGatsText(scope, cont)
	if e != nil {
		return e
	}
	e = handleGatsAttributes(scope, cont)
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
				e1 := fillInTemplate(instance, c) // TODO: stop ignoring the returned errors here
				if e1 != nil {
					panic(e1) // is this really panic worthy?
				}
				instance.RemoveAttr("gatsrepeatover")
			}
		}
	}
	sel.Remove()
}

func handleGatsRemoves(t *goquery.Selection) {
	t.Find("[gatsremove]").Remove()
}

func handleGats(t *goquery.Selection, cont *context, selector string,
	meat func(string, *goquery.Selection) error) (result error) {
	attribName := selector[1 : len(selector)-1]
	t.Find(selector).Each(func(_ int, sel *goquery.Selection) {
		fieldName, _ := sel.Attr(attribName)
		result = meat(fieldName, sel)
		sel.RemoveAttr(attribName)
	})
	return
}

func handleGatsIf(t *goquery.Selection, cont *context) (result error) {
	return handleGats(t, cont, "[gatsif]",
		func(fieldName string, sel *goquery.Selection) (res error) {
			show, res := getBool(fieldName, cont)
			if !show {
				sel.Remove()
			}
			return res
		})
	//t.Find("[gatsif]").Each(func(_ int, sel *goquery.Selection) {
	//	fieldName, _ := sel.Attr("gatsif")
	//	show, _ := getBool(fieldName, cont)
	//	if !show {
	//		sel.Remove()
	//	}
	//	sel.RemoveAttr("gatsif")
	//})
	//return result
}

func handleGatsContent(t *goquery.Selection, cont *context) (result error) {
	t.Find("[gatscontent]").Each(func(_ int, sel *goquery.Selection) {
		fieldName, _ := sel.Attr("gatscontent")
		if isString(fieldName, cont) {
			rawHtml, result := getString(fieldName, cont)
			if result == nil {
				sel.SetHtml(rawHtml)
			}
		} else {
			node, result := getHtmlNode(fieldName, cont)
			if result == nil {
				sel.AppendClones(node)
			}
		}
		sel.RemoveAttr("gatscontent")
	})
	return result
}

func handleGatsAttributes(t *goquery.Selection, cont *context) (result error) {
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

func handleGatsText(t *goquery.Selection, cont *context) (result error) {
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

func handleGatsOmitTag(t *goquery.Selection) {
	t.Find("[gatsomittag]").Each(func(_ int, parent *goquery.Selection) {
		parent.Contents().Remove().InsertBefore(parent)
		parent.Remove()
	})
}
