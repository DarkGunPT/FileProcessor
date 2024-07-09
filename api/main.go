package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"test/consumers"
)

func main() { // Added some tests to see some output; This should only be used to initialize db, config, consumers, publishers, cron jobs, etc...

	ctx := context.Background()

	consumer := consumers.NewConsumer()

	files, err := os.ReadDir("./files/") // Get every file presented in the folder files
	if err != nil {
		log.Fatalf("Failed to read directory: %s", err)
	}

	for _, file := range files { // Iterate each file
		if !file.IsDir() { // Check if is folder, if so just ignore
			consumeFile(ctx, "./files/"+file.Name(), consumer)
		}
	}

}

func consumeFile(ctx context.Context, filePath string, consumer *consumers.Consumer) { // Send the filePath to the consumer, receive the answer and output it
	var jsonData []byte
	output, err := consumer.Consume(ctx, filePath)

	if err != nil {
		return
	}

	if output != nil {
		jsonData, err = convertJson(output)
		log.Println(string(jsonData))
	}
}

func convertJson(output []consumers.JSONOutput) ([]byte, error) { // Convert the byte data received from the consumer to JSON
	jsonData, err := json.MarshalIndent(output, "", "  ")
	if err != nil {
		fmt.Printf("Error marshaling JSON: %v\n", err)
		return nil, err
	}
	return jsonData, nil
}
