// This file contains the entire program for create-gp4 since it's relatively trivial.

package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"time"
	"sort"
	"path/filepath"
	"encoding/xml"
)

// errorExit function will print the given formatted error to stdout and exit immediately after.
func errorExit(format string, params ...interface{}) {
	fmt.Printf(format, params...)
	os.Exit(-1)
}

type Dir struct {
	XMLName xml.Name  `xml:"dir"`
	TargName string   `xml:"targ_name,attr"`
	Dirs []Dir        `xml:"dir"`
}

type Rootdir struct {
	XMLName xml.Name  `xml:"rootdir"`
	Dirs []Dir        `xml:"dir"`
}

// check if slice contains specified string
func contains(s []string, e string) bool {
	for _, a := range s {
		if strings.Contains(a, e) {
			return true
		}
	}
	return false
}

// build rootdir tag
func buildRootDirTag(files []string) string {
	var paths []string
	var rootDir Rootdir;

	// sort files paths by len (to remove duplicate paths later)
	sort.Slice(files, func(i, j int) bool {
		return len(files[i]) > len(files[j])
	})

	// remove duplicate paths
	for _, file := range files {
		if file != "" && strings.Contains(file, "/") {
			if !contains(paths, filepath.Dir(file)) {
				paths = append(paths, filepath.Dir(file))
			}
		}
	}

	// parse paths
	for _, path := range paths {
		split := strings.Split(path, "/")
		var dir = Dir{TargName: split[0]}
		var dirPtr *Dir = &dir
		// parse childs paths
		for i := 1; i<len(split); i++ {
			dirPtr.Dirs = append(dirPtr.Dirs, Dir{TargName: split[i]})
			dirPtr = &dirPtr.Dirs[len(dirPtr.Dirs)-1]
		}
		rootDir.Dirs = append(rootDir.Dirs, dir)
	}

	out, _ := xml.MarshalIndent(rootDir, "\t", "\t")
	return string(out)
}

// build file list from path
func getFileList(filesPath string) []string {
	var files[] string

	// be sure path ends with a slash for strings.Replace
	if !strings.HasSuffix(filesPath, "/") {
		filesPath += "/"
	}

	// add files recursively
	filepath.Walk(filesPath, func(path string, _ os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		files = append(files, strings.Replace(path, filesPath, "", -1))
		return nil
	})

	return files
}

// parseFilesToTags takes a list of files as a space-deliminated string and parses it into a list of tags for the GP4 XML.
// Returns the list of XML tags for the files.
func parseFilesToTags(files []string) []string {
	var fileTags []string

	for _, file := range files {
		if file != "" {
			fileTags = append(fileTags, fmt.Sprintf("\t\t<file targ_path=\"%s\" orig_path=\"%s\" />", file, file))
		}
	}

	return fileTags
}

// createGP4 takes a set of values and constructs a .gp4 file and writes it to the given path. Returns an error if creation
// failed, nil otherwise.
func createGP4(path string, contentID string, files string, filesPath string) error {
	var fileList []string

	if files != "" {
		fileList = strings.Split(files, " ")
	} else {
		fileList = getFileList(filesPath)
	}
	fileTagList := parseFilesToTags(fileList)
	rootDir := buildRootDirTag(fileList)
	fileTags := strings.Join(fileTagList, "\n")
	currentTime := time.Now().Format("2006-01-02 15:04:05")

	gp4Contents := fmt.Sprintf("<?xml version=\"1.0\"?>\n"+
		"<psproject xmlns:xsd=\"http://www.w3.org/2001/XMLSchema\" xmlns:xsi=\"http://www.w3.org/2001/XMLSchema-instance\" fmt=\"gp4\" version=\"1000\">\n"+
		"\t<volume>\n"+
		"\t\t<volume_type>pkg_ps4_app</volume_type>\n"+
		"\t\t<volume_id>PS4VOLUME</volume_id>\n"+
		"\t\t<volume_ts>%s</volume_ts>\n"+
		"\t\t<package content_id=\"%s\" passcode=\"00000000000000000000000000000000\"\n"+
		"\t\t\tstorage_type=\"digital50\" app_type=\"full\" />\n"+
		"\t\t<chunk_info chunk_count=\"1\" scenario_count=\"1\">\n"+
		"\t\t\t<chunks>\n"+
		"\t\t\t\t<chunk id=\"0\" layer_no=\"0\" label=\"Chunk #0\" />\n"+
		"\t\t\t</chunks>\n"+
		"\t\t\t<scenarios default_id=\"0\">\n"+
		"\t\t\t\t<scenario id=\"0\" type=\"sp\" initial_chunk_count=\"1\" label=\"Scenario #0\">0</scenario>\n"+
		"\t\t\t</scenarios>\n"+
		"\t\t</chunk_info>\n"+
		"\t</volume>\n"+
		"\t<files img_no=\"0\">\n"+
		"%s"+
		"\n\t</files>\n"+
		"%s\n"+
		"</psproject>\n", currentTime, contentID, fileTags, rootDir)

	return ioutil.WriteFile(path, []byte(gp4Contents), 0644)
}

func main() {
	// Required flags
	outputFilePathPtr := flag.String("out", "homebrew.gp4", "`output gp4` to write to")
	contentIDPtr := flag.String("content-id", "", "content ID of the package")
	contentFilesPtr := flag.String("files", "", "list of files to pack into the package")
	contentPathPtr := flag.String("path", "", "path to files to pack into the package")

	flag.Parse()

	outputFilePath := *outputFilePathPtr
	contentID := *contentIDPtr
	contentFiles := *contentFilesPtr
	contentPath := *contentPathPtr

	if contentID == "" {
		errorExit("Content ID not specified, try -content-id=[content ID]\n")
	}

	if contentFiles == "" && contentPath == "" {
		errorExit("Content files or path not specified, try -files=\"[files, separated by spaces]\" or -path=\"[path/to/files]\"\n")
	}

	if err := createGP4(outputFilePath, contentID, contentFiles, contentPath); err != nil {
		errorExit("Error writing GP4: %s\n", err.Error())
	}
}
