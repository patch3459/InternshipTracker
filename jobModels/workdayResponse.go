package jobModels

type WorkDayResponse struct {
	Total       int                 `json:"total"`
	JobPostings []WorkDayJobPosting `json:"jobPostings"`
}
