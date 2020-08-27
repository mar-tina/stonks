# stonks

Take Home Challenge Interview.

## Languages Used:

- Go

Stonks is a CLI based application that allow stock data to be displayed according to the currency . It has a 3
step CSV update process to update the stock data .

### How To Run The App

- Clone the Repository

  `git clone https://github.com/mar-tina/stonks.git`

  `cd stonks`

  `go run main.go helpers.go -p stonks.csv`

  `or`

  `go build.`

  `./stonks -p stonks.csv`

### Environment Variables

    FIXER_KEY=

    CL_KEY=

### Running tests

To be able to run tests you need to have you API KEYS copied into the analyzer/analyzer.go file
respective variables.

![Keys ](https://github.com/mar-tina/stonks/blob/master/keys.png)

    Run : go test -v

### Modes

Display Mode :

    This is the default mode.  It displays data according to passed in user input.

![Default Mode](https://github.com/mar-tina/stonks/blob/master/defaultmode.png)

Update Mode :

    This mode allows you to update the CSV File.
    ./stonks -u on

The on parameter is necessary to switch to update mode.

![Update Mode](https://github.com/mar-tina/stonks/blob/master/updatemode.png)

GetCurrentPrice Mode :

    This mode allows you to view the current going price for a stock in relation to
    the base currently set as 'EUR'

    ./stonks -gp on -p stonks.csv

The on parameter is necessary to switch to update mode.

![GetCurrenPriceMode](https://github.com/mar-tina/stonks/blob/master/gp.png)

Conversion Mode :

    This mode allows you to convert from one currency to another.

    ./stonks -c on -lp lang.csv -p stonks.csv

There is a language picker that allows you to choose the language you would like to use:

![LanguagePrompt](https://github.com/mar-tina/stonks/blob/master/langpicker.png)

After picking the language the terminal switches to the default prompt style with the prompt labels in
the specified language.

![ConversionMode](https://github.com/mar-tina/stonks/blob/master/convert.png)

List All Mode :

    This mode allows you to list all the available stocks

    ./stonks -v on -p stonks.csv

![ListAllMode](https://github.com/mar-tina/stonks/blob/master/all.png)

### Flags

    -p "path" . Specifies the file path where the stock data is stored.

    -o "online hosted file" . Specifies a remote hosted file location.

    -u "update" . Drops the user into the update mode .
