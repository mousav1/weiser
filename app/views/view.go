package views

import (
	"errors"
	"fmt"
	"html/template"
	"path/filepath"

	"github.com/flosch/pongo2"
	"github.com/gofiber/fiber/v2"
	"github.com/spf13/viper"
)

// ViewData holds the data for rendering a view
type ViewData struct {
	Title string
	Data  interface{}
}

// TemplateEngine represents the supported template engines
type TemplateEngine string

const (
	HTMLTemplate   TemplateEngine = "html"
	Pongo2Template TemplateEngine = "pongo2"
)

// View is an interface that defines the methods for rendering a view
type View interface {
	Render(c *fiber.Ctx, data ViewData) error
}

// HTMLView renders HTML templates
type HTMLView struct {
	Template *template.Template
}

// Pongo2View renders Pongo2 templates
type Pongo2View struct {
	Template *pongo2.Template
}

// Render renders the HTML view with the given data
func (v *HTMLView) Render(c *fiber.Ctx, data ViewData) error {
	c.Set("Content-Type", "text/html; charset=utf-8")
	viewData := struct{ Data ViewData }{Data: data}
	if err := v.Template.Execute(c, viewData); err != nil {
		return fmt.Errorf("could not render HTML template: %w", err)
	}
	return nil
}

// Render renders the Pongo2 view with the given data
func (v *Pongo2View) Render(c *fiber.Ctx, data ViewData) error {
	c.Set("Content-Type", "text/html; charset=utf-8")
	if err := v.Template.ExecuteWriter(pongo2.Context{"Data": data}, c); err != nil {
		return fmt.Errorf("could not render Pongo2 template: %w", err)
	}
	return nil
}

// NewView creates a new view based on the template engine configuration
func NewView(pattern string) (View, error) {
	engine := viper.GetString("template_engine")
	if engine == "" {
		return nil, errors.New("template engine is not specified in configuration")
	}

	templateDir := viper.GetString("template_dir")
	if templateDir == "" {
		return nil, errors.New("template directory is not specified in configuration")
	}

	templatePath := filepath.Join(templateDir, pattern)

	switch TemplateEngine(engine) {
	case HTMLTemplate:
		tpl := template.Must(template.ParseFiles(templatePath))
		return &HTMLView{Template: tpl}, nil
	case Pongo2Template:
		tpl, err := pongo2.FromFile(templatePath)
		if err != nil {
			return nil, fmt.Errorf("could not parse Pongo2 template: %w", err)
		}
		return &Pongo2View{Template: tpl}, nil
	default:
		return nil, errors.New("unsupported template engine")
	}
}

// RenderView renders a view with the given data
func RenderView(c *fiber.Ctx, data ViewData, files string) error {
	view, err := NewView(files)
	if err != nil {
		return fmt.Errorf("could not create view: %w", err)
	}

	if err := view.Render(c, data); err != nil {
		return fmt.Errorf("could not render view: %w", err)
	}

	return nil
}
