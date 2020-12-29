package main

import (
	"database/sql"
	"fmt"
	"log"
	"testing"

	"github.com/comomac/furydb"
)

// TestCreate db, written this before the parser is made
func TestCreate(t *testing.T) {
	db, err := furydb.Create("tmp-db", "testme")
	if err != nil {
		t.Error(err)
	}

	db.Tables = []*furydb.Table{
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
					DefaultDataUUID: "gen_uuid_v4()",
				},
				{
					Name:       "cstr-unique-email",
					ColumnName: "email",
					IsUnique:   true,
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

	err = db.Save()
	if err != nil {
		t.Error(err)
	}
}

// // TestLoad db
// func TestLoad(t *testing.T) {
// 	db, err := furydb.Load("tmp-db")
// 	if err != nil {
// 		t.Error(err)
// 	}

// 	if db.Name != "testme" {
// 		t.Error(fmt.Errorf("name mismatch"))
// 	}
// 	if db.Tables == nil {
// 		t.Error(fmt.Errorf("tables is nil"))
// 	}
// 	if len(db.Tables) == 0 {
// 		t.Error(fmt.Errorf("tables is zero len"))
// 	}
// }

var db *sql.DB

// TestSqlDriverOpen
func TestSqlDriverOpen(t *testing.T) {
	var err error
	_, err = sql.Open("fury", "tmp-db")
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
}

// TestSqlDriverInsert
func TestSqlDriverInsert(t *testing.T) {
	var err error
	if db == nil {
		t.Error(fmt.Errorf("db not loaded"))
		return
	}

	query := `
	INSERT INTO users (email,password)
	VALUES ('bob@example.com','testpass');
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
	SELECT *
	FROM users;
	`

	// run select query
	rows, err := db.Query(query)
	if err != nil {
		t.Error(err)
		return
	}
	defer rows.Close()

	// get results
	for rows.Next() {
		var (
			id       [16]byte
			email    string
			password string
		)
		if err := rows.Scan(&id, &email, &password); err != nil {
			t.Error(err)
			return
		}
		log.Printf("email: %s    password: %s\n", email, password)
	}
	if !rows.NextResultSet() {
		t.Error(fmt.Errorf("expected more result sets: %v", rows.Err()))
		return
	}
}
