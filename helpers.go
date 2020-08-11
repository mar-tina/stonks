package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"os"
)

func ReadCSVFile(path string) error {
	f, err := os.Open(path)
	if err != nil {
		return fmt.Errorf("Something went wrong reading the file %v", err.Error())
	}

	data := csv.NewReader(f)

	for {
		r, err := data.Read()
		if err == io.EOF {
			break
		}

		if err != nil {
			return fmt.Errorf("Something went wrong reading the file %v", err.Error())
		}

		log.Printf("The records %v", r[0])
	}

	return nil

}
