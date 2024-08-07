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
	"net/http"
	"os"
	"strconv"
	"strings"
	"sync"

	"github.com/PuerkitoBio/goquery"
	"github.com/patch3459/InternshipTracker/config"
	"github.com/patch3459/InternshipTracker/jobModels"
)

/*
parseLeverCoJobsHtml

parses the Html of a lever co html response and gets jobs from it,
returnign it as a lever co response type using goquery library
*/
func parseLeverCoJobsHtml(html string) (jobModels.LeverCoResponse, error) {
	reader := strings.NewReader(html)

	doc, err := goquery.NewDocumentFromReader(reader)
	if err != nil {
		return jobModels.LeverCoResponse{}, errors.New("Error parsing html")
	}

	var postings jobModels.LeverCoResponse

	// grabbing each posting, which is wrapped in a div
	// with the classname posting
	doc.Find(".posting").Each(func(i int, s *goquery.Selection) {
		// for each one, instantitate  a job posting and add it to
		// the leverCo Object

		title := s.Find(".posting-name").Text()
		url, _ := s.Find(".posting-btn-submit").Attr("href")
		category := s.Find(".department").Text()
		location := s.Find(".location").Text()
		arrangement := s.Find(".commitment").Text()
		contractType := s.Find(".workplaceTypes").Text()

		job := jobModels.LeverCoJobPosting{title, location, category, contractType, arrangement, url}

		postings.JobPostings = append(postings.JobPostings, job)
	})

	postings.Total = len(postings.JobPostings)

	return postings, nil
}

/*
requestLeverCoJobs

makes a request from lever co jobs and parses it using html and returns an object of jobModels.LeverCoResponseType
*/
func RequestLeverCoJobs(url string) (jobModels.LeverCoResponse, error) {
	// making the html request

	var resp *http.Response

	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return jobModels.LeverCoResponse{}, errors.New("error with making request")
	}

	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,*/*;q=0.8")
	req.Header.Set("Accept_Language", "en-US,en;q=0.5")
	req.Header.Set("Accept-Encoding", "gzip, deflate, br")

	resp, err = client.Do(req)
	if err != nil {
		return jobModels.LeverCoResponse{}, errors.New("error making a request to " + url)
	}

	body, _ := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()

	var jobs jobModels.LeverCoResponse
	jobs, err = parseLeverCoJobsHtml(string(body))
	if err != nil {
		return jobModels.LeverCoResponse{}, errors.New("Error parsing Lever Co Html")
	}

	return jobs, nil
}

/*
requestGreenHouseJobs

makes a request for the jobs from a particular greenhouse link and returns them as an object of
jobModels.GreenHouseResponse Type
*/
func requestGreenHouseJobs(url string) (jobModels.GreenHouseResponse, error) {
	var resp *http.Response

	// via the url keyword
	reqUrl := fmt.Sprintf("https://boards-api.greenhouse.io/v1/boards/%s/jobs", url)
	// making the request
	resp, err := http.Get(reqUrl)
	if err != nil {
		return jobModels.GreenHouseResponse{}, errors.New("error making a request to " + reqUrl)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return jobModels.GreenHouseResponse{}, errors.New("error reading json response" + err.Error())
	}

	// parsing the JSON response
	var listings jobModels.GreenHouseResponse
	if err := json.Unmarshal(body, &listings); err != nil {
		return jobModels.GreenHouseResponse{}, errors.New("error parsing Json Response " + err.Error())
	}

	return listings, nil
}

/*
makeWorkdayAPILink

takes the page of the myworkdayjobs job page and generates a link for it's
jobs in json format

example : https://workday.wd5.myworkdayjobs.com/Workday

url workday.wd5.myworkdayjobs.com
company workday
page /Workday

resulting api link : workday.wd5.myworkdayjobs.com/Workday/wday/cxs/workday/Workday/jobs
*/
func makeWorkdayAPILink(url string) string {
	// pruning the string of www. and https://
	if strings.Contains(url, "www.") {
		ind := strings.Index(url, ".")
		url = url[ind+1:]
	} else if strings.Contains(url, "https://") {
		ind := strings.Index(url, "https://")
		url = url[ind+8:]
	}

	baseUrl := url[:strings.LastIndex(url, "/")]

	ind := strings.Index(url, ".")
	company := url[:ind]
	ind = strings.Index(url, "/")
	page := url[ind+1:]

	apiLink := fmt.Sprintf("https://%s/wday/cxs/%s/%s/jobs", baseUrl, company, page)

	return apiLink
}

/*
requestWorkDayJobs

Makes request to myworkdayjobs page and returns the jobs as a
jobModels.WorkDayResponse object
*/
func requestWorkDayJobs(url string) (jobModels.WorkDayResponse, error) {
	var resp *http.Response

	url = makeWorkdayAPILink(url)

	client := &http.Client{}
	req, err := http.NewRequest("POST", url, nil)
	if err != nil {
		return jobModels.WorkDayResponse{}, errors.New("error with making request")
	}

	req.Header.Set("Accept", "application/json")
	req.Header.Set("Accept_Language", "en-US")
	req.Header.Set("Content-Type", "application/json")

	resp, err = client.Do(req)
	if err != nil {
		return jobModels.WorkDayResponse{}, errors.New("error making a request to " + url)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return jobModels.WorkDayResponse{}, errors.New("error reading json response" + err.Error())
	}
	var listings jobModels.WorkDayResponse
	if err := json.Unmarshal(body, &listings); err != nil {
		return jobModels.WorkDayResponse{}, errors.New("error parsing Json Response " + err.Error())
	}
	return listings, nil
}

/*
writeInternshipToFile()

Writes a job listing into a csv found at found_internship_csv_path in the config.json
*/
func writeInternshipToFile(job *jobModels.JobListing, path string) (bool, error) {

	file, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return false, errors.New("error: Cannot open job list csv")
	}
	defer file.Close()
	csvWriter := csv.NewWriter(file)
	var line []string = []string{
		(*job).ID,
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
func grabJobs(entry []string, config config.Config) (bool, error) {
	url := entry[3]

	// converting jobType to Int
	jobType, err := strconv.Atoi(entry[2])
	if err != nil {
		return false, errors.New("error turning jobType to Int")
	}

	// if we're requesting from a GreenHouse io link
	if jobType == 1 {
		//TO:DO ADD WRITING SUPPORT FOR OTHER WEBSITE SUPPORTS OF INTERNSHIPS
		jobs, err := requestGreenHouseJobs(url)
		if err != nil {
			return false, errors.New("Error requesting job at" + url)
		}

		// conducting filter and converting those to listing

		for _, jobListing := range jobs.Jobs {
			var listing jobModels.JobListing

			if hasKeyword(jobListing.Title, config.Keywords) {
				listing = jobModels.GreenHouseJob_to_JobListing(&jobListing, entry[1])
				// writing filtered listings to the file
				_, err = writeInternshipToFile(&listing, config.JobListPath)
				if err != nil {
					return false, err
				}
				break
			}
		}

		// if we're requesting form a myworkdayjobs link
	} else {
		jobs, err := requestWorkDayJobs(url)
		if err != nil {
			return false, err
		}

		// conducting filter and converting those to listing
		for _, jobListing := range jobs.JobPostings {

			var listing jobModels.JobListing
			for _, keyword := range config.Keywords {
				if strings.Contains(jobListing.Title, keyword) {
					listing = jobModels.WorkDayJobPosting_to_JobListing(&jobListing, url, entry[1])
					// writing filtered listings to the file
					_, err = writeInternshipToFile(&listing, config.JobListPath)
					if err != nil {
						return false, err
					}
					break
				}
			}
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
	// reading the config file
	config := config.ReadConfigFile()

	// reading company list csv
	file, err := os.Open(config.CompanyListPath)
	if err != nil {
		return false, errors.New("error: unable to find job links csv file")
	}
	defer file.Close()

	// reading in the csv holding internship links to
	// aggregate from
	csvReader := csv.NewReader(file)
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
			_, err := grabJobs(entry, config)
			if err != nil {
				fmt.Println(err.Error())
			}
		}(entry)
	}

	wg.Wait()
	return true, nil
}
