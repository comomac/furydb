package main

import (
	"fmt"
	"log"
	"strings"
)

// Debug flag
const Debug bool = true

func main() {
	// qrySelect := `SELECT code, name, price
	// FROM products;`

	qryInsert := `INSERT INTO product (code, name, price)
	VALUES ('a-821', 'honey',18.50);`

	if Debug {
		fmt.Println(qryInsert)
	}

	parser := NewParser(strings.NewReader(qryInsert))
	stmt, err := parser.Parse()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%+v\n", stmt)
}
