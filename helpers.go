package main

import (
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/manifoldco/promptui"
	"github.com/mar-tina/stonks/analyzer"
	"github.com/olekukonko/tablewriter"
	"github.com/urfave/cli/v2"
)

type Stonk struct {
	Country  string
	Currency string
	Code     string
}

var Responses map[string][]string

func PopulateResponses() {
	Responses = make(map[string][]string)
	Responses["en"] = []string{"Currency ", "From ", "To ", "Amount "}
	Responses["fr"] = []string{"Devise ", "De ", "À ", "Prix "}
	Responses["sw"] = []string{"Fedha ", "Kutoka ", "Kuenda ", "Bei "}
	Responses["pt"] = []string{"Moeda ", "De ", "Para ", "Preço "}
}

type Language struct {
	Code string
	Lang string
}

type StonkBank struct {
	Stonks      map[string]Stonk
	Languages   map[string]Language
	DefaultLang Language
}

func (a *StonkBank) Insert(key string, value Stonk) {
	a.Stonks[key] = value
}

func (a *StonkBank) InsertLanguages(key string, value Language) {
	a.Languages[key] = value
}

// RenderStonk,  Renders the selected stock to the user.
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

func RenderCheckPrice(s Stonk, amount float64, base string) {
	data := [][]string{
		[]string{base, s.Country, s.Currency, s.Code, fmt.Sprintf("%f", amount), time.Now().Format("Mon Jan 2 15:04:05 -0700 MST 2006")},
	}

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Base", "Country", "Currency", "ISO 4217 Code", "Price", "Time"})

	for _, v := range data {
		table.Append(v)
	}
	table.Render()
}

func RenderAllStonks(c *cli.Context) error {
	sb, err := ReadCSVFile(c.String("p"))
	if err != nil {
		return err
	}
	for _, value := range sb.Stonks {
		RenderStonk(value)
	}

	return nil
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

func PopulateStocksForExistingStockBank(sb *StonkBank, path string) (*StonkBank, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("Something went wrong reading the file %v", err.Error())
	}

	data := csv.NewReader(f)
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

func ReadLanguagesFile(path string) (*StonkBank, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("Something went wrong reading the file %v", err.Error())
	}

	data := csv.NewReader(f)
	sb := &StonkBank{}
	counter := 0
	sb.Languages = make(map[string]Language)

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
			lang := Language{
				Lang: r[0],
				Code: r[1],
			}
			sb.InsertLanguages(lang.Code, lang)
		}

		counter++
	}
	return sb, nil
}

func CheckPrice(c *cli.Context) error {

	sb, err := ReadCSVFile(c.String("p"))
	if err != nil {
		sb, err = ReadHostedFile(c.String("o"))
		if err != nil {
			return err
		}
	}

	stonk := Stonk{}
	lyzer := analyzer.Init()

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

		currency, err := prompt.Run()

		if err != nil {
			fmt.Printf("Something went wrong %v\n", err)
			return err
		}

		res, err := lyzer.GetCurrentPrice([]string{currency})
		if err != nil {
			fmt.Printf("Something went wrong %v\n", err)
			return err
		}

		var base string
		var rate float64
		// CurrencyAnalyzer uses a different parameter for setting the base currency and rates parameter.
		if res["base"] == "" {
			base = res["source"].(string)
			rate = res["quotes"].(map[string]interface{})[stonk.Code].(float64)
		} else {
			base = res["base"].(string)
			rate = res["rates"].(map[string]interface{})[stonk.Code].(float64)
		}

		RenderCheckPrice(stonk, rate, base)
	}
}

func ListAllStocks(c *cli.Context) error {
	if c.String("p") == "" {
		log.Printf("Please provide a path for the stocks input")
		os.Exit(1)
	}

	RenderAllStonks(c)
	return nil
}

func ConversionMode(c *cli.Context) error {
	sb, err := ReadLanguagesFile(c.String("lp"))
	if err != nil {
		return errors.New("Something went wrong. Please try again later")
	}

	if c.String("p") == "" {
		log.Printf("Please provide a path for the stocks input")
		os.Exit(1)
	}

	if c.String("dl") == "" {
		prompt := promptui.Select{
			Label: "Select Language",
			Items: []string{"en", "fr", "pt", "sw"},
		}

		_, result, err := prompt.Run()

		if err != nil {
			return fmt.Errorf("Something went wrong. %v", err.Error())
		}

		sb.DefaultLang = sb.Languages[result]
		fmt.Printf("You choose %q\n", result)
	}

	PopulateStocksForExistingStockBank(sb, c.String("p"))
	lyzer := analyzer.Init()
	lang := Responses[sb.DefaultLang.Code]
	for {
		validate := func(input string) error {
			var ok bool
			_, ok = sb.Stonks[input]
			if !ok {
				return errors.New("We currently do not support the specified currency")
			}
			return nil
		}

		amountValidate := func(input string) error {
			_, err := strconv.ParseFloat(input, 64)

			if err != nil {
				return fmt.Errorf("Make sure input is a valid number %v", err.Error())
			}

			return nil
		}

		fromPrompt := promptui.Prompt{
			Label:    "" + lang[1] + "[" + lang[0] + "]",
			Validate: validate,
		}

		from, err := fromPrompt.Run()

		toPrompt := promptui.Prompt{
			Label:    "" + lang[2] + "[" + lang[0] + "]",
			Validate: validate,
		}

		to, err := toPrompt.Run()

		amountPrompt := promptui.Prompt{
			Label:    "" + lang[3] + "[" + lang[0] + "]",
			Validate: amountValidate,
		}

		var amount string
		amount, err = amountPrompt.Run()

		if err != nil {
			fmt.Printf("Something went wrong %v\n", err)
			return err
		}

		parsedAmount, err := strconv.ParseFloat(amount, 64)
		response, err := lyzer.Convert(to, from, parsedAmount)
		if err != nil {
			return fmt.Errorf("Failed to complete conversion %v ", err.Error())
		}

		log.Printf(" %v : %v -> %v : %v", amount, from, to, response["result"])
	}
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
	// a currency, validation is run and the user is shown all the information about a currency if the
	// currency is supported
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
			return errors.New("Please provide valid input")
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
