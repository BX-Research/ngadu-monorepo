// Command dburl supplies the local connection URL to development processes.
package main

import (
	"fmt"
	"github.com/BX-Research/ngadu-monorepo/backend/app/database"
	"os"
)

func main() {
	url, err := database.LocalURL()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	fmt.Print(url)
}
