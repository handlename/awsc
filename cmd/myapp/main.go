package main

import (
	"os"

	"github.com/handlename/my-golang-template/cli"
)

func main() {
	os.Exit(int(cli.Run()))
}
