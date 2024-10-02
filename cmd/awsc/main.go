package main

import (
	"os"

	"github.com/handlename/awsc/cli"
)

func main() {
	os.Exit(int(cli.Run()))
}
