package help

import (
	"fmt"
	"os"
	"strings"
)

// ShowHelp displays the help message for git-commit
func ShowHelp() {
	fmt.Println("git-commit - AI-powered git commit message generator")
	fmt.Println()
	fmt.Println("Usage:")
	fmt.Println("  git-commit [prompt-name]  Generate AI prompt with git diff and copy to clipboard")
	fmt.Println("  git-commit -h             Show this help message")
	fmt.Println("  git-commit -v             Enable verbose output")
	fmt.Println("  git-commit -generate-prompt  Generate prompt for current changes")
	fmt.Println()
	fmt.Println("Examples:")
	fmt.Println("  git-commit              # Generate prompt and copy to clipboard")
	fmt.Println("  git-commit mark         # Use custom prompt from custom-instructions/mark.md")
	fmt.Println("  git-commit -v           # Run with verbose logging")
	fmt.Println("  git-commit -generate-prompt  # Generate prompt without copying to clipboard")
	fmt.Println()
	fmt.Println("Configuration:")
	fmt.Println("  Create .git-commit/ignore file to specify patterns to ignore")
	fmt.Println("  Create .git-commit/prompt.md file for default custom AI prompt")
	fmt.Println("  Create .git-commit/custom-instructions/ folder with .md files for custom prompts")
	fmt.Println("    Example: .git-commit/custom-instructions/mark.md")
	fmt.Println("    Usage: git-commit mark")
	fmt.Println()
	
	// Show available custom prompts
	customPrompts := GetAvailableCustomPrompts()
	if len(customPrompts) > 0 {
		fmt.Println("Available custom prompts:")
		for _, promptName := range customPrompts {
			fmt.Printf("  git-commit %s\n", promptName)
		}
	} else {
		fmt.Println("No custom prompts found. Create .git-commit/custom-instructions/ folder with .md files.")
	}
}

// GetAvailableCustomPrompts scans the custom-instructions directory and returns available prompt names
func GetAvailableCustomPrompts() []string {
	var prompts []string
	
	// Check if custom-instructions directory exists
	customInstructionsPath := ".git-commit/custom-instructions"
	if _, err := os.Stat(customInstructionsPath); os.IsNotExist(err) {
		return prompts
	}
	
	// Read all .md files in the directory
	files, err := os.ReadDir(customInstructionsPath)
	if err != nil {
		return prompts
	}
	
	for _, file := range files {
		if !file.IsDir() && strings.HasSuffix(file.Name(), ".md") {
			// Remove .md extension to get the prompt name
			promptName := strings.TrimSuffix(file.Name(), ".md")
			prompts = append(prompts, promptName)
		}
	}
	
	return prompts
}