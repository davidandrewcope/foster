package projectree

import (
	"lib/FileStat"
	"strings"
	"os"
	"path/filepath"
	"bufio"
	"io"
	"net/http"
	"fmt"
	
	"github.com/cheggaaa/pb"
)


type ProjectTree struct {
	BasePath               string
	SourceFiles            []filestat.FileStat
	ExcludedSourcePatterns []string
	ProgressBar				*pb.ProgressBar
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
	pt.ProgressBar.Increment() 
	
	if f.IsDir() {
		return nil
	}

	if pt.isExcludedPath(filePath) {
		return nil
	}

	//Open the file
	fi, err := os.Open(filePath)
	if err != nil {
		panic(err)
	}
	
	//Initialize a filestat struct, this may be discarded if we are not reading a text type file
	fileStat := filestat.FileStat{} 
	
	// close fi on exit and check for its returned error
    defer func() {
        if err := fi.Close(); err != nil {
            panic(err)
        }
    }()
	
	// make a read buffer
    r := bufio.NewReader(fi)
	
	// make a buffer to keep chunks that are read
    buf := make([]byte, 1024)
	
	// make a chunk slice to hold the all the chunks as string
	chunks := []string{}
	
    for i := 0; ; i++ {
        // read a chunk
        n, err := r.Read(buf)
        if err != nil && err != io.EOF { panic(err) }
        if n == 0 { break }
		
		chunks = append(chunks, string(buf[:n]) )
		
		// We should only need the first chunk to detect the contentType
		if i == 0 {
			contentType := http.DetectContentType([]byte(buf[:n]))
			
			//Only search text type files
			if !strings.Contains(contentType, "text") {
				//Break out here if we are not reading a text file
				return nil
			} 
			
			//File was text so create the FileStat object
			fileStat = filestat.FileStat{f, filePath, contentType, []filestat.FileStat{}}

		}
        
    }
		

	//Check for matches in this chunk from our sourceFiles list
	// Anytime we find a reference to a file, add it to our referencedFiles slice
	for i, _ := range pt.SourceFiles {
		if strings.Contains(strings.Join(chunks, ""), pt.SourceFiles[i].FileInfo.Name()) {
			pt.SourceFiles[i].AppendReference(fileStat);
			
			//fmt.Printf("Matched string: %s in file: %s\n", pt.SourceFiles[i].FileInfo.Name(), string(b))
			//return nil  don't break on finding a reference
		}
	}

	return nil
}

/**
 * Kick off the directory walking function to build up a list of source files in the project
 */
func (pt *ProjectTree) BuildSourceFileList() error {
	walkResult := filepath.Walk(pt.BasePath, pt.appendSourceFile)
	
	pt.ProgressBar = pb.StartNew(len(pt.SourceFiles))
	pt.ProgressBar.ShowTimeLeft = false
	
	return walkResult
}

/**
 * End the progress bar timer, and print the used files in the project
 */
func (pt *ProjectTree) PrintUnUsedFiles() {
	pt.ProgressBar.Finish()
	
	fmt.Printf("UnUsed Files: \n")
	for _, sourceFile := range pt.SourceFiles {
		if (len(sourceFile.ReferencingFiles) == 0) && !strings.Contains(sourceFile.ContentType, "text") {  //Need a flag to print text files
			fmt.Printf("%s\n", sourceFile.Path)
		}
	}
}

/**
 * End the progress bar timer, and print the UNused files in the project
 */
func (pt *ProjectTree) PrintUsedFiles() {
	pt.ProgressBar.Finish()
	
	fmt.Printf("Used Files: \n")
	for _, sourceFile := range pt.SourceFiles {
		if (len(sourceFile.ReferencingFiles) > 0) {
			fmt.Printf("%s\n", sourceFile.Path)
		}
	}
}

/**
 * Initialize a ProjectTree struct
 */
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
