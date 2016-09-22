package main

import (
    "github.com/joho/godotenv"
    "log"
)

func LoadEnvironment() {
  err := godotenv.Load()
  if err != nil {
    log.Println("Did not load environment variables from .env file")
  }
}
