package main

import "go-websocket-sample/external"

func main() {
	router := external.NewRouter()
	router.Run(9080)
}
