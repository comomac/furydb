package main

// basic server for demo

import (
	"bytes"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"text/template"

	_ "github.com/comomac/furydb"
)

var db *sql.DB
var tplIndexHTML *template.Template
var indexHTML = `
<html>
<body>
<div><h1>Users:</h1></div>
<div>
	<form method="get" action="/api/add_user">
		email: <input name="email">
		password: <input name="password">
		<input type="submit" value="add">
	</form>
<div>
	<table>
	<thead>
		<td>id</td>
		<td>email</td>
		<td>password</td>
	</thead>
{{range .Users}}
	<tr>
		<td>{{.ID}}</td>
		<td>{{.Email}}</td>
		<td>{{.Password}}</td>
	</tr>
{{end}}
	</table>
</div>
</body>
<html>
`

func init() {
	tplIndexHTML = template.Must(template.New("indexHtml").Parse(indexHTML))
}

func main() {
	var err error
	db, err = sql.Open("fury", "tmp-db")
	if err != nil {
		log.Fatal(err)
	}

	http.HandleFunc("/", handlerIndex)
	http.HandleFunc("/api/add_user", handlerAddUser)

	fmt.Printf("Starting server at port 8080\n")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}

type results struct {
	Users []*User
}
type User struct {
	ID       string
	Email    string
	Password string
}

func handlerIndex(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "Method is not supported.", http.StatusNotFound)
	}

	query := `
	SELECT (id,email,password)
	FROM users;
	`

	// run select query
	rows, err := db.Query(query)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	res := results{}

	for rows.Next() {
		var (
			id       string
			email    string
			password string
		)
		if err := rows.Scan(&id, &email, &password); err != nil {
			res.Users = append(res.Users, &User{
				ID:       id,
				Email:    email,
				Password: password,
			})
			break
		}
	}

	buf := bytes.Buffer{}
	err = tplIndexHTML.Execute(&buf, res)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	fmt.Fprintf(w, buf.String())
}

func handlerAddUser(w http.ResponseWriter, r *http.Request) {
	var err error
	if r.Method != "GET" {
		http.Error(w, "Method is not supported.", http.StatusNotFound)
	}

	email := r.URL.Query().Get("email")
	password := r.URL.Query().Get("password")

	_, err = db.Query(`INSERT INTO users (email,password)
	VALUES ('` + email + `','` + password + `');`)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	fmt.Fprintf(w, "added")
}
