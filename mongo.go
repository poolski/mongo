package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/go-redis/redis"
	"github.com/tjvr/go-monzo"
)

type Config struct {
	AccessToken   string `json:"access_token"`
	ListenPort    string `json:"listen_port"`
	CallbackURL   string `json:"callback_url"`
	GeckoURL      string `json:"gecko_url"`
	GeckoAPIKey   string `json:"gecko_apikey"`
	MonzoURL      string `json:"monzo_url"`
	MonzoClientID string `json:"monzo_client_id"`
	MonzoSecret   string `json:"monzo_secret"`
}

func LoadConfig(file string) Config {
	var c Config
	cf, err := os.Open(file)
	defer cf.Close()
	if err != nil {
		fmt.Println(err.Error)
	}
	jp := json.NewDecoder(cf)
	jp.Decode(&c)
	return c
}

func main() {
	// c := LoadConfig("config.json")
	red := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})

	at, err := red.Get("monzoAuthToken").Result()
	if err != nil {
		panic(err)
	}
	cl := monzo.Client{
		BaseURL:     "https://api.monzo.com",
		AccessToken: at,
	}

	txs, err := cl.Transactions("acc_000096il5DQUqpBLrchOMb", false)
	tx := txs[1000]

	txd, err := cl.Transaction(tx.ID)
	fmt.Println(txd.Category)

	// auth := monzo.NewAuthenticator(c.MonzoClientID, c.MonzoSecret, c.CallbackURL)
	//
	// http.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
	// 	session := auth.GetSession(w, req)
	// 	at, err := red.Get("monzoAuthToken").Result()
	// 	if err != nil {
	// 		panic(err)
	// 	}
	// 	if at != "" {
	// 		session.Client.AccessToken = at
	// 	}
	// 	fmt.Println(session)
	// 	if !session.IsAuthenticated() {
	// 		http.Error(w, "Not Authenticated", 403)
	// 		return
	// 	}
	// })
	//
	// http.HandleFunc("/login", auth.Login)
	//
	// http.HandleFunc("/callback", func(w http.ResponseWriter, req *http.Request) {
	// 	session := auth.Callback(w, req)
	// 	if session == nil {
	// 		// TODO something went wrong
	// 	} else {
	// 		err := red.Set("monzoAuthToken", session.Client.AccessToken, time.Hour).Err()
	// 		if err != nil {
	// 			fmt.Println("ERROR!!!!!")
	// 			fmt.Println(err)
	// 		}
	// 		// http.Redirect(w, req, "/", http.StatusTemporaryRedirect)
	// 	}
	// })
	// http.HandleFunc("/logout", func(w http.ResponseWriter, req *http.Request) {
	// 	auth.Logout(w, req)
	// 	http.Redirect(w, req, "/", http.StatusTemporaryRedirect)
	// })
	//
	// http.ListenAndServe(":"+c.ListenPort, nil)
}

func must(err error) {
	if err != nil {
		panic(err)
	}
}
