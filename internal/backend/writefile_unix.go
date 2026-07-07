// Copyright 2022 Francisco Souza. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build !windows

package backend

import (
	"io"
	"os"

	"github.com/google/renameio/v2"
)

func writeFile(filename string, data []byte, perm os.FileMode) error {
	return renameio.WriteFile(filename, data, perm)
}

// writeFileStream atomically writes the contents of r to filename without
// buffering the whole payload in memory: r is copied straight into the pending
// temp file, which is then atomically renamed into place.
func writeFileStream(filename string, r io.Reader, perm os.FileMode) error {
	t, err := renameio.NewPendingFile(filename, renameio.WithPermissions(perm), renameio.WithExistingPermissions())
	if err != nil {
		return err
	}
	defer t.Cleanup()
	if _, err := io.Copy(t, r); err != nil {
		return err
	}
	return t.CloseAtomicallyReplace()
}
