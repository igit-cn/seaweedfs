package mount

import (
	"github.com/chrislusf/seaweedfs/weed/pb/filer_pb"
	"github.com/hanwen/go-fuse/v2/fuse"
)

func (wfs *WFS) AcquireHandle(inode uint64, uid, gid uint32) (fileHandle *FileHandle, status fuse.Status) {
	var entry *filer_pb.Entry
	_, _, entry, status = wfs.maybeReadEntry(inode)
	if status == fuse.OK {
		fileHandle = wfs.fhmap.AcquireFileHandle(wfs, inode, entry)
	}
	return
}

func (wfs *WFS) ReleaseHandle(handleId FileHandleId) {
	wfs.fhmap.ReleaseByHandle(handleId)
}

func (wfs *WFS) GetHandle(handleId FileHandleId) *FileHandle {
	return wfs.fhmap.GetFileHandle(handleId)
}
