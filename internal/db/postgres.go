package db

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

type DB struct{
	db *sql.DB 
	init bool
}

func NewDb(connStr string) (*DB, error){
	db := &DB{db: nil, init: false}
	if connStr == ""{
		fmt.Println("Empty address data base")
		return db, nil
	}
//	connStr := "user=postgres sslmode=disable"
	conn, err := sql.Open("postgres", connStr)
	if err != nil {
		return db, fmt.Errorf(" Can't open DB %w", err)
	}
	db.db = conn
	db.init = true
	return db, nil
}

func (db *DB) Close() error {
	fmt.Println("CLOSE : ",db)
	if !db.init{
		return fmt.Errorf(" Can't close DB : DataBase wasn't inited")
	}
	if err := db.db.Close(); err != nil{
		return fmt.Errorf(" Can't close DB %w", err)
	}
	return nil
}

func (db DB) Ping() error{
	fmt.Println("PING : ",db)
	if !db.init{
		fmt.Println(" Can't ping DB : DataBase wasn't inited")
		return fmt.Errorf(" Can't ping DB : DataBase wasn't inited")
	}
	if err := db.db.Ping(); err != nil{
		return fmt.Errorf(" Can't ping DB %w", err)
	}
	return nil
}
