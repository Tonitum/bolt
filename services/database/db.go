package database

// TODO: implement a real database interface
// store in memory for now

var database map[string]string = map[string]string{} // TODO: is this good / right syntax?

// TODO: return a struct not just a string? that's kinda OOP
func GetURL(id string) string {
	// how do I type with multiple returns?
	// I want to return nil if the provided id doesn't exist
	// TODO: This will always break if we get an ask for an ID that doesn't exist
	// 	options for handling non-existent ids
	// 		- return 404 < easier
	// 		- return a prefilled page (oops! that id doesn't exist)
	return database[id]
}

// TODO: return a struct not just a string? that's kinda OOP
func PutID(id string, url string) {
	// how do I type with multiple returns?
	// I want to return nil if the provided id doesn't exist
	database[id] = url
}


func Dump() map[string]string {
	return database
}
