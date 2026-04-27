package wiremock

import "github.com/wiremock/go-wiremock"

type MockClient struct {
	wm *wiremock.Client
}

func (client *MockClient) ResetAllScenarios() error {
	return client.wm.ResetAllScenarios()
}

func newMockClient(mockURL string) (*MockClient, error) {
	wm := wiremock.NewClient(mockURL)
	if err := wm.ResetAllScenarios(); err != nil {
		return nil, err
	}

	return &MockClient{
		wm: wm,
	}, nil
}
