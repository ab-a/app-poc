package main

import (
    "database/sql"
    "log"
    "net/http"
    "text/template"

    _ "github.com/go-sql-driver/mysql"
)

type App_table struct {
    ID     int
    Name   string 
    Color  string
    Pet    string
}

func dbConn() (db *sql.DB) {
    dbDriver := "mysql"
    dbUser := "root"
    dbPass := ""
    dbName := "app_db"
    db, err := sql.Open(dbDriver, dbUser+":"+dbPass+"@/"+dbName)
    if err != nil {
        panic(err.Error())
    }
    return db
}

var tmpl = template.Must(template.ParseGlob("form/*"))

func Index(w http.ResponseWriter, r *http.Request) {
    db := dbConn()
    selDB, err := db.Query("SELECT * FROM app_table ORDER BY id DESC")
    if err != nil {
        panic(err.Error())
    }
    emp := App_table{}
    res := []App_table{}
    for selDB.Next() {
        var id int
        var name, color, pet string
        err = selDB.Scan(&id, &name, &color, &pet)
        if err != nil {
            panic(err.Error())
        }
        emp.ID = id
        emp.Name = name
        emp.Color = color        
	emp.Pet = pet
	res = append(res, emp)
    }
    tmpl.ExecuteTemplate(w, "Index", res)
    defer db.Close()
}

func Show(w http.ResponseWriter, r *http.Request) {
    db := dbConn()
    nId := r.URL.Query().Get("id")
    selDB, err := db.Query("SELECT * FROM app_table WHERE id=?", nId)
    if err != nil {
        panic(err.Error())
    }
    emp := App_table{}
    for selDB.Next() {
        var id int
        var name, color, pet string
        err = selDB.Scan(&id, &name, &color, &pet)
        if err != nil {
            panic(err.Error())
        }
        emp.ID = id
        emp.Name = name
        emp.Color = color
        emp.Pet = pet
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
    selDB, err := db.Query("SELECT * FROM app_table WHERE id=?", nId)
    if err != nil {
        panic(err.Error())
    }
    emp := App_table{}
    for selDB.Next() {
        var id int
        var name, color, pet string
        err = selDB.Scan(&id, &name, &color, &pet)
        if err != nil {
            panic(err.Error())
        }
        emp.ID = id
        emp.Name = name
        emp.Color = color
        emp.Pet = pet
    }
    tmpl.ExecuteTemplate(w, "Edit", emp)
    defer db.Close()
}

func Insert(w http.ResponseWriter, r *http.Request) {
    db := dbConn()
    if r.Method == "POST" {
        name := r.FormValue("name")
        color := r.FormValue("color")
        pet := r.FormValue("pet")
        insForm, err := db.Prepare("INSERT INTO app_table(name, color, pet) VALUES(?,?,?)")
        if err != nil {
            panic(err.Error())
        }
        insForm.Exec(name, color, pet)
        log.Println("INSERT: Name: " + name + " | Color: " + color + " | Pet: " + pet)
    }
    defer db.Close()
    http.Redirect(w, r, "/", 301)
}

func Update(w http.ResponseWriter, r *http.Request) {
    db := dbConn()
    if r.Method == "POST" {
        name := r.FormValue("name")
        color := r.FormValue("color")
        pet := r.FormValue("pet")
        id := r.FormValue("uid")
        insForm, err := db.Prepare("UPDATE app_table SET name=?, color=?, pet=? WHERE id=?")
        if err != nil {
            panic(err.Error())
        }
        insForm.Exec(name, color, pet, id)
        log.Println("UPDATE: Name: " + name + " | Color: " + color + " | Pet: " + pet)
    }
    defer db.Close()
    http.Redirect(w, r, "/", 301)
}

func Delete(w http.ResponseWriter, r *http.Request) {
    db := dbConn()
    emp := r.URL.Query().Get("id")
    delForm, err := db.Prepare("DELETE FROM app_table WHERE id=?")
    if err != nil {
        panic(err.Error())
    }
    delForm.Exec(emp)
    log.Println("DELETE")
    defer db.Close()
    http.Redirect(w, r, "/", 301)
}

func main() {
    log.Println("Server started on: http://0.0.0.0:8080")
    http.HandleFunc("/", Index)
    http.HandleFunc("/show", Show)
    http.HandleFunc("/new", New)
    http.HandleFunc("/edit", Edit)
    http.HandleFunc("/insert", Insert)
    http.HandleFunc("/update", Update)
    http.HandleFunc("/delete", Delete)
    http.ListenAndServe(":8080", nil)
}
