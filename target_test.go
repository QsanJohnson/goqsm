package goqsm

import (
	"context"
	"fmt"
	"testing"
)

func TestTarget(t *testing.T) {
	fmt.Println("------------TestTarget--------------")

	ctx = context.Background()

	// createTargetTest(t)
}

func createTargetTest(t *testing.T) {
	fmt.Println("createTargetTest Enter")

	param := &CreateTargetParam{
		Type: "iSCSI",
		Iscsi: Iscsi{
			Eths: []string{"c0e1", "c0e2"},
		},
		HostGroup: []HostGroup{
			{
				Name:  "test_group6",
				Hosts: []Host{{Name: []string{"*"}}},
			},
		},
	}
	tgt, err := testConf.targetOp.CreateTarget(ctx, param)
	if err != nil {
		t.Fatalf("CreateTarget failed: %v", err)
	}
	fmt.Printf("  A Target was created. %+v\n", tgt)

	fmt.Println("createTargetTest Leave")
}
