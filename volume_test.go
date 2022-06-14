package goqsm

import (
	"context"
	"fmt"
	"testing"
	"time"
)

// var poolId string
// var volName string
// var volSize uint64

func TestVolume(t *testing.T) {
	fmt.Println("------------TestVolume--------------")

	now := time.Now()
	timeStamp := now.Format("20220101000000")
	volName := "gotest-vol-" + timeStamp
	poolId := "14595646129689353715"
	var volSize uint64 = 1024

	ctx = context.Background()

	listVolumesTest(t, poolId)

	volId := createVolumeTest(t, poolId, volName, volSize)
	exportVolumeTest(t, poolId, volId)
	unexportVolumeTest(t, poolId, volId)
	deleteVolumeTest(t, poolId, volId)

}

func createVolumeTest(t *testing.T, poolId, volName string, volSize uint64) string {
	fmt.Println("createVolumeTest Enter")

	vol, err := testConf.volumeOp.CreateVolume(ctx, poolId, volName, volSize)
	if err != nil {
		t.Fatalf("createVolume failed: %v", err)
	}

	fmt.Printf("  A volume was created. Id:%s, path: %s \n", vol.ID, vol.VMPath)
	fmt.Println("createVolumeTest Leave")

	return vol.ID
}

func deleteVolumeTest(t *testing.T, poolId, volId string) {
	fmt.Println("deleteVolumeTest Enter")

	err := testConf.volumeOp.DeleteVolume(ctx, poolId, volId)
	if err != nil {
		t.Fatalf("DeleteVolume failed: %v", err)
	}

	fmt.Printf("  A volume was deleted. Id:%s\n", volId)
	fmt.Println("deleteVolumeTest Leave")
}

func listVolumesTest(t *testing.T, poolId string) {
	fmt.Println("listVolumesTest Enter")

	vols, err := testConf.volumeOp.ListVolumes(ctx, poolId, "")
	if err != nil {
		t.Fatalf("ListVolumes failed: %v", err)
	}

	fmt.Printf("volume count: %d \n", len(*vols))
	for _, v := range *vols {
		fmt.Println("  ID:", v.ID, ", Name:", v.Name, ", VMPath:", v.VMPath)
		vol, err := testConf.volumeOp.ListVolumes(ctx, poolId, v.ID)
		if err != nil {
			t.Fatalf("ListVolumes failed with volId %s: %v", v.ID, err)
		}
		if len(*vol) != 1 {
			t.Fatalf("ListVolumes failed with volId %s: cnt=%d", v.ID, len(*vol))
		}
	}

	fmt.Println("listVolumesTest Leave")
}

func exportVolumeTest(t *testing.T, poolId, volId string) {
	fmt.Println("exportVolumeTest Enter")

	err := testConf.volumeOp.ExportVolume(ctx, poolId, volId)
	if err != nil {
		t.Fatalf("ExportVolume failed: %v", err)
	}

	fmt.Printf("  A volume was exported. Id:%s\n", volId)
	fmt.Println("exportVolumeTest Leave")
}

func unexportVolumeTest(t *testing.T, poolId, volId string) {
	fmt.Println("unexportVolumeTest Enter")

	err := testConf.volumeOp.UnexportVolume(ctx, poolId, volId)
	if err != nil {
		t.Fatalf("UnexportVolume failed: %v", err)
	}

	fmt.Printf("  A volume was unexported. Id:%s\n", volId)
	fmt.Println("unexportVolumeTest Leave")
}
