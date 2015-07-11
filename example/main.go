package main

import (
	"fmt"
	"github.com/captncraig/temple"
	"log"
)

func main() {
	templateManager, err := temple.New(false, myTemplates, "templates")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(templateManager.GetTemplate("a.tmpl"))
}
