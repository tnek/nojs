package serv

import (
	"log"
	"net/http"

	"github.com/google/uuid"
	"github.com/tnek/notes-site/model"
)

const (
	kDefaultProfilePath = "/static/img/default-profile.png"
)

type IndexPage struct {
	Title             string
	ProfilePictureURL string
	ReportInstance    string

	User  *model.User
	Notes []model.Note
}

func (a *App) IndexPage(w http.ResponseWriter, r *http.Request, u *model.User) {
	commonHeaders(w)

	notes, err := model.FetchNotes(u)
	if err != nil {
		log.Printf("failed to fetch notes: %v", err)
	}

	p := &IndexPage{
		Title:          a.AppName,
		User:           u,
		Notes:          notes,
		ReportInstance: uuid.New().String(),
	}
	if u.AvatarURL == "" {
		p.ProfilePictureURL = kDefaultProfilePath
	}

	if err := a.Tmpl.ExecuteTemplate(w, "index.html", p); err != nil {
		log.Printf("error rendering template for '/': %v", err)
	}
}

type ErrorPage struct {
	Title      string
	StatusText string
	Code       int
}

func (a *App) ErrorPage(w http.ResponseWriter, r *http.Request) {
	commonHeaders(w)

	if err := a.Tmpl.ExecuteTemplate(w, "error.html", &ErrorPage{
		Title:      "Error",
		StatusText: http.StatusText(http.StatusNotFound),
		Code:       http.StatusNotFound,
	}); err != nil {
		log.Printf("error rendering template for '/404': %v", err)
	}

}
