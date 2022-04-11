package main

import (
	"database/sql"
	"flag"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
)

type ignoreList map[string]struct{}

func (i *ignoreList) String() string {
	return fmt.Sprint(*i)
}

func (i *ignoreList) Set(value string) error {
	for _, v := range strings.Split(value, ",") {
		(*i)[v] = struct{}{}
	}

	return nil
}

func main() {
	var path, template string
	var ignore = make(ignoreList)

	flag.StringVar(&path, "db", "", "sqlite database path")
	flag.StringVar(&template, "template", "diagram.tpl.dot", "template file to use")
	flag.Var(&ignore, "ignore", "tables to ignore")
	flag.Parse()

	if path == "" && len(ignore) == 0 && len(os.Args) == 2 {
		path = os.Args[1]
	}

	if path == "" {
		fmt.Fprintf(os.Stderr, "Usage of %s:\n", os.Args[0])
		flag.PrintDefaults()
		os.Exit(1)
	}

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
	for _, table := range tables {
		if _, inList := ignore[table.Name]; inList {
			continue
		}

		// fmt.Println(table.Name)
		columns, err := Info(sqlxDB, table.Name)
		if err != nil {
			fmt.Printf("TableInfo failed: %v", err)
			os.Exit(1)
		}

		attrs := make([]*Attribute, len(columns))
		for i, col := range columns {
			// fmt.Println(i, col.Name)
			attrs[i] = &Attribute{
				Name:     col.Name,
				Primary:  col.Pk == 1,
				Type:     col.DataType,
				Nullable: !col.NotNull,
			}
		}

		entities = append(entities, &Entity{
			Name:       table.Name,
			Attributes: attrs,
		})
	}

	var rs []*ForeignKey
	for _, table := range tables {
		fks, err := ForeignKeys(sqlxDB, table.Name)
		if err != nil {
			fmt.Printf("ForeignKeys failed: %v", err)
			os.Exit(1)
		}

		for _, fk := range fks {
			setKey(entities, table.Name, fk.From)
			setKey(entities, fk.Table, fk.To)

			rs = append(rs, &ForeignKey{
				TargetTable:  table.Name,
				TargetColumn: fk.From,
				SourceTable:  fk.Table,
				SourceColumn: fk.To,
			})
		}
	}

	indices, err := Indices(sqlxDB)
	if err != nil {
		fmt.Printf("Error listing indices: %s\n", err)
		os.Exit(1)
	}

	for _, index := range indices {
		list, err := IndexInfo(sqlxDB, index.Name)
		if err != nil {
			os.Exit(1)
		}

		var compositeKey []string
		for _, key := range list {
			compositeKey = append(compositeKey, key.Name)
		}

		idx := fmt.Sprintf("%s (%s)", index.Name, strings.Join(compositeKey, ", "))
		addIndex(entities, index.Table, idx, index.Unique == 1)
	}

	err = RenderFromTemplate(&Diagram{
		Name:      path,
		Date:      time.Now().Format(time.RFC3339),
		Entities:  entities,
		Relations: rs,
	}, template)
	if err != nil {
		fmt.Fprint(os.Stderr, "failed to render digraph:", err)
	}
}
