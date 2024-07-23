
# Json Parser

Json Parser from scratch in go

## Test Locally

### Clone the project
`git clone https://github.com/KhushPatibandha/jsonParser.git`

### Navigate to the project directory
`cd .\jsonParser\`

### Test
`go run cmd/main.go`

And paste your json blob.

## Use as a package

### Get the latest package
`go get -u github.com/KhushPatibandha/jsonParser`

### Usage
`
    result, err := jsonparser.ParseIt(<your json string>)
`

The method `jsonparser.ParseIt(jsonString)` takes in a `string` and returns `interface{}, error`
