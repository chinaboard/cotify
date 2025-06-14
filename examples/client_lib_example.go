package main

import (
	"fmt"
	"log"

	"github.com/chinaboard/cotify/sdk/cotify"
)

func main() {
	c := cotify.NewClient("http://localhost:8080", nil)

	req := &cotify.StoreRequest{
		Url:      "https://example.com/video",
		Title:    "Example Video",
		Type:     "video",
		Metadata: "HD",
	}

	resp, err := c.Store(req)
	if err != nil {
		log.Fatalf("Failed to store item: %v", err)
	}

	fmt.Printf("Is new item: %v\n", resp.IsNew)
	fmt.Printf("Item details: %+v\n", resp.Item)
}
