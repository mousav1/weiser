package views

import (
	"errors"
	"html/template"
	"path/filepath"

	"github.com/flosch/pongo2"
	"github.com/gofiber/fiber/v2"
	"github.com/spf13/viper"
)

// ViewData is a struct that holds the data required for rendering a view
type ViewData struct {
	Title string
	Data  interface{}
}

// TemplateEngine is an enum that represents the supported template engines
type TemplateEngine int

const (
	HTMLTemplate TemplateEngine = iota
	Pongo2Template
)

// View is a struct that represents a template view
type ViewTemplate struct {
	Template *template.Template
}

// ViewPongo2 is a struct that represents a template view with Pongo2
type ViewPongo2 struct {
	Template *pongo2.Template
}

var templateEngineNames = map[TemplateEngine]string{
	HTMLTemplate:   "html",
	Pongo2Template: "pongo2",
	// Add other supported template engines here
}

// Render renders the view with the given data
func (v *ViewTemplate) Render(c *fiber.Ctx, data ViewData) error {
	c.Set("Content-Type", "text/html; charset=utf-8")
	viewData := struct{ Data ViewData }{Data: data}
	err := v.Template.Execute(c, viewData)
	if err != nil {
		return c.SendStatus(fiber.StatusInternalServerError)
	}
	return nil
}

// RenderPongo2 renders the view with the given data using Pongo2
func (v *ViewPongo2) RenderPongo2(c *fiber.Ctx, data ViewData) error {
	c.Set("Content-Type", "text/html; charset=utf-8")
	err := v.Template.ExecuteWriter(pongo2.Context{"Data": data}, c)
	if err != nil {
		return c.SendStatus(fiber.StatusInternalServerError)
	}
	return nil
}

// NewView creates a new view with the given template files and engine
func NewView(pattern string) (interface{}, error) {
	engine := viper.GetString("template_engine")
	if engine == "" {
		return nil, errors.New("template engine is not specified in configuration")
	}

	templateDir := viper.GetString("template_dir")
	if templateDir == "" {
		return nil, errors.New("template directory is not specified in configuration")
	}

	templatePath := filepath.Join(templateDir, pattern)

	var view interface{}
	var err error
	switch engine {
	case templateEngineNames[HTMLTemplate]:
		tpl := template.Must(template.ParseFiles(templatePath))
		view = &ViewTemplate{
			Template: tpl,
		}
	case templateEngineNames[Pongo2Template]:
		var tpl *pongo2.Template
		tpl, err := pongo2.FromFile(templatePath)
		if err != nil {
			return nil, err
		}
		view = &ViewPongo2{
			Template: tpl,
		}
	default:
		return nil, errors.New("unsupported template engine")
	}

	return view, err
}

// view renders a view with the given data
func View(c *fiber.Ctx, data ViewData, files string) error {
	viewI, err := NewView(files)
	if err != nil {
		return err
	}

	switch view := viewI.(type) {
	case *ViewTemplate:
		err = view.Render(c, data)
	case *ViewPongo2:
		err = view.RenderPongo2(c, data)
	default:
		return errors.New("unsupported view type")
	}

	if err != nil {
		return err
	}

	return nil
}
