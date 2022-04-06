package main

import (
	"fmt"

	"github.com/cloud66-oss/cloud66"
)

func main() {
	// Create a new cloud66 Go client, this will read the
	// cx.json token file from the root of the project.
	client := cloud66.GetClient("cx.json", "", "", cloud66.NewClientConfig("https://app.cloud66.com"))

	// Use client to get a list of stacks.
	// Error omited for brevity.
	stacks, _ := client.StackList()

	// Print the list of stacks.
	for _, stack := range stacks {
		fmt.Println(stack.Name)
	}
}
