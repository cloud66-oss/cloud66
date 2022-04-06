package main

import (
	"fmt"

	"github.com/cloud66-oss/cloud66"
)

func main() {
	// Create a new cloud66 Go client, this will read the
	// cx.json token file from the root of the project.
	client := cloud66.GetClient("cx.json", "", "", cloud66.NewClientConfig("https://app.cloud66.com"))

	// Use client to send an unauthorized ping.
	if err := client.UnauthenticatedPing(); err != nil {
		fmt.Println("Unauthorized ping failed:", err)
		return
	}

	fmt.Println("Unauthorized ping succeeded")

	// Use client to send an authorized ping.
	if err := client.AuthenticatedPing(); err != nil {
		fmt.Println("Authorized ping failed:", err)
		return
	}

	fmt.Println("Authorized ping succeeded")
}
