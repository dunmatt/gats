gats
====

Go Attribute Templating System

gats is exactly what it sounds like.  It is a templating system written in go that uses html/xml attributes for its control structures.

## Dear lord, not another templating system, why??

Two reasons, first and foremost it's because all the other templating systems I could find had some funky non-html syntax which means that the templates couldn't render in a browser without processing.
Also, because I thought it would be a good starter project for learning Go.

## Requirements

gats uses my own fork of goquery, which in turn requires Go's experimental html package and cascadia, so these are all required.
It just should be a simple matter of following this guide: http://code.google.com/p/go-wiki/wiki/InstallingExp and then running

`go get github.com/dunmatt/gats`

and if not, please please please file a bug so I can correct this doc!
