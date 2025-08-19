package main

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"

	"git-commit/git"
	"git-commit/prompt"
	"git-commit/utils"
)

func main() {
	// 1. Read the .gitdiffignore file and get ignore patterns
	patterns, err := git.ParseGitDiffIgnore()
	if err != nil {
		log.Fatalf("Error reading .gitdiffignore: %v", err)
	}
	// Log the obtained patterns
	log.Printf("Obtained patterns from .gitdiffignore: %v", patterns)

	// 2. Get the list of staged files
	stagedFiles, err := git.GetStagedFiles()
	if err != nil {
		log.Fatalf("Error getting staged files: %v", err)
	}

// Log the list of staged files
	log.Printf("Found staged files: %v", stagedFiles)

	// 3. Determine files that need to be ignored
	filesToIgnore := git.GetFilesToIgnore(patterns, stagedFiles)
	log.Printf("Files to ignore: %v", filesToIgnore)

	// 4. Remove ignored files from staged
	if len(filesToIgnore) > 0 {
		fmt.Printf("Removing %d files from staged to ignore:\n", len(filesToIgnore))
		for _, file := range filesToIgnore {
			fmt.Printf("  - %s\n", file)
		}

		err = git.RemoveFilesFromStaged(filesToIgnore)
		if err != nil {
			log.Fatalf("Error removing files from staged: %v", err)
		}
	}

	// 5. Perform main work - get git diff
	cmd := exec.Command("git", "diff", "--staged")
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err = cmd.Run()

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

	// 6. Get AI prompt (standard or custom)
	aiPrompt := prompt.GetAIPrompt()

	// 7. Copy the result to clipboard (first diff, then prompt)
	finalPrompt := fmt.Sprintf("<diff>\n%s\n</diff>\n\n%s", diffOutput, aiPrompt)
	utils.CopyToClipboard(finalPrompt)
	fmt.Println("AI prompt with git diff copied to clipboard.")

	// 8. Return ignored files back to staged
	if len(filesToIgnore) > 0 {
		fmt.Printf("Returning %d files to staged:\n", len(filesToIgnore))
		for _, file := range filesToIgnore {
			fmt.Printf("  + %s\n", file)
		}

		err = git.AddFilesToStaged(filesToIgnore)
		if err != nil {
			log.Printf("Error returning files to staged: %v", err)
		} else {
			fmt.Println("Files successfully returned to staged.")
		}
	}
}
