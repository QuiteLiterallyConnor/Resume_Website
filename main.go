package main

import (
	"flag"
	"fmt"
	"resume_website/app"
)

func main() {
	port := flag.String("port", "8080", "Port on which the server will run")
	flag.Parse()

	fmt.Printf("Starting server on port %s\n", *port)

	app.StartServer(*port)
}
