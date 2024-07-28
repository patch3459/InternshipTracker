package jobModels

type GreenHouseDataCompliance struct {
	Data_type                   string
	Requires_consent            bool
	Requires_processing_consent bool
	Requires_retention_consent  bool
	Retention_period            int
}
