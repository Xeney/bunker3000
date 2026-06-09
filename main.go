package main

import (
	"fmt"
	"main/player"
)

func main() {
	player := player.CreatePlayer()
	for {
		err := player.StartNewDay()
		if err != nil {
			fmt.Println(err)
			return
		}
	}
}
