# Git-Commit User Documentation

AI-powered Git commit message generator with customizable prompts and branch naming conventions.

---

## Table of Contents

- [Introduction](#introduction)
- [Quick Start](#quick-start)
- [Installation](#installation)
- [Command Usage](#command-usage)
  - [Basic Usage](#basic-usage)
  - [Advanced Options](#advanced-options)
- [Configuration](#configuration)
  - [Ignore Files](#ignore-files)
  - [Custom Prompts](#custom-prompts)
  - [Default Prompt Customization](#default-prompt-customization)
- [Working with Context](#working-with-context)
- [Troubleshooting](#troubleshooting)
- [Examples](#examples)
- [Appendix A: Default Prompt Template](#appendix-a-default-prompt-template)
- [Appendix B: Supported Platforms](#appendix-b-supported-platforms)

---

## Introduction

Git-Commit is a command-line tool that helps developers generate high-quality, standardized Git commit messages using AI. It analyzes your staged changes and creates structured prompts that can be used with AI models to generate appropriate commit messages following best practices like Conventional Commits.

Key features:

- Generates branch names and commit messages based on code changes
- Supports custom prompt templates
- Integrates with .gitdiffignore files to exclude specific files from analysis
- Copies generated prompts directly to clipboard for easy use with AI tools
- Cross-platform support (macOS, Linux, Windows)

---

## Quick Start

1. Stage your changes with `git add`
2. Run `git-commit` to generate a prompt and copy it to clipboard
3. Paste the prompt into your preferred AI tool to generate commit messages

```bash
# Stage your changes
git add src/main.go README.md

# Generate AI prompt and copy to clipboard
git-commit

# For verbose output
git-commit -v
```

---

## Installation

### Prerequisites

- Git must be installed and available in your PATH
- Go 1.19 or higher (for building from source)
- Clipboard utilities (pbcopy on macOS, xclip/xsel on Linux, clip on Windows)

### Building from Source

```bash
# Clone the repository
git clone https://github.com/your-username/git-commit.git
cd git-commit

# Build the binary
go build -o git-commit ./cmd/git-commit

# Move to a directory in your PATH
sudo mv git-commit /usr/local/bin/
```

### Verifying Installation

```bash
git-commit -h
```

---

## Command Usage

### Basic Usage

```bash
git-commit [prompt-name]
```

When run without arguments, git-commit:

1. Reads staged changes
2. Applies ignore rules from `.git-commit/ignore`
3. Generates an AI prompt based on the diff
4. Copies the prompt to clipboard

### Advanced Options

| Flag               | Description                                  | Example                       |
| ------------------ | -------------------------------------------- | ----------------------------- |
| `-h`               | Show help message                            | `git-commit -h`               |
| `-v`               | Enable verbose output                        | `git-commit -v`               |
| `-generate-prompt` | Generate prompt without copying to clipboard | `git-commit -generate-prompt` |

#### Examples

```bash
# Basic usage - generates prompt and copies to clipboard
git-commit

# Use a custom prompt template
git-commit conventional

# Enable verbose logging
git-commit -v

# Generate prompt without copying to clipboard
git-commit -generate-prompt
```

---

## Configuration

Git-Commit uses configuration files located in a `.git-commit` directory at the root of your repository.

### Ignore Files

Create `.git-commit/ignore` to specify files that should be excluded from the diff analysis:

```text
# Comments start with #
*.log
node_modules/
.env
dist/
test/fixtures/
```

The ignore file supports:

- Full file paths
- Directory patterns ending with `/`
- Wildcard patterns (limited support: `*.ext`, `prefix*`, `*suffix`)

### Custom Prompts

Create custom prompt templates in `.git-commit/custom-instructions/` with `.md` extensions:

```
.git-commit/
└── custom-instructions/
    ├── conventional.md
    ├── angular.md
    └── semantic.md
```

Use custom prompts by specifying their name:

```bash
git-commit conventional
```

#### Custom Prompt Directives

Custom prompts support special directives:

1. **@context:** - Include file or directory content in the prompt

   ```markdown
   @context: src/components/Button.js
   @context: src/utils/
   ```

2. **@diff** - Insert the git diff at that location
   ```markdown
   Please analyze these changes:
   @diff
   And provide a commit message.
   ```

### Default Prompt Customization

Override the default prompt by creating `.git-commit/prompt.md`. If this file exists and contains content, it will be used instead of the built-in prompt.

Example `.git-commit/prompt.md`:

```markdown
You are a senior software engineer reviewing code changes.

Please generate:

1. A branch name following kebab-case with appropriate prefix
2. A commit message following conventional commits format

@diff

Provide only the branch name and commit message with no additional explanation.
```

---

## Working with Context

Git-Commit can include relevant code context in your prompts to help AI models generate more accurate commit messages.

### Including File Context

In custom prompt files, use the `@context:` directive:

```markdown
# Code Review Prompt

Review the following component changes:
@context: src/components/UserProfile.jsx

Together with these updates:
@diff

Generate a concise commit message.
```

### Including Directory Context

Include entire directories for broader context:

```markdown
# Feature Implementation Prompt

Understand the application structure:
@context: src/api/

Review these specific changes:
@diff

Generate an appropriate commit message.
```

---

## Troubleshooting

### Common Issues

#### No Changes Detected

**Problem**: "No changes detected. Please use 'git add' to add files to staged."
**Solution**: Ensure you've staged your changes with `git add` before running git-commit.

#### Clipboard Not Working

**Problem**: "Clipboard utility not found."
**Solution**: Install appropriate clipboard utilities:

- macOS: Already included (pbcopy)
- Linux: Install xclip or xsel (`sudo apt install xclip` or `sudo yum install xclip`)
- Windows: Should work with built-in clip command

#### Custom Prompt Not Found

**Problem**: "Error reading custom prompt"
**Solution**: Verify the prompt file exists in `.git-commit/custom-instructions/` with the correct name and `.md` extension.

#### Permission Denied

**Problem**: "Permission denied" when running git-commit
**Solution**: Ensure the binary has execute permissions (`chmod +x git-commit`)

### Debugging with Verbose Mode

Use the `-v` flag to enable detailed logging:

```bash
git-commit -v
```

This will show:

- File operations
- Git command executions
- Configuration file processing
- Error details

---

## Examples

### Basic Workflow

```bash
# 1. Make changes to your code
echo "console.log('Hello World');" >> src/index.js

# 2. Stage your changes
git add src/index.js

# 3. Generate AI prompt
git-commit

# 4. Paste the prompt into your AI tool (e.g., ChatGPT, Claude)
# 5. Use the generated commit message
git commit -m "feat(client): add hello world logging"
```

### Using Custom Prompts

```bash
# Create a custom prompt
mkdir -p .git-commit/custom-instructions
cat > .git-commit/custom-instructions/angular.md << EOF
Generate a commit message following Angular's commit conventions.

Types: build, ci, docs, feat, fix, perf, refactor, test

@diff

Output only the commit message:
EOF

# Use the custom prompt
git-commit angular
```

### With Ignore Rules

```bash
# Create ignore file
mkdir -p .git-commit
echo "*.log" > .git-commit/ignore
echo "node_modules/" >> .git-commit/ignore

# These files will be ignored in the diff
echo "2023-01-01.log" >> debug.log
git add src/main.go debug.log

# Run git-commit (debug.log will be excluded from analysis)
git-commit
```

---

## Appendix A: Default Prompt Template

The default prompt template follows these conventions:

1. Branch naming with prefixes:

   - `feature/` - New functionality
   - `bugfix/` - Bug fixes
   - `hotfix/` - Urgent production fixes
   - `docs/` - Documentation changes
   - `refactor/` - Code refactoring
   - `test/` - Test-related changes
   - `chore/` - Maintenance tasks

2. Commit messages following Conventional Commits:
   - Format: `<type>[optional scope]: <description>`
   - Types: feat, fix, docs, refactor, test, chore, perf, style, build, ci
   - Body as markdown list with changes
   - Footer for references and breaking changes

Example output format:

```
feature/add-user-authentication
feat(auth): implement user login functionality

- Add authentication service module
- Create login form component
- Integrate with backend API
- Add unit tests for auth flows
```

---

## Appendix B: Supported Platforms

Git-Commit has been tested on:

| Platform | Clipboard Support | Notes                                    |
| -------- | ----------------- | ---------------------------------------- |
| macOS    | ✅ pbcopy         | Works out of the box                     |
| Linux    | ✅ xclip/xsel     | Requires installation of xclip or xsel   |
| Windows  | ✅ clip           | Works with Command Prompt and PowerShell |

### Linux Clipboard Setup

Install clipboard utilities:

```bash
# Ubuntu/Debian
sudo apt install xclip

# Or alternatively
sudo apt install xsel

# CentOS/RHEL/Fedora
sudo yum install xclip
# Or
sudo dnf install xclip
```

Check which utility is being used:

```bash
git-commit -v
# Output will show: "Using xclip for clipboard" or "Using xsel for clipboard"
```
