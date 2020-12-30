package main

import (
	"fmt"
	"log"
	"strings"
)

func main() {
	// qryInsert := `INSERT INTO product (code, name, price)
	// VALUES ('ab821', 'honey',18.50);`

	qrySelect := `SELECT code_, name, price
	FROM products;`

	parser := NewParser(strings.NewReader(qrySelect))
	stmt, err := parser.Parse()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%+v\n", stmt)
}
