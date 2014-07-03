package filestat

import (
	"testing"
	"syscall"
	"io/ioutil"
)

func TestFileStatStringPrintsPath(t *testing.T) {
	
	tempFile, err := ioutil.TempFile("", "TestFileStatStringPrintsPath.tmp")
	ioutil.WriteFile(tempFile.Name(), []byte("Just a test file\n"), 0644)
	defer syscall.Unlink(tempFile.Name())
	
	if err != nil {
		t.Fatalf("Cannot create temporary file for testing.")
	}
	
	tempFileInfo, err := tempFile.Stat()
	
	if err != nil {
		t.Fatalf("Cannot stat temporary file for testing.")
	}
	
	
	filestat := FileStat{tempFileInfo, tempFile.Name(), "", []FileStat{}}
	
	if len(filestat.String()) < 1 {
		t.Errorf("Expected FileStat.String() longer than 1 letter , was: %v", len(filestat.String()))
	}
}