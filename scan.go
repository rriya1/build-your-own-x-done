package main

import(
	"fmt"
	"strings"
	"log"
	"os"
	"bufio"
	"io"
	"os/user"
)

func scan(folder string) {
	print("Found folders:\n\n")
	repositories := recursiveScanFolders(folder)
	filePath := fileDotPath()
	addNewSliceElementsToFile(filePath, repositories)
	fmt.Printf("\n\nSucessfully added\n\n")
}

func recursiveScanFolders(folder string) []string{
	return scanGitFolders(make([]string, 0), folder)
}

func scanGitFolders(folders []string, folder string) []string{
	folder = strings.TrimSuffix(folder, "/")

	f, err := os.Open(folder)
	if err != nil{
		log.Fatal("error in opening the folder", err)
	}

	files, err := f.Readdir(-1)
	f.Close()
	if err != nil{
		log.Fatal("cannot read files in the folder", err)
	}

	var path string

	for _, file := range(files){
		if file.IsDir(){
			path = folder + "/" + file.Name()
			if file.Name() == ".git"{
				path = strings.TrimSuffix(path, "/.git")
				fmt.Print(path)
				folders = append(folders, path)
				continue
			}

			if file.Name() == ".vendor" || file.Name() == ".gitignore" || file.Name() == "node_modules" || file.Name() == ".claude"{
				continue
			}

			folders =  scanGitFolders(folders, path)
		}
	} 
	return folders
}

func fileDotPath() string{
	usr, err := user.Current()
	
	if err != nil{
		log.Fatal("cannot get the current session's user", err)
	}
	dotFile := usr.HomeDir + "/.gogitlocalstats"

	return dotFile
}

func addNewSliceElementsToFile(filepath string, newRepos []string){
	existingRepos := parseFilePathToSlice(filepath)
	repos := joinSlices(newRepos, existingRepos)
	dumpStringSliceToFiles(repos, filepath)
}

func parseFilePathToSlice(filepath string) []string{
	f := openFile(filepath)
	defer f.Close()

	var lines []string
	scanner := bufio.NewScanner(f)

	for scanner.Scan(){
		lines = append(lines, scanner.Text())
	}

	if err := scanner.Err(); err != nil{
		if err != io.EOF{
			panic(err)
		}
	}

	return lines
}

func openFile(filepath string) *os.File{
	f, err := os.OpenFile(filepath, os.O_APPEND|os.O_RDWR, 0755)
	if err != nil{
		if os.IsNotExist(err){
			f, err := os.Create(filepath)
			if err != nil{
				panic(err)
			}
		}

		panic(err)
	}
	return f
}

func joinSlices(newrepos []string, existingrepos []string) []string{
	for _,i := range newrepos{
		if !sliceContains(i, existingrepos){
			existingrepos = append(existingrepos, i)
		}
	}
	return existingrepos
}

func sliceContains(toCheckVal string, toCheckIn []string) bool{
	for _,i := range toCheckIn{
		if i == toCheckVal{
			return true
		}
	}
	return false
}

func dumpStringSliceToFiles(repos []string, filePath string){
	content := strings.Join(repos, "\n")
	os.WriteFile(filePath, []byte(content), 0755)
}