package main

import (
    "database/sql"
    "log"
    "net/http"
    "text/template"

    _ "github.com/go-sql-driver/mysql"
)

type Actor struct {
    Id    int
    FName  string
    LName string
    LUpdate string
}

func dbConn() (db *sql.DB) {
    dbDriver := "mysql"
    dbUser := "root"
    dbPass := "Pw_4087119317"
    dbName := "sakila"
    db, err := sql.Open(dbDriver, dbUser+":"+dbPass+"@/"+dbName)
    if err != nil {
        panic(err.Error())
    }
    return db
}

var tmpl = template.Must(template.ParseGlob("form/*"))

func Index(w http.ResponseWriter, r *http.Request) {
    db := dbConn()
    selDB, err := db.Query("SELECT * FROM actor ORDER BY actor_id DESC")
    if err != nil {
        panic(err.Error())
    }
    emp := Actor{}
    res := []Actor{}
    for selDB.Next() {
        var Id int
        var FName, LName, LUpdate string
        err = selDB.Scan(&Id, &FName, &LName, &LUpdate)
        if err != nil {
            panic(err.Error())
        }
        emp.Id = Id
        emp.FName = FName
        emp.LName = LName
        emp.LUpdate = LUpdate
        res = append(res, emp)
    }
    tmpl.ExecuteTemplate(w, "Index", res)
    defer db.Close()
}

func Show(w http.ResponseWriter, r *http.Request) {
    db := dbConn()
    nId := r.URL.Query().Get("id")
    selDB, err := db.Query("SELECT * FROM actor WHERE actor_id=?", nId)
    if err != nil {
        panic(err.Error())
    }
    emp := Actor{}
    for selDB.Next() {
        var Id int
        var FName, LName, LUpdate string
        err = selDB.Scan(&Id, &FName, &LName, &LUpdate)
        if err != nil {
            panic(err.Error())
        }
        emp.Id = Id
        emp.FName = FName
        emp.LName = LName
        emp.LUpdate = LUpdate
    }
    tmpl.ExecuteTemplate(w, "Show", emp)
    defer db.Close()
}

func New(w http.ResponseWriter, r *http.Request) {
    tmpl.ExecuteTemplate(w, "New", nil)
}

func Edit(w http.ResponseWriter, r *http.Request) {
    db := dbConn()
    nId := r.URL.Query().Get("id")
    selDB, err := db.Query("SELECT * FROM Actor WHERE actor_id=?", nId)
    if err != nil {
        panic(err.Error())
    }
    emp := Actor{}
    for selDB.Next() {
        var Id int
        var FName, LName, LUpdate string
        err = selDB.Scan(&Id, &FName, &LName, &LUpdate)
        if err != nil {
            panic(err.Error())
        }
        emp.Id = Id
        emp.FName = FName
        emp.LName = LName
        emp.LUpdate = LUpdate
    }
    tmpl.ExecuteTemplate(w, "Edit", emp)
    defer db.Close()
}

func Insert(w http.ResponseWriter, r *http.Request) {
    db := dbConn()
    if r.Method == "POST" {
        Id := r.FormValue("Id")
        FName := r.FormValue("FName")
        LName := r.FormValue("LName")
        LUpdate := r.FormValue("LUpdate")
        insForm, err := db.Prepare("INSERT INTO Actor(actor_id, first_name, last_name, last_update) VALUES(?,?,?,?)")
        if err != nil {
            panic(err.Error())
        }
        insForm.Exec(Id, FName, LName, LUpdate)
        log.Println("INSERT: ID : " + Id + " | First Name: " + FName + " | Last Name: " + LName + " | Last Update: " + LUpdate)
    }
    defer db.Close()
    http.Redirect(w, r, "/", 301)
}

func Update(w http.ResponseWriter, r *http.Request) {
    db := dbConn()
    if r.Method == "POST" {
        FName := r.FormValue("FName")
        LName := r.FormValue("LName")
        LUpdate := r.FormValue("LUpdate")
        Id := r.FormValue("Id")
        insForm, err := db.Prepare("UPDATE actor SET first_name=?, last_name=?, last_update=? WHERE actor_id=?")
        if err != nil {
            panic(err.Error())
        }
        insForm.Exec(FName, LName, LUpdate, Id)
        log.Println("UPDATE: First Name: " + FName + " | Last Name: " + LName + " | Last Update: " + LUpdate + " WHERE ID : " + Id )
    }
    defer db.Close()
    http.Redirect(w, r, "/", 301)
}

func Delete(w http.ResponseWriter, r *http.Request) {
    db := dbConn()
    emp := r.URL.Query().Get("id")
    delForm, err := db.Prepare("DELETE FROM actor WHERE actor_id=?")
    if err != nil {
        panic(err.Error())
    }
    delForm.Exec(emp)
    log.Println("DELETE")
    defer db.Close()
    http.Redirect(w, r, "/", 301)
}

func main() {
    log.Println("Server started on: http://localhost:8080")
    http.HandleFunc("/", Index)
    http.HandleFunc("/show", Show)
    http.HandleFunc("/new", New)
    http.HandleFunc("/edit", Edit)
    http.HandleFunc("/insert", Insert)
    http.HandleFunc("/update", Update)
    http.HandleFunc("/delete", Delete)
    http.ListenAndServe(":8080", nil)
}