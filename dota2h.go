package main

import (
    "flag"
    "math/rand"
    "time"

    "github.com/robotxet/dota2h/server"
)

func main() {
    configPath := flag.String("config", "./conf/config.json", "Server config")
    flag.Parse()
    config := server.ParseConfig(*configPath)

    rand.Seed(time.Now().UTC().UnixNano())

    s := server.New(config)

    s.Run()
}