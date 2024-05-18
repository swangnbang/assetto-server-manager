package views

import (
	"html/template"
	"io/fs"
	"path/filepath"
	"strings"

	"github.com/sirupsen/logrus"
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

		logrus.WithFields(logrus.Fields{
			"path": path,
		}).Info("loading template")

		switch {
		case strings.HasPrefix(path, "pages/"):
			t.pages = append(t.pages, path)
		case strings.HasPrefix(path, "partials/"):
			t.partials = append(t.partials, path)
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

	templateData, err := t.fileContents("layout/base.html")
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
