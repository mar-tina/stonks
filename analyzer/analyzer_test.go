package analyzer

import (
	"log"
	"testing"
)

func TestFixerGetCurrentPrice(t *testing.T) {
	// Init Analyzer
	lyzer := Init()

	res, err := lyzer.GetCurrentPrice([]string{"USD", "EUR"})
	if err != nil {
		log.Printf("Failed to get current price")
	}

	got := res["success"]
	want := true

	if got != want {
		t.Errorf("Got %v, want %v", got, want)
	}
}

func TestConvert(t *testing.T) {
	// Init Analyzer
	lyzer := Init()

	res, err := lyzer.Convert("USD", "KES", 10000)
	if err != nil {
		t.Errorf("Failed to complete conversion %v ", err.Error())
	}

	got := res["success"]
	want := true

	log.Printf("THE RES %v", res["result"])

	if got != want {
		t.Errorf("Got %v, want %v", got, want)
	}
}
