package filestat

import (
	"os"
	"fmt"
)

type FileStat struct {
	os.FileInfo
	Path string
	ContentType string
	//ReferencingFiles []FileStat
	
	
}


func (f *FileStat) String() string {
	return fmt.Sprintf("%s", f.Path)
}

func (a *FileStat) NameInSlice( list []FileStat) bool {
	for _, b := range list {
		if b.FileInfo.Name() == a.FileInfo.Name() {
			return true
		}
	}
	return false
}


