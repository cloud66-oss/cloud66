package main

import (
	"fmt"

	"github.com/cloud66-oss/cloud66"
)

func main() {
	// Create a new cloud66 Go client, this will read the
	// cx.json token file from the root of the project.
	client := cloud66.GetClient("cx.json", "", "", cloud66.NewClientConfig("https://app.cloud66.com"))

	// Use client to get a list of clouds.
	// Error omited for brevity.
	clouds, _ := client.GetCloudsInfo()

	// Print the list of clouds.
	for _, cloud := range clouds {
		fmt.Println(cloud.Name)

		// Print the list of regions.
		for _, region := range cloud.Regions {
			fmt.Println("\t", region.Name)
		}
	}
}
