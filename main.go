package main

import (
	"crypto/sha256"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"log"
	"os"
)

// Create a struct for storing CSV lines and annotate it with JSON struct field tags
type NamingRecord struct {
	Format            string            `json:"format,omitempty"`
	FileName          string            `json:"file_name,omitempty"`
	Description       string            `json:"description,omitempty"`
	MintingTool       string            `json:"minting_tool,omitempty"`
	Sensitive_content bool              `json:"Sensitive_content,omitempty"`
	Series_number     int64             `json:"series_number,omitempty"`
	Series_total      int64             `json:"series_total,omitempty"`
	Attributes        map[string]string `json:"attributes,omitempty"`
	Collection        Collection        `json:"collection,omitempty"`
	Data              map[string]string `json:"data,omitempty"`
	Hash              string            `json:"hash,omitempty"`
}

// Create a struct for type Collection
type Collection struct {
	Name       string              `json:"name,omitempty"`
	ID         string              `json:"id,omitempty"`
	Attributes []map[string]string `json:"attributes,omitempty"`
}

func main() {
	// open file
	f, err := os.Open("NFT Naming csv - Team Clutch.csv")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	csvReader := csv.NewReader(f)
	data, err := csvReader.ReadAll()
	if err != nil {
		log.Fatal(err)
	}
	//Assign sha256 value
	var list []NamingRecord
	for _, line := range data[2:] {
		attributes := map[string]string{"trait_type": "Gender", "value": line[3]}
		hash := sha256.New()
		format := "CHIP-0007"
		namingList := NamingRecord{
			Format:      format,
			FileName:    line[1],
			Description: line[2],
			Collection:  Collection{ID: line[4]},
			Attributes:  attributes,
		}
		hash.Write([]byte(fmt.Sprint(namingList)))
		hash1 := fmt.Sprintf("%x", hash.Sum(nil))
		namingList.Hash = hash1
		list = append(list, namingList)
	}
	// Create an output file to store hashed data
	csvFile, err := os.Create("filename.output.csv")
	if err != nil {
		log.Fatal(err)
	}
	defer csvFile.Close()
	w := csv.NewWriter(csvFile)
	defer w.Flush()
	//Create Headers
	header := []string{"Filename", "Name", "Description", "Gender", "Attributes", "UUID", "HASH"}
	err = w.Write(header)
	if err != nil {
		log.Fatal(err)
	}
	//Create rows and columns
	for _, r := range list {
		var csvRow []string
		csvRow = append(csvRow, r.FileName, r.Description, r.Attributes["value"], r.Collection.ID, r.Hash)
		err = w.Write(csvRow)
		if err != nil {
			log.Fatal(err)
		}
	}
	//Display JSON Output
	prettyJson, _ := json.MarshalIndent(list, "", "  ")
	fmt.Println(string(prettyJson))
}
