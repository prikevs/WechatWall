package main

import (
	"crawler/ucrawler"

	"fmt"
)

func main() {
	usersch := make(chan []ucrawler.User)
	go ucrawler.Run(usersch)

	for us := range usersch {
		fmt.Println(us)
	}
}
