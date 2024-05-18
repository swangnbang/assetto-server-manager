package views

import (
	"html/template"
	"io/fs"
	"path/filepath"
	"strings"
)

type TemplateLoader struct {
	pages, partials []string
}

func (t *TemplateLoader) Init() error {
	return fs.WalkDir(static, ".", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if d.IsDir() {
			return nil
		}

		n := d.Name()

		switch {
		case strings.HasPrefix(n, "/pages/"):
			t.pages = append(t.pages, n)
		case strings.HasPrefix(n, "/partials/"):
			t.partials = append(t.partials, n)
		}

		return nil
	})
}

func (t *TemplateLoader) fileContents(name string) (string, error) {
	if b, err := static.ReadFile(name); err != nil {
		return "", err
	} else {
		return string(b), err
	}
}

func (t *TemplateLoader) Templates(funcs template.FuncMap) (map[string]*template.Template, error) {
	templates := make(map[string]*template.Template)

	templateData, err := t.fileContents("/layout/base.html")
	if err != nil {
		return nil, err
	}

	for _, partial := range t.partials {
		contents, err := t.fileContents(partial)
		if err != nil {
			return nil, err
		}

		templateData += contents
	}

	for _, page := range t.pages {
		pageData := templateData

		pageText, err := t.fileContents(page)
		if err != nil {
			return nil, err
		}

		pageData += pageText

		t, err := template.New(page).Funcs(funcs).Parse(pageData)
		if err != nil {
			return nil, err
		}

		templates[strings.TrimPrefix(filepath.ToSlash(page), "/pages/")] = t
	}

	return templates, nil
}
