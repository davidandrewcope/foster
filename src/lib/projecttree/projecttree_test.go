package projectree

import (
	"testing"
	"syscall"
	"os"
	"io/ioutil"
	"path/filepath"
)

func TestIsExcludedPath(t *testing.T) {
	t.Skipf("TestIsExcludedPath is not yet implemented")
}

func TestAppendSourceFileAppendsTheSourceFile(t *testing.T) {
	projectTree := New(".", "");
	tempFile, err := ioutil.TempFile("", "TestAppendSourceFileAppendsTheSourceFile.tmp")
	ioutil.WriteFile(tempFile.Name(), []byte("Just a test file\n"), 0644)
	defer syscall.Unlink(tempFile.Name())
	
	if err != nil {
		t.Fatalf("Cannot create temporary file for testing.")
	}
	
	tempFileInfo, err := tempFile.Stat()
	
	if err != nil {
		t.Fatalf("Cannot stat temporary file for testing.")
	}
	
	returnVal := projectTree.appendSourceFile(tempFile.Name(), tempFileInfo, nil)
	t.Logf("appendSourceFile returned %v", returnVal)
	
	if err != nil {
		t.Fatalf("appendSourceFile returned error: %v", err)
	}
	
	t.Logf("Appended test file %v", tempFile.Name())
	
	//Assertions
	
	if len(projectTree.SourceFiles) != 1 {
		t.Errorf("Expected SourceFiles == 1, was: %v", len(projectTree.SourceFiles))
	}
	
}


func TestCheckUsage(t *testing.T) {
	tempDir, _ := ioutil.TempDir("", "ProjectTreeTest")
	defer os.RemoveAll(tempDir)
	
	projectTree := New(tempDir, "");
	
	//Create the first test file
	tempFile1, err := ioutil.TempFile(tempDir, "TestCheckUsageTestFile1.tmp")
	ioutil.WriteFile(tempFile1.Name(), []byte("Nothing to see here, carry on...\n"), 0644)
	defer syscall.Unlink(tempFile1.Name())
	
	if err != nil {
		t.Fatalf("Cannot create temporary file for testing.")
	}
	
	tempFileInfo1, err := tempFile1.Stat()
	
	if err != nil {
		t.Fatalf("Cannot stat temporary file for testing.")
	}
	
	projectTree.appendSourceFile(tempFile1.Name(), tempFileInfo1, nil)

	if err != nil {
		t.Fatalf("appendSourceFile returned error: %v", err)
	}

	
	//Create the second test file
	tempFile2, err := ioutil.TempFile(tempDir, "TestCheckUsageTestFile2.tmp")
	ioutil.WriteFile(tempFile2.Name(), []byte(tempFile1.Name()), 0644)//Write the name of file1 to file2
	t.Logf("Wrote this to tempFile2\n%v\n", tempFile1.Name())
	defer syscall.Unlink(tempFile2.Name())
	
	if err != nil {
		t.Fatalf("Cannot create temporary file for testing.")
	}
	
	tempFileInfo2, err := tempFile2.Stat()
	
	if err != nil {
		t.Fatalf("Cannot stat temporary file for testing.")
	}
	
	projectTree.appendSourceFile(tempFile2.Name(), tempFileInfo2, nil)

	if err != nil {
		t.Fatalf("appendSourceFile returned error: %v", err)
	}
	
	if len(projectTree.SourceFiles) != 2 {
		t.Fatalf("Test setup failed, expected len(projectTree.SourceFiles) == 2 but was: %v", len(projectTree.SourceFiles))
	}
	
	filepath.Walk(projectTree.BasePath, projectTree.CheckUsage)
	
	//Assertions
	// 
	// t.Logf("len(projectTree.SourceFiles[0].ReferencingFiles) = %v", len(projectTree.SourceFiles[0].ReferencingFiles))
// 	t.Logf("len(projectTree.SourceFiles[1].ReferencingFiles) = %v", len(projectTree.SourceFiles[1].ReferencingFiles))

	if len(projectTree.SourceFiles[1].ReferencingFiles) != 0 {
		t.Errorf("Expected projectTree.SourceFiles[0].ReferencingFiles to be 0 but was: %v", len(projectTree.SourceFiles[1].ReferencingFiles))
	}
	
	if len(projectTree.SourceFiles[0].ReferencingFiles) < 1 {
		t.Errorf("Expected projectTree.SourceFiles[0].ReferencingFiles longer than 1 , was: %v", len(projectTree.SourceFiles[0].ReferencingFiles))
	} else if (projectTree.SourceFiles[0].ReferencingFiles[0].Path != tempFile2.Name()) {
		t.Errorf("Expected projectTree.SourceFiles[0].ReferencingFiles[0].Path == %s, was: %v", tempFile2.Name(), projectTree.SourceFiles[0].ReferencingFiles[0].Path)
	}
	
}


func TestBuildSourceFileList(t *testing.T) {
	t.Skipf("TestBuildSourceFileList is not yet implemented")
}

func TestNew(t *testing.T) {
	t.Skipf("TestNew is not yet implemented")
}