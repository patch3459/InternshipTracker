/*
	Contains utilities for parsing job listings and updating the
	spread sheet data base.
*/

package parse

import (
	"encoding/csv"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"sync"
)

/*
readConfigFile

Reads the config.json file and will return it as a
config struct.

Will not return an error object. Will abort the program if the
file cannot be processed correctly.
*/
func readConfigFile() Config {
	file, err := os.Open("./config.json")
	if err != nil {
		log.Fatal("error: Could not load config.json file")
	}
	defer file.Close()

	bytes, err := ioutil.ReadAll(file)
	if err != nil {
		log.Fatal("error: Had trouble reading config.json as bytes")
	}

	var config Config
	err = json.Unmarshal(bytes, &config)
	if err != nil {
		log.Fatal("error: Could not unmarshall config.json")
	}

	return config
}

/*
requestGreenHouseJobs

makes a request for the jobs from a particular greenhouse link and returns them as an object of
GreenHouseResponse Type
*/
func requestGreenHouseJobs(url string) (GreenHouseResponse, error) {
	var resp *http.Response

	// via the url keyword
	reqUrl := fmt.Sprintf("https://boards-api.greenhouse.io/v1/boards/%s/jobs", url)
	// making the request
	resp, err := http.Get(reqUrl)
	if err != nil {
		return GreenHouseResponse{}, errors.New("error making a request to " + reqUrl)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return GreenHouseResponse{}, errors.New("error reading json response" + err.Error())
	}

	// parsing the JSON response
	var listings GreenHouseResponse
	if err := json.Unmarshal(body, &listings); err != nil {
		return GreenHouseResponse{}, errors.New("error parsing Json Response " + err.Error())
	}

	for _, jobListing := range listings.Jobs {
		// To Do : Add keywords
		if strings.Contains(jobListing.Title, "Intern") {
			fmt.Println(jobListing.Title)
		}
	}

	return listings, nil
}

/*
requestWorkDayJobs

Makes request to myworkdayjobs page and returns the jobs as a
workdayresponse object
*/
func requestWorkDayJobs(url string) (WorkDayResponse, error) {
	var resp *http.Response

	url = url + "/jobs"
	client := &http.Client{}
	req, err := http.NewRequest("POST", url, nil)
	if err != nil {
		return WorkDayResponse{}, errors.New("error with making request")
	}

	req.Header.Set("Accept", "application/json")
	req.Header.Set("Accept_Language", "en-US")
	req.Header.Set("Content-Type", "application/json")

	resp, err = client.Do(req)
	if err != nil {
		return WorkDayResponse{}, errors.New("error making a request to " + url)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return WorkDayResponse{}, errors.New("error reading json response" + err.Error())
	}
	var listings WorkDayResponse
	if err := json.Unmarshal(body, &listings); err != nil {
		return WorkDayResponse{}, errors.New("error parsing Json Response " + err.Error())
	}

	for _, jobListing := range listings.JobPostings {
		// To Do : Add keywords
		if strings.Contains(jobListing.Title, "Intern") {
			fmt.Println(jobListing.Title)
		}
	}

	return listings, nil
}

/*
writeInternshipToFile()

Writes a job listing into a csv found at found_internship_csv_path in the config.json
*/
func writeInternshipToFile(job *JobListing, path string) (bool, error) {

	file, err := os.OpenFile("/Users/pat/Desktop/InternshipTracker/JobList.csv", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return false, errors.New("error: Cannot open job list csv")
	}
	defer file.Close()
	csvWriter := csv.NewWriter(file)
	var line []string = []string{
		strconv.Itoa((*job).ID),
		(*job).Title,
		(*job).Company,
		(*job).DatePosted,
		(*job).Link,
		(*job).DateUploaded,
	}
	err = csvWriter.Write(line)
	csvWriter.Flush()

	if err != nil {
		return false, errors.New("error: CSV Writer could not write to JobList csv")
	}

	return true, nil
}

/*
grabJobs()

Makes a request to site and parses Internships in particular
*/
func grabJobs(entry []string) (bool, error) {
	url := entry[3]
	fmt.Println(url)
	// converting jobType to Int
	jobType, err := strconv.Atoi(entry[2])
	if err != nil {
		return false, errors.New("error turning jobType to Int")
	}

	if jobType == 1 {
		job, err := requestGreenHouseJobs(url)
		if err != nil {
			return false, errors.New("Error requesting job at" + url)
		}

		listing := GreenHouseJob_to_JobListing(&(job.Jobs[1]), entry[1])

		if err != nil {
			return false, errors.New("error converting workday job listing to generic job listing")
		}

		_, err = writeInternshipToFile(&listing, "./JobList.csv")
		if err != nil {
			return false, err
		}

	} else {
		_, err := requestWorkDayJobs(url)
		if err != nil {
			return false, errors.New("Error requesting job at" + url)
		}

	}

	return true, nil
}

/*
ScrapeNewInternships() -> bool

Will go through the JobLinks.csv database, locate job listings that have not yet
been located and then concurrently request and scrape them.

Once all the new listings have been found, it will write it into the newJobs database
(non concurrently as I don't think csv supports concurrent writing)

returns true if successful, false if not and an error object?
*/
func ScrapeNewInternships() (bool, error) {

	file, err := os.Open("./JobLinks.csv")
	if err != nil {
		return false, errors.New("error: unable to find job links csv file")
	}
	defer file.Close()

	csvReader := csv.NewReader(file)
	// ToDo: Check the space efficiency of using ReadAll()
	records, err := csvReader.ReadAll()
	if err != nil {
		return false, errors.New("error: Unable to read Job Links CSV File")
	}

	var wg sync.WaitGroup

	for _, entry := range records[1:] {

		wg.Add(1)
		go func(entry []string) {
			defer wg.Done()
			// ToDo : Think about how I should handle this if the function grabJobs
			// Returns an error ???
			_, err := grabJobs(entry)
			if err != nil {
				fmt.Println(err.Error())
			}
		}(entry)
	}

	wg.Wait()
	return true, nil
}
