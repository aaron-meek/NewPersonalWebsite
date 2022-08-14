package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/urhumantoast/NewPersonalWebsite/pkg/config"
	"github.com/urhumantoast/NewPersonalWebsite/pkg/handlers"
	"github.com/urhumantoast/NewPersonalWebsite/pkg/render"

	"github.com/alexedwards/scs/v2"
)

const portNumber = ":80"

var app config.AppConfig
var session *scs.SessionManager

func main() {
	app.InProduction = false

	session = scs.New()
	session.Lifetime = 24 * time.Hour
	session.Cookie.Persist = true
	session.Cookie.SameSite = http.SameSiteLaxMode
	session.Cookie.Secure = app.InProduction

	app.Session = session

	tc, err := render.CreateTemplateCache()
	if err != nil {
		log.Fatal("cannot create template cache")
	}

	app.TemplateCache = tc
	app.UseCache = false

	repo := handlers.NewRepo(&app)
	handlers.NewHandlers(repo)

	render.NewTemplates(&app)

	fmt.Printf("Starting application on port %s\n", portNumber)

	srv := &http.Server{
		Addr:    portNumber,
		Handler: routes(&app),
	}

	err = srv.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}

}
