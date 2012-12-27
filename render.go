package gats

import (
	"exp/html"
	"github.com/dunmatt/goquery"
	"io"
	"os"
)

func RenderTemplateFile(filename string, data interface{}, out io.Writer) error {
	f, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer f.Close()
	rootNode, err := html.Parse(f)
	if err != nil {
		return err
	}
	template := goquery.NewDocumentFromNode(rootNode)
	err = fillInTemplateTemplate(template, data)
	if err == nil {
		html.Render(out, rootNode)
	}
	return err
}

func fillInTemplate(t *goquery.Document, data interface{}) error {
	return nil
}
