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
	errors "github.com/libstorage/libstoragemgmt-golang/errors"
)

var URI = getEnv("LSM_GO_URI", "sim://")

const PASSWORD = ""
const TMO uint32 = 30000

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

func getEnv(variable string, defValue string) string {
	var p = os.Getenv(variable)
	if len(p) > 0 {
		return p
	}
	return defValue
}

func TestConnect(t *testing.T) {
	var c, libError = lsm.Client(URI, PASSWORD, TMO)
	assert.Nil(t, libError)
	assert.Equal(t, c.Close(), nil)
}

func TestSystems(t *testing.T) {
	var c, _ = lsm.Client(URI, PASSWORD, TMO)
	var systems, sysError = c.Systems()

	assert.Nil(t, sysError)
	assert.Equal(t, len(systems), 1)

	for _, s := range systems {
		t.Logf("%+v", s)
	}
	assert.Equal(t, c.Close(), nil)
}

func TestReadCachePercentSet(t *testing.T) {
	var c, _ = lsm.Client(URI, PASSWORD, TMO)
	var systems, sysError = c.Systems()
	assert.Nil(t, sysError)

	assert.Nil(t, c.SysReadCachePctSet(&systems[0], 0))
	assert.Nil(t, c.SysReadCachePctSet(&systems[0], 100))

	var expectedErr = c.SysReadCachePctSet(&systems[0], 101)
	assert.NotNil(t, expectedErr)
	var e = expectedErr.(*errors.LsmError)
	assert.Equal(t, e.Code, errors.InvalidArgument)

	assert.Equal(t, c.Close(), nil)
}

func TestIscsiChapSet(t *testing.T) {
	var c, _ = lsm.Client(URI, PASSWORD, TMO)

	var init = "iqn.1994-05.com.domain:01.89bd01"
	var e = c.IscsiChapAuthSet(init, nil, nil, nil, nil)

	var u = rs("user_", 3)
	var p = rs("password_", 3)
	e = c.IscsiChapAuthSet(init, &u, &p, nil, nil)
	assert.Nil(t, e)

	var outU = rs("outuser_", 3)
	var outP = rs("out_password_", 3)
	e = c.IscsiChapAuthSet(init, &u, &p, &outU, &outP)
	assert.Nil(t, e)

	assert.Equal(t, c.Close(), nil)
}

func TestVolumes(t *testing.T) {
	var c, _ = lsm.Client(URI, PASSWORD, TMO)
	var items, err = c.Volumes()

	assert.Nil(t, err)

	for _, i := range items {
		t.Logf("%+v", i)
	}
	assert.Equal(t, c.Close(), nil)
}

func TestPools(t *testing.T) {
	var c, _ = lsm.Client(URI, PASSWORD, TMO)
	var items, err = c.Pools()

	assert.Nil(t, err)
	assert.Greater(t, len(items), 0)

	for _, i := range items {
		t.Logf("%+v", i)
	}
	assert.Equal(t, c.Close(), nil)
}

func TestDisks(t *testing.T) {
	var c, _ = lsm.Client(URI, PASSWORD, TMO)
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
	assert.Equal(t, c.Close(), nil)
}

func TestFs(t *testing.T) {
	var c, _ = lsm.Client(URI, PASSWORD, TMO)
	var items, err = c.FileSystems()

	assert.Nil(t, err)
	assert.Greater(t, len(items), 0)

	for _, i := range items {
		t.Logf("%+v", i)
	}
	assert.Equal(t, c.Close(), nil)
}

func TestNfsExports(t *testing.T) {
	var c, _ = lsm.Client(URI, PASSWORD, TMO)
	var items, err = c.NfsExports()
	assert.Nil(t, err)

	if len(items) == 0 {
		var fs, err = c.FileSystems()
		assert.Nil(t, err)
		var access lsm.NfsAccess

		access.Root = []string{"192.168.1.1"}
		access.Rw = []string{"192.168.1.1"}
		access.AnonGID = lsm.AnonUIDGIDNotApplicable
		access.AnonUID = lsm.AnonUIDGIDNotApplicable

		var export lsm.NfsExport

		var exportPath = "/mnt/fubar"
		var auth = "standard"
		var exportErr = c.FsExport(&fs[0], &exportPath, &access, &auth, nil, &export)
		assert.Nil(t, exportErr)
		assert.Equal(t, export.ExportPath, exportPath)
	}

	for _, i := range items {
		t.Logf("%+v", i)

		var unExportErr = c.FsUnExport(&i)
		assert.Nil(t, unExportErr)
	}
	assert.Equal(t, c.Close(), nil)
}

func TestNfsAuthTypes(t *testing.T) {
	var c, _ = lsm.Client(URI, PASSWORD, TMO)
	var authTypes, err = c.NfsExportAuthTypes()
	assert.Nil(t, err)
	assert.Equal(t, 1, len(authTypes))
	assert.Equal(t, "standard", authTypes[0])

	fmt.Printf("%v", authTypes)
	assert.Equal(t, c.Close(), nil)
}

func TestAccessGroups(t *testing.T) {
	var c, _ = lsm.Client(URI, PASSWORD, TMO)
	var items, err = c.AccessGroups()
	var volumes, volErr = c.Volumes()
	assert.Nil(t, volErr)

	assert.Nil(t, err)

	if len(items) == 0 {
		var systems, sysErr = c.Systems()
		assert.Nil(t, sysErr)

		var ag lsm.AccessGroup
		var agCreateErr = c.AccessGroupCreate(rs("lsm_ag_", 4),
			"iqn.1994-05.com.domain:01.89bd01", lsm.InitiatorTypeIscsiIqn, &systems[0], &ag)
		assert.Nil(t, agCreateErr)

		var maskErr = c.VolumeMask(&volumes[0], &ag)
		assert.Nil(t, maskErr)

		var volsMasked, volsMaskedErr = c.VolsMaskedToAg(&ag)
		assert.Nil(t, volsMaskedErr)
		assert.Equal(t, len(volsMasked), 1)
		assert.Equal(t, volumes[0].Name, volsMasked[0].Name)

		var agsGranted, agsGrantedErr = c.AgsGrantedToVol(&volumes[0])
		assert.Nil(t, agsGrantedErr)
		assert.Equal(t, len(agsGranted), 1)
		assert.Equal(t, agsGranted[0].Name, ag.Name)

		var unmaskErr = c.VolumeUnMask(&volumes[0], &ag)
		assert.Nil(t, unmaskErr)

		volsMasked, volsMaskedErr = c.VolsMaskedToAg(&ag)
		assert.Nil(t, volsMaskedErr)
		assert.Equal(t, len(volsMasked), 0)

		agsGranted, agsGrantedErr = c.AgsGrantedToVol(&volumes[0])
		assert.Nil(t, agsGrantedErr)
		assert.Equal(t, len(agsGranted), 0)

		var agInitAdd lsm.AccessGroup
		var initAddErr = c.AccessGroupInitAdd(&ag, "iqn.1994-05.com.domain:01.89bd02", lsm.InitiatorTypeIscsiIqn, &agInitAdd)
		assert.Nil(t, initAddErr)
		assert.NotEqual(t, len(ag.InitIDs), len(agInitAdd.InitIDs))

		var agInitDel lsm.AccessGroup
		var initDelErr = c.AccessGroupInitDelete(&ag, "iqn.1994-05.com.domain:01.89bd02", lsm.InitiatorTypeIscsiIqn, &agInitDel)
		assert.Nil(t, initDelErr)
		assert.Equal(t, len(ag.InitIDs), len(agInitDel.InitIDs))

		items, err = c.AccessGroups()
		assert.Nil(t, err)
	}

	assert.Greater(t, len(items), 0)

	for _, i := range items {
		t.Logf("%+v", i)

		var agDelErr = c.AccessGroupDelete(&i)
		assert.Nil(t, agDelErr)
	}

	items, err = c.AccessGroups()
	assert.Equal(t, len(items), 0)
	assert.Equal(t, c.Close(), nil)
}

func TestTargetPorts(t *testing.T) {
	var c, _ = lsm.Client(URI, PASSWORD, TMO)
	var items, err = c.TargetPorts()

	assert.Nil(t, err)

	for _, i := range items {
		t.Logf("%+v", i)
	}

	assert.Greater(t, len(items), 0)
	assert.Equal(t, c.Close(), nil)
}

func TestBatteries(t *testing.T) {
	var c, _ = lsm.Client(URI, PASSWORD, TMO)
	var items, err = c.Batteries()

	assert.Nil(t, err)

	for _, i := range items {
		t.Logf("%+v", i)
	}

	assert.Greater(t, len(items), 0)
	assert.Equal(t, c.Close(), nil)
}

func TestCapabilities(t *testing.T) {
	var c, _ = lsm.Client(URI, PASSWORD, TMO)
	var systems, sysError = c.Systems()
	assert.Nil(t, sysError)

	var cap, capErr = c.Capabilities(&systems[0])
	assert.Nil(t, capErr)

	assert.True(t, cap.IsSupported(lsm.CapVolumeCreate))
	assert.Equal(t, c.Close(), nil)
}

func TestCapabilitiesSet(t *testing.T) {
	var c, _ = lsm.Client(URI, PASSWORD, TMO)
	var systems, sysError = c.Systems()
	assert.Nil(t, sysError)

	var cap, capErr = c.Capabilities(&systems[0])
	assert.Nil(t, capErr)

	var set = []lsm.CapabilityType{lsm.CapVolumeCreate, lsm.CapVolumeCResize}
	assert.True(t, cap.IsSupportedSet(set))
	assert.Equal(t, c.Close(), nil)
}

func TestRepBlockSize(t *testing.T) {
	var c, err = lsm.Client(URI, PASSWORD, TMO)
	assert.Nil(t, err)
	assert.NotNil(t, c)

	var systems, sysError = c.Systems()
	assert.Nil(t, sysError)

	var repRangeBlkSize, rpbE = c.VolumeRepRangeBlkSize(&systems[0])
	assert.Nil(t, rpbE)
	assert.Equal(t, repRangeBlkSize, uint32(512))
	assert.Equal(t, c.Close(), nil)
}

func createVolume(t *testing.T, c *lsm.ClientConnection, name string) *lsm.Volume {
	var pools, poolError = c.Pools()
	assert.Nil(t, poolError)

	var poolToUse = pools[3] // Arbitrary

	var volume lsm.Volume
	var jobID, errVolCreate = c.VolumeCreate(&poolToUse, name, 1024*1024*1, 2, true, &volume)
	assert.Nil(t, errVolCreate)
	assert.Nil(t, jobID)

	return &volume
}

func TestVolumeCreate(t *testing.T) {
	var volumeName = rs("lsm_go_vol_", 8)
	var c, err = lsm.Client(URI, PASSWORD, TMO)
	assert.Nil(t, err)

	var volume = createVolume(t, c, volumeName)

	assert.Equal(t, volume.Name, volumeName)

	// Try and clean-up
	var volDel, volDelErr = c.VolumeDelete(volume, true)
	assert.Nil(t, volDel)
	assert.Nil(t, volDelErr)
	assert.Equal(t, c.Close(), nil)
}

func TestVolumeEnableDisable(t *testing.T) {
	var volumeName = rs("lsm_go_vol_", 8)
	var c, err = lsm.Client(URI, PASSWORD, TMO)
	assert.Nil(t, err)

	var volume = createVolume(t, c, volumeName)
	assert.Equal(t, volume.Name, volumeName)

	var disableErr = c.VolumeDisable(volume)
	assert.Nil(t, disableErr)

	var enableErr = c.VolumeEnable(volume)
	assert.Nil(t, enableErr)

	// Try and clean-up
	var volDel, volDelErr = c.VolumeDelete(volume, true)
	assert.Nil(t, volDel)
	assert.Nil(t, volDelErr)
	assert.Equal(t, c.Close(), nil)
}

func TestScale(t *testing.T) {

	var c, err = lsm.Client("simc://", "", 30000)
	assert.Nil(t, err)

	var pools, poolError = c.Pools()
	assert.Nil(t, poolError)

	var poolToUse = pools[3]

	for i := 0; i < 10; i++ {
		var volumeName = rs("lsm_go_vol_", 8)

		var volume lsm.Volume
		var jobID, errVolCreate = c.VolumeCreate(&poolToUse, volumeName, 1024*1024*10, 2, true, &volume)
		if errVolCreate != nil {
			fmt.Printf("Created %d volume before we got error %s\n", i, errVolCreate)
			break
		}
		assert.Nil(t, jobID)
	}

	var volumes, vE = c.Volumes()
	assert.Nil(t, vE)
	assert.Greater(t, len(volumes), 10)

	assert.Equal(t, c.Close(), nil)
}

func TestVolumeDelete(t *testing.T) {
	var volumeName = rs("lsm_go_vol_", 8)

	var c, err = lsm.Client(URI, PASSWORD, TMO)
	assert.Nil(t, err)

	var volume = createVolume(t, c, volumeName)

	var _, errD = c.VolumeDelete(volume, true)
	assert.Nil(t, errD)
	assert.Equal(t, c.Close(), nil)
}

func TestJobWait(t *testing.T) {
	var c, err = lsm.Client(URI, PASSWORD, TMO)
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
	assert.Equal(t, c.Close(), nil)
}

func TestTmo(t *testing.T) {
	var c, err = lsm.Client(URI, PASSWORD, TMO)
	assert.Nil(t, err)

	var tmo uint32 = 32121

	assert.Nil(t, c.TimeOutSet(tmo))

	assert.Equal(t, tmo, c.TimeOutGet())
	assert.Equal(t, c.Close(), nil)
}

func TestVolumeResize(t *testing.T) {
	var volumeName = rs("lsm_go_vol_", 8)
	var c, err = lsm.Client(URI, PASSWORD, TMO)
	assert.Nil(t, err)
	assert.NotNil(t, c)

	var volume = createVolume(t, c, volumeName)
	var resized lsm.Volume
	var _, resizeErr = c.VolumeResize(volume, (volume.BlockSize*volume.NumOfBlocks)*2, true, &resized)
	assert.Nil(t, resizeErr)
	assert.NotEqual(t, volume.NumOfBlocks, resized.NumOfBlocks)

	c.VolumeDelete(&resized, true)
	assert.Equal(t, c.Close(), nil)
}

func TestVolumeRaidType(t *testing.T) {
	var c, err = lsm.Client(URI, PASSWORD, TMO)
	assert.Nil(t, err)
	assert.NotNil(t, c)

	var volumes, volErr = c.Volumes()
	assert.Nil(t, volErr)

	var raidInfo, raidInfoErr = c.VolRaidInfo(&volumes[0])
	assert.Nil(t, raidInfoErr)
	assert.NotNil(t, raidInfo)

	assert.Equal(t, c.Close(), nil)
}

func TestPoolMemberInfo(t *testing.T) {
	var c, err = lsm.Client(URI, PASSWORD, TMO)
	assert.Nil(t, err)
	assert.NotNil(t, c)

	var pools, volErr = c.Pools()
	assert.Nil(t, volErr)

	var poolInfo, poolInfoErr = c.PoolMemberInfo(&pools[0])
	assert.Nil(t, poolInfoErr)
	assert.NotNil(t, poolInfo)

	t.Logf("%+v", poolInfo)

	assert.Equal(t, c.Close(), nil)
}

func TestVolumeReplicate(t *testing.T) {
	var c, err = lsm.Client(URI, PASSWORD, TMO)
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
	assert.Equal(t, c.Close(), nil)

}

func TestVolumeReplicateRange(t *testing.T) {
	var c, err = lsm.Client(URI, PASSWORD, TMO)
	assert.Nil(t, err)
	assert.NotNil(t, c)

	var volumeName = rs("lsm_go_vol_", 8)
	var volume = createVolume(t, c, volumeName)

	var ranges []lsm.BlockRange
	ranges = append(ranges, lsm.BlockRange{BlkCount: 100, SrcBlkAddr: 10, DstBlkAddr: 400})

	var jobID, repErr = c.VolumeReplicateRange(lsm.VolumeReplicateTypeCopy, volume, volume, ranges, true)
	assert.Nil(t, jobID)
	assert.Nil(t, repErr)

	c.VolumeDelete(volume, true)
	assert.Equal(t, c.Close(), nil)
}

func TestFsCreateResizeCloneDelete(t *testing.T) {
	var c, err = lsm.Client(URI, PASSWORD, TMO)
	assert.Nil(t, err)
	assert.NotNil(t, c)

	var pools, pE = c.Pools()
	assert.Nil(t, pE)

	var newFs lsm.FileSystem
	var fsCreateJob, fsCreateErr = c.FsCreate(&pools[2], rs("lsm_go_pool_", 4), 1024*1024*100, true, &newFs)
	assert.Nil(t, fsCreateJob)
	assert.Nil(t, fsCreateErr)

	var resizedFs lsm.FileSystem
	var resizedJob, resizedErr = c.FsResize(&newFs, newFs.TotalSpace*2, true, &resizedFs)
	assert.Nil(t, resizedJob)
	assert.Nil(t, resizedErr)
	assert.NotEqual(t, newFs.TotalSpace, resizedFs)

	var cloned lsm.FileSystem
	var cloneFsJob, cloneErr = c.FsClone(&resizedFs, "lsm_go_cloned_fs", nil, true, &cloned)
	assert.Nil(t, cloneFsJob)
	assert.Nil(t, cloneErr)

	assert.Equal(t, cloned.Name, "lsm_go_cloned_fs")
	assert.Equal(t, resizedFs.TotalSpace, cloned.TotalSpace)

	var delcloneFsJob, delCloneFsErr = c.FsDelete(&cloned, true)
	assert.Nil(t, delcloneFsJob)
	assert.Nil(t, delCloneFsErr)

	var delFsJob, delFsErr = c.FsDelete(&resizedFs, true)
	assert.Nil(t, delFsJob)
	assert.Nil(t, delFsErr)

	assert.Equal(t, c.Close(), nil)
}

func TestFsFileClone(t *testing.T) {
	var c, err = lsm.Client(URI, PASSWORD, TMO)
	assert.Nil(t, err)
	assert.NotNil(t, c)

	var pools, pE = c.Pools()
	assert.Nil(t, pE)

	var newFs lsm.FileSystem
	var fsCreateJob, fsCreateErr = c.FsCreate(&pools[2], rs("lsm_go_pool_", 4), 1024*1024*100, true, &newFs)
	assert.Nil(t, fsCreateJob)
	assert.Nil(t, fsCreateErr)

	var fsFileCloneJob, fsFcErr = c.FsFileClone(&newFs, "some_file", "some_other_file", nil, true)
	assert.Nil(t, fsFileCloneJob)
	assert.Nil(t, fsFcErr)

	var delFsJob, delFsErr = c.FsDelete(&newFs, true)
	assert.Nil(t, delFsJob)
	assert.Nil(t, delFsErr)

	assert.Equal(t, c.Close(), nil)
}

func TestFsSnapShots(t *testing.T) {
	var c, err = lsm.Client(URI, PASSWORD, TMO)
	assert.Nil(t, err)
	assert.NotNil(t, c)

	var pools, pE = c.Pools()
	assert.Nil(t, pE)

	var newFs lsm.FileSystem
	var fsCreateJob, fsCreateErr = c.FsCreate(&pools[2], rs("lsm_go_pool_", 4), 1024*1024*100, true, &newFs)
	assert.Nil(t, fsCreateJob)
	assert.Nil(t, fsCreateErr)

	var ss lsm.FileSystemSnapShot
	var ssJob, ssE = c.FsSnapShotCreate(&newFs, "lsm_go_ss", true, &ss)

	assert.Nil(t, ssJob)
	assert.Nil(t, ssE)
	assert.Equal(t, ss.Name, "lsm_go_ss")

	var hasDep, depErr = c.FsHasChildDep(&newFs, make([]string, 0))
	assert.Nil(t, depErr)
	assert.True(t, hasDep)

	var snaps, snapsErr = c.FsSnapShots(&newFs)
	assert.Nil(t, snapsErr)

	assert.Equal(t, 1, len(snaps))

	for _, i := range snaps {
		t.Logf("%+v", i)
	}

	var ssDelJob, ssDelErr = c.FsSnapShotDelete(&newFs, &ss, true)
	assert.Nil(t, ssDelJob)
	assert.Nil(t, ssDelErr)

	c.FsDelete(&newFs, true)
	assert.Equal(t, c.Close(), nil)
}

func TestFsSnapShotRestore(t *testing.T) {
	var c, err = lsm.Client(URI, PASSWORD, TMO)
	assert.Nil(t, err)
	assert.NotNil(t, c)

	var pools, pE = c.Pools()
	assert.Nil(t, pE)

	var newFs lsm.FileSystem
	var fsCreateJob, fsCreateErr = c.FsCreate(&pools[2], rs("lsm_go_pool_", 4), 1024*1024*100, true, &newFs)
	assert.Nil(t, fsCreateJob)
	assert.Nil(t, fsCreateErr)

	var ss lsm.FileSystemSnapShot
	var ssName = rs("lsm_go_ss_", 4)
	var ssJob, ssE = c.FsSnapShotCreate(&newFs, ssName, true, &ss)

	assert.Nil(t, ssJob)
	assert.Nil(t, ssE)
	assert.Equal(t, ss.Name, ssName)

	var ssRestoreJob, ssRestoreErr = c.FsSnapShotRestore(
		&newFs, &ss, true, make([]string, 0), make([]string, 0), true)

	assert.Nil(t, ssRestoreJob)
	assert.Nil(t, ssRestoreErr)

	var org = []string{"/tmp/bar"}
	var rst = []string{"/tmp/fubar"}

	var ssRestoreJobF, ssRestoreErrF = c.FsSnapShotRestore(
		&newFs, &ss, false, org, rst, true)

	assert.Nil(t, ssRestoreJobF)
	assert.Nil(t, ssRestoreErrF)

	var ssDelJob, ssDelErr = c.FsSnapShotDelete(&newFs, &ss, true)
	assert.Nil(t, ssDelJob)
	assert.Nil(t, ssDelErr)

	c.FsDelete(&newFs, true)

	c.FsDelete(&newFs, true)
	assert.Equal(t, c.Close(), nil)
}

func TestTemplate(t *testing.T) {
	var c, err = lsm.Client(URI, PASSWORD, TMO)
	assert.Nil(t, err)
	assert.NotNil(t, c)
	assert.Equal(t, c.Close(), nil)
}

func setup() {
	var c, _ = lsm.Client(URI, PASSWORD, TMO)

	var pools, _ = c.Pools()
	var volumes, _ = c.Volumes()

	if len(volumes) == 0 {
		var volume lsm.Volume
		var _, _ = c.VolumeCreate(
			&pools[1], rs("lsm_go_vol_", 4),
			1024*1024*100,
			lsm.VolumeProvisionTypeDefault, true, &volume)
	}

	var fs, _ = c.FileSystems()
	if len(fs) == 0 {
		var fileSystem lsm.FileSystem
		var _, _ = c.FsCreate(
			&pools[1], rs("lsm_go_fs_", 4), 1024*1024*1000, true, &fileSystem)
	}
}

func TestMain(m *testing.M) {
	setup()

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
