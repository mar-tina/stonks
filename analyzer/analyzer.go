package analyzer

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

type Analyzer interface {
	Init()
	GetCurrentPrice(stonk []string) (map[string]interface{}, error)
	Convert(to, from string, price int) (map[string]interface{}, error)
}

type Fixer struct {
	client *http.Client
}

type CurrencyLayer struct {
	client *http.Client
}

func (f *Fixer) Init() {
	f.client = &http.Client{}
}

func GetEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

func (f *Fixer) makeRequest(req *http.Request) (map[string]interface{}, error) {
	resp, err := f.client.Do(req)
	if err != nil {
		log.Print("Failed to complete HTTP request . Please try again")
		return nil, fmt.Errorf("Could not complete request to Fixer analyzer %v", err.Error())
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Failed to read response body %v", err.Error())
		return nil, fmt.Errorf("Could not read response body %v", err.Error())
	}

	fix := make(map[string]interface{})
	json.Unmarshal(body, &fix)
	log.Printf("THE RESPONSE %v", fix)
	return fix, nil
}

func (f *Fixer) GetCurrentPrice(stonk []string) (map[string]interface{}, error) {
	key := GetEnv("FIXER_KEY", "INSERT YOUR API KEY for FIXER HERE FOR THE TESTS TO WORK")

	symbols := ""
	sep := ""
	for idx, val := range stonk {
		if idx >= 1 {
			sep = ","
		}
		symbols += fmt.Sprintf("%v%v", sep, val)
	}

	req, err := http.NewRequest("POST", fmt.Sprintf("http://data.fixer.io/api/latest?access_key=%v&symbols=%v", key, symbols), nil)
	if err != nil {
		return nil, fmt.Errorf("Could not complete request to Fixer analyzer %v", err.Error())
	}

	return f.makeRequest(req)
}

func (f *Fixer) Convert(to, from string, amount int) (map[string]interface{}, error) {
	key := GetEnv("FIXER_KEY", "INSERT YOUR API KEY HERE FOR THE TESTS TO WORK")

	req, err := http.NewRequest("POST", fmt.Sprintf("http://data.fixer.io/api/convert?access_key=%v&from=%v&to=%v&amount=%v", key, from, to, amount), nil)
	if err != nil {
		return nil, fmt.Errorf("Could not complete request to Fixer analyzer %v", err.Error())
	}

	return f.makeRequest(req)
}

func (c *CurrencyLayer) makeRequest(req *http.Request) (map[string]interface{}, error) {
	resp, err := c.client.Do(req)
	if err != nil {
		log.Print("Failed to complete HTTP request . Please restart the application")
		return nil, fmt.Errorf("Could not complete request to Fixer analyzer %v", err.Error())
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Failed to read response body %v", err.Error())
		return nil, fmt.Errorf("Could not read response body %v", err.Error())
	}

	fix := make(map[string]interface{})
	json.Unmarshal(body, &fix)
	log.Printf("THE RESPONSE %v", fix)
	return fix, nil
}

func (c *CurrencyLayer) Init() {
	c.client = &http.Client{}
}
