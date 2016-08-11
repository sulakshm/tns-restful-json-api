package main

import (
	"encoding/json"
	_ "fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

func AccountsIndex(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	res := GetDb("accounts").Show()
	if err := json.NewEncoder(w).Encode(res); err != nil {
		panic(err)
	}
}

func AccountShow(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	var id int
	var err error
	if id, err = strconv.Atoi(vars["Id"]); err != nil {
		panic(err)
	}

	var info DbEntry
	if info, err = GetDb("accounts").Find(id); err != nil {
		// If we didn't find it, 404
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusNotFound)
		if err := json.NewEncoder(w).Encode(jsonErr{Code: http.StatusNotFound, Text: "Not Found"}); err != nil {
			panic(err)
		}
	} else {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(info.String()); err != nil {
			panic(err)
		}
	}
}

/*
Test with this curl command:

curl -H "Content-Type: application/json" -d '{"name":"New Todo"}' http://localhost:8080/accounts

*/

func AccountCreate(w http.ResponseWriter, r *http.Request) {
	var act Account
	body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1048576))
	if err != nil {
		panic(err)
	}
	if err := r.Body.Close(); err != nil {
		panic(err)
	}
	if err := json.Unmarshal(body, &act); err != nil {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(422) // unprocessable entity
		if err := json.NewEncoder(w).Encode(err); err != nil {
			panic(err)
		}
	} else if err := Accounts.Create(&act); err != nil {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(422) // unprocessable entity
		if err := json.NewEncoder(w).Encode(err); err != nil {
			panic(err)
		}
	} else {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusCreated)
		if err := json.NewEncoder(w).Encode(act.String()); err != nil {
			panic(err)
		}
	}
}
