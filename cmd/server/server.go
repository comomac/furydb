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
	<form method="post" action="/add">
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

var submitted = `<!DOCTYPE html>
<html>
<head>
	<meta http-equiv="refresh" content="2; url='/'" />
</head>
<body>
	<p>Added, click <a href="/">here</a> to return</p>
</body>
</html>`

func init() {
	tplIndexHTML = template.Must(template.New("indexHtml").Parse(indexHTML))
}

func main() {
	var err error
	db, err = sql.Open("fury", "tmp-db")
	if err != nil {
		log.Fatal(err)
	}

	http.HandleFunc("/", handlerList)
	http.HandleFunc("/add", handlerAdd)

	fmt.Printf("Starting server at port 8080\n")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}

type results struct {
	Users []*user
}
type user struct {
	ID       string
	Email    string
	Password string
}

func handlerList(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "Method is not supported.", http.StatusMethodNotAllowed)
		return
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
		err := rows.Scan(&id, &email, &password)
		if err != nil {
			log.Println("Error rows.Scan() failed", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		res.Users = append(res.Users, &user{
			ID:       id,
			Email:    email,
			Password: password,
		})
	}

	buf := bytes.Buffer{}
	err = tplIndexHTML.Execute(&buf, res)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	fmt.Fprintf(w, buf.String())
}

func handlerAdd(w http.ResponseWriter, r *http.Request) {
	var err error
	if r.Method != "POST" {
		http.Error(w, "Method is not supported.", http.StatusMethodNotAllowed)
		return
	}

	err = r.ParseForm()
	if err != nil {
		http.Error(w, "Process form failed", http.StatusInternalServerError)
		return
	}

	email := r.PostForm.Get("email")
	password := r.PostForm.Get("password")

	_, err = db.Query(`INSERT INTO users (email,password)
	VALUES ('` + email + `','` + password + `');`)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	fmt.Fprintf(w, submitted)
}
