package gats

import (
	"exp/html"
	"fmt"
	"github.com/dunmatt/goquery"
	"io"
	"os"
	"strings"
)

func RenderTemplateFile(filename string, data interface{}, out io.Writer) error {
	rootNode, err := parseFile(filename)
	if err != nil {
		return err
	}
	cont := makeContext(data, nil)
	template := goquery.NewDocumentFromNode(rootNode).Find("html") // wrap goquery around the DOM

	handleGatsRemoves(template)           // take out everything that definitely won't be shown
	err = handleGatsTranscludes(template) // transclude in all the sub-templates
	if err != nil {
		return err
	}
	err = fillInTemplate(template, cont) // process the template
	handleGatsOmitTag(template)          // make sure this one comes last, it can interfere with other attributes
	if err == nil {
		html.Render(out, rootNode) // render the DOM back to html and send it off
	}
	return err
}

func parseFile(filename string) (*html.Node, error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	return html.Parse(f)
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
	e = handleGatsAttribute(scope, cont)
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

func handleGats(t *goquery.Selection, selector string, meat func(string, *goquery.Selection)) {
	attribName := selector[1 : len(selector)-1]
	t.Find(selector).Each(func(_ int, sel *goquery.Selection) {
		fieldName, _ := sel.Attr(attribName)
		meat(fieldName, sel)
		sel.RemoveAttr(attribName)
	})
	return
}

func splitString(val string) (string, string, error) {
	index := strings.Index(val, ";")
	if index == -1 {
		return "", "", fmt.Errorf("Invalid transclude string '%v', contains no semicolon.", val)
	}
	return val[:index], val[index+1:], nil
}

func handleGatsTranscludes(scope *goquery.Selection) (result error) {
	for scope.Find("[gatstransclude]").Length() > 0 && result == nil {
		handleGats(scope, "[gatstransclude]", func(ts string, sel *goquery.Selection) {
			if result != nil {
				return
			}
			filename, selector, res := splitString(ts)
			if res != nil {
				result = res
				return
			}
			rootNode, res := parseFile(filename)
			if res != nil {
				result = res
				return
			}
			newKids := goquery.NewDocumentFromNode(rootNode).Find(selector)
			sel.Empty().Append(newKids)
		})
	}
	return result
}

func handleGatsIf(t *goquery.Selection, cont *context) (result error) {
	handleGats(t, "[gatsif]", func(fieldName string, sel *goquery.Selection) {
		show, res := getBool(fieldName, cont)
		if !show {
			sel.Remove()
		}
		result = res
	})
	return result
}

func handleGatsContent(t *goquery.Selection, cont *context) (result error) {
	handleGats(t, "[gatscontent]", func(fieldName string, sel *goquery.Selection) {
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
	})
	return result
}

func handleGatsAttribute(t *goquery.Selection, cont *context) (result error) {
	handleGats(t, "[gatsattribute]", func(val string, sel *goquery.Selection) {
		attribute, fieldName, err := splitString(val)
		if err != nil {
			result = err
			return
		}
		attributeValue, err := getString(fieldName, cont)
		if err != nil {
			result = err
			return
		}
		sel.SetAttr(attribute, attributeValue)
	})
	return result
}

func handleGatsAttributes(t *goquery.Selection, cont *context) (result error) {
	handleGats(t, "[gatsattributes]", func(fieldName string, sel *goquery.Selection) {
		attribs, result := getStringMap(fieldName, cont)
		if result == nil {
			for k, v := range attribs {
				sel.SetAttr(k, v)
			}
		}
	})
	return result
}

func handleGatsText(t *goquery.Selection, cont *context) (result error) {
	handleGats(t, "[gatstext]", func(fieldName string, sel *goquery.Selection) {
		text, result := getString(fieldName, cont)
		if result == nil {
			sel.SetText(text)
		}
	})
	return result
}

func handleGatsOmitTag(t *goquery.Selection) {
	t.Find("[gatsomittag]").Each(func(_ int, parent *goquery.Selection) {
		parent.Contents().Remove().InsertBefore(parent)
		parent.Remove()
	})
}
