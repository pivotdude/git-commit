# git-commit

A Go tool that helps you generate AI-powered commit messages by analyzing your staged changes. It automatically creates well-formatted commit messages and branch names following best practices, then copies them to your clipboard for easy use.

## Features

- **AI-Powered Generation**: Uses AI to generate meaningful commit messages and branch names based on your code changes
- **File Ignoring**: Supports `ignore` file to exclude specific files from analysis with advanced pattern matching
- **Custom Prompts**: Allows custom AI prompts via `.git-commit/prompt.md` with intelligent fallback to default prompt
- **Multiple Custom Prompt Support**: Supports multiple custom prompts in `.git-commit/custom-instructions/` directory
- **Context File Support**: Allows including additional context files in prompts using `@context:filename.md` syntax
- **Command Line Interface**: Enhanced CLI with flags for various operations
- **Cross-Platform**: Works on macOS, Linux, and Windows
- **Clipboard Integration**: Automatically copies the generated prompt to your clipboard

## Installation

### Prerequisites

- Go 1.16+
- Git

### Method 1: Using Go install (Recommended)

```bash
go install github.com/your-username/git-commit@latest
```

### Method 2: Build from source

```bash
# Clone the repository
git clone https://github.com/your-username/git-commit.git
cd git-commit

# Build the binary
go build -o git-commit

# Move to a directory in your PATH
sudo mv git-commit /usr/local/bin/
```

## Usage

### Basic Usage

```bash
# Stage your changes
git add .

# Generate AI-powered commit message
git-commit

# The generated prompt is now in your clipboard, paste it into git commit
git commit
```

### Using Custom Prompts

Create custom prompts in the `.git-commit/custom-instructions/` directory:

```bash
# Create the custom instructions directory
mkdir -p .git-commit/custom-instructions

# Create a custom prompt for documentation changes
cat > .git-commit/custom-instructions/docs.md << EOF
You are a technical writer. Generate commit messages that focus on documentation improvements and user experience.
EOF
```

Use your custom prompt:

```bash
git-commit docs
```

### Available Commands

```bash
git-commit [prompt-name]           # Generate AI prompt with git diff and copy to clipboard
git-commit -h                      # Show help message
git-commit -v                      # Enable verbose output
git-commit -generate-prompt        # Generate prompt for current changes without copying to clipboard
```

### Configuration

#### Custom Default Prompt

Create a `.git-commit/prompt.md` file to override the default AI prompt:

```bash
mkdir -p .git-commit
cat > .git-commit/prompt.md << EOF
You are a senior software engineer. Generate commit messages that follow our team's specific guidelines...
EOF
```

#### File Ignoring

Create a `.git-commit/ignore` file to specify patterns of files to exclude from analysis:

```bash
# .git-commit/ignore
*.test.js
test/
*.log
config/local.*
```

#### Multiple Custom Prompts

Create multiple prompt files in the `.git-commit/custom-instructions/` directory:

```bash
.git-commit/
├── ignore
├── prompt.md
└── custom-instructions/
    ├── api.md
    ├── ui.md
    ├── docs.md
    └── refactor.md
```

Use them with:

```bash
git-commit api
git-commit ui
```

### Context File Support

You can include additional context files in your custom prompts using the `@context:` syntax. This is useful when you need to provide AI with additional context such as:

- Project documentation
- API specifications
- Configuration files
- Database schemas
- Style guides

**File Context Syntax:**

```bash
@context:filename.md
```

This will read the specified file and include its content in the prompt.

**Git Diff Context Syntax:**

```bash
@diff
```

This will include the git diff output in the prompt. Use this when you want the AI to have access to the actual code changes.

**Example:**

Create a custom prompt with context:

```bash
mkdir -p .git-commit/custom-instructions
cat > .git-commit/custom-instructions/api-feature.md << EOF
You are a senior API developer. When working on API changes, consider the following context:

@context:docs/api-spec.md
@context:config/database.yml
@diff

Generate API-focused commit messages that follow REST conventions and include proper versioning based on the provided API spec, database configuration, and code changes.
```

When you run `git-commit api-feature`, the tool will:

1. Read `docs/api-spec.md` and `config/database.yml` and include their contents
2. Include the git diff output (because of `@diff`)
3. Format the result as:

```markdown
<context file="docs/api-spec.md">
[Content of api-spec.md]
</context file>

<context file="config/database.yml">
[Content of database.yml]
</context file>

<diff>
[Git diff output here]
</diff>
```

**Error Handling:** If a context file cannot be read, it will be replaced with an error message in the prompt.

**Note:** Git diff is only included when `@diff` is explicitly specified in the custom prompt.

## How It Works

1. **Parse Git Diff Ignore**: Reads `.git-commit/ignore` for file patterns to exclude with enhanced wildcard support
2. **Get Staged Files**: Retrieves the list of currently staged files from Git
3. **Filter Files**: Removes files matching ignore patterns from the staged area temporarily
4. **Generate Diff**: Creates a git diff of the remaining staged changes
5. **Create AI Prompt**: Combines the diff with the appropriate AI prompt (default, custom, or specific named prompt)
6. **Copy to Clipboard**: Places the complete prompt in your system clipboard
7. **Restore Staged Files**: Returns the previously ignored files back to the staged area

## Custom Prompt Best Practices

When creating custom prompts, consider these guidelines:

- Be specific about the role you want the AI to assume
- Include any project-specific conventions or requirements
- Specify the desired format and structure of the output
- Reference relevant documentation or specifications
- Use clear, concise language

Example custom prompt:

```markdown
You are a frontend specialist working on a React application. When generating commit messages:

1. Focus on user-facing changes and component improvements
2. Reference the component hierarchy when relevant
3. Mention any UI/UX improvements
4. Follow our style guide at docs/frontend-style.md

@context:docs/frontend-style.md
@diff
```

## License

MIT License

## Contributing

Contributions are welcome! Please feel free to submit issues, feature requests, or pull requests.
