package main

import (
	"context"
	"log"

	"github.com/BX-Research/ngadu-monorepo/backend/app/server"
)

func main() {
	if err := server.Run(context.Background()); err != nil {
		log.Fatal(err)
	}
}
