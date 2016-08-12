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

// Understand http authorization here...

// Endpoint: POST /v1/account/register, AuthToken: none
// Endpoint: POST /v1/account/login, Authorization: Basic <>
//	 Respond with Token.
// Should carry: Authorization: Token <key>
// Endpoint: POST /v1/account/logout, AuthToken: yes
// Endpoint: POST /v1/account/setpasswd, AuthToken: yes
// Endpoint: PUT /v1/account/resetpasswd, AuthToken: yes, -d "email=value"
// Endpoint: PUT /v1/account/users, AuthToken: yes, "update entries of Account"
// Endpoint: GET /v1/account/users, AuthToken: yes, "get entries of account"
// Endpoint: DELETE /v1/account/users, AuthToken: yes, "delete entry of account"
// Endpoint: POST /v1/account/auth/google, -- dropped


var acctRoutes = []Route{
	Route{
		"AccountCreate",
		"POST",
		"/v1/account/register",
		AccountCreate,
	},
	Route{
		"AccountLogin",
		"POST",
		"/v1/account/login",
		AccountLogin,
	},
	Route{
		"AccountLogout",
		"POST",
		"/v1/account/logout",
		AccountLogout,
	},
	Route{
		"AccountsIndex",
		"GET",
		"/v1/account/users",
		AccountsIndex,
	},
	Route{
		"AccountShow",
		"GET",
		"/v1/account/users/{Id}",
		AccountShow,
	},
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

func AccountLogin(w http.ResponseWriter, r *http.Request) {
	var act *Account

	// Header: Authorization: Basic <key> 
	username, passwd, ok := r.BasicAuth()
	if !ok {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		http.Error(w, "Authorization failed", http.StatusUnauthorized)
		return
	}

	/*
	if v, ok := r.Header['Device_Type']; ok {
		devType = v
		fmt.Printf("Got DeviceType: %s\n", devType)
	}

	if v, ok := r.Header['Device_Id']; ok {
		fmt.Printf("Got DeviceId: %s\n", v)
	} */
	// XXX: why should login carry body, should be part of header.
	// Header: device_type: "Android" | "iOS"
	// Header: deveice_id: "devId"
	body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1048576))
	if err != nil {
		panic(err)
	}
	if err := r.Body.Close(); err != nil {
		panic(err)
	}
	type Info struct {
		DeviceId	string `json:"device_id"`
		DeviceType	string `json:"device_type"`
	}
	var inf Info

	if err := json.Unmarshal(body, &inf); err != nil {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(422) // unprocessable entity
		if err := json.NewEncoder(w).Encode(err); err != nil {
			panic(err)
		}
	} 


	db := GetDb("accounts")
	act, err = db.(*AccountsDb).CredCheck(username, passwd)
	if err != nil {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		http.Error(w, "Authorization failed", http.StatusUnauthorized)
		return
	}

	act.DeviceType = inf.DeviceType
	act.DeviceId = inf.DeviceId
	act.active = true

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
}

func AccountLogout(w http.ResponseWriter, r *http.Request) {
	var act *Account

	// Header: Authorization: Basic <key> 
	username, passwd, ok := r.BasicAuth()
	if !ok {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		http.Error(w, "Authorization failed", http.StatusUnauthorized)
		return
	}

	/*
	if v, ok := r.Header['Device_Type']; ok {
		devType = v
		fmt.Printf("Got DeviceType: %s\n", devType)
	}

	if v, ok := r.Header['Device_Id']; ok {
		fmt.Printf("Got DeviceId: %s\n", v)
	} */
	// XXX: why should login carry body, should be part of header.
	// Header: device_type: "Android" | "iOS"
	// Header: deveice_id: "devId"
	body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1048576))
	if err != nil {
		panic(err)
	}
	if err := r.Body.Close(); err != nil {
		panic(err)
	}
	type Info struct {
		DeviceId	string `json:"device_id"`
		DeviceType	string `json:"device_type"`
	}
	var inf Info

	if err := json.Unmarshal(body, &inf); err != nil {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(422) // unprocessable entity
		if err := json.NewEncoder(w).Encode(err); err != nil {
			panic(err)
		}
	} 


	db := GetDb("accounts")
	act, err = db.(*AccountsDb).CredCheck(username, passwd)
	if err != nil { // ignore error
		goto out
	}

	act.DeviceType = inf.DeviceType
	act.DeviceId = inf.DeviceId
	act.active = false

out:
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
}


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

