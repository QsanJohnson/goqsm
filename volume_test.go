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

	resizeVolumeTest(t)

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

func resizeVolumeTest(t *testing.T) {
	fmt.Println("resizeVolumeTest Enter")
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

	volSize = 5120
	err = testConf.volumeOp.ResizeVolume(ctx, scId, vol.ID, volSize)
	if err != nil {
		t.Fatalf("resizeVolumeTest failed: %v", err)
	}
	vols, err := testConf.volumeOp.ListVolumes(ctx, scId, vol.ID)
	if err != nil {
		t.Fatalf("ListVolumes failed: %v", err)
	}
	if len(*vols) != 1 {
		t.Fatalf("Volume %s not found.", vol.ID)
	}
	if (*vols)[0].SizeMB != volSize {
		t.Fatalf("resizeVolumeTest failed: size is not match (%d vs %d)", (*vols)[0].SizeMB, volSize)
	}
	fmt.Printf("  A volume with ID %s was resize to %d MB\n", vol.ID, volSize)

	volSize = 10240
	err = testConf.volumeOp.ResizeVolume(ctx, scId, vol.ID, volSize)
	if err != nil {
		t.Fatalf("resizeVolumeTest failed: %v", err)
	}
	vols, err = testConf.volumeOp.ListVolumes(ctx, scId, vol.ID)
	if err != nil {
		t.Fatalf("ListVolumes failed: %v", err)
	}
	if len(*vols) != 1 {
		t.Fatalf("Volume %s not found.", vol.ID)
	}
	if (*vols)[0].SizeMB != volSize {
		t.Fatalf("resizeVolumeTest failed: size is not match (%d vs %d)", (*vols)[0].SizeMB, volSize)
	}
	fmt.Printf("  A volume with ID %s was resize to %d MB\n", vol.ID, volSize)

	err = testConf.volumeOp.DeleteVolume(ctx, scId, vol.ID)
	if err != nil {
		t.Fatalf("DeleteVolume failed: %v", err)
	}
	fmt.Printf("  A volume was deleted. Id:%s\n", vol.ID)

	fmt.Println("resizeVolumeTest Leave")
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
