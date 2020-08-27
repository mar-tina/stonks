package main

import (
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/urfave/cli/v2"
)

func main() {

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	PopulateResponses()

	app := cli.App{
		Name:  "stonks",
		Usage: "Find out whether we support your stonks",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "path",
				Aliases: []string{"p"},
				Usage:   "Filepath for the stocks input",
			},
			&cli.StringFlag{
				Name:    "online",
				Aliases: []string{"o"},
				Usage:   "Hosted Filepath for the stocks input",
			},
			&cli.StringFlag{
				Name:    "update",
				Aliases: []string{"u"},
				Usage:   "Update mode that allows editing CSV File",
			},
			&cli.StringFlag{
				Name:    "conversion",
				Aliases: []string{"c"},
				Usage:   "Conversion mode that allows users to query the API and convert between currencies",
			},
			&cli.StringFlag{
				Name:    "Language File Path",
				Aliases: []string{"lp"},
				Usage:   "File path for the languages input",
			},
			&cli.StringFlag{
				Name:    "Default Language",
				Aliases: []string{"dl"},
				Usage:   "Sets the default language for the running instance",
			},

			&cli.StringFlag{
				Name:    "Get Current Price",
				Aliases: []string{"gp"},
				Usage:   "Get the current stock price",
			},
			&cli.StringFlag{
				Name:    "List All Stocks",
				Aliases: []string{"v"},
				Usage:   "Lists all the available stocks",
			},
		},
		Action: func(c *cli.Context) error {
			if c.String("u") == "on" {
				return UpdateCSV(c)
			} else if c.String("c") == "on" {
				return ConversionMode(c)
			} else if c.String("gp") == "on" {
				return CheckPrice(c)
			} else if c.String("v") == "on" {
				return ListAllStocks(c)
			}
			return PromptDisplay(c)
		},
	}

	err = app.Run(os.Args)
	if err != nil {
		log.Printf("Oops something went wrong . %v", err.Error())
	}

}
