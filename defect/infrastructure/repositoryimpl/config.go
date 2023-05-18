package repositoryimpl

type Config struct {
	Table Table `json:"table" required:"true"`
}

type Table struct {
	Defect string `json:"defect_manager" required:"true"`
}
