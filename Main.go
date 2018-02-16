package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"

	_ "github.com/go-sql-driver/mysql"
)

func main() {

	router := mux.NewRouter().StrictSlash(true)

	router.HandleFunc("/", Index).Methods("GET")
	router.HandleFunc("/store", SavePerson).Methods("POST")
	router.HandleFunc("/update/{id}", UpdatePerson).Methods("PUT")
	router.HandleFunc("/destroy/{id}", DeletePerson).Methods("DELETE")
	router.HandleFunc("/show/{id}", ShowPerson).Methods("GET")

	log.Fatal(http.ListenAndServe(":8080", router))
}

// Person Data model
type Person struct {
	Id          int    `json:id`
	Firstname   string `json:firstname`
	Lastname    string `json:lastname`
	Email       string `json:email`
	Phonenumber string `json:phonenumber`
	Birth       string `json:birth`
}

// Index Index page
func Index(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	fmt.Fprintln(w, "SPARK!\n")

	db := dbconn()

	rows, err := db.Query("SELECT * FROM person")
	checkErr(err)

	var p Person
	var ppl []Person

	for rows.Next() {
		err = rows.Scan(&p.Id, &p.Firstname, &p.Lastname, &p.Email, &p.Phonenumber, &p.Birth)
		checkErr(err)
		ppl = append(ppl, p)
		// json.NewEncoder(w).Encode(p)
	}

	json.NewEncoder(w).Encode(ppl)

	db.Close()
}

// PersonStore STORE
func SavePerson(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var p Person
	err := decoder.Decode(&p)
	checkErr(err)
	defer r.Body.Close()

	db := dbconn()

	stmt, err := db.Prepare("INSERT INTO person(firstname, lastname, email, phonenumber, birth) VALUES(?,?,?,?,?)")
	checkErr(err)

	res, err := stmt.Exec(p.Firstname, p.Lastname, p.Email, p.Phonenumber, p.Birth)
	checkErr(err)

	id, err := res.LastInsertId()
	checkErr(err)

	fmt.Fprintln(w, id)

	db.Close()
}

// PersonUpdate UPDATE
func UpdatePerson(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var p Person
	err := decoder.Decode(&p)
	checkErr(err)
	defer r.Body.Close()
	vars := mux.Vars(r)

	db := dbconn()

	stmt, err := db.Prepare("UPDATE person SET firstname=?, lastname=?, email=?, phonenumber=?, birth=? WHERE id=?")
	checkErr(err)

	res, err := stmt.Exec(p.Firstname, p.Lastname, p.Email, p.Phonenumber, p.Birth, vars["id"])
	checkErr(err)

	affect, err := res.RowsAffected()
	checkErr(err)

	fmt.Println(affect)
}

// PersonShow SHOW
func ShowPerson(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	vars := mux.Vars(r)

	var person Person

	db := dbconn()

	rows, err := db.Query("SELECT * FROM person WHERE id=" + vars["id"])
	checkErr(err)

	for rows.Next() {
		err = rows.Scan(&person.Id, &person.Firstname, &person.Lastname, &person.Email, &person.Phonenumber, &person.Birth)
		checkErr(err)
	}

	j, _ := json.Marshal(person)
	w.Write(j)
}

// PersonDelete DESTROY
func DeletePerson(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	db := dbconn()

	stmt, err := db.Prepare("DELETE FROM person WHERE id=?")
	checkErr(err)

	res, err := stmt.Exec(vars["id"])
	checkErr(err)

	affect, err := res.RowsAffected()
	checkErr(err)

	fmt.Println(affect)

	db.Close()

}

func dbconn() *sql.DB {
	db, err := sql.Open("mysql", "root:@tcp(127.0.0.1:3306)/spark")
	checkErr(err)
	return db
}

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}
