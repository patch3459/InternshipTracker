package parse

import "time"

/*
	Contains methods appertaining to the models in models.go
*/

/*
GreenHouseJob_to_JobListing

Converts a greenhouse job listing to generic joblisting for writing
to the CSV
*/
func GreenHouseJob_to_JobListing(gh *GreenHouseJob, company string) JobListing {
	t := time.Now()

	return JobListing{
		1,
		(*gh).Title,
		company,
		(*gh).Updated_at,
		(*gh).Absolute_url,
		t.String(),
	}
}

/*
WorkDayJobPosting_to_JobListing

Converts a WorkDay Job POsting to generic joblisting for writing to the csv
*/

func WorkDayJobPosting_to_JobListing(wd *WorkDayJobPosting, baseUrl string, company string) JobListing {
	return JobListing{}
}
