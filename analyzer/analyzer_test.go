package analyzer

import (
	"log"
	"testing"
)

func TestGetCurrentPrice(t *testing.T) {
	// Init Analyzer
	var lyzer Analyzer
	lyzer = &Fixer{}
	lyzer.Init()

	res, err := lyzer.GetCurrentPrice([]string{"USD", "EUR"})
	if err != nil {
		log.Printf("Failed to get current price")
	}

	stonk := make(map[string]interface{})
	stonk["success"] = true
	got := res["success"]
	want := stonk["success"]

	if got != want {
		t.Errorf("Got %v, want %v", got, want)
	}
}

func TestConvert(t *testing.T) {
	// Init Analyzer
	var lyzer Analyzer
	lyzer = &Fixer{}
	lyzer.Init()

	res, err := lyzer.Convert("USD", "EUR", 25)
	if err != nil {
		log.Printf("Failed to get current price")
	}

	stonk := make(map[string]interface{})
	stonk["success"] = true
	got := res["success"]
	want := stonk["success"]

	if got != want {
		t.Errorf("Got %v, want %v", got, want)
	}
}
