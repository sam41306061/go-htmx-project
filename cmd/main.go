package main

import (
	"html/template"
	"io"

	"github.com/labstack/echo/v4"
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

type Contacts = []Contact

func(d Data) hasEmail(email string) bool {
    for _, contact := range d.Contacts {
        if contact.Email == email {
         return true
        }
    }
    return false
}

type Data struct {
    Contacts Contacts
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


type FormData struct {
    Values map[string]string
    Errors map[string]string
}

func newFormData() FormData {
    return FormData {
        Values: make(map[string]string),
        Errors: make(map[string]string),
    }
}

// index page 
type Page struct{
    Data Data 
    Form FormData
}

func newPage() Page {
    return Page {
        Data: *NewData(),
        Form: newFormData(),
    }
}

func main() {

    e := echo.New()
    page := newPage()
    e.Renderer = newTemplate()

    e.GET("/", func(c echo.Context) error {
        return c.Render(200, "index", page)
    }) 
    // handles contact updating 
    e.POST("/contacts", func(c echo.Context) error {
        name := c.FormValue("name")
        email := c.FormValue("email")

        if page.Data.hasEmail(email) {
            formData := newFormData()
            formData.Values["name"] = name 
            formData.Values["email"] = email
            formData.Errors["email"] = "Email already exists"

            return c.Render(400, "form", formData)
        }

        page.Data.Contacts = append(page.Data.Contacts, NewContact(name, email))
        return c.Render(200, "display", page)
    })

    e.Logger.Fatal(e.Start(":42069"))
}