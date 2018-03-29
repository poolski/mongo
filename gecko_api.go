package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

type Client struct {
	BaseURL  string
	Dataset  string
	Username string
	Password string `default:""`
}

type MonzoDataset struct {
	Fields struct {
		ID struct {
			Type     string `json:"type"`
			Name     string `json:"name"`
			Optional bool   `json:"optional"`
		} `json:"id"`
		Amount struct {
			Type     string `json:"type"`
			Name     string `json:"name"`
			Optional bool   `json:"optional"`
		} `json:"amount"`
		Category struct {
			Type     string `json:"type"`
			Name     string `json:"name"`
			Optional bool   `json:"optional"`
		} `json:"category"`
		Created struct {
			Type     string `json:"type"`
			Name     string `json:"name"`
			Optional bool   `json:"optional"`
		} `json:"created"`
		Merchant struct {
			Type     string `json:"type"`
			Name     string `json:"name"`
			Optional bool   `json:"optional"`
		} `json:"merchant"`
	} `json:"fields"`
	UniqueBy []string `json:"unique_by"`
}

func NewClient(base, dataset, apikey string) *Client {
	return &Client{
		BaseURL:  base,
		Dataset:  dataset,
		Username: apikey,
	}
}

// This needs to be implemented yet.
func (s *Client) CreateDataset(ds string) error {
	fmt.Println("Creating Dataset")
	// url := s.BaseURL + "datasets/" + ds
	tpl := []byte(`{
	"fields": {
		"id": {
			"type": "string",
			"name": "Transaction ID",
			"optional": false
		},
		"amount": {
			"type": "number",
			"name": "Amount",
			"optional": false
		},
		"category": {
			"type": "string",
			"name": "Category",
			"optional": false
		},
		"created": {
			"type": "datetime",
			"name": "Transaction Date",
			"optional": false
		},
		"merchant": {
			"type": "string",
			"name": "Merchant",
			"optional": false
		}
	},
	"unique_by": ["id"]
}`)
	j, err := json.Marshal(tpl)
	// req, err := http.NewRequest("PUT", url, nil)
	if err != nil {
		return err
	}

	var md MonzoDataset
	err = json.Unmarshal(j, &md)
	fmt.Printf("%+v", md)
	return nil
	// _, err = s.doRequest(req)
	// return err
}

func (s *Client) UpdateDataset(data Results) error {
	url := s.BaseURL + "datasets/" + s.Dataset + "/data"
	fmt.Println(url)
	j, err := json.Marshal(data)
	if err != nil {
		return err
	}
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(j))
	if err != nil {
		return err
	}
	_, err = s.doRequest(req)
	return err
}

func (s *Client) doRequest(req *http.Request) ([]byte, error) {
	req.SetBasicAuth(s.Username, s.Password)
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Print("Something broked")
		return nil, err
	}
	switch resp.StatusCode {
	case 200:
		fmt.Println("Success!")
		return body, nil
	case 404:
		fmt.Println("Not found. Creating")
		s.CreateDataset(s.Dataset)
	default:
		return nil, fmt.Errorf("%s", body)
	}
	return body, nil
}
