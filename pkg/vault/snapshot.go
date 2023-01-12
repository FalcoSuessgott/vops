package vault

import (
	"bytes"
	"io"
)

// SnapshotBackup takes a snapshot.
func (v *Vault) SnapshotBackup() (*bytes.Buffer, error) {
	var writer bytes.Buffer

	err := v.Client.Sys().RaftSnapshot(&writer)
	if err != nil {
		return &writer, err
	}

	return &writer, nil
}

// SnapshotBackup restores a snapshot.
func (v *Vault) SnapshotRestore(reader io.Reader, force bool) error {
	err := v.Client.Sys().RaftSnapshotRestore(reader, force)
	if err != nil {
		return err
	}

	return nil
}
