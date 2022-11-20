package serv

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"text/template"

	"github.com/gorilla/mux"
	"github.com/tnek/notes-site/assets"
	"github.com/tnek/notes-site/config"

	// OMIT
	"github.com/tnek/notes-site/admin"
	// OMITEND
)

type App struct {
	AppName string
	Host    string
	Port    int

	Tmpl *template.Template

	// OMIT
	Admin *admin.Admin
	// OMITEND
}

func NewApp(c *config.AppConfig) (*App, error) {
	a := &App{}

	if c.AppName == "" {
		a.AppName = "NoJS"
	} else {
		a.AppName = c.AppName
	}

	a.Host = c.Host
	a.Port = c.Port
	a.Tmpl = template.Must(template.ParseFiles(c.Templates...))
	return a, nil
}

func (a *App) Run(c *config.AppConfig) error {
	r := mux.NewRouter().StrictSlash(true)
	r.PathPrefix("/static/").Handler(assets.StaticHTTPServ).Methods("GET")

	r.HandleFunc("/login", a.LoginPage).Methods("GET", "POST")
	r.HandleFunc("/register", a.RegisterPage).Methods("GET", "POST")
	r.HandleFunc("/logout", a.Logout).Methods("GET")

	r.Handle("/", NewEnsureAuth(a.IndexPage)).Methods("GET")
	r.Handle("/delete", NewEnsureAuth(a.DeleteNoteRequest)).Methods("GET")
	r.Handle("/post", NewEnsureAuth(a.CreateNoteRequest)).Methods("POST")

	s := http.Server{
		Addr:    fmt.Sprintf("%v:%v", a.Host, a.Port),
		Handler: r,
	}

	return s.ListenAndServe()
}

func parseBody(r *http.Request, req interface{}) error {
	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return err
	}
	defer r.Body.Close()

	if err := json.Unmarshal(b, req); err != nil {
		return err
	}
	return nil
}
