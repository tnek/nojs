package serv

import (
	"log"
	"net/http"

	"github.com/tnek/notes-site/model"
)

func (a *App) CreateNoteRequest(w http.ResponseWriter, r *http.Request, u *model.User) {
	r.ParseForm()
	title := r.Form.Get("submission-title")
	contents := r.Form.Get("submission-text")
	recipient := r.Form.Get("recipient")

	id, err := model.NewNote(u, title, contents)
	if err != nil {
		log.Printf("note creation failed: %v\n", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if recipient != "" {
		to, err := model.ShareNote(u, id, recipient)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		// OMIT
		if to.IsAdmin && a.Admin != nil {
			if err := a.Admin.Visit(u, recipient); err != nil {
				log.Println("visit err: %v", err)
				http.Error(w, "something went wrong with the challenge, retry or contact an admin on slack", http.StatusInternalServerError)
				return
			}
		}
		// OMITEND
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (a *App) DeleteNoteRequest(w http.ResponseWriter, r *http.Request, u *model.User) {
	noteID := r.URL.Query().Get("id")
	log.Printf("wtf %v; %v", noteID, r.URL.Query())

	ok, err := model.DeleteNote(u, noteID)
	if err != nil {
		log.Printf("deleteNote by %v of %v failed: %v\n", u.Name, noteID, err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if !ok {
		http.Error(w, "you don't own this note", http.StatusUnauthorized)
		return
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}
