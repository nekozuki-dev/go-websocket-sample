package main

import "go_websocket/external"

func main() {
	router := external.NewRouter()
	router.Run(9080)
}
