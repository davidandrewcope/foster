package main

import "fmt"
import "flag"
import "os"
import "path/filepath"
import "strings"
import "net/http"
import "io/ioutil"
import "github.com/cheggaaa/pb"


var sourceFiles []FileStat

var referencedFiles []FileStat

var bar *pb.ProgressBar

var ignoreFolders []string

type FileStat struct {
	os.FileInfo
	path string
	contentType string
}

func (f *FileStat) String() string {
	return fmt.Sprintf("%s", f.path)
}

func appendIfMissing(slice []FileStat, i FileStat) []FileStat {
	for _, ele := range slice {
		if ele.path == i.path {
			return slice
		}
	}
	return append(slice, i)
}

func isExcludedPath(filePath string) bool{
	return strings.Contains( filePath, ".git")
}

func addToSourceList(fileName string, f os.FileInfo, err error) error {
	if f.IsDir() {
		return nil
	}

	if isExcludedPath(fileName) {
		return nil
	}

	sourceFiles = append(sourceFiles, FileStat{f, fileName, ""})

	return nil
}

func checkUsage(filePath string, f os.FileInfo, err error) error {

	bar.Increment();

	if f.IsDir() {
		return nil
	}

	if isExcludedPath(filePath) {
		return nil
	}

	//Skip any folders passed in the ignoreFolders args
	for _, ignoredFolder := range ignoreFolders {
		if (strings.Contains(filePath, "/" + ignoredFolder )) {
			return nil;
		}
	}

	// read the whole file
	b, err := ioutil.ReadFile(filePath)
	if err != nil { panic(err) }

	fileStat := FileStat{f, filePath, ""}

	// This should be pretty efficient since we only sniff the first chunk
	fileStat.contentType = http.DetectContentType([]byte(b))

	//Only search text type files
	if ( !strings.Contains( fileStat.contentType , "text") ) {

		return nil
	}

	//Check for matches in this chunk from our sourceFiles list
	// Anytime we find a reference to a file, add it to our referencedFiles slice
	for _, sourceFile := range sourceFiles {
		if strings.Contains(string(b), sourceFile.FileInfo.Name()) {
			referencedFiles = appendIfMissing(referencedFiles, fileStat)
			//fmt.Printf("Matched string: %s in file: %s\n", sourceFile, filePath)
			//return nil  don't break on finding a reference
		}
	}

	return nil
}

func fileInSlice(a FileStat, list []FileStat) bool {
	for _, b := range list {
		if b.FileInfo.Name() == a.FileInfo.Name() {
			return true
		}
	}
	return false
}

func main() {


	//Collect the arguments
	basePath := flag.String("base", ".", "Site root folder")
	showUsed := *flag.Bool("used", false, "Show used files instead of unused files")
	ignoreFoldersArgs := strings.Split(*flag.String("ignoreFolders", ".", "CSV list of folders to skip"), ",")

	//Cleanup formatting of ignored folders
	for _, ignoreFoldersArg := range ignoreFoldersArgs {
		ignoreFolders = append(ignoreFolders, strings.TrimSpace(ignoreFoldersArg))
	}

	flag.Parse()

	fmt.Printf("Crawling directory: %s\n", *basePath)


	err := filepath.Walk(*basePath, addToSourceList)

	if err != nil {
		fmt.Printf("Unable to crawl directory.\n Exit %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("\n")

	bar = pb.StartNew(len(sourceFiles))
	bar.ShowTimeLeft = false

	fmt.Printf("Found %d files, searching for usages.\n", len(sourceFiles))


	//The parsing is done via a channel
	done := make(chan bool)
	go func () {
		err := filepath.Walk(*basePath, checkUsage)

		if err != nil {
			fmt.Printf("Unable to read file.\n Exit %v\n", err)
			os.Exit(2)
		}

		close(done)
	}()

	<-done

	bar.Finish()

	fmt.Printf("\n\n")


	if (!showUsed) {
		fmt.Printf("Unused Files: \n")

		for _, sourceFile := range sourceFiles {
			if !fileInSlice(sourceFile, referencedFiles) {
				fmt.Printf("%s\n", sourceFile.path)
			}
		}
	}

	if (showUsed) {
		fmt.Printf("Used Files: \n")


		for _, referencedFile := range referencedFiles {
			fmt.Printf("%s\n", referencedFile.path)
		}
	}

}
