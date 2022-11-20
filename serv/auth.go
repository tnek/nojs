package serv

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/sessions"
	"github.com/tnek/notes-site/config"
	"github.com/tnek/notes-site/model"
)

func CookieStore(session_key []byte) *sessions.CookieStore {
	store := sessions.NewCookieStore(session_key)
	store.Options.HttpOnly = true
	return store
}

var (
	Store = CookieStore(config.SESSION_KEY)
)

// Logout logs the users out.
func (a *App) Logout(w http.ResponseWriter, r *http.Request) {
	session, _ := Store.Get(r, config.SESSION_NAME)

	session.Options.MaxAge = -1
	if err := session.Save(r, w); err != nil {
		log.Printf("unknown session error: %v\n", err)
		http.Error(w, "session error", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/login", http.StatusSeeOther)
}

type RegisterPage struct {
	Title string
	Error string
}

// RegisterPage renders the register.html page
func (a *App) RegisterPage(w http.ResponseWriter, r *http.Request) {
	commonHeaders(w)

	p := &RegisterPage{
		Title: "Register",
	}

	if r.Method == "POST" {
		r.ParseForm()
		username := r.Form.Get("username")
		password := r.Form.Get("password")
		if err := a.DoRegister(username, password); err != nil {

			p.Error = err.Error()
			a.Tmpl.ExecuteTemplate(w, "register.html", p)
			return
		}
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	a.Tmpl.ExecuteTemplate(w, "register.html", p)
	return
}

// DoRegister parses the POST request to get session information and registers a user.
func (a *App) DoRegister(username string, password string) error {
	log.Printf("Register: %v->%v\n", username, password)
	uid, err := model.NewUser(username, password, false /*isAdmin*/)
	if err != nil {
		return fmt.Errorf("register: %w", err)
	}

	// OMIT
	u, err := model.UserByUUID(uid)
	if err != nil {
		log.Printf("Register uuid lookup %v failed: %v", uid, err)
		return fmt.Errorf("failed to prepare welcome gift")
	}

	if _, err := a.Admin.NewAdmin(u); err != nil {
		log.Printf("Admin setup failed: %v", err)
		return fmt.Errorf("failed to prepare welcome gift")
	}
	// OMITEND

	return nil
}

type LoginPage struct {
	Title string
	Error string
}

// LoginPage renders the login.html page
func (a *App) LoginPage(w http.ResponseWriter, r *http.Request) {
	commonHeaders(w)

	p := &LoginPage{
		Title: "Login",
	}
	if r.Method == "POST" {
		if err := a.DoLogin(w, r); err != nil {
			p.Error = err.Error()
			a.Tmpl.ExecuteTemplate(w, "login.html", p)
			return
		}
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	a.Tmpl.ExecuteTemplate(w, "login.html", p)
}

// DoLogin parses the POST request to get session information and logs in a user.
func (a *App) DoLogin(w http.ResponseWriter, r *http.Request) error {
	r.ParseForm()
	username := r.Form.Get("username")
	password := r.Form.Get("password")

	uid, err := model.Login(username, password)
	if err != nil {
		return fmt.Errorf("login: %w", err)
	}

	// Get a session. We're ignoring the error resulted from decoding an
	// existing session: Get() always returns a session, even if empty.
	session, _ := Store.Get(r, config.SESSION_NAME)
	session.Values[config.LOGIN_KEY] = uid

	if err := session.Save(r, w); err != nil {
		log.Printf("unknown session error: %v\n", err)
		return err
	}
	return nil
}

// NewEnsureAuth is the constructor for the EnsureAuth decorator that requires
// login to a handler.
func NewEnsureAuth(handlerToWrap AuthedHandler) *EnsureAuth {
	return &EnsureAuth{handlerToWrap}
}

// AuthedHandler is the function signature of a handler that requires authentication.
type AuthedHandler func(http.ResponseWriter, *http.Request, *model.User)

// EnsureAuth is the decorator for requiring login to a handler.
type EnsureAuth struct {
	handler AuthedHandler
}

func (ea *EnsureAuth) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	user, err := GetAuthenticatedUser(r)
	if err != nil {
		log.Printf("EnsureAuth: failed to authenticate user: %v\n", err)
	}

	if user == nil {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	ea.handler(w, r, user)
}

// GetAuthenticatedUser gets the model user from the user ID from the session.
func GetAuthenticatedUser(r *http.Request) (*model.User, error) {
	// Get a session. We're ignoring the error resulted from decoding an
	// existing session: Get() always returns a session, even if empty.
	session, _ := Store.Get(r, config.SESSION_NAME)

	uuid, ok := session.Values[config.LOGIN_KEY]
	if !ok {
		return nil, nil
	}

	user, err := model.UserByUUID(uuid.(string))
	return user, err
}
