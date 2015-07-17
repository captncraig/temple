package temple

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"html/template"
	"io"
	"path/filepath"
)

type TemplateStore interface {
	GetTemplate(name string) (*template.Template, error)
	Execute(w io.Writer, ctx interface{}, name string) error
	ExecuteMaster(w io.Writer, ctx interface{}, master, child string) error
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

func (s *staticTemplateStore) Execute(w io.Writer, ctx interface{}, name string) error {
	return execute(s, w, ctx, name)
}

func (s *staticTemplateStore) ExecuteMaster(w io.Writer, ctx interface{}, master, child string) error {
	return executeMaster(s, w, ctx, master, child)
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

func (d *devTemplateStore) Execute(w io.Writer, ctx interface{}, name string) error {
	return execute(d, w, ctx, name)
}

func (d *devTemplateStore) ExecuteMaster(w io.Writer, ctx interface{}, master, child string) error {
	return executeMaster(d, w, ctx, master, child)
}

const (
	// number of buffers to keep in rotation
	numBuffers = 100
	// initial size to allocate for new buffers
	initialSize = 1024 * 64 // 64KB
	// size at which we don't put buffers back in pool
	maxBufferSize = 1024 * 1024 // 1M
)

var buffers = make(chan *bytes.Buffer, numBuffers)

func getBuffer() *bytes.Buffer {
	var b *bytes.Buffer
	select {
	case b = <-buffers:
		return b
	default:
	}
	arr := make([]byte, initialSize)
	b = bytes.NewBuffer(arr)
	return b
}

func putBuffer(b *bytes.Buffer) {
	b.Reset()
	if cap(b.Bytes()) > maxBufferSize {
		return
	}
	select {
	case buffers <- b:
		return
	default:
	}
}

func execute(store TemplateStore, w io.Writer, ctx interface{}, name string) error {
	tpl, err := store.GetTemplate(name)
	if err != nil {
		return err
	}
	buf := getBuffer()
	defer putBuffer(buf)
	err = tpl.Execute(buf, ctx)
	if err != nil {
		return err
	}
	_, err = io.Copy(w, buf)
	return err
}

func executeMaster(store TemplateStore, w io.Writer, ctx interface{}, master, child string) error {
	masterTpl, err := store.GetTemplate(master)
	if err != nil {
		return err
	}
	childTpl, err := store.GetTemplate(child)
	if err != nil {
		return err
	}
	// Render child template
	childBuffer := getBuffer()
	defer putBuffer(childBuffer)
	err = childTpl.Execute(childBuffer, ctx)
	if err != nil {
		return err
	}
	masterBuffer := getBuffer()
	defer putBuffer(masterBuffer)
	// Wrap up original ctx and child content for master template
	masterCtx := &struct {
		Data    interface{}
		Content string
	}{ctx, childBuffer.String()}
	err = masterTpl.Execute(w, masterCtx)
	if err != nil {
		return err
	}
	_, err = io.Copy(w, masterBuffer)
	return err
}
