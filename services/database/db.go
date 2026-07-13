package database

import (
	"database/sql"
	"fmt"

	_ "modernc.org/sqlite"
)

// TODO: implement a real database interface
// store in memory for now

var conn *sql.DB

func InitDB() {
	var err error
	conn, err = sql.Open("sqlite", "./bolt.db")
	if err != nil {
		fmt.Println("Failed to connect to database")
		fmt.Printf("err: %v\n", err)
		return
	}
	sql := `CREATE TABLE IF NOT EXISTS aliases (
		alias TEXT PRIMARY KEY NOT NULL,
		url TEXT NOT NULL
	)`
	_, err = conn.Exec(sql)
	if err != nil {
		fmt.Println("failed to initialize database")
		fmt.Printf("err: %v\n", err)
		return
	}
}

// TODO: return a struct not just a string? that's kinda OOP
func GetURL(id string) string {
	// how do I type with multiple returns?
	// I want to return nil if the provided id doesn't exist
	// TODO: This will always break if we get an ask for an ID that doesn't exist
	// 	options for handling non-existent ids
	// 		- return 404 < easier
	// 		- return a prefilled page (oops! that id doesn't exist)

	stmt, err := conn.Prepare("SELECT url FROM aliases WHERE alias = ?")
	if err != nil {
		fmt.Println("Failed to prepare statement")
		fmt.Printf("err: %v\n", err)
		return ""
	}
	defer stmt.Close()
	var resulturl string
	err = stmt.QueryRow(id).Scan(&resulturl)
	if err != nil {
		fmt.Println("failed to load results")
		fmt.Printf("err: %v\n", err)
		return ""
	}
	return resulturl
}

// TODO: return a struct not just a string? that's kinda OOP
func PutID(id string, url string) (bool, error) {
	// how do I type with multiple returns?
	// I want to return nil if the provided id doesn't exist
	stmt, err := conn.Prepare("INSERT INTO aliases (alias,url) VALUES (?,?)")
	if err != nil {
		fmt.Println("Failed to prepare statement")
		fmt.Printf("err: %v\n", err)
		var retErr error
		return false, retErr
	}
	_, err = stmt.Exec(id, url)
	if err != nil {
		fmt.Println("Failed to prepare statement")
		fmt.Printf("err: %v\n", err)
		var retErr error
		return false, retErr
	}
	return true, nil
}

func Dump() map[string]string {
	stmt, err := conn.Prepare("SELECT * FROM aliases")
	if err != nil {
		fmt.Println("Failed to prepare statement")
		fmt.Printf("err: %v\n", err)
		return nil
	}
	defer stmt.Close()
	rows, err := stmt.Query()
	if err != nil {
		fmt.Println("failed to load results")
		fmt.Printf("err: %v\n", err)
		return nil
	}
	defer rows.Close()
	var database map[string]string = map[string]string{} // TODO: is this good / right syntax?
	var alias string
	var url string

	for rows.Next() {
		rows.Scan(&alias, &url)
		database[alias] = url
	}
	return database
}
