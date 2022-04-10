# SQLiteViz

Go application that creates simple diagrams from your SQLite database schemas.

Install with:

    go install github.com/n0x1m/sqliteviz@latest

Cli usage:

```sh
# to dotfile
sqliteviz db.sqlite3 > output.dot
# to png
sqliteviz db.sqlite3 | dot -Tpng > output.png
# to svg
sqliteviz db.sqlite3 | dot -Tsvg > output.svg
```
