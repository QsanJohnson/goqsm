// @2022 QSAN Inc. All rights reserved

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

// The response data of volume related methods, ex ListVolumes and CreateVolume.
type VolumeData struct {
	ID      string `json:"id"`
	Name    string `json:"name"`
	IntName string `json:"intName"`
	NaaID   string `json:"naaId"`
	VMPath  string `json:"vmPath"`
	SizeMB  int    `json:"sizeMB"`
	UsedMB  int    `json:"usedMB"`
}

// NewVolume returns volume operation
func NewVolume(client *AuthClient) *VolumeOp {
	return &VolumeOp{client}
}

// ListVolumes list all volumes or a dedicated volume with volId
func (v *VolumeOp) ListVolumes(ctx context.Context, scId, volId string) (*[]VolumeData, error) {
	req, err := v.client.NewRequest(ctx, http.MethodGet, "/rest/internal/cloud/containers/"+scId+"/vols/"+volId, nil)

	if err != nil {
		return nil, err
	}

	res := []VolumeData{}
	if err := v.client.sendRequest(ctx, req, &res); err != nil {
		return nil, err
	}

	return &res, nil
}

// CreateVolume create a volume on a storage container
func (v *VolumeOp) CreateVolume(ctx context.Context, scId, name string, size uint64) (*VolumeData, error) {
	params := url.Values{}
	params.Add("name", name)
	params.Add("sizeMB", strconv.FormatUint(size, 10))

	req, err := v.client.NewRequest(ctx, http.MethodPost, "/rest/internal/cloud/containers/"+scId+"/vols", params)
	if err != nil {
		return nil, err
	}

	res := VolumeData{}
	if err := v.client.sendRequest(ctx, req, &res); err != nil {
		return nil, err
	}

	return &res, nil
}

// DeleteVolume delete a volume from a storage container
func (v *VolumeOp) DeleteVolume(ctx context.Context, scId, volId string) error {
	req, err := v.client.NewRequest(ctx, http.MethodDelete, "/rest/internal/cloud/containers/"+scId+"/vols/"+volId, nil)
	if err != nil {
		return err
	}

	res := EmptyData{}
	if err := v.client.sendRequest(ctx, req, &res); err != nil {
		return err
	}

	return nil
}

// ExportVolume export a NFS volume
func (v *VolumeOp) ExportVolume(ctx context.Context, scId, volId string) error {
	req, err := v.client.NewRequest(ctx, http.MethodPost, "/rest/internal/cloud/containers/"+scId+"/vols/"+volId+"/share", nil)
	if err != nil {
		return err
	}

	res := EmptyData{}
	if err := v.client.sendRequest(ctx, req, &res); err != nil {
		return err
	}

	return nil
}

// UnexportVolume unexport a NFS volume
func (v *VolumeOp) UnexportVolume(ctx context.Context, scId, volId string) error {
	req, err := v.client.NewRequest(ctx, http.MethodDelete, "/rest/internal/cloud/containers/"+scId+"/vols/"+volId+"/share", nil)
	if err != nil {
		return err
	}

	res := EmptyData{}
	if err := v.client.sendRequest(ctx, req, &res); err != nil {
		return err
	}

	return nil
}
