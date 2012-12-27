package gats

import (
	"exp/html"
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
	return nil
}

func handleGatsRemoves(t *goquery.Document) {
	t.Find("[gatsremove]").Remove()
}
