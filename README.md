gats
====

Go Attribute Templating System

gats is exactly what it sounds like.  It is a templating system written in go that uses html/xml attributes for its control structures.

IMPORTANT NOTE:  This is my "hey lets learn Go" project.  Drop me a line and I'll see what I can do to support your use case.  Or just drop me a line to say you like it, that's fine too :-D .

## Dear lord, not another templating system, why??

Two reasons, first and foremost it's because all the other templating systems I could find had some funky non-html syntax which means that the templates couldn't render in a browser without processing.
Also, because I thought it would be a good starter project for learning Go.

## Requirements

gats uses my own fork of goquery (TODO: issue pull request once it all works), which in turn requires Go's experimental html package and cascadia, so these are all required.
It just should be a simple matter of following this guide: http://code.google.com/p/go-wiki/wiki/InstallingExp and then running

`go get github.com/dunmatt/gats`

and if not, please please please file a bug so I can correct this doc!

## Example Usage

Just decorate some html however you want (yay, it even supports nesting!):

```HTML
<html>
  <head>
    <title>Yo Dawg</title>
  </head>
  <body>
    Yo dawg, stuff:
    <ul gatsif="showul">
      <li>things</li>
      <li>misc</li>
    </ul>
    <table>
      <tr>
        <th gatsattributes="Titleattrs">Title</th>
        <th>Author</th>
        <th>Year</th>
        <th>Bibtex</th>
      </tr>
      <tr gatsrepeatover="Entries">
        <td gatstext="title">Doing stuff with items</td>
        <td gatsrepeatover="Entries">Me</td>
        <td gatstext="year">Soon</td>
        <td gatstext="bibtex"></td>
      </tr>
      <tr gatsremove="true">
        <td>This</td>
        <td>row</td>
        <td>should</td>
        <td>disappear</td>
      </tr>
    </table>
  </body>
</html>
```

Then build a struct with names that match the elements in the template and call RenderTemplateFile.  Note that under the covers GATS uses reflection, so any field that gets iterated over (gatsrepeatover and gatsattributes) absolutely must be public.

```Go
package main

import (
	"fmt"
	"github.com/dunmatt/gats"
	"os"
)

type pub struct {
	title  string
	bibtex string
}

type data struct {
	showul     bool
	Entries    []pub
	Titleattrs map[string]string
	year       string
}

func main() {
	d := data{
		showul: true,
		Entries: []pub{
			{title: "first", bibtex: "meh"},
			{title: "the matrix", bibtex: "look over there ---->"},
			{title: "the three amigos", bibtex: "a plethora of laughs"}},
		Titleattrs: make(map[string]string),
		year:       "2013",
	}
	d.Titleattrs["hi"] = "there"
	d.Titleattrs["test"] = "data"
	e := gats.RenderTemplateFile("exampleGood.html", &d, os.Stdout)
	if e != nil {
		fmt.Println(e)
	}
}
```

## All Attributes (and their semantics)

* **gatsattributes** : A map\[string\]string of attributes to programmatically give the element.
* **gatsif** : If the name given as a value is in the data and evaluates to true, show this element, otherwise remove it (and all its kids).
* **gatsremove** : Remove the attributed element and all of its children from the DOM.
* **gatsrepeatover** : Populate a copy of the attributed element (and its children) with each item in the named array/slice (in order).
* **gatstext** : Replace the children of the attributed element with the named string.
