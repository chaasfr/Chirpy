package main

import (
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"//this tells go to import the driver for side effects

) 
func main() {
	godotenv.Load()
	InitServer()
}