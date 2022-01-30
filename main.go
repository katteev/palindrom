package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"strings"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
)

type User struct {
	ID        int    `json:"Id"`
	Fornavn   string `json:"fornavn"`
	Etternavn string `json:"etternavn"`
}

type ResultatPalindrom struct {
	ID                 int    `json: "Id"`
	Fornavn            string `json:"fornavn"`
	FornavnPalindrom   string `json:"fornavnpalindrom`
	Etternavn          string `json:"etternavn"`
	EtternavnPalindrom string `json:"etternavnpalindrom"`
}

var db *sql.DB
var err error
var users []User
var prevID = 0

// boolsk funksjon som sjekker om en tekst er palindrom
func isPalindrom(text string) bool {
	for i := 0; i < len(text); i++ {
		j := len(text) - 1 - i
		if text[i] != text[j] {
			return false
		}
	}
	return true
}

//hent bruker
func getUsers(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var user User
	var intern_users []User

	rows, e := db.Query(
		`SELECT ID, fornavn, etternavn FROM users;`)
	if e != nil {
		log.Println(e)
		http.Error(w, e.Error(), http.StatusInternalServerError)
		return
	}

	for rows.Next() {
		rows.Scan(&user.ID, &user.Fornavn, &user.Etternavn)
		intern_users = append(intern_users, user)
	}
	users = intern_users //oppdaterer global variabel
	json.NewEncoder(w).Encode(users)
}

//lag bruker
func createUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var user User
	var ID int
	json.NewDecoder(r.Body).Decode(&user)

	row := db.QueryRow(`SELECT COALESCE(MAX(ID), 0) FROM users;`)
	switch err := row.Scan(&ID); err {
	case sql.ErrNoRows:
		log.Println("No rows were returned!")
		prevID = 0
	case nil:
		prevID = ID
	default:
		panic(err)
	}

	prevID++
	user.ID = prevID

	_, err = db.Exec(
		"INSERT INTO users (ID, fornavn, etternavn ) VALUES (?, ?, ?)",
		user.ID,
		user.Fornavn,
		user.Etternavn,
	)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	//vis alle brukere
	getUsers(w, r)
}

//Vis en bruker
func getUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	var user User
	var inputID = params["ID"]

	var statement string = "SELECT fornavn, etternavn from users where ID = " + inputID + ";"

	row := db.QueryRow(statement)
	switch err := row.Scan(&user.Fornavn, &user.Etternavn); err {
	case sql.ErrNoRows:
		log.Println("No rows were returned!")
	case nil:
		user.ID, _ = strconv.Atoi(inputID)
		json.NewEncoder(w).Encode(&user)
	default:
		panic(err)
	}

}

//oppdaterer en bruker
func updateUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var user User
	params := mux.Vars(r)
	var inputID, _ = strconv.Atoi(params["ID"])
	json.NewDecoder(r.Body).Decode(&user)
	var inputFornavn = user.Fornavn
	var inputEtternavn = user.Etternavn

	// update user i DB
	// get users

	_, er := db.Exec(
		"UPDATE users SET fornavn = ?, etternavn = ? WHERE ID = ? ",
		inputFornavn, inputEtternavn, inputID)
	if er != nil {
		log.Println("Error updating...")
		http.Error(w, er.Error(), http.StatusInternalServerError)
		return
	}
	getUsers(w, r)

}

//sletter en bruker
func deleteUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	inputID, _ := strconv.Atoi(params["ID"])

	//ta i mot brukerID, slett fra DB; hent alle brukere i DB etterpå
	_, e := db.Exec("DELETE FROM users WHERE ID = ?", inputID)
	if e != nil {
		http.Error(w, e.Error(), http.StatusBadRequest)
		return
	}
	getUsers(w, r)
}

//sjekker om fornavn og etternavn til en bruker er palindrom
//getUsers må være kjørt først
func getPalindrom(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var sjekk ResultatPalindrom
	params := mux.Vars(r)
	inputID, _ := strconv.Atoi(params["ID"])

	for _, user := range users {
		if user.ID == inputID {
			json.NewDecoder(r.Body).Decode(&user)

			var resultatpalindrom string
			sjekk.ID = user.ID
			sjekk.Fornavn = user.Fornavn

			if isPalindrom(strings.ToLower(user.Fornavn)) {
				resultatpalindrom = "ja"
			} else {
				resultatpalindrom = "nei"
			}
			sjekk.FornavnPalindrom = resultatpalindrom
			sjekk.Etternavn = user.Etternavn

			if isPalindrom(strings.ToLower(user.Etternavn)) {
				resultatpalindrom = "ja"
			} else {
				resultatpalindrom = "nei"
			}
			sjekk.EtternavnPalindrom = resultatpalindrom
			json.NewEncoder(w).Encode(sjekk)
			return
		}
	}
	//hvis bruker ikke finnes i lista returneres kun innsendt ID
	if sjekk.ID == 0 {
		sjekk.ID = inputID
		json.NewEncoder(w).Encode(sjekk)
	}

}

func main() {
	router := mux.NewRouter()

	//åpner til DB
	//db, err = sql.Open("mysql", "root:admin@tcp(127.0.0.1:3306)/test")
	db, err = sql.Open("mysql", "root:admin@tcp(palindromDB)/test")
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()

	router.HandleFunc("/users", createUser).Methods("POST")
	router.HandleFunc("/users/{ID}", getUser).Methods("GET")
	router.HandleFunc("/users", getUsers).Methods("GET")
	router.HandleFunc("/users/{ID}", updateUser).Methods("PUT")
	router.HandleFunc("/users/{ID}", deleteUser).Methods("DELETE")
	router.HandleFunc("/users/palindrom/{ID}", getPalindrom).Methods("GET")

	log.Fatal(http.ListenAndServe(":8080", router))

}
