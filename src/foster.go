package main

import "fmt"
import "flag"
import "os"
import "path/filepath"

import "lib/ProjectTree"


//Declare the ProjectTree
var projectTree *projectree.ProjectTree

func main() {
	
	
	//Collect the arguments
	basePath 		:= flag.String("base", "", "The root folder of your project tree")
	ignoreFolders 	:= flag.String("ignoreFolders", "", "CSV list of folders to ignore")
	showUsed 		:= *flag.Bool("used", false, "Show used files instead of unused files")

	flag.Parse()
	
	projectTree = projectree.New(*basePath, *ignoreFolders);

	//Reject and empty basePath
	if len(projectTree.BasePath) < 1 {
		flag.PrintDefaults()
		os.Exit(2)
	}

	fmt.Printf("Crawling directory: %s\n", projectTree.BasePath)
	err := projectTree.BuildSourceFileList()
	
	if err != nil {
		fmt.Printf("Unable to crawl directory.\n Exit %v\n", err)
		os.Exit(1)
	}


	fmt.Printf("Found %d files, searching for usages.\n", len(projectTree.SourceFiles))

	//The parsing is done via a channel
	done := make(chan bool)
	go func() {
		err := filepath.Walk(projectTree.BasePath, projectTree.CheckUsage)

		if err != nil {
			fmt.Printf("Unable to read file.\n Exit %v\n", err)
			os.Exit(2)
		}

		close(done)
	}()

	<-done

	

	fmt.Printf("\n\n")
	if showUsed {
		projectTree.PrintUsedFiles();
	} else {
		projectTree.PrintUnUsedFiles();
	}


}
