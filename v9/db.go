package main

import (
	"fmt"
	"sync"
	"errors"
	_ "log"
	_ "os"
	"encoding/gob"
)

type DbEntry interface {
	Id()		int
	String()	string
	Delete()	error
	Update(DbEntry)	error
	Encode(encoder *gob.Encoder)
}

type Database interface {
	Init(string)
	Store() string
	Create(DbEntry)	error
	Update(DbEntry)	error
	Find(int)		(DbEntry, error)
	Delete(int)		error
	Show()		[]string
	Name()		string
	Commit() 
	//Decode(*gob.Decoder) (DbEntry, error)
}

type DbTemplate struct {
	Entry	map[int]DbEntry
	sync.Mutex
	currentId int
	name    string
	concrete interface{}
}

func(t* DbTemplate) Init(n string, con interface{}) {
	t.Entry = make(map[int]DbEntry)
	t.name = n
	t.concrete = con
	//t.Load()
}
func(t *DbTemplate) NextId() int {
	t.Lock()
	defer t.Unlock()
	t.currentId += 1
	return t.currentId
}

func(t* DbTemplate) Name() string {
	return t.name
}

func (t* DbTemplate) Set(k int, v DbEntry) {
	t.Lock()
	defer t.Unlock()
	t.Entry[k] = v
	if t.currentId < k {
		t.currentId = k
	}
}

func (t* DbTemplate) Find(k int) (DbEntry, error) {
	t.Lock()
	defer t.Unlock()

	if v, ok := t.Entry[k]; ok {
		return v, nil
	}
	return nil, errors.New("key not found")
}

func (t* DbTemplate) Store() string {
	return fmt.Sprintf("%s.gob", t.Name())
}

func (t* DbTemplate) Delete(k int) error {
	t.Lock()
	defer t.Unlock()

	if _, ok := t.Entry[k]; ok {
		delete(t.Entry, k)
		return nil
	}
	return errors.New("key not found")
}

func (t* DbTemplate) Show() []string {
	t.Lock()
	defer t.Unlock()
	var res []string
	fmt.Printf("DB: Length %d\n", len(t.Entry))
	for _, val := range t.Entry {
		res = append(res, val.String())
	}
	fmt.Printf("DB: Show returned %v\n", res)
	return res
}

/*
func (t* DbTemplate) Load() {
	fmt.Printf("Loading from store: %s\n", t.Store())
	// Open a RO file
	decodeFile, err := os.Open(t.Store())
	if err != nil {
		fmt.Printf("Db: %s, no store file found - %s\n",
				t.Name(), err)
		return
	}
	defer decodeFile.Close()

	// Create a decoder
	decoder := gob.NewDecoder(decodeFile)
	// decoder.Decode(&t.Entry)
	for {
		item, err := t.concrete.Decode(decoder)
		if err != nil {
			break
		}
		t.Entry[item.Id()] = item
	}
	fmt.Printf(".. got inited to %d entries.\n", len(t.Entry))
}

func (t* DbTemplate) Commit() {
	fmt.Printf("Committing to store: %s, len %d\n", t.Store(), len(t.Entry))
	file, err := os.Create(t.Store())
	if err != nil {
		panic(err)
	}
	defer file.Close()

	encoder := gob.NewEncoder(file)
	// Write to the file
	// if err := encoder.Encode(t.Entry); err != nil {
	// 	panic(err)
	// } 
	// for _, v := range t.Entry {
	// 	v.Encode(encoder)
	// }
}
*/
