# File Processor

## The Goal

The aim of this service is to receive csv and xml files and process them, removing the data from them and transforming it into JSON format.

### CSV Example

```
id;price;sku;
123;22.0;3123
asdv;23.0;7651
```

### OUTPUT Expected

```json
[{
    "id": "some-random-id",
    "entity_id": "123",
    "data": {
        "price": "22.0",
        "sku": "3123"
    }
}, {
    "id": "some-random-id",
    "entity_id": "asdv",
    "data": {
        "price": "23.0",
        "sku": "7651"
    }
}]
```

### XML Example

```xml
<products>
    <product>
        <id>123</id>
        <price>22.0</price>
        <sku>3123</sku>
        <name>some name</name>
    <product>
</products>
```

### OUTPUT Expected

```json
[{
    "id": "some-random-id",
    "entity_id": "123",
    "data": {
        "price": "22.0",
        "sku": "3123"
    }
}]
```
## The project structure

The project is made up of 4 folders: 
* ./api, which contains our main 
* ./consumers, which contains our consumer 
* ./files, which contains the files to be processed
* ./tests, which contains the tests and files for those tests

## The Main

The main function is the initial function of our project and is presented in ./api/main.go It is the function that the compiler runs first and is responsible for calling our file processor.

When main receives the file processed by the consumer, in []bytes format, converts it to JSON and displays it.

### Important

The files need to be in the ./files/ folder so that main can fetch them.

```go
files, err := os.ReadDir("./files/")
```

### How to run

In the terminal navigate to the project folder, then run the main.go file

```
cd ...\challenges\file-processor
go run api/main.go
```

## The File-Processor

A consumer has been created that receives the path to the file in the system. Is presented in ./consumers/fileConsumer.go

The consumer has an array of strings with the name and extension of the files that have already been processed, validating whether or not they have been processed.

```go
type Consumer struct {
	processedFiles []string // Array of string, example: {"xml.xml","csv.csv"}, so we can check if the file was already processed or not
}
```

If they haven't been processed, it processes the documents, performing 2 different functions depending on whether the file extension is .csv or .xml. 

```go
func (c *Consumer) readCSV(file *os.File) ([]JSONOutput, error)
func (c *Consumer) readXML(file *os.File) ([]JSONOutput, error)
```

If the file has already been processed, a log is sent informing you that the file has already been processed.

```
File name.extension already processed, will not process again
```

If the file has an extension that is neither of the two allowed, consumer returns null and outputs that the file extension is not allowed.

```
Extension .extension of file name.extension is not supported
```

Finally, when the processing is finished and there has been no error, the output is returned to main().

## Tests

To test our consumer, we created a tests folder and imported the testing package. This package allows you to run automatic tests and return the results.

### Important

For these tests to work, the files csv.csv, test.csv, test.html and xml.xml must be present in the ./tests/ folder.

* func TestDuplicateFiles(t *testing.T)

Responsible for testing the sending of duplicated files

* func TestUnavailableDirectory(t *testing.T)

Responsible for testing the sending of files with an invalid directory

* func TestUnsupportedFile(t *testing.T)

Responsible for testing the sending of files with an unsupported extension

* func TestEmptyDirectory(t *testing.T)

Responsible for testing what happen is the directory received is empty

* func TestEmptyFile(t *testing.T)

Responsible for testing if the file to process is empty, without data

## Output

```
PS ...\challenges\file-processor> go run api/main.go
2024/07/08 14:31:14 Going to process file csv.csv
2024/07/08 14:31:14 Finished to process file csv.csv
2024/07/08 14:31:14 [
  {
    "id": "53418832-6ded-491d-899a-2eb04502547c",
    "entity_id": "123",
    "data": {
      "price": "22.0",
      "sku": "3123"
    }
  },
  {
    "id": "556ab05f-4764-472d-98d4-950bcf63ec76",
    "entity_id": "asdv",
    "data": {
      "price": "23.0",
      "sku": "7651"
    }
  }
]
2024/07/08 14:31:14 Going to process file test.csv
2024/07/08 14:31:14 Finished to process file test.csv
2024/07/08 14:31:14 File test.csv is empty
2024/07/08 14:31:14 Going to process file xml.xml
2024/07/08 14:31:14 Finished to process file xml.xml
2024/07/08 14:31:14 [
  {
    "id": "4c2acd62-0363-43f6-9039-b141b9ccb1cc",
    "entity_id": "123",
    "data": {
      "price": 22,
      "sku": 3123
    }
  },
  {
    "id": "9d533319-2378-4336-8d3e-ac2dad03bf73",
    "entity_id": "asdv",
    "data": {
      "price": 23,
      "sku": 7651
    }
  }
]
```