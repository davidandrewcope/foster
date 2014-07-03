package projectree

import (
	"lib/FileStat"
	"strings"
	"os"
	"path/filepath"
	"io/ioutil"
	"net/http"
)

type ProjectTree struct {
	BasePath               string
	SourceFiles            []filestat.FileStat
	ExcludedSourcePatterns []string
}


/**
 * Reject files matching the exclusion list
 */
func (pt *ProjectTree) isExcludedPath(fileName string) bool {
	for _, pat := range pt.ExcludedSourcePatterns {
		
		if (len(pat) < 1) { //Ignore empty patterns
			return false
		} 
		//TODO: add some robustness here... first convert to an os.FileInfo then get the basename and then check the pattern
		if strings.Contains(fileName, pat) {
			return true
		}
	}
	return false
}

/**
 * Append sorce file if they meet our criteria, this is called by a directory walker
 */
func (pt *ProjectTree) appendSourceFile(fileName string, f os.FileInfo, err error) error {

	if err != nil {
		return err
	}
	
	if f.IsDir() {
		return nil
	}

	if pt.isExcludedPath(fileName) {
		return nil
	}

	pt.SourceFiles = append(pt.SourceFiles, filestat.FileStat{f, fileName, "", []filestat.FileStat{}})

	return nil
}

/**
 * Append sorce file if they meet our criteria, this is called by a directory walker
 */
func (pt *ProjectTree) CheckUsage(filePath string, f os.FileInfo, err error) error {

	//bar.Increment() //TODO: need a way to call this, ideas?
	
	if f.IsDir() {
		return nil
	}

	if pt.isExcludedPath(filePath) {
		return nil
	}

	// read the whole file
	b, err := ioutil.ReadFile(filePath)
	if err != nil {
		panic(err)
	}

	filestat := filestat.FileStat{f, filePath, "", []filestat.FileStat{}}

	filestat.ContentType = http.DetectContentType([]byte(b))

	//Only search text type files
	if !strings.Contains(filestat.ContentType, "text") {

		return nil
	}

	//Check for matches in this chunk from our sourceFiles list
	// Anytime we find a reference to a file, add it to our referencedFiles slice
	for _, sourceFile := range pt.SourceFiles {
		if strings.Contains(string(b), sourceFile.FileInfo.Name()) {
			sourceFile.AppendReference(filestat);
			
			//fmt.Printf("Matched string: %s in file: %s\n", sourceFile, filePath)
			//return nil  don't break on finding a reference
		}
	}

	return nil
}




func (pt *ProjectTree) BuildSourceFileList() error {
	return filepath.Walk(pt.BasePath, pt.appendSourceFile)
}

func New(basePath string, ignoreFolderArgs string) *ProjectTree {
	pt := new(ProjectTree)
	pt.BasePath = basePath;
	pt.ExcludedSourcePatterns = []string{};
	
	//Cleanup formatting of ignored folders
	for _, ignoreFoldersArg := range strings.Split(ignoreFolderArgs, ",") {
		pt.ExcludedSourcePatterns = append(pt.ExcludedSourcePatterns, strings.TrimSpace(ignoreFoldersArg))
	}
	
	pt.SourceFiles = []filestat.FileStat{};
	
	
	return pt;
}
