package main

var Accounts	AccountsDb
var dbMap	map[string]Database

func RegisterDb(db Database) {
	dbMap[db.Name()] = db
}

func GetDb(db string) Database {
	return dbMap[db]
}

func init() {
	dbMap = make(map[string]Database)
	Accounts.Init("accounts")
	RegisterDb(&Accounts)

	// create few accounts
	// NewAccount("Chokkar", "Guruswamy", "chokkar@gmail.com", "1234512345", "password", "Android", "ABCDEF12345")
	// NewAccount("Ravi", "Karunas", "ravik@gmail.com", "5432154321", "password", "Android", "ABCDEF54321")
}


