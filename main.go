package main

import (
    "bytes"
    "fmt"
    "log"
    "os"
    "os/exec"
    "runtime"
    "strings"
)

const aiPrompt = `
You are a senior software engineer who specializes in Git. Your task is to generate a branch name and a commit message based on the provided code changes (git diff). Follow these rules:
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
<generated_commit_message>

**Here is the 'git diff' output to analyze:**`

func main() {
    // Get staged git diff
    cmd := exec.Command("git", "diff", "--staged")
    var stdout, stderr bytes.Buffer
    cmd.Stdout = &stdout
    cmd.Stderr = &stderr
    err := cmd.Run()

    if err != nil {
        log.Fatalf("Error running git diff: %v\n%s", err, stderr.String())
    }

    diffOutput := strings.TrimSpace(stdout.String())
    if diffOutput == "" {
        fmt.Println("No changes detected. Please use 'git add' to stage files.")
        os.Exit(1)
    }

    finalPrompt := fmt.Sprintf("%s\n<diff>\n%s\n</diff>", aiPrompt, diffOutput)
    copyToClipboard(finalPrompt)
    fmt.Println("AI prompt with git diff has been copied to clipboard.")
}

func copyToClipboard(text string) {
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
        fmt.Printf("Clipboard utility not found. Please copy manually:\n---\n%s\n---\n", text)
        return
    }

    cmd.Stdin = strings.NewReader(text)
    if err := cmd.Run(); err != nil {
        fmt.Printf("Clipboard copy error: %v\n", err)
        fmt.Printf("Please copy manually:\n---\n%s\n---\n", text)
    }
}
