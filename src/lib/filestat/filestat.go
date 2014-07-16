package filestat

import (
	"fmt"
	"os"
)

/**
 * Extend FileInfo to maintain the full Path, ContentType, and any Referencing files 
 */
type FileStat struct {
	os.FileInfo
	Path        string
	ContentType string
	ReferencingFiles []FileStat

}

/**
 * We typically only need to print the path of a file
 */
func (f *FileStat) String() string {
	return fmt.Sprintf("%v", f.Path)
}

/**
 * Determine if this file is in the given slice of other files
 */
func (a *FileStat) NameInSlice(list []FileStat) bool {
	for _, b := range list {
		if b.FileInfo.Name() == a.FileInfo.Name() {
			return true
		}
	}
	return false
}

/**
 * Append the given FileStat to the ReferencingFiles for this FileStat
 */
func (f *FileStat) AppendReference(i FileStat) bool {
	for _, ele := range f.ReferencingFiles {
		if ele.Path == i.Path {
			return false
		}
	}
	f.ReferencingFiles = append(f.ReferencingFiles, i)
	return true
}

