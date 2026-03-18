package report

type allureResult struct {
	Name          string       `json:"name"`
	FullName      string       `json:"fullName"`
	Status        string       `json:"status"`
	StatusDetails allureDetail `json:"statusDetails"`
	Start         int64        `json:"start"`
	Stop          int64        `json:"stop"`
	Steps         []allureStep `json:"steps"`
}

type allureDetail struct {
	Message string `json:"message"`
	Trace   string `json:"trace"`
}

type allureStep struct {
	Name          string       `json:"name"`
	Status        string       `json:"status"`
	StatusDetails allureDetail `json:"statusDetails"`
	Start         int64        `json:"start"`
	Stop          int64        `json:"stop"`
	Steps         []allureStep `json:"steps"`
}
