# Git Commit Message Generator

A Go tool that helps you generate AI-powered commit messages by analyzing your staged git changes and preparing prompts for AI assistants.

## Features

- **Automatic Git Diff Analysis**: Analyzes staged files and generates git diff output
- **AI Prompt Generation**: Creates structured prompts for AI assistants to generate commit messages
- **File Ignoring**: Supports `.gitdiffignore` file to exclude specific files from analysis
- **Custom Prompts**: Allows custom AI prompts via `.git-commit/prompt.md`
- **Cross-Platform**: Works on macOS, Linux, and Windows
- **Clipboard Integration**: Automatically copies the generated prompt to your clipboard

## Installation

### Prerequisites

- Go 1.24.6 or higher
- Git installed on your system

### Build from Source

```bash
# Clone the repository
git clone <repository-url>
cd git-commit-message

# Build the application
go build -o git-commit-message

# Or run directly
go run main.go
```

## Usage

### Basic Usage

1. Stage your files for commit:

   ```bash
   git add <your-files>
   ```

2. Run the tool:

   ```bash
   ./git-commit-message
   ```

3. The tool will:
   - Analyze your staged files
   - Generate a git diff
   - Create an AI prompt with the diff
   - Copy the prompt to your clipboard

### Configuration

#### Ignoring Files

Create a `.git-commit/ignore` file to specify patterns of files to ignore:

```bash
mkdir -p .git-commit
echo "# Ignore generated files" >> .git-commit/ignore
echo "*.log" >> .git-commit/ignore
echo "temp" >> .git-commit/ignore
```

#### Custom AI Prompts

Create a custom AI prompt by creating `.git-commit/prompt.md`:

```bash
mkdir -p .git-commit
cat > .git-commit/prompt.md << EOF
You are a senior software engineer who specializes in Git...

[Your custom prompt here]
EOF
```

## How It Works

1. **Parse Git Diff Ignore**: Reads `.git-commit/ignore` for file patterns to exclude
2. **Get Staged Files**: Retrieves files added to git staging area
3. **Filter Files**: Removes files matching ignore patterns
4. **Generate Diff**: Creates git diff output for remaining files
5. **Create AI Prompt**: Combines diff with AI instructions
6. **Copy to Clipboard**: Places the complete prompt in clipboard
7. **Restore Files**: Returns ignored files to staging area

## Project Structure

```
git-commit-message/
├── main.go              # Main application entry point
├── go.mod               # Go module file
├── git/                 # Git-related functionality
│   └── git.go          # Git operations and file management
├── prompt/             # AI prompt generation
│   └── prompt.go       # Prompt handling and customization
└── utils/              # Utility functions
    └── utils.go        # Helper functions and clipboard operations
```

## API Reference

### Main Functions

- `ParseGitDiffIgnore()` - Parses ignore patterns from `.git-commit/ignore`
- `GetStagedFiles()` - Retrieves staged files from git
- `GetFilesToIgnore()` - Filters files based on ignore patterns
- `GetAIPrompt()` - Gets AI prompt (default or custom)

### Utility Functions

- `GlobMatch()` - Pattern matching for file paths
- `CopyToClipboard()` - Cross-platform clipboard copying

## Contributing

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Changelog

### v1.0.0

- Initial release
- Basic git diff analysis
- AI prompt generation
- File ignoring functionality
- Cross-platform clipboard support
- Custom prompt support

## Support

If you encounter any issues or have questions:

1. Check the [Issues](https://github.com/your-repo/issues) page
2. Create a new issue with detailed description
3. Include your OS and Go version

## Acknowledgments

- Built with Go programming language
- Uses standard git command-line interface
- Inspired by modern AI-assisted development workflows
