package goqsm

import (
	"context"
	"net/http"
	"net/url"
	"strconv"
)

// VolumeOp handles volume related methods of the QSM storage.
type VolumeOp struct {
	client *AuthClient
}

type VolumeData struct {
	ID      string `json:"id"`
	Name    string `json:"name"`
	IntName string `json:"intName"`
	NaaID   string `json:"naaId"`
	VMPath  string `json:"vmPath"`
	SizeMB  int    `json:"sizeMB"`
	UsedMB  int    `json:"usedMB"`
}

func NewVolume(client *AuthClient) *VolumeOp {
	return &VolumeOp{client}
}

func (v *VolumeOp) ListVolumes(ctx context.Context, poolId, volId string) (*[]VolumeData, error) {
	req, err := v.client.NewRequest(ctx, http.MethodGet, "/rest/internal/cloud/containers/"+poolId+"/vols/"+volId, nil)

	if err != nil {
		return nil, err
	}

	res := []VolumeData{}
	if err := v.client.sendRequest(ctx, req, &res); err != nil {
		return nil, err
	}

	return &res, nil
}

func (v *VolumeOp) CreateVolume(ctx context.Context, poolId, name string, size uint64) (*VolumeData, error) {
	params := url.Values{}
	params.Add("name", name)
	params.Add("sizeMB", strconv.FormatUint(size, 10))

	req, err := v.client.NewRequest(ctx, http.MethodPost, "/rest/internal/cloud/containers/"+poolId+"/vols", params)
	if err != nil {
		return nil, err
	}

	res := VolumeData{}
	if err := v.client.sendRequest(ctx, req, &res); err != nil {
		return nil, err
	}

	return &res, nil
}

func (v *VolumeOp) DeleteVolume(ctx context.Context, poolId, volId string) error {
	req, err := v.client.NewRequest(ctx, http.MethodDelete, "/rest/internal/cloud/containers/"+poolId+"/vols/"+volId, nil)
	if err != nil {
		return err
	}

	res := EmptyData{}
	if err := v.client.sendRequest(ctx, req, &res); err != nil {
		return err
	}

	return nil
}

func (v *VolumeOp) ExportVolume(ctx context.Context, poolId, volId string) error {
	req, err := v.client.NewRequest(ctx, http.MethodPost, "/rest/internal/cloud/containers/"+poolId+"/vols/"+volId+"/share", nil)
	if err != nil {
		return err
	}

	res := EmptyData{}
	if err := v.client.sendRequest(ctx, req, &res); err != nil {
		return err
	}

	return nil
}

func (v *VolumeOp) UnexportVolume(ctx context.Context, poolId, volId string) error {
	req, err := v.client.NewRequest(ctx, http.MethodDelete, "/rest/internal/cloud/containers/"+poolId+"/vols/"+volId+"/share", nil)
	if err != nil {
		return err
	}

	res := EmptyData{}
	if err := v.client.sendRequest(ctx, req, &res); err != nil {
		return err
	}

	return nil
}
