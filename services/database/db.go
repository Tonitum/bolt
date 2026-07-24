package database

import (
	"database/sql"
	"fmt"
	_ "modernc.org/sqlite"
)

type Database interface {
	// todo: define get/set/delete/list functions
	Init() error
	GetURL(alias string) (string, error)
	PutAlias(alias string, url string) (bool, error)
	DeleteAlias(alias string) (bool, error)
	ListAliases() (map[string]string, error)
}

type SQLITE struct {
	conn *sql.DB
}

var DB Database

func InitDB() {
	// todo: determine database type based on configuration
	DB = &SQLITE{}
	DB.Init()
}

func (db *SQLITE) Init() error {
	var err error
	db.conn, err = sql.Open("sqlite", "./bolt.db")
	if err != nil {
		fmt.Println("Failed to connect to database")
		fmt.Printf("err: %v\n", err)
		return err
	}
	sql := `CREATE TABLE IF NOT EXISTS aliases (
		alias TEXT PRIMARY KEY NOT NULL,
		url TEXT NOT NULL
	)`
	_, err = db.conn.Exec(sql)
	if err != nil {
		fmt.Println("failed to initialize database")
		fmt.Printf("err: %v\n", err)
		return err
	}
	return nil
}

func (db *SQLITE) GetURL(alias string) (string, error) {
	// how do I type with multiple returns?
	// I want to return nil if the provided id doesn't exist
	// TODO: This will always break if we get an ask for an ID that doesn't exist
	// 	options for handling non-existent ids
	// 		- return 404 < easier
	// 		- return a prefilled page (oops! that id doesn't exist)

	stmt, err := db.conn.Prepare("SELECT url FROM aliases WHERE alias = ?")
	if err != nil {
		fmt.Println("Failed to prepare statement")
		fmt.Printf("err: %v\n", err)
		return "", err
	}
	defer stmt.Close()
	var resulturl string
	err = stmt.QueryRow(alias).Scan(&resulturl)
	if err != nil {
		fmt.Println("failed to load results")
		fmt.Printf("err: %v\n", err)
		return "", err
	}
	return resulturl, nil
}

// TODO: return a struct? not just a string? that's kinda OOP
// TODO: we should conditionally upsert, but that will require some user PutID
// also, maybe another flag in the incoming request to overwrite existing?
// also, maybe store a database column to indicate team/group?
func (db SQLITE) PutAlias(alias string, url string) (bool, error) {
	// we want to do an upsert
	// if the record exists, update
	// if it does not, insert
	var command = `INSERT INTO aliases (alias, url) 
	VALUES(?,?)
	ON CONFLICT(alias)
	DO
		UPDATE SET url = ?
		WHERE alias = ?
	`
	stmt, err := db.conn.Prepare(command)
	if err != nil {
		fmt.Println("Failed to prepare upsert statement")
		fmt.Printf("err: %v\n", err)
		return false, err
	}
	res, err := stmt.Exec(alias, url, url, alias)
	if err != nil {
		fmt.Println("Failed to execute upsert statement")
		fmt.Printf("err: %v\n", err)
		return false, err
	}
	fmt.Println(res.RowsAffected())

	return true, nil
}

func (db *SQLITE) DeleteAlias(alias string) (bool, error) {
	return true, nil
}

func (db *SQLITE) ListAliases() (map[string]string, error) {
	stmt, err := db.conn.Prepare("SELECT * FROM aliases")
	if err != nil {
		fmt.Println("Failed to prepare statement")
		fmt.Printf("err: %v\n", err)
		return nil, err
	}
	defer stmt.Close()
	rows, err := stmt.Query()
	if err != nil {
		fmt.Println("failed to load results")
		fmt.Printf("err: %v\n", err)
		return nil, err
	}
	defer rows.Close()
	var database map[string]string = map[string]string{} // TODO: is this good / right syntax?
	var alias string
	var url string

	for rows.Next() {
		rows.Scan(&alias, &url)
		database[alias] = url
	}
	return database, nil
}
