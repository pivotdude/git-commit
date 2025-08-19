# Git Commit Message Generator

A Go tool that helps you generate AI-powered commit messages by analyzing your staged git changes and preparing prompts for AI assistants.

## Features

- **Automatic Git Diff Analysis**: Analyzes staged files and generates git diff output
- **AI Prompt Generation**: Creates structured prompts for AI assistants to generate commit messages
- **File Ignoring**: Supports `.gitdiffignore` file to exclude specific files from analysis with advanced pattern matching
- **Custom Prompts**: Allows custom AI prompts via `.git-commit/prompt.md` with intelligent fallback to default prompt
- **Multiple Custom Prompt Support**: Supports multiple custom prompts in `.git-commit/custom-instructions/` directory
- **Command Line Interface**: Enhanced CLI with flags for various operations
- **Cross-Platform**: Works on macOS, Linux, and Windows
- **Clipboard Integration**: Automatically copies the generated prompt to your clipboard
- **Enhanced Pattern Matching**: Improved wildcard and directory pattern support for file filtering

## Installation

### Prerequisites

- Go 1.24.6 or higher
- Git installed on your system

### Build from Source

```bash
# Clone the repository
git clone <repository-url>
cd git-commit

# Build the application
go build -o git-commit

# Or run directly
go run main.go
```

````

## Usage

### Basic Usage

1. Stage your files for commit:

   ```bash
   git add <your-files>
   ```

2. Run the tool:

   ```bash
   ./git-commit
   ```

3. The tool will:
   - Analyze your staged files
   - Generate a git diff
   - Create an AI prompt with the diff
   - Copy the prompt to your clipboard

### Command Line Options

```bash
git-commit [prompt-name]  Generate AI prompt with git diff and copy to clipboard
git-commit -h             Show help message
git-commit -v             Enable verbose output
git-commit -generate-prompt  Generate prompt for current changes
```

### Examples

```bash
git-commit              # Generate prompt and copy to clipboard
git-commit mark         # Use custom prompt from custom-instructions/mark.md
git-commit -v           # Run with verbose logging
git-commit -generate-prompt  # Generate prompt without copying to clipboard
```

### Configuration

#### Ignoring Files

Create a `.git-commit/ignore` file to specify patterns of files to ignore:

```bash
mkdir -p .git-commit
echo "# Ignore generated files" >> .git-commit/ignore
echo "*.log" >> .git-commit/ignore
echo "temp" >> .git-commit/ignore
```

#### Default Custom AI Prompt

Create a default custom AI prompt by creating `.git-commit/prompt.md`:

```bash
mkdir -p .git-commit
cat > .git-commit/prompt.md << EOF
You are a senior software engineer who specializes in Git...

[Your custom prompt here]
EOF
```

#### Multiple Custom Prompts

Create multiple custom prompts by adding `.md` files to `.git-commit/custom-instructions/`:

```bash
mkdir -p .git-commit/custom-instructions
cat > .git-commit/custom-instructions/mark.md << EOF
You must update the README.md file with the latest information about your project. Please update the file.
EOF
```

Use custom prompts by specifying their name:

```bash
git-commit mark
```

**Note**: If the custom prompt file is empty or contains only whitespace, the system will automatically use the default prompt.

## How It Works

1. **Parse Git Diff Ignore**: Reads `.git-commit/ignore` for file patterns to exclude with enhanced wildcard support
2. **Get Staged Files**: Retrieves files added to git staging area
3. **Filter Files**: Removes files matching ignore patterns using improved pattern matching
4. **Generate Diff**: Creates git diff output for remaining files
5. **Create AI Prompt**: Combines diff with AI instructions, using custom prompt if available and valid
6. **Copy to Clipboard**: Places the complete prompt in clipboard
7. **Restore Files**: Returns ignored files to staging area

## Project Structure

```
git-commit/
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
- `GetAIPrompt()` - Gets AI prompt (default or custom) with intelligent fallback
- `ParseGitCustomCommit()` - Reads and validates custom prompt from `.git-commit/prompt.md`

### Utility Functions

- `GlobMatch()` - Enhanced pattern matching for file paths with wildcard and directory support
- `CopyToClipboard()` - Cross-platform clipboard copying with error handling

## Contributing

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Changelog

### v1.1.0

- Added command line interface with flags
- Implemented multiple custom prompt support
- Added verbose mode for debugging
- Added generate-prompt mode for viewing without clipboard copy
- Enhanced help system with examples
- Added automatic detection of available custom prompts
- Improved error handling and user feedback

### v1.0.0

- Initial release
- Basic git diff analysis
- AI prompt generation
- File ignoring functionality
- Cross-platform clipboard support
- Custom prompt support
- Enhanced pattern matching with wildcard support
- Improved custom prompt handling with intelligent fallback
- Comprehensive test coverage

## Support

If you encounter any issues or have questions:

1. Check the [Issues](https://github.com/your-repo/issues) page
2. Create a new issue with detailed description
3. Include your OS and Go version

## Acknowledgments

- Built with Go programming language
- Uses standard git command-line interface
- Inspired by modern AI-assisted development workflows
````
