package main

import (
	"html/template"
	"io"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

// calling templates
type Templates struct {
	templates *template.Template
}

// echo boiler plate 
func(t *Templates) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return t.templates.ExecuteTemplate(w, name, data)
}

// template call 
func newTemplate () *Templates {
	return &Templates{
		templates: template.Must(template.ParseGlob("views/*.html")),
	}
}

type Contact struct {
    Name  string
    Email string
}

type Data struct {
    Contacts []Contact
}

func NewData() *Data {
	return &Data{
		Contacts: []Contact{
			NewContact("Clara", "cd@gmail.com"),
			NewContact("John", "jd@gmail.com"),
		},
	}
}

func NewContact(name, email string) Contact {
    return Contact{
        Name: name,
        Email: email,
    }
}



func main() {

    e := echo.New()

    data := NewData()

    e.Renderer = newTemplate()
    e.Use(middleware.Logger())

    e.GET("/", func(c echo.Context) error {
        return c.Render(200, "index.html", data)
    })

    e.POST("/contacts", func(c echo.Context) error {
        name := c.FormValue("name")
        email := c.FormValue("email")

        data.Contacts = append(data.Contacts, NewContact(name, email))
        return c.Render(200, "index.html", data)
    })

    e.Logger.Fatal(e.Start(":42069"))
}