package main

import (
    "math/rand"
    "time"

    "github.com/robotxet/dota2h/server"
)

func main() {
    config := server.ParseConfig("./conf/config.json")

    rand.Seed(time.Now().UTC().UnixNano())

    s := server.New(config)

    s.Run()
}