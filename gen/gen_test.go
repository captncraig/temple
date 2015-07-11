package gen

import (
	"testing"
)

func TestRender(t *testing.T) {
	P := GenParameters{"pkg1", "myTemplates", "", ""}
	D := map[string]string{"a": "qwerty", "B": "qqqqqqq"}
	_, err := render(P, D)
	if err != nil {
		t.Fatal(err)
	}
}

func TestFullGen(t *testing.T) {
	P := GenParameters{"main", "myTemplates", "../example/templates.go", "../example/templates"}
	Generate(P)
}
