package main

import (
	"log"
	"os"
)

func getProjectPath() string {
	dir, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	return dir
}
