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
		},
		Action: func(c *cli.Context) error {
			if c.String("u") == "on" {
				return UpdateCSV(c)
			}
			return PromptDisplay(c)
		},
	}

	err = app.Run(os.Args)
	if err != nil {
		log.Printf("Oops something went wrong . %v", err.Error())
	}

}
