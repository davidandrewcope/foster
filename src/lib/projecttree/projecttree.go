package projectree

import (
	"lib/filestat"
	"strings"
)

type ProjectTree struct {
	BasePath               string
	SourceFiles            []filestat.FileStat
	ExcludedSourcePatterns []string
}

func (pt *ProjectTree) IsExcludedPath(fileName string) bool {
	for _, pat := range pt.ExcludedSourcePatterns {
		//TODO: add some robustness here... first convert to an os.FileInfo then get the basename and then check the pattern
		if strings.Contains(fileName, pat) {
			return true
		}
	}
	return false
}

func (pt *ProjectTree) AppendSourceFile(fileName string, f os.FileInfo, err error) error {
	if f.IsDir() {
		return nil
	}

	if pt.IsExcludedPath(fileName) {
		return nil
	}

	pt.SourceFiles = append(pt.SourceFiles, filestat.FileStat{f, fileName, ""})

	return nil
}
