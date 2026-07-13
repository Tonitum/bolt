package database

import (
	"database/sql"
	"fmt"
	_ "modernc.org/sqlite"
)

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

// TODO: return a struct? not just a string? that's kinda OOP
// TODO: we should conditionally upsert, but that will require some user PutID
// also, maybe another flag in the incoming request to overwrite existing?
// also, maybe store a database column to indicate team/group?
func PutID(id string, url string) (bool, error) {
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
	stmt, err := conn.Prepare(command)
	if err != nil {
		fmt.Println("Failed to prepare upsert statement")
		fmt.Printf("err: %v\n", err)
		return false, err
	}
	res, err := stmt.Exec(id, url, url, id)
	if err != nil {
		fmt.Println("Failed to execute upsert statement")
		fmt.Printf("err: %v\n", err)
		return false, err
	}
	fmt.Println(res.RowsAffected())

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
