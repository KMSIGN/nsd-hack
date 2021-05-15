package loader

import (
	"errors"
	"nsd-hack/server/app/file"
)

func StartUploading(addr string, name string) error {
	if file.CheckFileExists(name) { return errors.New("no such file")}




	return nil
}
