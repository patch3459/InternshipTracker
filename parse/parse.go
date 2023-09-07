/*
	Contains utilities for parsing job listings and updating the
	spread sheet data base.
*/

package parse

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"
)

/*
ScrapeNewInternships() -> bool

Will go through the JobLinks.csv database, locate job listings that have not yet
been located and then concurrently request and scrape them.

Once all the new listings have been found, it will write it into the newJobs database
(non concurrently as I don't think csv supports concurrent writing)

returns true if successful, false if not and an error object?
*/
func ScrapeNewInternships() bool {

	file, err := os.Open("../JobLinks.csv")
	if err != nil {
		log.Fatal("Unable to find job links csv file")
	}
	defer file.Close()

	csvReader := csv.NewReader(file)
	// ToDo: Check the space efficiency of using ReadAll()
	records, err := csvReader.ReadAll()
	if err != nil {
		log.Fatal("Error: Unable to read Job Links CSV File")
	}

	for _, entry := range records {
		fmt.Println(entry)
	}

	return true
}
