package consumers

import (
	"context"
	"encoding/csv"
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/google/uuid"
)

type JSONOutput struct {
	ID       string                 `json:"id"`
	EntityID string                 `json:"entity_id"`
	Data     map[string]interface{} `json:"data"`
}

type Consumer struct {
	processedFiles map[string]bool
}

func NewConsumer() *Consumer {
	return &Consumer{
		processedFiles: make(map[string]bool),
	}
}

func (c *Consumer) Consume(ctx context.Context, filePath string) ([]JSONOutput, error) {

	fileName := filepath.Base(filePath)
	output := make([]JSONOutput, 0)

	if !c.isProcessed(fileName) {

		//Get file extension so we can choose which method we will use to extract the info from the file
		fileExtension := strings.ToLower(filepath.Ext(filePath))

		switch fileExtension {
		case ".csv":
			log.Printf("Going to process file %s\n", fileName)
			file, err := c.openFile(filePath, fileName)
			if err != nil {
				if errors.Is(err, os.ErrNotExist) { // This isn't error, the file just do not exist in the directory so we wil return to main and process the next file
					return nil, nil
				}
				return nil, err
			}
			output, err = c.readCSV(file)
			if err != nil {
				return nil, err
			}
			defer file.Close()
			log.Printf("Finished to process file %s\n", fileName)
		case ".xml":
			log.Printf("Going to process file %s\n", fileName)
			file, err := c.openFile(filePath, fileName)
			if err != nil {
				if errors.Is(err, os.ErrNotExist) { // This isn't error, the file just do not exist in the directory so we wil return to main and process the next file
					return nil, nil
				}
				return nil, err
			}
			output, err = c.readXML(file)
			if err != nil {
				return nil, err
			}
			defer file.Close()
			log.Printf("Finished to process file %s\n", fileName)
		default: // If not .csv or .xml will always give this error
			log.Printf("Extension of file %s is not supported", fileName)
			return nil, fmt.Errorf("extension of file is not supported")
		}

		// If we got 0 reads from the files the file is empty
		if len(output) == 0 {
			log.Printf("File %s is empty", fileName)
			return nil, fmt.Errorf("file is empty")
		}

		// Save the processeced name file to the array saved in memory
		c.processedFiles[fileName] = true

	} else {
		log.Printf("File %s already processed, will not process again", fileName)
		return nil, nil
	}

	return output, nil
}

// Function responsible for receiving the .csv file and reading it, as well as extracting the necessary data from it
func (c *Consumer) readCSV(file *os.File) ([]JSONOutput, error) {

	reader := csv.NewReader(file)
	reader.Comma = ';'
	records, err := reader.ReadAll() // Read the records from the file with delimiter ;
	if err != nil {
		return nil, err
	}

	output := make([]JSONOutput, 0)

	// Responsible to remove from the records readed from the file
	for _, record := range records[1:] { // Skip the first row since is the name of the column, each record is 1 product
		output = append(output, JSONOutput{
			ID:       uuid.New().String(), // Autogenerated ID, library from google
			EntityID: record[0],
			Data: map[string]interface{}{
				"price": record[1],
				"sku":   record[2],
			},
		})
	}

	return output, nil
}

// Tags xml so we can use unmarshal
type Product struct {
	EntityID string  `xml:"id"`
	Price    float64 `xml:"price"`
	SKU      int64   `xml:"sku"`
}

type Products struct {
	Products []Product `xml:"product"`
}

// Function responsible for receiving the .xml file and reading it, as well as extracting the necessary data from it
func (c *Consumer) readXML(file *os.File) ([]JSONOutput, error) {

	bytes, err := io.ReadAll(file) // Read the data from the file
	if err != nil {
		return nil, err
	}

	var products Products

	err = xml.Unmarshal(bytes, &products) // Exports the data to the struct Products

	if err != nil { // Check if we got any error doing the export of the data
		return nil, err
	}

	output := make([]JSONOutput, 0)

	for _, product := range products.Products {
		output = append(output, JSONOutput{
			ID:       uuid.New().String(), // Autogenerated ID, library from google
			EntityID: product.EntityID,
			Data: map[string]interface{}{
				"price": product.Price,
				"sku":   product.SKU,
			},
		})
	}

	return output, nil
}

// Function responsible to check if we already got that file at our memory storage; return true if exists, false if not;
func (c *Consumer) isProcessed(fileName string) bool {
	_, processed := c.processedFiles[fileName]
	return processed
}

func (c *Consumer) openFile(filePath string, fileName string) (*os.File, error) {
	file, err := os.Open(filePath)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			log.Printf("File %s does not exist in this directory", fileName)
		}
		return nil, err
	}
	return file, nil
}
