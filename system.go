package goqsm

import (
	"context"
	"net/http"
)

// SystemOp handles system related methods of the QSM storage.
type SystemOp struct {
	client *Client
}

type AboutData struct {
	Addresses []struct {
		Address string `json:"address"`
		Online  bool   `json:"online"`
	} `json:"addresses"`
	SystemName   string `json:"systemName"`
	FirmwareVer  string `json:"firmwareVer"`
	ModelName    string `json:"modelName"`
	ModelType    string `json:"modelType"`
	SerialNumber string `json:"serialNumber"`
	Wwn          string `json:"wwn"`
}

//NewSystem function returns system operation
func NewSystem(client *Client) *SystemOp {
	return &SystemOp{client}
}

func (s *SystemOp) GetAbout(ctx context.Context) (*AboutData, error) {
	req, err := s.client.NewRequest(ctx, http.MethodGet, "/rest/v1/about", nil)
	if err != nil {
		return nil, err
	}

	res := AboutData{}
	if err := s.client.sendRequest(ctx, req, &res); err != nil {
		return nil, err
	}

	return &res, nil
}
