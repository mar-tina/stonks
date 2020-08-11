package main

import (
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/manifoldco/promptui"
	"github.com/olekukonko/tablewriter"
	"github.com/urfave/cli/v2"
)

type Stonk struct {
	Country  string
	Currency string
	Code     string
}

type StonkBank struct {
	Stonks map[string]Stonk
}

func (a *StonkBank) Insert(key string, value Stonk) {
	a.Stonks[key] = value
}

func RenderStonk(s Stonk) {
	data := [][]string{
		[]string{s.Country, s.Currency, s.Code},
	}

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Country", "Currency", "ISO 4217 Code"})

	for _, v := range data {
		table.Append(v)
	}
	table.Render()
}

func ReadCSVFile(path string) (*StonkBank, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("Something went wrong reading the file %v", err.Error())
	}

	data := csv.NewReader(f)
	sb := &StonkBank{}
	counter := 0
	sb.Stonks = make(map[string]Stonk)

	for {
		r, err := data.Read()
		if err == io.EOF {
			break
		}

		if err != nil {
			return nil, fmt.Errorf("Something went wrong reading the file %v", err.Error())
		}

		// Skip the first CSV column
		if counter >= 1 {
			st := Stonk{
				Country:  r[0],
				Currency: r[1],
				Code:     r[2],
			}
			sb.Insert(st.Code, st)
		}

		counter++
	}

	return sb, nil

}

func PromptDisplay(c *cli.Context) error {
	sb, err := ReadCSVFile(c.String("p"))
	if err != nil {
		sb, err = ReadHostedFile(c.String("o"))
		if err != nil {
			return err
		}
	}

	stonk := Stonk{}

	// Modeling the prompt after the structure of a simple guessing game. User inputs
	// a currency, validation is run and the user is shown all the information about a currencyF
	for {
		validate := func(input string) error {
			var ok bool
			stonk, ok = sb.Stonks[input]
			if !ok {
				return errors.New("We currently do not support the specified currency")
			}
			return nil
		}

		prompt := promptui.Prompt{
			Label:    "Currency",
			Validate: validate,
		}

		_, err := prompt.Run()

		if err != nil {
			fmt.Printf("Something went wrong %v\n", err)
			return err
		}

		RenderStonk(stonk)
	}
}

func ReadHostedFile(url string) (*StonkBank, error) {
	if err := DownloadFile(url); err != nil {
		return nil, err
	}

	return ReadCSVFile("stonks.csv")
}

func UpdateCSV(c *cli.Context) error {

	validate := func(input string) error {
		if input == "" {
			return errors.New("We currently do not support the specified currency")
		}
		return nil
	}

	prompt := promptui.Prompt{
		Label:    "Country",
		Validate: validate,
	}

	countryRes, err := prompt.Run()

	currprompt := promptui.Prompt{
		Label:    "Currency",
		Validate: validate,
	}

	currRes, err := currprompt.Run()

	codeprompt := promptui.Prompt{
		Label:    "ISO Code ",
		Validate: validate,
	}

	codeRes, err := codeprompt.Run()

	var data = []string{countryRes, currRes, codeRes}
	err = WriteToFile(data)

	if err != nil {
		fmt.Printf("Something went wrong %v\n", err)
		return err
	}

	return nil

}

func WriteToFile(data []string) error {
	f, err := os.OpenFile("stonks.csv", os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		return fmt.Errorf("Could not open file %v", err.Error())
	}

	defer f.Close()

	writer := csv.NewWriter(f)
	defer writer.Flush()

	f.Write([]byte("\n"))
	err = writer.Write(data)
	if err != nil {
		return fmt.Errorf("Cannot write to file %v", err.Error())
	}

	return nil
}


func DownloadFile(url string) error {
	resp, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("Something went wrong fetching the file %v", err.Error())
	}

	defer resp.Body.Close()

	out, err := os.Create("stonks.csv")
	if err != nil {
		return fmt.Errorf("Failed to save file %v", err.Error())
	}

	_, err = io.Copy(out, resp.Body)
	return err

}
