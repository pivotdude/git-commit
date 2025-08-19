package prompt

import (
	"os"
	"testing"
)

func setupTestDir(t *testing.T) string {
	// Create a temporary directory in /tmp for testing
	tempDir, err := os.MkdirTemp("/tmp", "prompt-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	
	// Save original working directory
	originalDir, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get original working directory: %v", err)
	}
	
	// Change to temp directory
	err = os.Chdir(tempDir)
	if err != nil {
		t.Fatalf("Failed to change to temp dir: %v", err)
	}
	
	// Create .git-commit directory
	err = os.Mkdir(".git-commit", 0755)
	if err != nil {
		t.Fatalf("Failed to create .git-commit dir: %v", err)
	}
	
	// Cleanup function
	t.Cleanup(func() {
		// Change back to original directory
		os.Chdir(originalDir)
		// Remove temp directory
		os.RemoveAll(tempDir)
	})
	
	return tempDir
}

func TestParseGitCustomCommit_FileExists(t *testing.T) {
	setupTestDir(t)
	
	// Create custom prompt file with some content
	customContent := `# Custom AI Prompt

You are a custom AI assistant for generating commit messages.

Please follow these rules:
1. Use conventional commits format
2. Keep messages concise
3. Include scope when applicable

Example: feat(auth): add login functionality
`
	err := os.WriteFile(".git-commit/prompt.md", []byte(customContent), 0644)
	if err != nil {
		t.Fatalf("Failed to create custom prompt file: %v", err)
	}
	
	result, err := ParseGitCustomCommit()
	if err != nil {
		t.Fatalf("ParseGitCustomCommit() error = %v", err)
	}
	
	if result != customContent {
		t.Errorf("Expected custom prompt content, got different content")
	}
}

func TestParseGitCustomCommit_FileDoesNotExist(t *testing.T) {
	setupTestDir(t)
	
	// Ensure .git-commit/prompt.md doesn't exist
	os.Remove(".git-commit/prompt.md")
	
	result, err := ParseGitCustomCommit()
	if err != nil {
		t.Fatalf("ParseGitCustomCommit() error = %v", err)
	}
	
	if result != "" {
		t.Errorf("Expected empty string, got %q", result)
	}
}

func TestParseGitCustomCommit_EmptyFile(t *testing.T) {
	setupTestDir(t)
	
	// Create empty prompt file
	err := os.WriteFile(".git-commit/prompt.md", []byte(""), 0644)
	if err != nil {
		t.Fatalf("Failed to create empty prompt file: %v", err)
	}
	
	result, err := ParseGitCustomCommit()
	if err != nil {
		t.Fatalf("ParseGitCustomCommit() error = %v", err)
	}
	
	if result != "" {
		t.Errorf("Expected empty string, got %q", result)
	}
}

func TestParseGitCustomCommit_WhitespaceOnly(t *testing.T) {
	setupTestDir(t)
	
	// Create prompt file with only whitespace
	whitespaceContent := "   \n\t   \n   "
	err := os.WriteFile(".git-commit/prompt.md", []byte(whitespaceContent), 0644)
	if err != nil {
		t.Fatalf("Failed to create whitespace-only prompt file: %v", err)
	}
	
	result, err := ParseGitCustomCommit()
	if err != nil {
		t.Fatalf("ParseGitCustomCommit() error = %v", err)
	}
	
	if result != whitespaceContent {
		t.Errorf("Expected whitespace content, got different content")
	}
}

func TestGetAIPrompt_CustomPromptExists(t *testing.T) {
	setupTestDir(t)
	
	// Create custom prompt file
	customContent := `# Custom AI Prompt
This is a custom prompt for testing.`
	err := os.WriteFile(".git-commit/prompt.md", []byte(customContent), 0644)
	if err != nil {
		t.Fatalf("Failed to create custom prompt file: %v", err)
	}
	
	result := GetAIPrompt("")
	if result != customContent {
		t.Errorf("Expected custom prompt, got default prompt")
	}
}

func TestGetAIPrompt_CustomPromptEmpty(t *testing.T) {
	setupTestDir(t)
	
	// Create empty custom prompt file
	err := os.WriteFile(".git-commit/prompt.md", []byte(""), 0644)
	if err != nil {
		t.Fatalf("Failed to create empty custom prompt file: %v", err)
	}
	
	result := GetAIPrompt("")
	if result != defaultAIPrompt {
		t.Errorf("Expected default prompt, got custom prompt")
	}
}

func TestGetAIPrompt_CustomPromptWhitespaceOnly(t *testing.T) {
	setupTestDir(t)
	
	// Create custom prompt file with only whitespace
	whitespaceContent := "   \n\t   \n   "
	err := os.WriteFile(".git-commit/prompt.md", []byte(whitespaceContent), 0644)
	if err != nil {
		t.Fatalf("Failed to create whitespace-only custom prompt file: %v", err)
	}
	
	result := GetAIPrompt("")
	if result != defaultAIPrompt {
		t.Errorf("Expected default prompt, got custom prompt")
	}
}

func TestGetAIPrompt_NoCustomPromptFile(t *testing.T) {
	setupTestDir(t)
	
	// Ensure .git-commit/prompt.md doesn't exist
	os.Remove(".git-commit/prompt.md")
	
	result := GetAIPrompt("")
	if result != defaultAIPrompt {
		t.Errorf("Expected default prompt, got something else")
	}
}

func TestGetAIPrompt_CustomPromptReadError(t *testing.T) {
	setupTestDir(t)
	
	// Create a directory instead of a file to cause read error
	err := os.Mkdir(".git-commit/prompt.md", 0755)
	if err != nil {
		t.Fatalf("Failed to create directory instead of file: %v", err)
	}
	
	result := GetAIPrompt("")
	if result != defaultAIPrompt {
		t.Errorf("Expected default prompt on read error, got something else")
	}
}

func TestDefaultAIPromptFormat(t *testing.T) {
	// Test that the default prompt has the expected format and content
	if defaultAIPrompt == "" {
		t.Error("Default prompt should not be empty")
	}
	
	// Check for key sections
	expectedSections := []string{
		"You are a senior software engineer",
		"generate a branch name",
		"Branch Name:",
		"Commit Message:",
		"<generated_branch_name>",
		"<generated_commit_message>",
	}
	
	for _, section := range expectedSections {
		if !contains(defaultAIPrompt, section) {
			t.Errorf("Default prompt missing expected section: %q", section)
		}
	}
}

func TestDefaultAIPromptBranchRules(t *testing.T) {
	// Test that the default prompt includes branch naming rules
	expectedPrefixes := []string{
		"feature/",
		"fix/",
		"hotfix/",
		"refactor/",
		"chore/",
		"docs/",
	}
	
	for _, prefix := range expectedPrefixes {
		if !contains(defaultAIPrompt, prefix) {
			t.Errorf("Default prompt missing branch prefix: %q", prefix)
		}
	}
	
	// Check for branch format description
	if !contains(defaultAIPrompt, "lowercase, hyphen-separated") {
		t.Error("Default prompt missing branch format description")
	}
}

func TestDefaultAIPromptCommitRules(t *testing.T) {
	// Test that the default prompt includes commit message rules
	expectedRules := []string{
		"Conventional Commits",
		"imperative, present tense",
		"under 50 characters",
		"lowercase letter",
	}
	
	for _, rule := range expectedRules {
		if !contains(defaultAIPrompt, rule) {
			t.Errorf("Default prompt missing commit rule: %q", rule)
		}
	}
}

// Helper function to check if a string contains a substring
func contains(s, substr string) bool {
	return len(s) >= len(substr) && indexOf(s, substr) >= 0
}

// Helper function to find the index of a substring
func indexOf(s, substr string) int {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return i
		}
	}
	return -1
}