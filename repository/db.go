package repository

import (
	"database/sql"
	"log"

	_ "github.com/lib/pq"
)

const db_query = `CREATE TABLE IF NOT EXISTS requests (id SERIAL PRIMARY KEY,field1 TEXT NOT NULL,field2 TEXT NOT NULL,field3 TEXT NOT NULL,field4 TEXT NOT NULL,field5 TEXT NOT NULL,created_at TIMESTAMP NOT NULL DEFAULT NOW());`

func NewDatabase(connStr string) (*sql.DB, error) {
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Printf("error opening the database ---> %v", err.Error())
		return nil, err
	}
	log.Println("DATABASE OPENED")
	if err := db.Ping(); err != nil {
		log.Printf("error pinging the database ---> %v", err.Error())
		return nil, err
	}
	if _, err = db.Exec(db_query); err != nil {
		log.Printf("error making the table ---> %v", err.Error())
		return nil, err
	}
	log.Println("TABLE CREATED")
	return db, nil
}
