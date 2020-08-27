package analyzer

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
)

var (
	FIXER_KEY = ""
	CL_KEY    = "8b4fcd671369c6bf9c332d88843333bf"
)

type Analyzer interface {
	GetCurrentPrice(stonk []string) (map[string]interface{}, error)
	Convert(to, from string, amount float64) (map[string]interface{}, error)
}

type Fixer struct {
	client *http.Client
}

type CurrencyLayer struct {
	client *http.Client
}

// GenericAnalyzer holds a reference to both Fixer and CurrencyAnalyzer structs and
//fallsback on either or function calls when one fails.
type GenericAnalyzer struct {
	fixer *Fixer
	cl    *CurrencyLayer
}

func GetEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

//Init initializes the GenericAnalyzer providing references to both
//Fixer and CurrencyAnalyzer structs
func Init() Analyzer {
	lyzer := &GenericAnalyzer{
		fixer: &Fixer{},
		cl:    &CurrencyLayer{},
	}
	return lyzer
}

func (g *GenericAnalyzer) Convert(to, from string, amount float64) (map[string]interface{}, error) {
	fixerResponse, err := g.fixer.Convert(to, from, amount)
	if err != nil || fixerResponse["success"] != true {
		return g.cl.Convert(to, from, amount)
	}
	return fixerResponse, nil
}

func (g *GenericAnalyzer) GetCurrentPrice(stonk []string) (map[string]interface{}, error) {
	fixerResponse, err := g.fixer.GetCurrentPrice(stonk)
	if err != nil || fixerResponse["success"] != true {
		return g.cl.GetCurrentPrice(stonk)
	}
	return fixerResponse, nil
}

//GetCurrentPrice takes in a uri string that holds the endpoint and an array of
//stock names , Creates a new http request and returns it to the caller.
func GetCurrentPrice(uri string, stonk []string) (*http.Request, error) {
	var key string
	symbols := ""
	sep := ""
	// Loop over the stocks array and create a comma separated string to pass as a
	// query parameter
	for idx, val := range stonk {
		if idx >= 1 {
			sep = ","
		}
		symbols += fmt.Sprintf("%v%v", sep, val)
	}

	//switch the query param name depending on the endpoint being hit
	var currQueryParam string
	if strings.Contains(uri, "fixer") {
		currQueryParam = "symbols"
		key = GetEnv("FIXER_KEY", FIXER_KEY)
	} else {
		currQueryParam = "currencies"
		key = GetEnv("CL_KEY", CL_KEY)
	}

	req, err := http.NewRequest("POST", fmt.Sprintf("%v?access_key=%v&%v=%v", uri, key, currQueryParam, symbols), nil)
	if err != nil {
		return nil, fmt.Errorf("Could not complete request to Fixer analyzer %v", err.Error())
	}

	return req, nil
}

//makeRequest intiates the http request and returns the response to the caller.
func makeRequest(req *http.Request) (map[string]interface{}, error) {
	client := http.Client{}
	raw, err := client.Do(req)
	if err != nil {
		log.Print("Failed to complete HTTP request . Please restart the application")
		return nil, fmt.Errorf("Could not complete request to Fixer analyzer %v", err.Error())
	}

	body, err := ioutil.ReadAll(raw.Body)
	if err != nil {
		log.Printf("Failed to read response body %v", err.Error())
		return nil, fmt.Errorf("Could not read response body %v", err.Error())
	}

	response := make(map[string]interface{})
	json.Unmarshal(body, &response)
	return response, nil
}

//Convert is the default convert function that is ran. If access is limited the conversion is
//done manually by the fallbackConvert function
func Convert(uri, to, from string, amount float64) (map[string]interface{}, error) {
	var key string
	//Set the access key according to the endpoint being hit.
	if strings.Contains(uri, "fixer") {
		key = GetEnv("FIXER_KEY", FIXER_KEY)
	} else {
		key = GetEnv("CL_KEY", CL_KEY)
	}

	req, err := http.NewRequest("POST", fmt.Sprintf("%v?access_key=%v&from=%v&to=%v&amount=%v", uri, key, from, to, amount), nil)
	if err != nil {
		return nil, fmt.Errorf("Could not complete request to Fixer analyzer %v", err.Error())
	}

	return makeRequest(req)
}

func BaseCoversion(base, rate, amount float64) float64 {
	temp := base * amount
	return temp / rate
}

//GetCurrentPrice calls the generic GetCurrentPrice function for the fixer API
func (f *Fixer) GetCurrentPrice(stonk []string) (map[string]interface{}, error) {
	req, err := GetCurrentPrice("http://data.fixer.io/api/latest", stonk)
	if err != nil {
		return nil, err
	}
	return makeRequest(req)
}

func (f *Fixer) Convert(to, from string, amount float64) (map[string]interface{}, error) {
	res, err := Convert("http://data.fixer.io/api/convert", to, from, amount)
	if err != nil || res["success"] != true {
		return f.FallbackConvert(to, from, amount)
	}
	return res, nil
}

func (f *Fixer) FallbackConvert(to, from string, amount float64) (map[string]interface{}, error) {
	fallbackResponse := make(map[string]interface{})
	res, err := f.GetCurrentPrice([]string{to, from, "EUR"})
	if err != nil {
		return nil, fmt.Errorf("Conversion failed. Please try again later")
	}

	if res["rates"] != nil {
		// base is the base currency set by the API provider .
		// Conversion is done like this: FROM --> BASE --> TO
		base := res["rates"].(map[string]interface{})["EUR"].(float64)
		fromConverted := BaseCoversion(base, res["rates"].(map[string]interface{})[from].(float64), amount)
		toConverted := BaseCoversion(base, res["rates"].(map[string]interface{})[to].(float64), fromConverted)

		fallbackResponse["success"] = true
		fallbackResponse["result"] = toConverted
		return fallbackResponse, nil
	}

	fallbackResponse["success"] = false
	fallbackResponse["message"] = "Conversion failed. Please try again later"
	return fallbackResponse, nil
}

func (c *CurrencyLayer) GetCurrentPrice(stonk []string) (map[string]interface{}, error) {
	req, err := GetCurrentPrice("http://api.currencylayer.com/live", stonk)
	if err != nil {
		return nil, err
	}
	return makeRequest(req)
}

func (c *CurrencyLayer) Convert(to, from string, amount float64) (map[string]interface{}, error) {
	res, err := Convert("http://api.currencylayer.com/convert", to, from, amount)
	if err != nil || res["success"] != true {
		return c.FallbackConvert(to, from, amount)
	}
	return res, nil
}

func (c *CurrencyLayer) FallbackConvert(to, from string, amount float64) (map[string]interface{}, error) {
	fallbackResponse := make(map[string]interface{})
	res, err := c.GetCurrentPrice([]string{to, from, "EUR"})
	if err != nil {
		return nil, fmt.Errorf("Conversion failed. Please try again later")
	}
	if res["quotes"] != nil {
		// base is the base currency set by the API provider .
		// Conversion is done like this: FROM --> BASE --> TO . Chose to use EUR consistently for both FIXER and
		// CURRENCYANALYZER Endpoints
		base := res["quotes"].(map[string]interface{})["USDEUR"].(float64)
		fromConverted := BaseCoversion(base, res["quotes"].(map[string]interface{})["USD"+from].(float64), amount)
		toConverted := BaseCoversion(base, res["quotes"].(map[string]interface{})["USD"+to].(float64), fromConverted)

		fallbackResponse["success"] = true
		fallbackResponse["result"] = toConverted
		return fallbackResponse, nil
	}

	fallbackResponse["success"] = false
	fallbackResponse["message"] = "Conversion failed. Please try again later"
	return fallbackResponse, nil
}
