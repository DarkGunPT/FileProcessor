package tests

import (
	"context"
	"test/consumers"
	"testing"
)

// FILES NEED TO BE IN FOLDER ./file-processor/tests/

// Test sending duplicated files to the consumer
func TestDuplicateFiles(t *testing.T) {

	ctx := context.Background()

	fileProcessorConsumer := consumers.NewConsumer()

	// Test 1

	_, err := fileProcessorConsumer.Consume(ctx, "./xml.xml")

	if err != nil {
		t.Errorf("Error: %s", err.Error())
	}

	// Test 2 - repeated xml file

	_, err = fileProcessorConsumer.Consume(ctx, "./xml.xml")

	if err != nil {
		t.Errorf("Error: %s", err.Error())
	}

	// Test 3

	_, err = fileProcessorConsumer.Consume(ctx, "./csv.csv")

	if err != nil {
		t.Errorf("Error: %s", err.Error())
	}

	// Test 4 - repeated file

	_, err = fileProcessorConsumer.Consume(ctx, "./csv.csv")

	if err != nil {
		t.Errorf("Error: %s", err.Error())
	}

}

// Test sending unavailable directory to the consumer
func TestUnavailableDirectory(t *testing.T) {
	ctx := context.Background()

	fileProcessorConsumer := consumers.NewConsumer()

	// Test 5 - file not existent

	_, err := fileProcessorConsumer.Consume(ctx, "./re.csv")

	if err != nil {
		t.Errorf("Error: %s", err.Error())
	}
}

// Test sending unsupported extension file to the consumer
func TestUnsupportedFile(t *testing.T) {
	ctx := context.Background()

	fileProcessorConsumer := consumers.NewConsumer()

	// Test 6 - Unsupported File

	_, err := fileProcessorConsumer.Consume(ctx, "./test.html")

	if err != nil {
		t.Errorf("Error: %s", err.Error())
	}
}

// Test sending empty directory to the consumer
func TestEmptyDirectory(t *testing.T) {
	ctx := context.Background()

	fileProcessorConsumer := consumers.NewConsumer()

	// Test 7 - Empty directory

	_, err := fileProcessorConsumer.Consume(ctx, "")

	if err != nil {
		t.Errorf("Error: %s", err.Error())
	}
}

// Test empty file
func TestEmptyFile(t *testing.T) {
	ctx := context.Background()

	fileProcessorConsumer := consumers.NewConsumer()

	// Test 7 - Empty directory

	_, err := fileProcessorConsumer.Consume(ctx, "test.csv")

	if err != nil {
		t.Errorf("Error: %s", err.Error())
	}
}
