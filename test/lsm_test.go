package libstoragemgmt

import (
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	lsm "github.com/libstorage/libstoragemgmt-golang"
)

// Running these tests requires lsmd up and running with the simulator
// plugin available.  In addition a number of tests require things to
// exists which don't exist by default and need to be created.  At the
// moment this needs to be done via lsmcli.  As functionality evolves
// these requirements will be reduced as the unit tests will create
// things as needed.

func rs(pre string, n int) string {
	var l = []byte("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
	b := make([]byte, n)
	for i := range b {
		b[i] = l[rand.Intn(len(l))]
	}
	return fmt.Sprintf("%s%s", pre, string(b))
}

func TestConnect(t *testing.T) {
	var c, libError = lsm.Client("sim://", "", 30000)
	assert.Nil(t, libError)
	assert.Equal(t, c.Close(), nil)
}

func TestSystems(t *testing.T) {
	var c, _ = lsm.Client("sim://", "", 30000)
	var systems, sysError = c.Systems()

	assert.Nil(t, sysError)
	assert.Equal(t, len(systems), 1)

	for _, s := range systems {
		t.Logf("%+v", s)
	}
}

func TestVolumes(t *testing.T) {
	var c, _ = lsm.Client("sim://", "", 30000)
	var items, err = c.Volumes()

	assert.Nil(t, err)

	for _, i := range items {
		t.Logf("%+v", i)
	}
}

func TestPools(t *testing.T) {
	var c, _ = lsm.Client("sim://", "", 30000)
	var items, err = c.Pools()

	assert.Nil(t, err)
	assert.Greater(t, len(items), 0)

	for _, i := range items {
		t.Logf("%+v", i)
	}
}

func TestDisks(t *testing.T) {
	var c, _ = lsm.Client("sim://", "", 30000)
	var items, err = c.Disks()

	assert.Nil(t, err)
	assert.Greater(t, len(items), 0)

	for _, s := range items {

		if s.DiskType == lsm.DiskTypeSata {
			t.Logf("Got the sata disk!")
		}

		if s.Status&(lsm.DiskStatusOk|lsm.DiskStatusFree) == lsm.DiskStatusOk|lsm.DiskStatusFree {
			t.Logf("Disk OK and FREE %x\n", lsm.DiskStatusOk|lsm.DiskStatusFree)
		}

		t.Logf("%+v", s)
	}
}

func TestFs(t *testing.T) {
	var c, _ = lsm.Client("sim://", "", 30000)
	var items, err = c.FileSystems()

	assert.Nil(t, err)
	assert.Greater(t, len(items), 0)

	for _, i := range items {
		t.Logf("%+v", i)
	}
}

func TestNfsExports(t *testing.T) {
	var c, _ = lsm.Client("sim://", "", 30000)
	var items, err = c.NfsExports()

	assert.Nil(t, err)
	assert.Greater(t, len(items), 0)

	for _, i := range items {
		t.Logf("%+v", i)
	}
}

func TestAccessGroups(t *testing.T) {
	var c, _ = lsm.Client("sim://", "", 30000)
	var items, err = c.AccessGroups()

	assert.Nil(t, err)

	for _, i := range items {
		t.Logf("%+v", i)
	}

	assert.Greater(t, len(items), 0)
}

func TestTargetPorts(t *testing.T) {
	var c, _ = lsm.Client("sim://", "", 30000)
	var items, err = c.TargetPorts()

	assert.Nil(t, err)

	for _, i := range items {
		t.Logf("%+v", i)
	}

	assert.Greater(t, len(items), 0)
}

func TestBatteries(t *testing.T) {
	var c, _ = lsm.Client("sim://", "", 30000)
	var items, err = c.Batteries()

	assert.Nil(t, err)

	for _, i := range items {
		t.Logf("%+v", i)
	}

	assert.Greater(t, len(items), 0)
}

func TestCapabilities(t *testing.T) {
	var c, _ = lsm.Client("sim://", "", 30000)
	var systems, sysError = c.Systems()
	assert.Nil(t, sysError)

	var cap, capErr = c.Capabilities(&systems[0])
	assert.Nil(t, capErr)

	assert.True(t, cap.IsSupported(lsm.CapVolumeCreate))
}

func TestRepBlockSize(t *testing.T) {
	var c, err = lsm.Client("sim://", "", 30000)
	assert.Nil(t, err)
	assert.NotNil(t, c)

	var systems, sysError = c.Systems()
	assert.Nil(t, sysError)

	var repRangeBlkSize, rpbE = c.VolumeRepRangeBlkSize(&systems[0])
	assert.Nil(t, rpbE)
	assert.Equal(t, repRangeBlkSize, uint32(512))
}

func createVolume(t *testing.T, c *lsm.ClientConnection, name string) *lsm.Volume {
	var pools, poolError = c.Pools()
	assert.Nil(t, poolError)

	var poolToUse = pools[2] // Arbitrary

	var volume lsm.Volume
	var jobID, errVolCreate = c.VolumeCreate(&poolToUse, name, 1024*1024*100, 2, true, &volume)
	assert.Nil(t, errVolCreate)
	assert.Nil(t, jobID)

	return &volume
}

func TestVolumeCreate(t *testing.T) {
	var volumeName = rs("lsm_go_vol_", 8)
	var c, err = lsm.Client("sim://", "", 30000)
	assert.Nil(t, err)

	var volume = createVolume(t, c, volumeName)

	assert.Equal(t, volume.Name, volumeName)

	// Try and clean-up
	c.VolumeDelete(volume, true)
}

func TestVolumeDelete(t *testing.T) {
	var volumeName = rs("lsm_go_vol_", 8)

	var c, err = lsm.Client("sim://", "", 30000)
	assert.Nil(t, err)

	var volume = createVolume(t, c, volumeName)

	var _, errD = c.VolumeDelete(volume, true)
	assert.Nil(t, errD)
}

func TestJobWait(t *testing.T) {
	var c, err = lsm.Client("sim://", "", 30000)
	assert.Nil(t, err)

	var pools, poolError = c.Pools()
	assert.Nil(t, poolError)

	var poolToUse = pools[2] // Arbitrary

	var volumeName = rs("lsm_go_vol_async_", 8)

	var volume lsm.Volume
	var jobID, errCreate = c.VolumeCreate(&poolToUse, volumeName, 1024*1024*100, 2, false, &volume)
	assert.Nil(t, errCreate)
	assert.NotNil(t, jobID)

	var waitForIt = c.JobWait(*jobID, &volume)
	assert.Nil(t, waitForIt)

	assert.Equal(t, volumeName, volume.Name)

	c.VolumeDelete(&volume, true)
}

func TestVolumeResize(t *testing.T) {
	var volumeName = rs("lsm_go_vol_", 8)
	var c, err = lsm.Client("sim://", "", 30000)
	assert.Nil(t, err)
	assert.NotNil(t, c)

	var volume = createVolume(t, c, volumeName)
	var resized lsm.Volume
	var _, resizeErr = c.VolumeResize(volume, (volume.BlockSize*volume.NumOfBlocks)*2, true, &resized)
	assert.Nil(t, resizeErr)
	assert.NotEqual(t, volume.NumOfBlocks, resized.NumOfBlocks)

	c.VolumeDelete(&resized, true)
}

func TestVolumeReplicate(t *testing.T) {
	var c, err = lsm.Client("sim://", "", 30000)
	assert.Nil(t, err)
	assert.NotNil(t, c)

	var volumeName = rs("lsm_go_vol_", 8)
	var repName = rs("lsm_go_rep_", 4)

	var repVol lsm.Volume
	var srcVol = createVolume(t, c, volumeName)
	var jobID, errRep = c.VolumeReplicate(nil, lsm.VolumeReplicateTypeCopy, srcVol, repName, true, &repVol)
	assert.Nil(t, jobID)
	assert.Nil(t, errRep)

	assert.Equal(t, repVol.Name, repName)

	c.VolumeDelete(&repVol, true)
	c.VolumeDelete(srcVol, true)

}
func TestTemplate(t *testing.T) {
	var c, err = lsm.Client("sim://", "", 30000)
	assert.Nil(t, err)
	assert.NotNil(t, c)
}

func TestMain(m *testing.M) {
	//setup()

	// This will allow us to reproduce the same sequence if needed
	// if we encounter an error.
	var seed = os.Getenv("LSM_GO_SEED")
	if len(seed) > 0 {
		var sInt, err = strconv.ParseInt(seed, 10, 64)
		if err != nil {
			os.Exit(1)
		}
		rand.Seed(sInt)
	} else {
		var s = time.Now().UnixNano()
		rand.Seed(s)
		fmt.Printf("export LSM_GO_SEED=%v\n", s)
	}

	code := m.Run()
	//shutdown()
	os.Exit(code)
}
