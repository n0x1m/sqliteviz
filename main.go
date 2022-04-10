package main

import (
	"database/sql"
	"fmt"
	"os"
	"time"

	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
)

func main() {
	if len(os.Args) == 1 {
		fmt.Printf("No database specified\n")
		os.Exit(1)
	}

	path := os.Args[1]

	db, err := sql.Open("sqlite3", path)
	if err != nil {
		fmt.Printf("Error opening database '%s': %s\n", path, err)
		os.Exit(1)
	}
	defer db.Close()

	sqlxDB := sqlx.NewDb(db, "sqlite3")

	tables, err := Tables(sqlxDB)
	if err != nil {
		fmt.Printf("Error listing tables: %s\n", err)
		os.Exit(1)
	}

	var entities []*Entity
	var rs []*ForeignKey

	for _, table := range tables {
		// TODO: cmd line arg
		if table.Name == "schema_migrations" {
			continue
		}

		// fmt.Println(table.Name)
		columns, err := Info(sqlxDB, table.Name)
		if err != nil {
			fmt.Printf("TableInfo failed: %v", err)
			os.Exit(1)
		}

		attrs := make([]*Column, len(columns))
		for i, col := range columns {
			// fmt.Println(i, col.Name)
			attrs[i] = &Column{
				Name:    col.Name,
				Primary: col.Pk == 1,
				Type:    col.DataType,
				NotNull: col.NotNull,
			}
		}

		entities = append(entities, &Entity{
			Name:       table.Name,
			Attributes: attrs,
		})
	}

	for _, table := range tables {
		fks, err := ForeignKeys(sqlxDB, table.Name)
		if err != nil {
			fmt.Printf("ForeignKeys failed: %v", err)
			os.Exit(1)
		}

		for _, fk := range fks {
			setKey(entities, table.Name, fk.From)
			setKey(entities, fk.Table, fk.To)

			// TODO: composite keys when seq increments with the same id
			rs = append(rs, &ForeignKey{
				TargetTable:  table.Name,
				TargetColumn: fk.From,
				SourceTable:  fk.Table,
				SourceColumn: fk.To,
			})
		}
	}

	RenderFromTemplate(&Diagram{
		Name:      os.Args[1],
		Date:      time.Now().Format(time.RFC3339),
		Entities:  entities,
		Relations: rs,
	}, "diagram.tpl.dot")
}
