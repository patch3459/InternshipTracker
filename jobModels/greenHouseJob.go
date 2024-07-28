package jobModels

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
