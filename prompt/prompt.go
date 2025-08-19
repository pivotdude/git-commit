package prompt

import (
	"fmt"
	"os"
	"strings"
)

const defaultAIPrompt = `You are a senior software engineer who specializes in Git. Your task is to generate a branch name and a commit message based on the provided code changes (git diff). Follow these rules:
1.  **Branch Name:**
    * Use a lowercase, hyphen-separated format (e.g., 'feature/add-login-button').
    * Start with a prefix indicating the type of change:
        * 'feature/' for new features.
        * 'fix/' for bug fixes.
        * 'hotfix/' for critical bug fixes.
        * 'refactor/' for code refactoring.
        * 'chore/' for routine tasks or maintenance.
        * 'docs/' for documentation changes.
    * The name should be concise and clearly describe the change.
2.  **Commit Message:**
    * Follow the Conventional Commits specification.
    * **Subject Line:**
        * Start with a type prefix (e.g., 'feat', 'fix', 'refactor', 'chore', 'docs', 'perf').
        * Use imperative, present tense: "Add," not "Added" or "Adds."
        * Keep it brief (under 50 characters).
        * Use a lowercase letter for the first word.
    * **Body (optional):**
        * Provide a more detailed explanation of the change. Explain the "what" and "why," not the "how."
        * Separate the subject from the body with a blank line.
    * **Footer (optional):**
        * Include a reference to a bug or issue tracker (e.g., 'Closes #123', 'Fixes #456').
3.  **Output Format:**
    * Provide the branch name and commit message separately, formatted as follows:

Branch Name: <generated_branch_name>
Commit Message:
<generated_commit_message>`

// ParseGitCustomCommit reads the prompt.md file from .git-commit and returns a custom prompt
func ParseGitCustomCommit() (string, error) {
	// Check if the prompt.md file exists in the .git-commit folder
	if _, err := os.Stat(".git-commit/prompt.md"); os.IsNotExist(err) {
		// If the file doesn't exist, return an empty string
		return "", nil
	}

	// Open the file for reading
	file, err := os.Open(".git-commit/prompt.md")
	if err != nil {
		return "", fmt.Errorf("failed to open file .git-commit/prompt.md: %v", err)
	}
	defer file.Close()

	// Read the entire file into a string
	content, err := os.ReadFile(".git-commit/prompt.md")
	if err != nil {
		return "", fmt.Errorf("error reading file .git-commit/prompt.md: %v", err)
	}

	return string(content), nil
}

// GetAIPrompt returns the AI prompt (standard or custom)
func GetAIPrompt(promptName string) string {
	// If a specific prompt name is provided, try to load it from custom-instructions
	if promptName != "" {
		customPrompt, err := loadCustomPrompt(promptName)
		if err != nil {
			fmt.Printf("Error reading custom prompt '%s': %v, using standard\n", promptName, err)
			return defaultAIPrompt
		}
		if strings.TrimSpace(customPrompt) != "" {
			return customPrompt
		}
	}

	// Try to read the default custom prompt
	customPrompt, err := ParseGitCustomCommit()
	if err != nil {
		fmt.Printf("Error reading custom prompt: %v, using standard\n", err)
		return defaultAIPrompt
	}

	// If the custom prompt is empty, use the standard one
	if strings.TrimSpace(customPrompt) == "" {
		return defaultAIPrompt
	}

	return customPrompt
}

// loadCustomPrompt loads a custom prompt from the custom-instructions folder
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