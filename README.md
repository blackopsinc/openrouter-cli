# OpenRouter CLI

A powerful command-line interface for sending prompts to OpenRouter's API with advanced features and configuration management.

## Overview

This CLI tool allows you to send prompts to various language models via OpenRouter's API. It supports multiple input methods, model aliases, configuration management, and flexible API key handling.

## Features

- 🚀 **Multiple Input Methods**: Direct prompts, file content, or interactive stdin
- 🎯 **Model Aliases**: Use short aliases instead of full model names
- ⚙️ **Configuration Management**: Save and load settings from config file
- ⏱️ **Timeout Control**: Configurable request timeouts
- 📝 **Verbose Logging**: Optional detailed output for debugging
- 🔒 **Secure API Key Handling**: Environment variables or command-line flags
- 📏 **File Size Validation**: Prevents uploading oversized files
- 🛡️ **Robust Error Handling**: Clear error messages and proper HTTP handling

## Installation

### Prerequisites

- Go 1.16 or later
- An OpenRouter API key (get one from [OpenRouter](https://openrouter.ai/))

### Building from Source

```bash
# Clone the repository
git clone https://github.com/yourusername/openrouter-cli.git

# Navigate to the project directory
cd openrouter-cli

# Build the project
go build
```

## Usage

### Basic Usage

```bash
# Using with an API key and prompt directly as an argument
./openrouter-cli -key "your-api-key-here" "What is the capital of France?"

# Using with an API key and prompt from stdin (when no arguments are provided)
./openrouter-cli -key "your-api-key-here"

# Using with an API key and prompt from a file
./openrouter-cli -key "your-api-key-here" -file "/path/to/file.txt"

# Using model aliases for convenience
./openrouter-cli -key "your-api-key-here" -model gpt-4 "Explain quantum computing"
```

### Advanced Usage

```bash
# Enable verbose output
./openrouter-cli -v -key "your-api-key" "What is Go programming?"

# Set custom timeout
./openrouter-cli -timeout 30s -key "your-api-key" "Complex question here"

# List available model aliases
./openrouter-cli --list-models

# Save current configuration
./openrouter-cli --save-config
```

## Available Flags

- `-key`: Your OpenRouter API key (can also be provided via environment variable)
- `-file`: Path to a file whose content will be used as the prompt
- `-model`: Model to use for the request (supports aliases, default: claude-thinking)
- `-verbose`, `-v`: Enable verbose output for debugging
- `-timeout`: Request timeout duration (default: 60s)
- `--list-models`: List all available model aliases
- `--save-config`: Save current configuration to file

## Environment Variables

- `OPENROUTER_API_KEY`: You can set this environment variable instead of using the `-key` flag

## Model Aliases

The CLI includes convenient aliases for popular models:

| Alias | Full Model Name |
|-------|----------------|
| `claude-3-opus` | `anthropic/claude-3-opus` |
| `claude-3-sonnet` | `anthropic/claude-3-sonnet` |
| `claude-3-haiku` | `anthropic/claude-3-haiku` |
| `gpt-4` | `openai/gpt-4` |
| `gpt-4-turbo` | `openai/gpt-4-turbo` |
| `gpt-3.5-turbo` | `openai/gpt-3.5-turbo` |
| `claude-thinking` | `anthropic/claude-3.7-sonnet:thinking` |

## Configuration

The CLI automatically creates a configuration file at `~/.config/openrouter-cli/config.json` with default settings. You can customize:

- Default model
- Request timeout
- Maximum file size
- Model aliases

Example configuration:
```json
{
  "default_model": "anthropic/claude-3.7-sonnet:thinking",
  "timeout_seconds": 60,
  "max_file_size_bytes": 10485760,
  "models": {
    "gpt-4": "openai/gpt-4",
    "claude-3-opus": "anthropic/claude-3-opus"
  }
}
```

## Input Priorities

The CLI prioritizes input sources as follows:

1. If a file is provided with `-file`, it will use the file content as the prompt
2. If no file is provided but a direct prompt is given, it will use the direct prompt
3. If neither a file nor direct prompt is provided, it will read from standard input (stdin)

## Examples

### Basic Examples

```bash
# Using a direct prompt
./openrouter-cli -key "your-api-key-here" "Explain quantum computing in simple terms"

# Using a file as the prompt source
./openrouter-cli -key "your-api-key-here" -file "code.go"

# Using stdin (interactive mode)
./openrouter-cli -key "your-api-key-here"
# Then type your prompt and press Ctrl+D (Unix/Linux) or Ctrl+Z followed by Enter (Windows) when done
```

### Advanced Examples

```bash
# Using environment variable for API key with model alias
export OPENROUTER_API_KEY="your-api-key-here"
./openrouter-cli -model gpt-4 "Explain the differences between Go and Python"

# Verbose mode with custom timeout
./openrouter-cli -v -timeout 2m -model claude-3-opus "Write a detailed analysis of..."

# List available models
./openrouter-cli --list-models

# Process a large file with verbose output
./openrouter-cli -v -file large_document.txt -model gpt-4-turbo
```

## Error Handling

The CLI provides detailed error messages for common issues:

- Invalid API keys
- Network timeouts
- File size limits exceeded
- Unsupported file formats
- API rate limits

## Security

- API keys are never logged or stored in configuration files
- Use environment variables for API keys in production environments
- File size limits prevent accidental upload of large files

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## License

MIT License - see LICENSE file for details.

## Support

For issues and questions:
- Check the [OpenRouter documentation](https://openrouter.ai/docs)
- Open an issue on GitHub
- Review the verbose output with `-v` flag for debugging
