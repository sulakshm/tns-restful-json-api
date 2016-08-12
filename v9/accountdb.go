package main

import (
	_ "time"
	"encoding/json"
	"log"
	"fmt"
	"errors"
	"os"
	"encoding/gob"
)

// Account - a DbEntry
type Account struct {
	AccId	int `json:"id"`
	FirstName string `json:"first_name"`
	LastName string `json:"last_name"`
	Email	string	`json:"email"`
	Mobile	string	`json:"mobile"`
	Password	string `json:"password"`

	DeviceType	string `json:"device_type"`
	DeviceId	string `json:"device_id"`

	active	bool	// currently logged 
}

func (e* Account) Id() int {
	return e.AccId
}

func (e* Account) String() string {
	v, err := json.Marshal(e)	
	if err != nil {
		log.Fatal(err)
	}
	return string(v)
}

func (e* Account) Update(o DbEntry) error {
	n, ok := o.(*Account)
	if !ok {
		return errors.New("Type assertion failed")
	}
	e.FirstName = n.FirstName
	e.LastName = n.LastName
	e.Email = n.Email
	e.Mobile = n.Mobile
	e.Password = n.Password
	e.DeviceType = n.DeviceType
	e.DeviceId = n.DeviceId

	return nil
}


func (e* Account) Delete() error {
	return Accounts.Delete(e.Id())
}

func (e* Account) Encode(encoder *gob.Encoder) {
	err := encoder.Encode(e)
	if err != nil {
		log.Fatal("encode:" , err)
	}
}

func (e* Account) Match(username, passwd string) bool {
	if e.Email == username && e.Password == passwd {
		return true
	}
	return false
}

func NewAccount(first, last, email, mobile, passwd, devtype, devid string) (Account, error) {
	new := &Account{FirstName: first, LastName: last, 
			Email: email, Mobile: mobile,
			Password: passwd,
			DeviceType: devtype,
			DeviceId: devid}
	if err := Accounts.Create(new); err != nil {
		return Account{}, err
	}
	return *new, nil
}

type AccountsDb struct {
	DbTemplate
}

func (db* AccountsDb) Init(n string) {
	fmt.Printf("AccountsDb.Init(%s) called\n", n)
	// gob.Register(AccountsDb{}) // 
	gob.Register(Account{}) // 
	db.DbTemplate.Init(n, db)
	db.Load()

	// finally register handlers
	RegisterRoutes(acctRoutes)
}

func (db* AccountsDb) Decode(decoder *gob.Decoder) (*Account, error) {
	var act Account
	err := decoder.Decode(&act)
	if err != nil {
		fmt.Printf("decode: %s", err)
		return nil, err
	}
	return &act, nil
}

func (db* AccountsDb) Create(e DbEntry) error {
	n, ok := e.(*Account)
	if !ok {
		fmt.Printf("AccountsDB: Create() failed %s\n", e.String())
		return errors.New("Type assertion failure.")
	}

	n.AccId = db.NextId()
	db.Set(n.AccId, n)
	fmt.Printf("AccountsDB: Create() %s\n", e.String())
	return nil
}

func (db* AccountsDb) Update(e DbEntry) error {
	v, err := db.Find(e.Id())
	if err != nil {
		return err
	}
	return v.Update(e)
}

func (db* AccountsDb) Load() {
	// Open a RO file
	decodeFile, err := os.Open(db.Store())
	if err != nil {
		fmt.Printf("Db: %s, no store file found - %s\n",
				db.Name(), err)
		return
	}
	defer decodeFile.Close()

	// Create a decoder
	var mydb map[int]*Account 
	gob.Register(Account{})
	decoder := gob.NewDecoder(decodeFile)
	decoder.Decode(&mydb)

	for k, v := range mydb {
		db.DbTemplate.Set(k, v)
	}
}


func (db* AccountsDb) Commit() {
	file, err := os.Create(db.Store())
	if err != nil {
		panic(err)
	}
	defer file.Close()

	mydb := make(map[int]*Account)
	for k, v := range db.DbTemplate.Entry {
		mydb[k] = v.(*Account)
	}
	gob.Register(Account{})
	encoder := gob.NewEncoder(file)
	// Write to the file
	if err := encoder.Encode(mydb); err != nil {
		panic(err)
	}
}

func (db* AccountsDb) CredCheck(username, passwd string) (*Account, error) {
	db.DbTemplate.Lock()
	defer db.DbTemplate.Unlock()

	for _, v := range db.DbTemplate.Entry {
		e := v.(*Account)
		if e.Match(username, passwd) {
			return e, nil
		}
	}

	return nil, errors.New("Invalid credentials")
}

// Group

// Member


// Node
