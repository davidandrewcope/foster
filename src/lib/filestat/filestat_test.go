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

func TestAppendReference(t *testing.T) {

	filestat1 := FileStat{nil, "/test/file/1", "", []FileStat{}}
	filestat2 := FileStat{nil, "/test/file/2", "", []FileStat{}}
	
	filestat1.AppendReference(filestat2);
	
	if len(filestat1.ReferencingFiles) < 1 {
		t.Errorf("Expected filestat1.ReferencingFiles longer than 1 , was: %v", len(filestat1.ReferencingFiles))
	}
	
	if (filestat1.ReferencingFiles[0].Path != "/test/file/2") {
		t.Errorf("Expected filestat1.ReferencingFiles[0].Path == '/test/file/2', was: %v", filestat1.ReferencingFiles[0].Path)
	}
	
	
}