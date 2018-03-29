package main

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/go-redis/redis"
	"github.com/tjvr/go-monzo"
)

type Config struct {
	AccountID       string        `json:"account_id"`
	AccessToken     string        `json:"access_token"`
	ListenPort      string        `json:"listen_port"`
	CallbackURL     string        `json:"callback_url"`
	GeckoURL        string        `json:"gecko_url"`
	GeckoAPIKey     string        `json:"gecko_apikey"`
	MonzoURL        string        `json:"monzo_url"`
	MonzoClientID   string        `json:"monzo_client_id"`
	MonzoSecret     string        `json:"monzo_secret"`
	RefreshInterval time.Duration `json:"refresh_interval"`
}

type GeckoData struct {
	ID       string  `json:"id"`
	Amount   float64 `json:"amount"`
	Category string  `json:"category"`
	Created  string  `json:"created"`
	Merchant string  `json:"merchant"`
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

type Results struct {
	Data []GeckoData `json:"data"`
}

func main() {
	c := LoadConfig("config.json")
	fmt.Println(c.RefreshInterval)
	for x := range time.NewTicker(time.Duration(c.RefreshInterval) * time.Second).C {
		fmt.Println(string(x.Format(time.RFC822Z)) + " - Running updates")
		getTransactions()
	}
}

func getTransactions() {
	c := LoadConfig("config.json")
	gc := NewClient(c.GeckoURL, "monzo.transactions", c.GeckoAPIKey)
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

	// t, err := time.Parse(time.RFC3339, "2018-03-01T23:00:00Z")
	// s := t.Format(time.RFC3339)

	b := time.Now().Format(time.RFC3339)

	mp := monzo.Parameters{
		AccountID:      c.AccountID,
		ExpandMerchant: true,
		Before:         b,
	}

	txs, err := cl.Transactions(mp)
	if err != nil {
		panic(err)
	}
	var res Results
	for i := range txs {
		if txs[i].Scheme != "payport_faster_payments" {
			var gd GeckoData
			gd.ID = txs[i].ID
			gd.Amount = (float64(txs[i].Amount) / 100) * -1
			gd.Category = txs[i].Category
			gd.Created = txs[i].Created
			gd.Merchant = txs[i].Merchant.Name
			res.Data = append(res.Data, gd)
		}
	}
	err = gc.UpdateDataset(res)
	if err != nil {
		fmt.Println(err)
	}
}

func must(err error) {
	if err != nil {
		panic(err)
	}
}
