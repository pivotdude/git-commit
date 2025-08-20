package prompt

import (
	"fmt"
	"os"
	"strings"
)

const defaultAIPrompt = "Generate a branch name and git commit message based on the provided git diff. Strictly follow these rules:\n" +
	"**Branch Name Rules:**\n\n" +
	"1. Use lowercase (e.g., \"feature/new-sidebar\")\n" +
	"2. Separate words with hyphens (\"-\")\n" +
	"3. Add a prefix at the start:\n" +
	"   - \"feature/\" for new functionality\n" +
	"   - \"bugfix/\" for bug fixes\n" +
	"   - \"hotfix/\" for urgent production fixes\n" +
	"   - \"docs/\" for documentation changes\n" +
	"   - \"refactor/\" for refactoring\n" +
	"   - \"test/\" for tests\n" +
	"   - \"chore/\" for maintenance tasks\n" +
	"4. Optionally add a ticket number after the prefix (e.g., \"feature/T-123-add-filters\")\n" +
	"5. Do not use: PascalCase, camelCase, snake_case, dots, special characters, or repeated hyphens\n" +
	"   **Git Commit Rules (Conventional Commits):**\n" +
	"6. **Message Structure:**\n" +
	"   - Header (mandatory): \"<type>[scope]: <description>\"\n" +
	"   - Body (optional): detailed description of changes\n" +
	"   - Footer (optional): metadata (ticket number, breaking changes)\n" +
	"7. **Header Rules:**\n" +
	"   - Use imperative mood (e.g., \"Add\", \"Fix\", \"Update\")\n" +
	"   - Start with a capital letter\n" +
	"   - No period at the end\n" +
	"   - Max length: 50 characters\n" +
	"   - Format: \"<type>[scope]: <description>\"\n" +
	"8. **Commit Types:**\n" +
	"   - \"feat\" - new functionality\n" +
	"   - \"fix\" - bug fixes\n" +
	"   - \"docs\" - documentation changes\n" +
	"   - \"refactor\" - code refactoring\n" +
	"   - \"test\" - adding/fixing tests\n" +
	"   - \"chore\" - maintenance tasks (configs, scripts)\n" +
	"   - \"perf\" - performance optimization\n" +
	"   - \"style\" - code style changes\n" +
	"   - \"build\" - build system changes\n" +
	"   - \"ci\" - CI configuration changes\n" +
	"9. **Scope:**\n" +
	"   - Enclosed in parentheses after the type\n" +
	"   - Indicates affected module (e.g., \"auth\", \"ui\", \"api\")\n" +
	"   - Example: \"feat(auth): add user login\"\n" +
	"10. **Body Message:**\n" +
	"    - Separated from header by a blank line\n" +
	"    - Each line ≤ 72 characters\n" +
	"    - Use markdown list for changes:\n" +
	"      ```\n" +
	"      - Change 1\n" +
	"      - Change 2\n" +
	"      ```\n" +
	"    - Each item starts with a capital letter\n" +
	"    - Use imperative mood\n" +
	"    - No periods at the end of items\n" +
	"11. **Footer:**\n" +
	"    - Separated from body by a blank line\n" +
	"    - Format: \"BREAKING CHANGE:\" for critical changes\n" +
	"    - Add ticket references: \"Refs: #123\"\n" +
	"      **Correct Examples:**\n" +
	"12. For a new feature:\n\n" +
	"```\n" +
	"feature/add-user-auth\n" +
	"feat(auth): implement user login\n\n" +
	"- Add authentication service\n" +
	"- Create login component\n" +
	"- Integrate OAuth provider\n" +
	"```\n\n" +
	"2. For a bug fix:\n" +
	"```\n" +
	"bugfix/T-456-fix-login-bug\n" +
	"fix(auth): resolve login validation\n\n" +
	"- Fix password validation\n" +
	"- Update error messages\n" +
	"- Add unit tests\n" +
	"```\n\n" +
	"3. For documentation:\n" +
	"```\n" +
	"docs/update-readme\n" +
	"docs: add installation guide\n\n" +
	"- Add setup instructions\n" +
	"- Update dependencies list\n" +
	"- Include configuration examples\n" +
	"```\n\n" +
	"**Git Diff for Analysis:**\n" +
	"[git diff will be inserted here]\n" +
	"**Output Format:**\n" +
	"The response must contain ONLY:\n" +
	"1. Branch name (one line)\n" +
	"2. Git commit in the format:\n" +
	"```\n" +
	"   <commit header>\n\n" +
	"   - Change 1\n" +
	"   - Change 2\n" +
	"   - Change 3\n" +
	"```\n" +
	"No explanations, additional text, or formatting.\n" +
	"```\n" +
	"**Key Prompt Features:**\n" +
	"1. **Clear structure** - separate rules for branch and commit\n" +
	"2. **Conventional Commits** - full specification with types, scopes, and format\n" +
	"3. **Markdown list for body** - explicit change formatting requirement\n" +
	"4. **Examples** - correct variants for different task types\n" +
	"5. **Output restrictions** - strict response format\n" +
	"6. **Ticket integration** - compatibility with tracking systems\n" +
	"7. **Formatting rules** - case, punctuation, line length constraints\n" +
	"Example model output:\n" +
	"```\n" +
	"feature/add-payment-gateway\n" +
	"feat(payment): integrate stripe api\n" +
	"- Add Stripe service module\n" +
	"- Create payment component\n" +
	"- Implement transaction handling\n" +
	"- Add error handling logic\n" +
	"```"

// ParseGitCustomCommit reads the prompt.md file from .git-commit and returns a custom prompt
func parseGitCustomCommitMessage() (string, error) {
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

func GetChangesAiPrompt() string {
	rawPrompt := ""
	customPrompt, err := parseGitCustomCommitMessage()
	if err != nil {
		fmt.Printf("Error reading custom prompt: %v, using standard\n", err)
		rawPrompt = defaultAIPrompt
	} else if strings.TrimSpace(customPrompt) == "" {
		rawPrompt = defaultAIPrompt
	} else {
		rawPrompt = customPrompt
	}
	
	return rawPrompt
}
