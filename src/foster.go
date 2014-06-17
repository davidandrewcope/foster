package main

import "fmt"
import "flag"
import "os"
import "path/filepath"
import "strings"
import "net/http"
import "io/ioutil"
import "github.com/cheggaaa/pb"

var sourceFiles []string

var referencedFiles []string

var bar *pb.ProgressBar

var ignoreFolders []string


func appendIfMissing(slice []string, i string) []string {
	for _, ele := range slice {
		if ele == i {
			return slice
		}
	}
	return append(slice, i)
}

func isExcludedPath(path string) bool{
	return strings.Contains( path, ".git")
}

func addToSourceList(path string, f os.FileInfo, err error) error {
	if f.IsDir() {
		return nil
	}

	if isExcludedPath(path) {
		return nil
	}

	sourceFiles = append(sourceFiles, f.Name())

	return nil
}

func checkUsage(path string, f os.FileInfo, err error) error {

	bar.Increment();

	if f.IsDir() {
		return nil
	}

	if isExcludedPath(path) {
		return nil
	}

	//Skip any folders passed in the ignoreFolders args
	for _, ignoredFolder := range ignoreFolders {
		if (strings.Contains(path, "/" + ignoredFolder )) {
			return nil;
		}
	}

	// read the whole file
	b, err := ioutil.ReadFile(path)
	if err != nil { panic(err) }

	//Only search text type files
	// This should be pretty efficient since we only sniff the first chunk
	if ( !strings.Contains( http.DetectContentType( []byte(b) ), "text") ) {
		return nil
	}

	//Check for matches in this chuck from our sourceFiles list
	// Anytime we find a reference to a file, add it to our referencedFiles slice
	for _, sourceFile := range sourceFiles {
		if strings.Contains(string(b), sourceFile) {
			referencedFiles = appendIfMissing(referencedFiles, sourceFile)
			//fmt.Printf("Matched string: %s in file: %s\n", sourceFile, path)
			//return nil  don't break on finding a reference
		}
	}

	return nil
}

func stringInSlice(a string, list []string) bool {
	for _, b := range list {
		if b == a {
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
			if !stringInSlice(sourceFile, referencedFiles) {
				fmt.Printf("%s\n", sourceFile)
			}
		}
	}

	if (showUsed) {
		fmt.Printf("Used Files: \n")

		for _, referencedFile := range referencedFiles {
			fmt.Printf("%s\n", referencedFile)
		}
	}

}
