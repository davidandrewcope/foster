package main

import "fmt"
import "flag"
import "os"
import "path/filepath"
import "github.com/cheggaaa/pb"

import "lib/ProjectTree"


var bar *pb.ProgressBar


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

	fmt.Printf("\n")
	
	

	bar = pb.StartNew(len(projectTree.SourceFiles))
	bar.ShowTimeLeft = false

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

	bar.Finish()

	fmt.Printf("\n\n")

	if !showUsed {
		fmt.Printf("Unused Files: \n")

		for _, sourceFile := range projectTree.SourceFiles {
			if (len(sourceFile.ReferencingFiles) > 0) {
			//if !(sourceFile.NameInSlice(referencedFiles)) {
				fmt.Printf("%s\n", sourceFile.Path)
			}
		}
	}

	if showUsed {
		fmt.Printf("Used Files: \n")
		for _, sourceFile := range projectTree.SourceFiles {
			if (len(sourceFile.ReferencingFiles) == 0) {
				fmt.Printf("%s\n", sourceFile.Path)
			}
		}
		
		//for _, referencedFile := range projectTree referencedFiles {
		//	fmt.Printf("%s\n", referencedFile.Path)
		//}
	}
	
	for _, sourceFile := range projectTree.SourceFiles {
		
		
			fmt.Printf("%v\n", sourceFile)
		
	}

}
