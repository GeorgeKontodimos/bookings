package render

import (
	"bytes"
	"html/template"
	"log"
	"net/http"
	"path/filepath"

	"github.com/georgekontodimos/bookings/pkg/config"
	"github.com/georgekontodimos/bookings/pkg/models"
)

var app *config.AppConfig

// NewTemplates sets the config for the template package
func NewTemplates(a *config.AppConfig) {
	app = a
}

func AddDefaultData(td*models.TemplateData) *models.TemplateData{
	// here add default data
	return td
}

// RenderTemplate renders template
func RenderTemplate(w http.ResponseWriter, tmpl string, td *models.TemplateData) {
	// get the template cache from the app config
	var tc map[string]*template.Template
	if app.UseCache {
		tc = app.TemplateCache
	} else {
		tc, _ = CreateTemplateCache()
	}

	// get requested template from cache
	t, ok := tc[tmpl]
	if !ok {
		log.Fatal("Could not get template from template cache")
	}

	buf := new(bytes.Buffer)
	td = AddDefaultData(td)
	err := t.Execute(buf, td)
	if err != nil {
		log.Println(err)
	}
	// render the template
	_, err = buf.WriteTo(w)
	if err != nil {
		log.Println(err)
	}
}

func CreateTemplateCache() (map[string]*template.Template, error) {
	myCache := map[string]*template.Template{}

	// get all of the files named *.page.tmpl from the ./templates
	pages, err := filepath.Glob("./templates/*.page.tmpl")
	if err != nil {
		return myCache, err
	}

	// range through all files ending with *.page.tmpl
	for _, page := range pages {
		name := filepath.Base(page)                    // get only the filename
		ts, err := template.New(name).ParseFiles(page) // give a name and parse
		if err != nil {
			return myCache, err
		}
		matches, err := filepath.Glob("./templates/*.layout.tmpl")
		if err != nil {
			return myCache, err
		}

		if len(matches) > 0 {
			ts, err = ts.ParseGlob("./templates/*.layout.tmpl")
			if err != nil {
				return myCache, err
			}
		}
		myCache[name] = ts
	}
	return myCache, nil
}

// old code
/*

func RenderTemplate(w http.ResponseWriter, tmpl string) {
	old code
	parsedTemplate, _ := template.ParseFiles("./templates/"+tmpl, "./templates/base.layout.tmpl")
	err := parsedTemplate.Execute(w, nil)
	if err != nil {
		fmt.Println("error parsing template:", err)
		return
	}
}


var tc = make(map[string]*template.Template)

func RenderTemplate(w http.ResponseWriter, t string) {
	var tmpl *template.Template
	var err error
	// check if we already have the temlate in our cache
	_, inMap := tc[t]
	if !inMap {
		// need to create template
		log.Println("creating temlate and adding to cache")
		err = createTemplateCache(t)
		if err != nil{
			log.Panicln(err)
		}
	} else {
		// we have the template in the cache
		log.Println("using cached temlate")
	}

	tmpl = tc[t]
	err = tmpl.Execute(w, nil)
	if err != nil{
			log.Panicln(err)
		}
}
func createTemplateCache(t string) error {
	templates := []string{
		fmt.Sprintf("./templates/%s", t),
		"./templates/base.layout.tmpl",
	}
	// parse the template
	tmpl, err := template.ParseFiles(templates...)
	if err !=nil{
		return err
	}
	// add template to cache (map)
	tc[t] = tmpl

	return nil
}
*/
