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

### Modes

Display Mode :

    This is the default mode.  It displays data according to the passed in user input prompt .

![Default Mode](https://github.com/mar-tina/stonks/blob/master/defaultmode.png)

Update Mode :

    This mode allows you to update the CSV File.
    ./stonks -u on

The on parameter is necessary to switch to update mode.

![Update Mode](https://github.com/mar-tina/stonks/blob/master/updatemode.png)

### Flags

    -p "path" . Specifies the file path where the stock data is stored.

    -o "online hosted file" . Specifies a remote hosted file location.

    -u "update" . Drops the user into the update mode .
