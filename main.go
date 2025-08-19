package main

import (
	"bytes"
	"flag"
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
	// Define command line flags
	help := flag.Bool("h", false, "Show help message")
	verbose := flag.Bool("v", false, "Enable verbose output")
	generatePrompt := flag.Bool("generate-prompt", false, "Generate prompt for current changes")
	flag.Parse()

	// Get custom prompt name if provided (skip flags)
	customPromptName := ""
	for _, arg := range os.Args[1:] {
		if !strings.HasPrefix(arg, "-") {
			customPromptName = arg
			break
		}
	}

	// If no arguments provided, show help
	if len(os.Args) == 1 || *help {
		showHelp()
		return
	}

	// Handle help flag
	if *help {
		showHelp()
		return
	}

	// Handle verbose mode
	if *verbose {
		log.SetOutput(os.Stdout)
		log.SetFlags(log.LstdFlags | log.Lshortfile)
	} else {
		log.SetOutput(os.Stderr)
		log.SetFlags(0)
	}

	// Handle generate prompt command
	if *generatePrompt {
		generatePromptForChanges(customPromptName)
		return
	}

	// Default behavior - run the main workflow
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
	aiPrompt := prompt.GetAIPrompt(customPromptName)

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

func showHelp() {
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
	customPrompts := getAvailableCustomPrompts()
	if len(customPrompts) > 0 {
		fmt.Println("Available custom prompts:")
		for _, promptName := range customPrompts {
			fmt.Printf("  git-commit %s\n", promptName)
		}
	} else {
		fmt.Println("No custom prompts found. Create .git-commit/custom-instructions/ folder with .md files.")
	}
}

// getAvailableCustomPrompts scans the custom-instructions directory and returns available prompt names
func getAvailableCustomPrompts() []string {
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

func generatePromptForChanges(customPromptName string) {
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

	// 4. Remove ignored files from staged
	if len(filesToIgnore) > 0 {
		fmt.Printf("Ignoring %d files based on .gitdiffignore:\n", len(filesToIgnore))
		for _, file := range filesToIgnore {
			fmt.Printf("  - %s\n", file)
		}

		err = git.RemoveFilesFromStaged(filesToIgnore)
		if err != nil {
			log.Printf("Error removing files from staged: %v", err)
		}
	}

	// 5. Get git diff
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
	aiPrompt := prompt.GetAIPrompt(customPromptName)

	// 7. Generate and display the final prompt
	finalPrompt := fmt.Sprintf("<diff>\n%s\n</diff>\n\n%s", diffOutput, aiPrompt)
	fmt.Println("Generated AI prompt:")
	fmt.Println(finalPrompt)

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
