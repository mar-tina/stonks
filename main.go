package main

import (
	"errors"
	"fmt"
	"log"
	"os"

	"github.com/manifoldco/promptui"
	"github.com/urfave/cli/v2"
)

func main() {

	app := cli.App{
		Name:  "stonks",
		Usage: ". ",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "path",
				Aliases: []string{"p"},
				Usage:   "Look up the Stonk you need for the currency you desire",
			},
		},
		Action: func(c *cli.Context) error {
			log.Printf("The file path %v", c.String("p"))
			ReadCSVFile(c.String("p"))
			for {
				validate := func(input string) error {
					if len(input) < 3 {
						return errors.New("Please input a valid currency")
					}
					return nil
				}

				prompt := promptui.Prompt{
					Label:    "Currency",
					Validate: validate,
				}

				result, err := prompt.Run()

				if err != nil {
					fmt.Printf("Prompt failed %v\n", err)
					return err
				}

				fmt.Printf("You choose %q\n", result)

			}
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal("")
	}

	// Modeling the prompt after the structure of a simple guessing game. User inputs
	// a currency, validation is run and the user is shown all the information about a currencyF

}
