package temple

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"html/template"
	"path/filepath"
)

type TemplateStore interface {
	GetTemplate(name string) (*template.Template, error)
	Execute(ctx interface{}, templates ...string) ([]byte, error)
}

func New(devMode bool, storedTemplates map[string]string, dir string) (TemplateStore, error) {
	if !devMode {
		return newStatic(storedTemplates)
	}
	return &devTemplateStore{dir}, nil
}

type staticTemplateStore struct {
	templates *template.Template
}

func (s *staticTemplateStore) GetTemplate(name string) (*template.Template, error) {
	if t := s.templates.Lookup(name); t != nil {
		return t, nil
	}
	return nil, fmt.Errorf("Template `%s` not found", name)
}

func (s *staticTemplateStore) Execute(ctx interface{}, templates ...string) ([]byte, error) {
	return execute(s, ctx, templates...)
}

func newStatic(storedTemplates map[string]string) (*staticTemplateStore, error) {
	t, err := mapToTemplate(storedTemplates)
	if err != nil {
		return nil, err
	}
	return &staticTemplateStore{t}, nil
}

func mapToTemplate(storedTemplates map[string]string) (*template.Template, error) {
	var t *template.Template
	for name, text := range storedTemplates {
		var tpl *template.Template
		if t == nil {
			tpl = template.New(name)
			t = tpl
		} else {
			tpl = t.New(name)
		}
		decoded, err := base64.StdEncoding.DecodeString(text)
		if err != nil {
			return nil, err
		}
		tpl, err = tpl.Parse(string(decoded))
		if err != nil {
			return nil, err
		}
	}
	return t, nil
}

type devTemplateStore struct {
	dir string
}

func (d *devTemplateStore) GetTemplate(name string) (*template.Template, error) {
	filePath := filepath.Join(d.dir, "*")
	tpl, err := template.ParseGlob(filePath)
	if err != nil {
		return nil, err
	}
	if t := tpl.Lookup(name); t != nil {
		return t, nil
	}
	return nil, fmt.Errorf("Template `%s` not found", name)
}

func (d *devTemplateStore) Execute(ctx interface{}, templates ...string) ([]byte, error) {
	return execute(d, ctx, templates...)
}

func execute(store TemplateStore, ctx interface{}, templates ...string) ([]byte, error) {
	buf := &bytes.Buffer{}
	var tpl *template.Template
	for _, name := range templates {
		var thisTpl *template.Template
		if tpl == nil {
			tpl, err := store.GetTemplate(name)
			if err != nil {
				return nil, err
			}
			thisTpl = tpl
		} else {
			thisTpl = tpl.Lookup(name)
			if thisTpl == nil {
				return nil, fmt.Errorf("Template `%s` not found", name)
			}
		}
		err := thisTpl.Execute(buf, ctx)
		if err != nil {
			return nil, err
		}
	}
	return buf.Bytes(), nil
}
