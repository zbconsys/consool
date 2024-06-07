package main

import (
	"github.com/zbconsys/consool/internal/app"
	"log"
)

func main() {
	a, err := app.New()
	if err != nil {
		log.Fatalln(err.Error())
	}

	if err = a.Run(); err != nil {
		log.Fatalln(err.Error())
	}
}
