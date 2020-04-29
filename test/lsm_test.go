package libstoragemgmt

import (
	"testing"

	"github.com/stretchr/testify/assert"

	lsm "github.com/libstorage/libstoragemgmt-golang"
)

// Running these tests requires lsmd up and running with the simulator
// plugin available.  In addition a number of tests require things to
// exists which don't exist by default and need to be created.  At the
// moment this needs to be done via lsmcli.  As functionality evolves
// these requirements will be reduced as the unit tests will create
// things as needed.

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

func TestTemplate(t *testing.T) {
	var c, err = lsm.Client("sim://", "", 30000)
	assert.Nil(t, err)
	assert.NotNil(t, c)
}
