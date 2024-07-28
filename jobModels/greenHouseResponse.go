package jobModels

type GreenHouseResponse struct {
	Jobs []GreenHouseJob            `json:"jobs"`
	Meta GreenHouseResponseMetaData `json:"meta"`
}
