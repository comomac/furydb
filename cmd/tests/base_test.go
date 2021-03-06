package main

import (
	"database/sql"
	"fmt"
	"testing"

	"github.com/comomac/furydb"
)

// holder so tests can use same connection without reconnect to db
var (
	db  *sql.DB
	fdb *furydb.Database
)

// TestCreate db, written this before the parser is made
func TestCreate(t *testing.T) {
	var err error
	fdb, err = furydb.Create("tmp-db", "testme")
	if err != nil {
		t.Error(err)
	}

	fdb.Tables = []*furydb.Table{
		{
			Name: "users",
			Columns: []*furydb.Column{
				{
					Name: "id",
					Type: furydb.ColumnTypeUUID,
				},
				{
					Name: "email",
					Type: furydb.ColumnTypeString,
				},
				{
					Name: "password",
					Type: furydb.ColumnTypeString,
				},
				{
					Name: "email_verified",
					Type: furydb.ColumnTypeString,
				},
				{
					Name: "created_at",
					Type: furydb.ColumnTypeTime,
				},
			},
			Constraints: []*furydb.Constraint{
				{
					Name:            "cstr-pk",
					ColumnName:      "id",
					IsPrimaryKey:    true,
					IsUnique:        true,
					IsNotNull:       true,
					DefaultDataUUID: "gen_uuid_v4()",
				},
				{
					Name:       "cstr-unique-email",
					ColumnName: "email",
					IsUnique:   true,
					IsNotNull:  true, // this will blow up if not nullable
				},
				{
					Name:            "cstr-created_at",
					ColumnName:      "created_at",
					DefaultDataUUID: "now()",
				},
			},
		},
		{
			Name: "customers",
			Columns: []*furydb.Column{
				{
					Name: "id",
					Type: furydb.ColumnTypeUUID,
				},
				{
					Name: "first_name",
					Type: furydb.ColumnTypeString,
				},
				{
					Name: "last_name",
					Type: furydb.ColumnTypeString,
				},
				{
					Name: "avatar",
					Type: furydb.ColumnTypeBytes,
				},
				{
					Name: "year_born",
					Type: furydb.ColumnTypeInt,
				},
				{
					Name: "credit",
					Type: furydb.ColumnTypeFloat,
				},
				{
					Name: "user_id",
					Type: furydb.ColumnTypeUUID,
				},
			},
			Constraints: []*furydb.Constraint{
				{
					Name:            "cstr-pk",
					ColumnName:      "id",
					IsPrimaryKey:    true,
					IsUnique:        true,
					DefaultDataUUID: "gen_uuid_v4()",
				},
				{
					Name:       "cstr-first_name",
					ColumnName: "first_name",
					IsNotNull:  true,
				},
				{
					Name:       "cstr-last_name",
					ColumnName: "last_name",
					IsNotNull:  true,
				},
				{
					Name:          "user_id",
					Type:          furydb.ColumnTypeUUID,
					IsNotNull:     true,
					IsForeignKey:  true,
					ForeignTable:  "users",
					ForeignColumn: "id",
				},
			},
		},
	}

	err = fdb.Save()
	if err != nil {
		t.Error(err)
	}
}

// TestLoad db
func TestLoad(t *testing.T) {
	var err error
	fdb, err = furydb.Load("tmp-db")
	if err != nil {
		t.Error(err)
	}

	if fdb.Name != "testme" {
		t.Error(fmt.Errorf("name mismatch"))
	}
	if fdb.Tables == nil {
		t.Error(fmt.Errorf("tables is nil"))
	}
	if len(fdb.Tables) == 0 {
		t.Error(fmt.Errorf("tables is zero len"))
	}
}

// TestSqlDriverOpen
func TestSqlDriverOpen(t *testing.T) {
	var err error
	db, err = sql.Open("fury", "tmp-db")
	if err != nil {
		t.Error(err)
		return
	}
}

// TestSqlDriverTableCreate
func TestSqlDriverTableCreate(t *testing.T) {
	if db == nil {
		t.Error(fmt.Errorf("db not loaded"))
		return
	}
	// todo finish me
}

// TestSqlDriverInsert
func TestSqlDriverInsert(t *testing.T) {
	var err error
	if db == nil {
		t.Error(fmt.Errorf("db not loaded"))
		return
	}

	query := `
	INSERT INTO users (id,email,password)
	VALUES ('0583e443-015a-0b3e-0f19-6bce4dfdd638','bob@example.com','testpass');
	`

	// run insert query
	_, err = db.Query(query)
	if err != nil {
		t.Error(err)
		return
	}
}

// TestSqlDriverSelect
func TestSqlDriverSelect(t *testing.T) {
	var err error
	if db == nil {
		t.Error(fmt.Errorf("db not loaded"))
		return
	}

	query := `
	SELECT (id,email,password)
	FROM users;
	`

	// run select query
	rows, err := db.Query(query)
	if err != nil {
		t.Error(err)
		return
	}
	defer rows.Close()

	var (
		// id       [16]byte
		id       string
		email    string
		password string
		// email    sql.NullString
		// password sql.NullString
	)
	// get results
	for rows.Next() {
		if err := rows.Scan(&id, &email, &password); err != nil {
			t.Error(err)
			return
		}
		// log.Printf("id: %+v   email: %+v   password: %+v\n", id, email, password)
	}
	if id != "0583e443-015a-0b3e-0f19-6bce4dfdd638" {
		t.Error(fmt.Errorf("invalid id"))
	}
	if email != "bob@example.com" {
		t.Error(fmt.Errorf("invalid email"))
	}
	if password != "testpass" {
		t.Error(fmt.Errorf("invalid password"))
	}
	// if !rows.NextResultSet() {
	// 	t.Error(fmt.Errorf("expected more result sets: %v", rows.Err()))
	// 	return
	// }
}
