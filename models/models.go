package models

import (
	"database/sql"
	"fmt"
	"html/template"
	"log"
	"net"
	"net/http"

	"github.com/gorilla/sessions"
	_ "github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"
)

var db *sql.DB
var store = sessions.NewCookieStore([]byte("EindiaBusiness"))

func InitDB(dataSourceName string) {
	var err error
	db, err = sql.Open("postgres", dataSourceName)
	if err != nil {
		panic(err)
	}

	if err = db.Ping(); err != nil {
		panic(err)
	}
	fmt.Println("Database connection established")
}

func InitSession(dataSourceName string) {
	// var err error
	// db, err = sql.Open("postgres", dataSourceName)
	// if err != nil {
	// 	panic(err)
	// }

	// if err = db.Ping(); err != nil {
	// 	panic(err)
	// }
	fmt.Println("Session Start")
}

var tpl *template.Template

func init() {
	//templates = template.Must(template.ParseGlob("templates/*.html"))
	tpl, _ = template.ParseGlob("static/*.html")
}

type User struct {
	ClientID   int
	UserName   string
	FullName   string
	Tokenid    int
	Clientid   int
	Logintime  string
	Logouttime string
	Loginip    string
}

type Result struct {
	Name   string
	ID     int
	Email  string
	Status int
	Alert  string
}

type Login struct {
	Tokenid    int
	Clientid   int
	Logintime  string
	Logouttime string
	Loginip    string
}

type Transaction struct {
	TransID        int
	ClientID       int
	BillAmount     float64
	BillCurrency   string
	TransAmount    float64
	TransCurrency  string
	TransType      string
	TransStatus    string
	Timestamp      string
	SettlementDate string
}
type Profile struct {
	MID          int
	Title        string
	Gender       string
	BirthDate    string
	CountryCode  int
	Mobile       string
	AddressLine1 string
	AddressLine2 string
	Alert        string
	SesName      string
}

// Fetch user list from client master
func GetUsers() ([]User, error) {
	rows, err := db.Query("SELECT client_id, username, full_name FROM client_master ORDER BY client_id DESC")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	users := []User{}
	for rows.Next() {
		var user User
		if err := rows.Scan(&user.ClientID, &user.UserName, &user.FullName); err != nil {
			return nil, err
		}
		users = append(users, user)
	}
	//fmt.Println(users)
	return users, nil
}

// Fetch user list from client master
func GetLoginlist() ([]Login, error) {
	rows, err := db.Query("SELECT token_id,client_id,login_time,logout_time,login_ip FROM login_history ORDER BY token_id DESC")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	logins := []Login{}

	for rows.Next() {
		var login Login
		if err := rows.Scan(&login.Tokenid, &login.Clientid, &login.Logintime, &login.Logouttime, &login.Loginip); err != nil {
			return nil, err
		}
		logins = append(logins, login)
	}

	//fmt.Println(logins)
	return logins, nil
}

func GetLoggedDetails(login_username, login_password string) (Result, error) {
	//fmt.Println(login_username)
	//fmt.Println(login_password)

	sqlStatement := `SELECT client_id, full_name, password, status FROM client_master WHERE username = $1;`
	var client_id int
	var full_name string
	var password string
	var status int

	// Replace 3 with an ID from your database or another random
	// value to test the no rows use case.
	row := db.QueryRow(sqlStatement, login_username)
	switch err := row.Scan(&client_id, &full_name, &password, &status); err {
	case sql.ErrNoRows:
		//fmt.Println("Data Not Found")
		message := Result{
			Alert: "Data Not Found",
		}
		return message, nil
	case nil:
		//fmt.Println(client_id, full_name, password, status)

		if status != 1 {
			message := Result{
				Alert: "Account Not",
			}
			return message, nil
		}

		// func CompareHashAndPassword(hashedPassword, password []byte) error
		err = bcrypt.CompareHashAndPassword([]byte(password), []byte(login_password))

		// returns nill
		if err == nil {
			//fmt.Println("You have successfully logged in :")
			message := Result{
				Name:   full_name,
				ID:     client_id,
				Email:  login_username,
				Status: status,
			}
			// manage login history
			//fmt.Println(GetLocalIP())

			var ip = "192.168.29.4"
			//var ip = string(GetLocalIP())
			//fmt.Println(ip)
			sqlStatement := `INSERT INTO login_history (client_id, login_ip) VALUES ($1,  $2);`
			db.QueryRow(sqlStatement, client_id, ip)

			return message, nil

		} else {

			message := Result{
				Alert: "incorrect password",
			}

			return message, nil

		}

	default:
		panic(err)
	}

}

func UsersRegistration(name, email string) (Result, error) {
	//fmt.Println(name)
	//fmt.Println(email)

	// create hash from password
	var password = "India123"
	var msg = "Success"
	/// For Data Validation  ///

	/// End Data Validation  ///

	var hash []byte
	// func GenerateFromPassword(password []byte, cost int) ([]byte, error)
	hash, _ = bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	sqlStatement := `
INSERT INTO client_master (username, full_name, password, status)
VALUES ($1, $2, $3, $4)
RETURNING client_id`
	id := 0
	err := db.QueryRow(sqlStatement, email, name, hash, 1).Scan(&id)
	if err != nil {
		//panic(err)
		msg = "Fail - please check your data"
		//fmt.Println(msg)
	}
	///
	sqlStatementNew := `INSERT INTO client_details (client_id) VALUES ($1)`
	errr := db.QueryRow(sqlStatementNew, id)
	//fmt.Println(errr)
	if errr == nil {

		msg = "client details Insert Issue"
		//fmt.Println(msg)
	}
	///

	message := Result{Name: name, ID: id, Email: email, Status: 1, Alert: msg}

	return message, nil
}

func GetLocalIP() net.IP {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	localAddress := conn.LocalAddr().(*net.UDPAddr)

	return localAddress.IP
}

// fetch login history
func LoginhistoryHandler(w http.ResponseWriter, r *http.Request) {

	session, _ := store.Get(r, "merchant")

	if session.Values["Merchant-ID"] == nil {
		var msg = "Session expired - please relogin"
		tpl.ExecuteTemplate(w, "login.html", msg)
		return
	}

	rows, err := db.Query("SELECT token_id,client_id,login_time,logout_time,login_ip FROM login_history ORDER BY token_id DESC")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var users []User
	for rows.Next() {
		var user User
		err := rows.Scan(&user.Tokenid, &user.Clientid, &user.Logintime, &user.Logouttime, &user.Loginip)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		users = append(users, user)
	}
	if err = rows.Err(); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	//fmt.Println(users)
	data := struct {
		Title string
		Users []User
	}{
		Title: "Users",
		Users: users,
	}
	//fmt.Println(data)
	tpl.ExecuteTemplate(w, "login-history.html", data)

}

// fetch Transactions

func TransactionsList() ([]Transaction, error) {
	rows, err := db.Query("SELECT transid, clientid, billamount, billcurrency, transamount, transcurrency, transtype, transstatus, timestamp, settlementdate FROM master_trans_table ORDER BY transid DESC ")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	transList := []Transaction{}
	for rows.Next() {
		var trans Transaction
		if err := rows.Scan(&trans.TransID, &trans.ClientID, &trans.BillAmount, &trans.BillCurrency, &trans.TransAmount, &trans.TransCurrency, &trans.TransType, &trans.TransStatus, &trans.Timestamp, &trans.SettlementDate); err != nil {
			return nil, err
		}
		transList = append(transList, trans)
	}

	//fmt.Println(transList)
	return transList, nil
} // fetch Transactions

func ProfileDetails(MID int) (Profile, error) {
	//fmt.Println(MID)

	sqlStatement := `SELECT title, gender, birth_date, country_code, mobile, address_line1, address_line2  FROM client_details WHERE client_id = $1;`
	var title sql.NullString
	var gender sql.NullString
	var birth_date sql.NullString
	var country_code sql.NullInt16
	var mobile sql.NullString
	var address_line1 sql.NullString
	var address_line2 sql.NullString

	// Replace 3 with an ID from your database or another random
	// value to test the no rows use case.
	row := db.QueryRow(sqlStatement, MID)
	err := row.Scan(&title, &gender, &birth_date, &country_code, &mobile, &address_line1, &address_line2)

	fmt.Println(title, gender, birth_date, country_code, mobile, address_line1, address_line2)

	// returns nill
	if err != nil {
		message := Profile{
			Alert: "Data Not Found",
		}
		return message, nil
	}

	message := Profile{
		Title:        title.String,
		Gender:       gender.String,
		BirthDate:    birth_date.String,
		CountryCode:  int(country_code.Int16),
		Mobile:       mobile.String,
		AddressLine1: address_line1.String,
		AddressLine2: address_line2.String,
		Alert:        "Data Found",
	}
	return message, nil
}

func UpdateProfile(title, gender, birth_date, country_code, mobile, address_line1, address_line2 string) (Result, error) {
	//fmt.Println(title)
	//fmt.Println(gender)

	var msg = "Success"

	_, err := db.Exec("UPDATE client_details SET title=$1, gender=$2, birth_date=$3, country_code=$4, mobile=$5, address_line1=$6, address_line2=$7 WHERE client_id=$8", title, gender, birth_date, country_code, mobile, address_line1, address_line2, 108)
	//fmt.Println(err)
	if err != nil {
		//panic(err)
		msg = "Update Failed"
		fmt.Println(msg)
	}
	//fmt.Println(msg)
	message := Result{Alert: msg}

	return message, nil
}
