package mock

import (
	"github.com/eu-sovereign-cloud/conformance/secalib"
	"github.com/wiremock/go-wiremock"
)

type Scenarios struct {
	params  secalib.GeneralParams
	mockURL string
}

func (scenario *Scenarios) newClient() (*wiremock.Client, error) {
	wm := wiremock.NewClient(scenario.mockURL)
	if err := wm.ResetAllScenarios(); err != nil {
		return nil, err
	}
	return wm, nil
}
