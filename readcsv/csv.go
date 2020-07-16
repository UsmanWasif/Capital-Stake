package csv

import (
	"encoding/csv"
	"encoding/json"
	"io"
	"log"
	"os"
	"strings"
)

// Covid ...
type Covid struct {
	CumulativeTestPositive  string `json:"cumulative_test_positive"`
	CumulativeTestPerformed string `json:"cumulative_test_performed"`
	Region                  string `json:"region"`
	Date                    string `json:"Date"`
	Discharged              string `json:"discharged"`
	Expired                 string `json:"expired"`
	Admitted                string `json:"admitted"`
}

// Loadcsvfile ...
func Loadcsvfile(str string) []Covid {
	// Open the file
	var table = make([]Covid, 0) // creating slices of dynamic size array just like list
	csvfile, err := os.Open(str)
	if err != nil {
		log.Fatalln("Couldn't open the csv file", err)
	}

	// Parse the file
	file := csv.NewReader(csvfile)

	// Iterate through the records
	for {
		// Read each record from csv
		record, err := file.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err)
		}
		var c = Covid{
			CumulativeTestPositive:  record[0],
			CumulativeTestPerformed: record[1],
			Date:                    record[2],
			Discharged:              record[3],
			Expired:                 record[4],
			Region:                  record[5],
			Admitted:                record[6],
		}
		table = append(table, c)
	}
	return table
}

// Search ...
func Search(table []Covid, find string) []byte {
	find1 := strings.Title(find)
	find2 := strings.ToUpper(find)
	var response []byte
	var searchedData []Covid
	for _, lim := range table { // one of the syntax of for loop that works for each elements/slices
		if lim.Region == find1 || lim.Date == find || lim.Region == find2 {
			var c = []Covid{
				{CumulativeTestPositive: lim.CumulativeTestPositive,
					CumulativeTestPerformed: lim.CumulativeTestPerformed,
					Date:                    lim.Date,
					Discharged:              lim.Discharged,
					Expired:                 lim.Expired,
					Region:                  lim.Region,
					Admitted:                lim.Admitted},
			}
			searchedData = append(searchedData, c...)

		}
	}
	response, err := json.MarshalIndent(searchedData, "", "")
	if err != nil {
		log.Println(err)
	}
	return response
}
