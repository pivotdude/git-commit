package prompt

import (
	"fmt"
	"git-commit/diff"
	"git-commit/utils"
	"os"
	"path/filepath"
	"strings"
)

// loadCustomPrompt loads a custom prompt from the custom-instructions folder
// processDirectory recursively processes all files in a directory
func processDirectory(dirPath string) (string, error) {
	var result []string

	err := filepath.Walk(dirPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		
		// Skip the root directory itself
		if path == dirPath {
			return nil
		}

		if info.IsDir() {
			// For directories, add directory tag
			relPath, err := filepath.Rel(dirPath, path)
			if err != nil {
				return err
			}
			result = append(result, fmt.Sprintf("<directory name=\"%s\" path=\"%s\">\n", info.Name(), relPath))
		} else {
			// For files, add context tag with content
			content, err := utils.ReadFileContent(path)
			if err != nil {
				return fmt.Errorf("error reading file %s: %v", path, err)
			}
			relPath, err := filepath.Rel(dirPath, path)
			if err != nil {
				return err
			}
			result = append(result, fmt.Sprintf("<context file=\"%s\">\n%s\n</context>\n", relPath, content))
		}
		return nil
	})

	if err != nil {
		return "", err
	}

	return strings.Join(result, "\n"), nil
}

func loadCustomPrompt(promptName string) (string, error) {
	// Construct the path to the custom prompt file
	customPromptPath := fmt.Sprintf(".git-commit/custom-instructions/%s.md", promptName)
	
	// Check if the file exists
	if _, err := os.Stat(customPromptPath); os.IsNotExist(err) {
		return "", fmt.Errorf("custom prompt file '%s' not found", customPromptPath)
	}

	// Read the file content
	content, err := os.ReadFile(customPromptPath)
	if err != nil {
		return "", fmt.Errorf("failed to read custom prompt file '%s': %v", customPromptPath, err)
	}

	return string(content), nil
}

// ProcessMarkdownDirectives processes special directives in markdown content
func ProcessMarkdownDirectives(content string) (string, error) {
	lines := strings.Split(content, "\n")
	var result []string

	for _, line := range lines {
		if strings.Contains(line, "@context:") {
			// Extract file path after @context:
			filePath := strings.TrimSpace(strings.Split(line, "@context:")[1])
			
			// Get file info to check if it's a directory
			fileInfo, err := os.Stat(filePath)
			if err != nil {
				return "", fmt.Errorf("error getting file info for %s: %v", filePath, err)
			}

			if fileInfo.IsDir() {
				// Handle directory recursively
				dirContent, err := processDirectory(filePath)
				if err != nil {
					return "", fmt.Errorf("error processing directory %s: %v", filePath, err)
				}
				replacement := fmt.Sprintf("<directory name=\"%s\" path=\"%s\">\n%s</directory>", 
					fileInfo.Name(), filePath, dirContent)
				result = append(result, replacement)
			} else {
				// Handle single file
				fileContent, err := utils.ReadFileContent(filePath)
				if err != nil {
					return "", fmt.Errorf("error reading context file %s: %v", filePath, err)
				}
				replacement := fmt.Sprintf("<context file=\"%s\">\n%s\n</context>", filePath, fileContent)
				result = append(result, replacement)
			}
		} else if strings.Contains(line, "@diff") {
			diffOutput := diff.GetDiffOutputWithoutIgnoresFiles()
			result = append(result, diffOutput)
		} else {
			result = append(result, line)
		}
	}

	return strings.Join(result, "\n"), nil
}

// GetAIPrompt returns the AI prompt (standard or custom) with context files processed
func GetAIPrompt(promptName string) string {
	var rawPrompt string
	
	// If a specific prompt name is provided, try to load it from custom-instructions
	if promptName != "" {
		customPrompt, err := loadCustomPrompt(promptName)
		if err != nil {
			fmt.Printf("Error reading custom prompt '%s': %v, using standard\n", promptName, err)
			rawPrompt = ""
		} else if strings.TrimSpace(customPrompt) != "" {
			rawPrompt = customPrompt
		} else {
			rawPrompt = ""
		}
	}
	
	processedPrompt, err := ProcessMarkdownDirectives(rawPrompt)
	if err != nil {
		fmt.Printf("Error processing prompt directives: %v\n", err)
		return rawPrompt
	}
	
	return processedPrompt
}
