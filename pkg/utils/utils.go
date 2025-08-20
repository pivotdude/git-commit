package utils

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"
)

// GlobMatch checks if a path matches a pattern (simplified implementation)
func GlobMatch(pattern, path string) bool {
	// Handle wildcard (*) matching first.
	if strings.Contains(pattern, "*") {
		// This simplified logic handles a single wildcard '*' at the beginning, middle, or end.
		// It splits the pattern by the wildcard.
		parts := strings.Split(pattern, "*")
		if len(parts) == 2 {
			// If pattern is "*.go", parts are ["", ".go"]. Path must have an empty prefix and a ".go" suffix.
			// If pattern is "test.*", parts are ["test.", ""]. Path must have a "test." prefix and an empty suffix.
			// If pattern is "test*.log", parts are ["test", ".log"]. Path must have a "test" prefix and a ".log" suffix.
			return strings.HasPrefix(path, parts[0]) && strings.HasSuffix(path, parts[1])
		}
		// Note: This simplified implementation does not handle multiple wildcards like "test*/*.go".
	}

	// If the pattern ends with /, it's a directory pattern.
	if strings.HasSuffix(pattern, "/") {
		// Check if the path is the directory itself or a file/subdir within it.
		return strings.HasPrefix(path, pattern) || path == strings.TrimSuffix(pattern, "/")
	}

	// If the pattern doesn't contain path separators, treat it as a directory name.
	if !strings.ContainsAny(pattern, "/\\") {
		// Check if the path is the directory itself or a file within that directory.
		return strings.HasPrefix(path, pattern+"/") || path == pattern
	}

	// Fallback to simple exact match.
	return pattern == path
}

// CopyToClipboard copies text to clipboard
func CopyToClipboard(text string) {
	var cmd *exec.Cmd

	switch runtime.GOOS {
	case "darwin":
		cmd = exec.Command("pbcopy")
	case "linux":
		if _, err := exec.LookPath("xclip"); err == nil {
			cmd = exec.Command("xclip", "-selection", "clipboard")
			fmt.Println("Using xclip for clipboard")
			} else if _, err := exec.LookPath("xsel"); err == nil {
			cmd = exec.Command("xsel", "--clipboard", "--input")
		}
	case "windows":
		cmd = exec.Command("clip")
	}

	if cmd == nil {
		printf("Clipboard utility not found. Please copy manually:\n---\n%s\n---\n", text)
		return
	}

	cmd.Stdin = strings.NewReader(text)
	if err := cmd.Run(); err != nil {
		printf("Clipboard copy error: %v\n", err)
		printf("Please copy manually:\n---\n%s\n---\n", text)
	}
}

// Printf - simplified output function
func printf(format string, a ...interface{}) {
	fmt.Printf(format, a...)
}


// readFileContent reads and returns the content of a file
func ReadFileContent(filePath string) (string, error) {
	// Check if the file exists
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return "", fmt.Errorf("file not found")
	}

	// Read the file content
	content, err := os.ReadFile(filePath)
	if err != nil {
		return "", fmt.Errorf("failed to read file: %v", err)
	}

	return string(content), nil
}
