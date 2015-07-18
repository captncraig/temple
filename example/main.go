package main

//go:generate templeGen -dir templates -pkg main -var myTemplates -o templates.go

import (
	"flag"
	"log"
	"net/http"

	"github.com/captncraig/temple"
)

var devMode = flag.Bool("dev", false, "activate dev mode for templates")

func main() {
	flag.Parse()
	templateManager, err := temple.New(*devMode, myTemplates, "templates")
	if err != nil {
		log.Fatal(err)
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		ctx := homePageContext{
			baseContext{"Home"},
			[]string{"foo", "bar", "baz"},
		}
		// homepage composes shared templates by directly referencing "header" and "footer" templates.
		err := templateManager.Execute(w, ctx, "homepage")
		if err != nil {
			http.Error(w, err.Error(), 500)
		}
	})

	http.HandleFunc("/mc", func(w http.ResponseWriter, r *http.Request) {
		ctx := masterChildContext{
			baseContext{"Master / Child"},
			"AAAAA",
			"BBBBB",
		}
		// master / child templates let the library compose things for you
		// so no explicit reference to other templates is needed.
		err := templateManager.ExecuteMaster(w, ctx, "master", "child")
		if err != nil {
			http.Error(w, err.Error(), 500)
		}
	})
	http.ListenAndServe(":5555", nil)

}

// base context that every page needs because the header/footer reference it.
// every one of my single-page contexts will embed this so that the fields are always availible to the header.
type baseContext struct {
	PageTitle string
}

type homePageContext struct {
	baseContext
	Items []string
}

type masterChildContext struct {
	baseContext
	Foobar string
	Xyz    string
}
