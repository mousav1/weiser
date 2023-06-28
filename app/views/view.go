package views

import (
	"fmt"
	"log"
	"text/template"

	"github.com/flosch/pongo2"
	"github.com/gofiber/fiber/v2"
)

// View is a struct that represents a template view
type View struct {
	Template *template.Template
	Layout   string
	Data     interface{}
}

// View is a struct that represents a template view
type Viewpongo2 struct {
	Template *pongo2.Template
	Layout   string
	Data     interface{}
}

// Render renders the view with the given data
func (v *View) Render(w *fiber.Ctx, data interface{}) error {
	v.Data = data

	if v.Layout != "" {
		err := v.Template.ExecuteTemplate(w, v.Layout, v)
		if err != nil {
			return err
		}
	} else {
		err := v.Template.Execute(w, v)
		if err != nil {
			return err
		}
	}

	return nil
}

// Render renders the view with the given data
func (v *Viewpongo2) Renderpongo2(w *fiber.Ctx, data interface{}) error {
	v.Data = data

	if v.Layout != "" {
		t, err := pongo2.FromFile("views/layouts/" + v.Layout + ".html")
		if err != nil {
			return err
		}
		err = t.ExecuteWriter(pongo2.Context{"View": v, "Data": data}, w)
		if err != nil {
			return err
		}
	} else {
		err := v.Template.ExecuteWriter(pongo2.Context{"View": v, "Data": data}, w)
		if err != nil {
			return err
		}
	}

	return nil
}

// NewView creates a new view with the given base template and file names
func NewView(layout string, files ...string) *View {
	// Add the layout file to the list of files
	files = append(files, fmt.Sprintf("resources/layouts/%s.html", layout))

	// Create a new template with the given files
	t, err := template.ParseFiles(files...)
	if err != nil {
		log.Fatalf("Error parsing template files: %v", err)
	}

	return &View{
		Template: t,
		Layout:   layout,
	}
}
