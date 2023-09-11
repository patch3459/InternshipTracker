/*
	Has Structs for the different kinds of response shapes we will get
*/

package parse

// Models for job listing responses

type GreenHouseResponseMetaData struct {
	Total int
}
type GreenHouseDataCompliance struct {
	Data_type                   string
	Requires_consent            bool
	Requires_processing_consent bool
	Requires_retention_consent  bool
	Retention_period            int
}

type GreenHouseLocation struct {
	Name string
}

type GreenHouseJobListingMetaData struct {
	Id         int
	Name       string
	Value      string
	Value_type string
}

type GreenHouseJob struct {
	Absolute_url    string
	Data_compliance []GreenHouseDataCompliance
	Internal_job_id int
	Location        GreenHouseLocation
	//ToDo : Something is horribly off with metadata
	metadata       []GreenHouseJobListingMetaData
	Id             int
	Updated_at     string
	Requisition_id string
	Title          string
}

type GreenHouseResponse struct {
	Jobs []GreenHouseJob            `json:"jobs"`
	Meta GreenHouseResponseMetaData `json:"meta"`
}

type WorkDayJobPosting struct {
	Title         string
	ExternalPath  string
	LocationsText string
	PostedOn      string
}

type WorkDayResponse struct {
	Total       int                 `json:"total"`
	JobPostings []WorkDayJobPosting `json:"jobPostings"`
}

// Other models

type JobListing struct {
	ID           string
	Title        string
	Company      string
	DatePosted   string
	Link         string
	DateUploaded string
}

// for the config.json
type Config struct {
	CompanyListPath string   `json:"company_list_csv_path"`
	JobListPath     string   `json:"job_list_csv_path"`
	Keywords        []string `json:"keywords"`
}
