package main

import (
	"fmt"
	"testing"

	"github.com/comomac/furydb"
)

// TestCreate db
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

// TestLoad db
func TestLoad(t *testing.T) {
	db, err := furydb.Load("tmp-db")
	if err != nil {
		t.Error(err)
	}

	if db.Name != "testme" {
		t.Error(fmt.Errorf("name mismatch"))
	}
	if db.Tables == nil {
		t.Error(fmt.Errorf("tables is nil"))
	}
	if len(db.Tables) == 0 {
		t.Error(fmt.Errorf("tables is zero len"))
	}
}
