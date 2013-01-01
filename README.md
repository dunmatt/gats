gats
====

Go Attribute Templating System

gats is exactly what it sounds like.  It is a templating system written in go that uses html/xml attributes for its control structures.

IMPORTANT NOTE:  This is my "hey lets learn Go" project.  Drop me a line and I'll see what I can do to support your use case.  Or just drop me a line to say you like it, that's fine too :-D .

## Dear lord, not another templating system, why??

Two reasons, first and foremost it's because all the other templating systems I could find had some funky non-html syntax which means that the templates couldn't render in a browser without processing.
Also, because I thought it would be a good starter project for learning Go.

The point has been raised that gats is remarkably similar to TAL (http://wiki.zope.org/ZPT/TALSpecification14), this was originally accidental, but now that I'm aware of TAL I've gone ahead and added the two basic constructs that TAL has that gats didn't: omit-tag and content.

As of now the only fundamental difference between TAL and gats is the complexity of the expressions.  gats is very simple, all expressions must be the name of a field in the data that the template will be rendered with.  Beyond that, the order of operations is different than TAL: in TAL variables are bound in ancestor first order, in gats it's just the opposite, there are no variables so it binds in descendents first order.  The attributes map as follows:

* tal:define == not necessary in gats
* tal:condition == gatsif
* tal:repeat == gatsrepeatover
* tal:content == gatscontent
* tal:replace == gatscontent && gatsomittag
* tal:attributes == gatsattributes
* tal:omit-tag == gatsomittag

## Requirements

gats uses my own fork of goquery (TODO: issue pull request once it all works), which in turn requires Go's experimental html package and cascadia, so these are all required.
It just should be a simple matter of following this guide: http://code.google.com/p/go-wiki/wiki/InstallingExp and then running

`go get github.com/dunmatt/gats`

and if not, please please please file a bug so I can correct this doc!

## All Attributes (and their semantics)

* **gatsattributes** : A map\[string\]string of attributes to programmatically give the element.
* **gatscontent** : Replace the children of the attributed element with either the raw string, or the html parse tree.
* **gatsif** : If the name given as a value is in the data and evaluates to true, show this element, otherwise remove it (and all its kids).
* **gatsomittag** : Replace the attributed element with its children.
* **gatsremove** : Remove the attributed element and all of its children from the DOM.
* **gatsrepeatover** : Populate a copy of the attributed element (and its children) with each item in the named array/slice (in order).
* **gatstext** : Replace the children of the attributed element with the named string.  This is much like gatscontent, except that it html escapes everything, so the string will display to the user instead of potentially becomming part of the DOM.

## Changelog

*    **v0.2.0** : Added gatsomittag and gatscontent
*    **v0.1.0** : Initial release.

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
    <div gatscontent="cont"></div>
    <table>
      <tr>
        <th gatsattributes="Titleattrs">Title</th>
        <th>Author</th>
        <th>Year</th>
        <th>Bibtex</th>
      </tr>
      <tr gatsrepeatover="Entries">
        <td><b gatsomittag="unr" gatstext="title">Doing stuff with items<b></td>
        <td gatsrepeatover="Entries">Me</td>
        <td gatstext="year">Soon</td>
        <td gatstext="bibtex">the</td>
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
	"exp/html"
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
	//cont       string
	cont *html.Node
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
		//cont:       "<hr/>",
		cont: &html.Node{
			Data: "hr",
			Type: html.ElementNode,
		},
	}
	d.Titleattrs["hi"] = "there"
	d.Titleattrs["test"] = "data"
	e := gats.RenderTemplateFile("exampleGood.html", &d, os.Stdout)
	if e != nil {
		fmt.Println(e)
	}
}
```

To yield:

```HTML
<html><head>
    <title>Yo Dawg</title>
  </head>
  <body>
    Yo dawg, stuff:
    <ul>
      <li>things</li>
      <li>misc</li>
    </ul>
    <div><hr/></div>
    <table>
      <tbody><tr>
        <th hi="there" test="data">Title</th>
        <th>Author</th>
        <th>Year</th>
        <th>Bibtex</th>
      </tr>
      <tr>
        <td>first</td>
        <td>Me</td><td>Me</td><td>Me</td>
        <td>2013</td>
        <td>meh</td>
      </tr><tr>
        <td>the matrix</td>
        <td>Me</td><td>Me</td><td>Me</td>
        <td>2013</td>
        <td>look over there ----&amp;gt;</td>
      </tr><tr>
        <td>the three amigos</td>
        <td>Me</td><td>Me</td><td>Me</td>
        <td>2013</td>
        <td>a plethora of laughs</td>
      </tr>

    </tbody></table>


</body></html>
```
