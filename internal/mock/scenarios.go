package mock

import (
	"github.com/wiremock/go-wiremock"
)

func newClient(mockURL string) (*wiremock.Client, error) {
	wm := wiremock.NewClient(mockURL)
	if err := wm.ResetAllScenarios(); err != nil {
		return nil, err
	}
	return wm, nil
}
