package main

import (
	"context"
	"log"
	"reporting/server"
)

func main() {
	if err := server.NewHTTPServer().Run(context.TODO()); err != nil {
		log.Fatalf("failed starting server, err : \n%+v", err)
	}
}
