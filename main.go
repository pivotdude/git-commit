package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	defaultprompt "git-commit/default_prompt"
	"git-commit/diff"
	"git-commit/git"
	"git-commit/help"
	"git-commit/prompt"
	"git-commit/utils"
)

func main() {
	// Define command line flags
	showHelp := flag.Bool("h", false, "Show help message")
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

	// Handle verbose mode
	if *verbose {
		log.SetOutput(os.Stdout)
		log.SetFlags(log.LstdFlags | log.Lshortfile)
	} else {
		log.SetOutput(os.Stderr)
		log.SetFlags(0)
	}

	// Route to appropriate handler based on flags
	if len(os.Args) == 1 || *showHelp {
		help.ShowHelp()
		return
	}

	if *generatePrompt {
		handleGeneratePrompt()
		return
	}

	// Default behavior - run the main workflow
	handleDefaultWorkflow(customPromptName)
}




func removeFilesFromStaged(filesToIgnore []string) {
	if len(filesToIgnore) > 0 {
		fmt.Printf("Removing %d files from staged to ignore:\n", len(filesToIgnore))
		for _, file := range filesToIgnore {
			fmt.Printf("  - %s\n", file)
		}

		err := git.RemoveFilesFromStaged(filesToIgnore)
		if err != nil {
			log.Printf("Error removing files from staged: %v", err)
		}
	}
}


func handleGeneratePrompt() {
	diffOutput := diff.GetDiffOutputWithoutIgnoresFiles()

	aiPrompt := defaultprompt.GetChangesAiPrompt()
	finalPrompt := fmt.Sprintf("<diff>\n%s\n</diff>\n\n%s", diffOutput, aiPrompt)

	utils.CopyToClipboard(finalPrompt)
}

func handleDefaultWorkflow(customPromptName string) {
	aiPrompt := prompt.GetAIPrompt(customPromptName)
	utils.CopyToClipboard(aiPrompt)
}


