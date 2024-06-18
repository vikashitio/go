package routes

import (
	"ebank/handlers"
	"ebank/models"
	"html/template"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
)

var tpl *template.Template
var store = sessions.NewCookieStore([]byte("EindiaBusiness"))

func init() {
	//templates = template.Must(template.ParseGlob("templates/*.html"))
	tpl, _ = template.ParseGlob("static/*.html")
}

func InitRoutes() *mux.Router {

	router := mux.NewRouter()
	//router.HandleFunc("/users", handlers.UsersHandler).Methods("GET")
	router.HandleFunc("/login", loginHandler)
	router.HandleFunc("/loginPost", handlers.UsersLogin).Methods("POST")
	router.HandleFunc("/profile", handlers.ProfileHandler).Methods("GET")
	router.HandleFunc("/profilePost", handlers.ProfilePost).Methods("Post")
	//router.HandleFunc("/profilePost", handlers.UsersLogin).Methods("POST")
	router.HandleFunc("/loginhistory", models.LoginhistoryHandler)
	router.HandleFunc("/registration", registrationHandler)
	router.HandleFunc("/registrationPost", handlers.UsersRegistration).Methods("POST")
	router.HandleFunc("/logout", logoutHandler)
	router.HandleFunc("/", handlers.IndexHandler)

	// Add more routes here

	// Serve static files
	//router.PathPrefix("/").Handler(http.FileServer(http.Dir("./static/")))

	return router
}

type logindata struct {
	Title    string
	Email    string
	Contact  string
	Message  string
	Name     string
	ClientID int
	Status   int
}

var logindata1 logindata

// loginHandler serves form for users to login with
func loginHandler(w http.ResponseWriter, r *http.Request) {

	session, _ := store.Get(r, "merchant")
	//fmt.Println(session)

	if session.Values["Login-Message"] == "" {
		var msg = ""
		tpl.ExecuteTemplate(w, "login.html", msg)
	} else {
		var msg = session.Values["Login-Message"]
		session.Values["Login-Message"] = ""
		session.Save(r, w)
		tpl.ExecuteTemplate(w, "login.html", msg)
	}

}

// logoutHandler serves form for users to login with
func logoutHandler(w http.ResponseWriter, r *http.Request) {

	session, _ := store.Get(r, "merchant") // Get all Session
	//fmt.Println("session:", session)
	//var client_id = session.Values["Merchant-ID"]
	//fmt.Println(client_id)
	//UpdateLogout(client_id)
	session.Options.MaxAge = -1 //destroy all session
	// Save it before we write to the response/return from the handler.
	session.Save(r, w)

	var msg = "Logout Successfully"
	tpl.ExecuteTemplate(w, "login.html", msg)

}

// Registration serves
func registrationHandler(w http.ResponseWriter, r *http.Request) {

	logindata1 = logindata{
		Title:   "Login Form",
		Email:   "vikashg@itio.in",
		Contact: "+ 977 9852 5862 55",
		Message: "",
	}

	tpl.ExecuteTemplate(w, "registration.html", logindata1)
}
