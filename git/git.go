package git

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"git-commit/utils"
)

// parseGitDiffIgnore reads the ignore file from .git-commit and returns a list of patterns to ignore
func ParseGitDiffIgnore() ([]string, error) {
	// Check if the ignore file exists in the .git-commit folder
	if _, err := os.Stat(".git-commit/ignore"); os.IsNotExist(err) {
		// If the file doesn't exist, return an empty list
		return []string{}, nil
	}

	// Open the file for reading
	file, err := os.Open(".git-commit/ignore")
	if err != nil {
		return nil, fmt.Errorf("failed to open file .git-commit/ignore: %v", err)
	}
	defer file.Close()

	var patterns []string
	scanner := bufio.NewScanner(file)
	
	// Read the file line by line
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		// Skip empty lines and comments
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		patterns = append(patterns, line)
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("error reading file .git-commit/ignore: %v", err)
	}

	return patterns, nil
}


// GetStagedFiles gets the list of files added to staged
func GetStagedFiles() ([]string, error) {
	cmd := exec.Command("git", "diff", "--staged", "--name-only")
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err := cmd.Run()

	if err != nil {
		return nil, fmt.Errorf("error getting staged files list: %v\n%s", err, stderr.String())
	}

	output := strings.TrimSpace(stdout.String())
	if output == "" {
		return []string{}, nil
	}

	// Split the output into lines
	files := strings.Split(output, "\n")
	return files, nil
}

// GetFilesToIgnore returns the list of files that should be ignored
func GetFilesToIgnore(patterns []string, files []string) []string {
	var ignoredFiles []string
	
	for _, file := range files {
		for _, pattern := range patterns {
			if utils.GlobMatch(pattern, file) {
				ignoredFiles = append(ignoredFiles, file)
				break // File matched the pattern, no need to check further
			}
		}
	}
	
	return ignoredFiles
}

// RemoveFilesFromStaged removes the specified files from staged
func RemoveFilesFromStaged(files []string) error {
	if len(files) == 0 {
		return nil
	}
	
	args := append([]string{"reset", "--"}, files...)
	cmd := exec.Command("git", args...)
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	err := cmd.Run()
	
	if err != nil {
		return fmt.Errorf("error removing files from staged: %v\n%s", err, stderr.String())
	}
	
	return nil
}

// AddFilesToStaged adds the specified files to staged
func AddFilesToStaged(files []string) error {
	if len(files) == 0 {
		return nil
	}
	
	args := append([]string{"add", "--"}, files...)
	cmd := exec.Command("git", args...)
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	err := cmd.Run()
	
	if err != nil {
		return fmt.Errorf("error adding files to staged: %v\n%s", err, stderr.String())
	}
	
	return nil
}