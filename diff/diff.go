package diff

import (
	"bytes"
	"fmt"
	"git-commit/git"
	"log"
	"os"
	"os/exec"
	"strings"
)

//
//	getFilesToIgnore reads the .gitdiffignore file, get staged files, and determines which files need to be ignored
//
func getFilesToIgnore() []string {
	// 1. Read the .gitdiffignore file and get ignore patterns
	patterns, err := git.ParseGitDiffIgnore()
	if err != nil {
		log.Printf("Error reading .gitdiffignore: %v", err)
	}

	// 2. Get the list of staged files
	stagedFiles, err := git.GetStagedFiles()
	if err != nil {
		log.Fatalf("Error getting staged files: %v", err)
	}

	// 3. Determine files that need to be ignored
	filesToIgnore := git.GetFilesToIgnore(patterns, stagedFiles)
	return filesToIgnore
}

func removeFilesFromStaged(filesToIgnore []string) {
		// 4. Remove ignored files from staged
	if len(filesToIgnore) > 0 {
		fmt.Printf("Ignoring %d files based on .gitdiffignore:\n", len(filesToIgnore))
		for _, file := range filesToIgnore {
			fmt.Printf("  - %s\n", file)
		}

		err := git.RemoveFilesFromStaged(filesToIgnore)
		if err != nil {
			log.Printf("Error removing files from staged: %v", err)
		}
	}
}


func addedFilesToStaged(filesToIgnore []string) []string {
	// 8. Return ignored files back to staged
	if len(filesToIgnore) > 0 {
		fmt.Printf("Returning %d files to staged:\n", len(filesToIgnore))
		for _, file := range filesToIgnore {
			fmt.Printf("  + %s\n", file)
		}

		err := git.AddFilesToStaged(filesToIgnore)
		if err != nil {
			log.Printf("Error returning files to staged: %v", err)
		} else {
			fmt.Println("Files successfully returned to staged.")
		}
	}

	return filesToIgnore
}

func parseGitDiff(filesToIgnore []string) string {
	// 5. Get git diff
	cmd := exec.Command("git", "diff", "--staged")
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err := cmd.Run()

	if err != nil {
		// If an error occurred, return the ignored files back to staged
		if len(filesToIgnore) > 0 {
			fmt.Println("An error occurred, returning files back to staged...")
			addErr := git.AddFilesToStaged(filesToIgnore)
			if addErr != nil {
				log.Printf("Error returning files to staged: %v", addErr)
			}
		}
		log.Fatalf("Error executing git diff: %v\n%s", err, stderr.String())
	}

	diffOutput := strings.TrimSpace(stdout.String())
	if diffOutput == "" {
		// If there are no changes, return the ignored files back to staged
		if len(filesToIgnore) > 0 {
			fmt.Println("No changes found, returning files back to staged...")
			addErr := git.AddFilesToStaged(filesToIgnore)
			if addErr != nil {
				log.Printf("Error returning files to staged: %v", addErr)
			}
		}
		fmt.Println("No changes detected. Please use 'git add' to add files to staged.")
		os.Exit(1)
	}
	return diffOutput
}

func GetDiffOutputWithoutIgnoresFiles() string {
	ignoredFiles := getFilesToIgnore()
	removeFilesFromStaged(ignoredFiles)
	diffOutput := parseGitDiff(ignoredFiles)
	addedFilesToStaged(ignoredFiles)

	return diffOutput
}