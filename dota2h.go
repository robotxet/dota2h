package main

import (
    "math/rand"
    "time"

    "dota2_server/server"
)

func main() {
    config := server.ParseConfig("./conf/config.json")

    rand.Seed(time.Now().UTC().UnixNano())

    s := server.New(config)

    s.Run()
}