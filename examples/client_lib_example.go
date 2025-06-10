package main

import (
	"fmt"
	"log"

	"github.com/chinaboard/cotify/sdk/httpclient"
)

func main() {
	c := httpclient.New("http://localhost:8080")

	req := httpclient.CotifyItemRequest{
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
