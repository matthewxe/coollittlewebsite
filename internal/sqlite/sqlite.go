package sqlite

import (
	"database/sql"
	"fmt"
	"log"
	"math/rand"
	"os"

	_ "modernc.org/sqlite"
)

var ApiKey = random_key(32)
var DB *sql.DB

func Hello() {
	fmt.Println("boolin")
}

func random_key(len int) (key []byte) {
	for i := 0; i < len; i++ {
		key = append(key, byte((rand.Int()%97)+33))
	}
	return key
}

func setup_sql(db *sql.DB) {
	table_create := `CREATE TABLE blogs (
		id INTEGER PRIMARY KEY,
		thumbnail BLOB,
		content TEXT NOT NULL,
		timestamp DATETIME DEFAULT CURRENT_TIMESTAMP)`
	db.Exec(table_create)
}

func Setup() {
	db, err := sql.Open("sqlite", "./data/data.db")
	if err != nil {
		return
	}
	defer db.Close()
	setup_sql(db)

	log.Printf("The key is\n%s", ApiKey)
	os.WriteFile("/tmp/coollittlewebsite_key", []byte(ApiKey), 0666)
}
