package utils

import (
	"fmt"
	"os/exec"
	"runtime"
	"strings"
)

// GlobMatch checks if a path matches a pattern (simplified implementation)
func GlobMatch(pattern, path string) bool {
	// If the pattern ends with /, it's a directory
	if strings.HasSuffix(pattern, "/") {
		// Check if the path is in this directory
		return strings.HasPrefix(path, pattern) || path == strings.TrimSuffix(pattern, "/")
	}
	
	// If the pattern doesn't contain separator characters (/, \), consider it a directory
	if !strings.ContainsAny(pattern, "/\\") {
		// Check if the path is a file in this directory or the directory itself
		return strings.HasPrefix(path, pattern+"/") || path == pattern || strings.HasPrefix(path, pattern+"/")
	}
	
	// If the pattern contains *, it's a wildcard
	if strings.Contains(pattern, "*") {
		// Simple implementation for * at the end (e.g., *.log)
		if strings.HasSuffix(pattern, "*") {
			prefix := strings.TrimSuffix(pattern, "*")
			return strings.HasPrefix(path, prefix)
		}
		// Simple implementation for * at the beginning (e.g., *.log)
		if strings.HasPrefix(pattern, "*") {
			suffix := strings.TrimPrefix(pattern, "*")
			return strings.HasSuffix(path, suffix)
		}
	}
	
	// Simple exact match
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