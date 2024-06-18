package handlers

import (
	"ebank/function"
	"ebank/models"
	"encoding/json"
	"fmt"
	"net/http"
	"text/template"

	"github.com/gorilla/sessions"
)

var store = sessions.NewCookieStore([]byte("EindiaBusiness"))

var tpl *template.Template

func init() {
	//templates = template.Must(template.ParseGlob("templates/*.html"))
	tpl, _ = template.ParseGlob("static/*.html")
}

type Sub struct {
	Username string
	Data     string
}

type PageData struct {
	Title       string
	Email       string
	Contact     string
	Message     string
	Name        string
	ClientID    int
	Status      int
	StatusText  string
	ListData    []models.Transaction
	DisplayData models.Profile
}

func UsersHandler(w http.ResponseWriter, r *http.Request) {
	users, err := models.GetUsers()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(users)
	//fmt.Println(users)

}

func UsersLogin(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()                      // Parses the request body
	username := r.Form.Get("username") // Fetch Value of username
	password := r.Form.Get("password") // Fetch Value of password

	message, _ := models.GetLoggedDetails(username, password)

	//fmt.Println(message)
	session, err := store.Get(r, "merchant")
	if message.Alert == "" {
		//fmt.Println("Client ID:", message.ID)
		//fmt.Println("Name:", message.Name)
		//fmt.Println("Email:", message.Email)
		//fmt.Println("Status:", message.Status)
		// Store Session Variable

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Set some session values.
		session.Values["Merchant-Name"] = message.Name
		session.Values["Merchant-ID"] = message.ID
		session.Values["Merchant-Email"] = message.Email
		session.Values["Merchant-Status"] = message.Status
		//session.Values["Login-Message"] = "Login Done"

		//fmt.Println("Client ID 11:", session.Values["Merchant-ID"])
		//fmt.Println("Name 22:", session.Values["Merchant-Name"])
		//fmt.Println("Email 33:", session.Values["Merchant-Email"])
		//fmt.Println("Status 44:", session.Values["Merchant-Status"])

		// Save it before we write to the response/return from the handler.
		err = session.Save(r, w)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		http.Redirect(w, r, "/", 302)
	} else {
		session.Values["Login-Message"] = message.Alert
		// Save it before we write to the response/return from the handler.
		err = session.Save(r, w)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		//fmt.Println("Status 44:", session.Values["Login-Message"])
		http.Redirect(w, r, "/login", 302)

	}

}

func UsersRegistration(w http.ResponseWriter, r *http.Request) {

	r.ParseForm()                // Parses the request body
	name := r.Form.Get("name")   // Fetch Value of name
	email := r.Form.Get("email") // Fetch Value of email

	message, err := models.UsersRegistration(name, email)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	//fmt.Println(message)
	//fmt.Println("Client ID:", message.ID)
	//fmt.Println("Name:", message.Name)
	//fmt.Println("Email:", message.Email)
	//fmt.Println("Status:", message.Status)

	//var client_id = message.ID

	if message.Alert == "Success" {

		//  Email///
		var domName = "http://localhost:8080"
		var subject = "Test Message"
		//var HTMLbody = "Hi this is message"
		HTMLbody :=
			`<html>
			<p><strong>Hello , ` + message.Name + `</strong></p>
			<br/>
			<p>Welcome to Golang Bank! We are pleased to inform that your account has been created.</p>
			<br/>
			<strong>Login Details for Your Account:<br/>=====================<br/><strong>
			<p>Username :  ` + message.Email + `</p>
			<span>Password : </span> 
			<a href="` + domName + `/emailver/` + message.Email + `">click to change password</a>
			<br/>
			Cheers,
			<br/>
            <strong>Golang Bank</strong>
		</html>`
		err = function.SendEmail(subject, HTMLbody)
		if err != nil {
			fmt.Println("issue sending verification email")
		}
		session, err := store.Get(r, "merchant")
		// Set some session values.
		session.Values["Merchant-Name"] = message.Name
		session.Values["Merchant-ID"] = message.ID
		session.Values["Merchant-Email"] = message.Email
		session.Values["Merchant-Status"] = message.Status
		// Save it before we write to the response/return from the handler.
		err = session.Save(r, w)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		///
		http.Redirect(w, r, "/", 302)
		//fmt.Println("111 :", message.Alert)

	} else {
		//fmt.Println(message)
		//fmt.Println("2222 :", message.Alert)

		tpl.ExecuteTemplate(w, "registration.html", struct {
			Title   string
			Email   string
			Contact string
			Message string
		}{
			Title:   "Login Form",
			Email:   "vikashg@itio.in",
			Contact: "+ 977 9852 5862 55",
			Message: message.Alert,
		})
	}

}

// Index data serves
func IndexHandler(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "merchant")

	if session.Values["Merchant-ID"] == nil {
		var msg = "Session expired - please Re login"
		tpl.ExecuteTemplate(w, "login.html", msg)
		return
	}
	//fmt.Println(session)
	SessionMerchantID := session.Values["Merchant-ID"].(int)          // Define session variable with data type
	SessionMerchantName := session.Values["Merchant-Name"].(string)   // Define session variable with data type
	SessionMerchantEmail := session.Values["Merchant-Email"].(string) // Define session variable with data type
	SessionMerchantStatus := session.Values["Merchant-Status"].(int)  // Define session variable with data type

	transList, err := models.TransactionsList()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	//fmt.Println(transList)

	MerchantStatus, err := function.GetStatus(SessionMerchantStatus)
	//fmt.Println(MerchantStatus.Status) //Fetch Status
	var MerchantTextStatus = MerchantStatus.Status
	if err != nil {
		fmt.Println("Status not found")
	}

	var data = PageData{
		Title:      "My Dashboard",
		Email:      SessionMerchantEmail,
		Name:       SessionMerchantName,
		ClientID:   SessionMerchantID,
		StatusText: MerchantTextStatus,
		Message:    "My Dashboard",
		ListData:   transList,
	}
	//fmt.Println(data)
	tpl.ExecuteTemplate(w, "index.html", data)
}

// Profile
func ProfileHandler(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "merchant")

	if session.Values["Merchant-ID"] == nil {
		var msg = "Session expired - please Re Login"
		tpl.ExecuteTemplate(w, "login.html", msg)
		return
	}

	MID := session.Values["Merchant-ID"].(int)                      // Define data type
	SessionMerchantName := session.Values["Merchant-Name"].(string) // Define data type
	//MID := session.Values["Merchant-Name"]
	//fmt.Println("===> ", session)

	message, err := models.ProfileDetails(MID)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var msg = ""

	if session.Values["Login-Message"] != "" {
		//fmt.Println("===> ===> ", MID)
		msg = "Profile Update Successfully"
		// 	msg = session.Values["Login-Message"].(string)
		session.Values["Login-Message"] = ""
		session.Save(r, w)
	}

	var data = models.Profile{
		Title:        message.Title,
		Gender:       message.Gender,
		BirthDate:    message.BirthDate,
		CountryCode:  message.CountryCode,
		Mobile:       message.Mobile,
		AddressLine1: message.AddressLine1,
		AddressLine2: message.AddressLine2,
		SesName:      SessionMerchantName, //Pass session Merchant Name
		Alert:        msg,
	}
	//fmt.Println("->", data)
	tpl.ExecuteTemplate(w, "profile.html", data)
}

func ProfilePost(w http.ResponseWriter, r *http.Request) {

	r.ParseForm()                                // Parses the request body
	title := r.Form.Get("title")                 // Fetch Value of title
	gender := r.Form.Get("gender")               // Fetch Value of gender
	birth_date := r.Form.Get("birth_date")       // Fetch Value of birth_date
	country_code := r.Form.Get("country_code")   // Fetch Value of country_code
	mobile := r.Form.Get("mobile")               // Fetch Value of mobile
	address_line1 := r.Form.Get("address_line1") // Fetch Value of address_line1
	address_line2 := r.Form.Get("address_line2") // Fetch Value of address_line2

	// fmt.Println("title:", title)
	// fmt.Println("gender:", gender)
	// fmt.Println("birth_date:", birth_date)
	// fmt.Println("country_code:", country_code)
	// fmt.Println("mobile:", mobile)
	// fmt.Println("address_line1:", address_line1)
	// fmt.Println("address_line2:", address_line2)

	message, err := models.UpdateProfile(title, gender, birth_date, country_code, mobile, address_line1, address_line2)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	//fmt.Println(message.Alert)
	if message.Alert == "Success" {
		session, _ := store.Get(r, "merchant")
		session.Values["Login-Message"] = "Profile Update Successfully"
		// Save it before we write to the response/return from the handler.
		err = session.Save(r, w)

		http.Redirect(w, r, "/profile", 302)

	} else {
		//fmt.Println(message)
		//fmt.Println("2222 :", message.Alert)

		tpl.ExecuteTemplate(w, "profile.html", struct {
			Title        string
			Gender       string
			BirthDate    string
			CountryCode  string
			Mobile       string
			AddressLine1 string
			AddressLine2 string
			Message      string
		}{
			Title:        title,
			Gender:       gender,
			BirthDate:    birth_date,
			CountryCode:  country_code,
			Mobile:       mobile,
			AddressLine1: address_line1,
			AddressLine2: address_line2,
			Message:      message.Alert,
		})
	}

}
