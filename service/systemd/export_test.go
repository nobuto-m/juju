// Copyright 2015 Canonical Ltd.
// Licensed under the AGPLv3, see LICENCE file for details.

package systemd

import (
	"github.com/juju/testing"
)

type patcher interface {
	PatchValue(interface{}, interface{})
}

func PatchFindDataDir(patcher patcher, dataDir string) {
	patcher.PatchValue(&findDataDir, func() (string, error) {
		return dataDir, nil
	})
}

func PatchNewChan(patcher patcher) chan string {
	ch := make(chan string, 1)
	patcher.PatchValue(&newChan, func() chan string { return ch })
	return ch
}

func PatchNewConn(patcher patcher, stub *testing.Stub) *StubDbusAPI {
	conn := &StubDbusAPI{Stub: stub}
	patcher.PatchValue(&newConn, func() (dbusAPI, error) { return conn, nil })
	return conn
}

func PatchFileOps(patcher patcher, stub *testing.Stub) *StubFileOps {
	fops := &StubFileOps{Stub: stub}
	patcher.PatchValue(&removeAll, fops.RemoveAll)
	patcher.PatchValue(&mkdirAll, fops.MkdirAll)
	patcher.PatchValue(&createFile, fops.CreateFile)
	return fops
}
