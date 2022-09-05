// @2022 QSAN Inc. All rights reserved

package goqsm

import (
	"context"
	"encoding/json"
	"net/http"
)

// TargetOp handles target related methods of the QSM storage.
type TargetOp struct {
	client *AuthClient
}

type Iscsi struct {
	Eths []string `json:"eths"`
}

type Host struct {
	Name []string `json:"name"`
}

type HostGroup struct {
	Name  string `json:"name"`
	Hosts []Host `json:"hosts"`
}

type CreateTargetParam struct {
	Type       string `json:"type"`
	Iscsi      `json:"iscsi"`
	HostGroups []HostGroup `json:"hostGroup"`
}

type TargetData struct {
	ID    string `json:"id"`
	Type  string `json:"type"`
	Iscsi struct {
		Iqn   string      `json:"iqn"`
		Name  string      `json:"name"`
		Alias interface{} `json:"alias"`
		Eths  []string    `json:"eths"`
	} `json:"iscsi"`
	Luns       []interface{} `json:"luns"`
	HostGroups []HostGroup   `json:"hostGroup"`
}

// NewTarget returns volume operation
func NewTarget(client *AuthClient) *TargetOp {
	return &TargetOp{client}
}

// CreateTarget create a target on a storage server
func (v *TargetOp) CreateTarget(ctx context.Context, param *CreateTargetParam) (*TargetData, error) {
	rawdata, _ := json.Marshal(param)
	req, err := v.client.NewRequest(ctx, http.MethodPost, "/rest/v2/dataTransfer/targets", string(rawdata))
	if err != nil {
		return nil, err
	}

	res := TargetData{}
	if err := v.client.SendRequest(ctx, req, &res); err != nil {
		return nil, err
	}

	return &res, nil
}
