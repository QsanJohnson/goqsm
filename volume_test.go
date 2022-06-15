package goqsm

import (
	"context"
	"fmt"
	"testing"
	"time"
)

func TestVolume(t *testing.T) {
	fmt.Println("------------TestVolume--------------")

	now := time.Now()
	timeStamp := now.Format("20060102150405")
	volName := "gotest-vol-" + timeStamp
	scId := testConf.scId
	var volSize uint64 = 1024
	fmt.Printf("TestConf: scId=%s, volName=%s\n", scId, volName)

	ctx = context.Background()

	listVolumesTest(t, scId)

	volId := createVolumeTest(t, scId, volName, volSize)
	exportVolumeTest(t, scId, volId)
	unexportVolumeTest(t, scId, volId)
	deleteVolumeTest(t, scId, volId)

}

func createVolumeTest(t *testing.T, scId, volName string, volSize uint64) string {
	fmt.Println("createVolumeTest Enter")

	vol, err := testConf.volumeOp.CreateVolume(ctx, scId, volName, volSize)
	if err != nil {
		t.Fatalf("createVolume failed: %v", err)
	}

	fmt.Printf("  A volume was created. Id:%s, path: %s \n", vol.ID, vol.VMPath)
	fmt.Println("createVolumeTest Leave")

	return vol.ID
}

func deleteVolumeTest(t *testing.T, scId, volId string) {
	fmt.Println("deleteVolumeTest Enter")

	err := testConf.volumeOp.DeleteVolume(ctx, scId, volId)
	if err != nil {
		t.Fatalf("DeleteVolume failed: %v", err)
	}

	fmt.Printf("  A volume was deleted. Id:%s\n", volId)
	fmt.Println("deleteVolumeTest Leave")
}

func listVolumesTest(t *testing.T, scId string) {
	fmt.Println("listVolumesTest Enter")

	vols, err := testConf.volumeOp.ListVolumes(ctx, scId, "")
	if err != nil {
		t.Fatalf("ListVolumes failed: %v", err)
	}

	fmt.Printf("volume count: %d \n", len(*vols))
	for _, v := range *vols {
		fmt.Println("  ID:", v.ID, ", Name:", v.Name, ", VMPath:", v.VMPath)
		vol, err := testConf.volumeOp.ListVolumes(ctx, scId, v.ID)
		if err != nil {
			t.Fatalf("ListVolumes failed with volId %s: %v", v.ID, err)
		}
		if len(*vol) != 1 {
			t.Fatalf("ListVolumes failed with volId %s: cnt=%d", v.ID, len(*vol))
		}
	}

	fmt.Println("listVolumesTest Leave")
}

func exportVolumeTest(t *testing.T, scId, volId string) {
	fmt.Println("exportVolumeTest Enter")

	err := testConf.volumeOp.ExportVolume(ctx, scId, volId)
	if err != nil {
		t.Fatalf("ExportVolume failed: %v", err)
	}

	fmt.Printf("  A volume was exported. Id:%s\n", volId)
	fmt.Println("exportVolumeTest Leave")
}

func unexportVolumeTest(t *testing.T, scId, volId string) {
	fmt.Println("unexportVolumeTest Enter")

	err := testConf.volumeOp.UnexportVolume(ctx, scId, volId)
	if err != nil {
		t.Fatalf("UnexportVolume failed: %v", err)
	}

	fmt.Printf("  A volume was unexported. Id:%s\n", volId)
	fmt.Println("unexportVolumeTest Leave")
}
