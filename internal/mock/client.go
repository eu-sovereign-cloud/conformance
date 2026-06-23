package mock

import "github.com/wiremock/go-wiremock"

type MockClient struct {
	Wiremock *wiremock.Client
}

func (client *MockClient) ResetAllScenarios() error {
	return client.Wiremock.ResetAllScenarios()
}

func NewMockClient(mockURL string) (*MockClient, error) {
	wm := wiremock.NewClient(mockURL)
	if err := wm.ResetAllScenarios(); err != nil {
		return nil, err
	}

	return &MockClient{
		Wiremock: wm,
	}, nil
}
