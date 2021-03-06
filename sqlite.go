package main

import (
	"github.com/guregu/null"
	"github.com/jmoiron/sqlx"
)

/*
sqlite> .headers ON
sqlite> select * from pragma_table_info('sqlite_master');
cid|name|type|notnull|dflt_value|pk
0|type|text|0||0
1|name|text|0||0
2|tbl_name|text|0||0
3|rootpage|int|0||0
4|sql|text|0||0
*/

type Table struct {
	Type     string      `db:"type"`
	Name     string      `db:"name"`
	TblName  string      `db:"tbl_name"`
	Rootpage int         `db:"rootpage"`
	SQL      null.String `db:"sql"`
}

func Tables(db *sqlx.DB) (results []Table, err error) {
	rows, err := db.Queryx("SELECT name,type,tbl_name,rootpage FROM sqlite_master WHERE type = 'table' AND name NOT IN ('sqlite_sequence')")
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		var t Table
		err = rows.StructScan(&t)
		if err != nil {
			return nil, err
		}

		results = append(results, t)
	}

	return results, nil
}

type TableInfo struct {
	ID        int         `db:"cid"`
	Name      string      `db:"name"`
	DataType  string      `db:"type"`
	NotNull   bool        `db:"notnull"`
	DfltValue null.String `db:"dflt_value"`
	Pk        int         `db:"pk"`
}

func Info(db *sqlx.DB, table string) (results []TableInfo, err error) {
	rows, err := db.Queryx("SELECT * FROM pragma_table_info(?)", table)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		var t TableInfo
		err = rows.StructScan(&t)
		if err != nil {
			return nil, err
		}

		results = append(results, t)
	}

	return results, nil
}

/*
sqlite> select * from pragma_foreign_key_list('fills');
id|seq|table|from|to|on_update|on_delete|match
*/

type TableForeignKey struct {
	ID       int    `db:"id"`
	Seq      int    `db:"seq"`
	Table    string `db:"table"`
	From     string `db:"from"`
	To       string `db:"to"`
	OnUpdate string `db:"on_update"`
	OnDelete string `db:"on_delete"`
	Match    string `db:"match"`
}

func ForeignKeys(db *sqlx.DB, table string) (results []TableForeignKey, err error) {
	rows, err := db.Queryx("SELECT * FROM pragma_foreign_key_list(?)", table)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		var t TableForeignKey
		err = rows.StructScan(&t)
		if err != nil {
			return nil, err
		}

		results = append(results, t)
	}

	return results, nil
}

type TableIndex struct {
	Table  string `db:"tbl"`
	Name   string `db:"name"`
	Unique int    `db:"uniq"`
}

func Indices(db *sqlx.DB) (results []TableIndex, err error) {
	rows, err := db.Queryx(`select m.name as tbl, il.name as name, il.[unique] as uniq from sqlite_master as m, pragma_index_list(m.name) as il, pragma_index_info(il.name) as ii where m.type='table' group by il.name, m.name order by m.name, il.name`)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		var t TableIndex
		err = rows.StructScan(&t)
		if err != nil {
			return nil, err
		}

		results = append(results, t)
	}

	return results, nil
}

type TableIndexInfo struct {
	Seq  int    `db:"seqno"`
	CID  int    `db:"cid"`
	Name string `db:"name"`
}

func IndexInfo(db *sqlx.DB, idx string) (results []TableIndexInfo, err error) {
	rows, err := db.Queryx(`select * from pragma_index_info(?)`, idx)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		var t TableIndexInfo
		err = rows.StructScan(&t)
		if err != nil {
			return nil, err
		}

		results = append(results, t)
	}

	return results, nil
}
