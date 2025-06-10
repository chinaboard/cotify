package main

import (
	"fmt"
	"log"

	"github.com/chinaboard/cotify/sdk/rpc"
)

func main() {
	c := rpc.New("http://localhost:8080")

	req := rpc.CotifyItemRequest{
		Url:      "https://example.com/video",
		Title:    "Example Video",
		Type:     "video",
		Metadata: "HD",
	}

	resp, err := c.StoreItem(req)
	if err != nil {
		log.Fatalf("Failed to store item: %v", err)
	}

	fmt.Printf("Is new item: %v\n", resp.IsNew)
	fmt.Printf("Item details: %+v\n", resp.Item)
}
