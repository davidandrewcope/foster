package projectree

import (
	"fmt"
	"testing"
	"syscall"
	"io/ioutil"
)

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
	
	for _, sourceFile := range projectTree.SourceFiles {
		
		fmt.Printf("%s\n", sourceFile.Path)
	}
	
	//fmt.Print("%v", projectTree)
	
}