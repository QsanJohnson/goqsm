package goqsm

import (
	"context"
	"fmt"
	"testing"
	"time"
)

func TestVolume(t *testing.T) {
	fmt.Println("------------TestVolume--------------")
	ctx = context.Background()

	options1 := VolumeCreateOptions{
		BlockSize: 4096,
		Provision: "thin",
		Compress:  "off",
		Dedup:     true,
	}
	options2 := VolumeCreateOptions{
		BlockSize: 65536,
		Provision: "thick",
		Compress:  "lz4",
		Dedup:     false,
	}

	createDeleteVolumeTest(t, 5120, &options1)
	createDeleteVolumeTest(t, 10240, &options2)

	exportUnexportVolumeTest(t)
}

func createDeleteVolumeTest(t *testing.T, volSize uint64, options *VolumeCreateOptions) {
	fmt.Printf("createDeleteVolumeTest Enter (volSize: %d,  %+v )\n", volSize, *options)

	now := time.Now()
	timeStamp := now.Format("20060102150405")
	volName := "gotest-vol-" + timeStamp
	scId := testConf.scId

	vol, err := testConf.volumeOp.CreateVolume(ctx, scId, volName, volSize, options)
	if err != nil {
		t.Fatalf("createVolume failed: %v", err)
	}
	fmt.Printf("  A volume was created. Id:%s, path: %s\n", vol.ID, vol.VMPath)

	vols, err := testConf.volumeOp.ListVolumes(ctx, scId, vol.ID)
	if err != nil {
		t.Fatalf("ListVolumes failed: %v", err)
	}
	if len(*vols) != 1 {
		t.Fatalf("Volume %s not found.", vol.ID)
	}

	err = testConf.volumeOp.DeleteVolume(ctx, scId, vol.ID)
	if err != nil {
		t.Fatalf("DeleteVolume failed: %v", err)
	}
	fmt.Printf("  A volume was deleted. Id:%s\n", vol.ID)

	fmt.Println("createDeleteVolumeTest Leave")
}

func exportUnexportVolumeTest(t *testing.T) {
	fmt.Println("exportUnexportVolumeTest Enter")
	now := time.Now()
	timeStamp := now.Format("20060102150405")
	volName := "gotest-vol-" + timeStamp
	scId := testConf.scId
	var volSize uint64 = 1024

	options := VolumeCreateOptions{}
	vol, err := testConf.volumeOp.CreateVolume(ctx, scId, volName, volSize, &options)
	if err != nil {
		t.Fatalf("createVolume failed: %v", err)
	}
	fmt.Printf("  A volume was created. Id:%s, path: %s \n", vol.ID, vol.VMPath)

	err = testConf.volumeOp.ExportVolume(ctx, scId, vol.ID)
	if err != nil {
		t.Fatalf("ExportVolume failed: %v", err)
	}
	fmt.Printf("  A volume was exported. Id:%s\n", vol.ID)

	err = testConf.volumeOp.UnexportVolume(ctx, scId, vol.ID)
	if err != nil {
		t.Fatalf("UnexportVolume failed: %v", err)
	}
	fmt.Printf("  A volume was unexported. Id:%s\n", vol.ID)

	err = testConf.volumeOp.DeleteVolume(ctx, scId, vol.ID)
	if err != nil {
		t.Fatalf("DeleteVolume failed: %v", err)
	}
	fmt.Printf("  A volume was deleted. Id:%s\n", vol.ID)

	fmt.Println("exportUnexportVolumeTest Leave")
}
