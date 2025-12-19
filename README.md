# OpenRouter CLI

A simple command-line interface for sending prompts to OpenRouter's API via stdin. Perfect for piping data from other commands.

## Overview

This CLI tool reads input from stdin and sends it to OpenRouter's API. It's designed to work seamlessly with Unix pipes, making it easy to process command output through AI models.

## Features

- 🚀 **Piped Input**: Read from stdin for seamless integration with other commands
- 🔒 **Environment-Based Configuration**: All settings via environment variables
- 🛡️ **Robust Error Handling**: Clear error messages and proper HTTP handling

## Installation

### Prerequisites

- Go 1.16 or later
- An OpenRouter API key (get one from [OpenRouter](https://openrouter.ai/))

### Building from Source

```bash
# Navigate to the project directory
cd openrouter-cli

# Build the project
go build -o openrouter-cli
```

## Usage

### Basic Usage

The tool reads from stdin and requires the API key to be set via environment variable:

```bash
# Set your API key
export OPENROUTER_API_KEY="your-api-key-here"

# Pipe command output to OpenRouter
ps aux | ./openrouter-cli

# Pipe file content
cat file.txt | ./openrouter-cli

# Pipe command output with a pre-prompt
echo "Hello world" | ./openrouter-cli
```

### Environment Variables

The tool uses the following environment variables:

- **`OPENROUTER_API_KEY`** (required): Your OpenRouter API key
- **`OPENROUTER_MODEL`** (optional): Model to use (default: `openai/gpt-oss-20b:free`)
- **`OPENROUTER_PRE_PROMPT`** (optional): Text to prepend to the stdin input

### Examples

#### Basic Piping

```bash
# Analyze process list
ps aux | ./openrouter-cli

# Analyze log file
tail -n 100 app.log | ./openrouter-cli

# Analyze command output
df -h | ./openrouter-cli
```

#### With Model Selection

```bash
export OPENROUTER_API_KEY="your-api-key"
export OPENROUTER_MODEL="openai/gpt-4"

ps aux | ./openrouter-cli
```

#### With Pre-Prompt

```bash
export OPENROUTER_API_KEY="your-api-key"
export OPENROUTER_PRE_PROMPT="Analyze the following process list and identify any suspicious processes:"

ps aux | ./openrouter-cli
```

#### Combined Configuration

```bash
export OPENROUTER_API_KEY="your-api-key"
export OPENROUTER_MODEL="anthropic/claude-3-opus"
export OPENROUTER_PRE_PROMPT="Summarize the following:"

cat document.txt | ./openrouter-cli
```

## Error Handling

The CLI provides clear error messages for common issues:

- Missing API key
- Empty input
- Network errors
- API errors

## Security

- API keys are only read from environment variables (never from command-line arguments)
- API keys are never logged or stored
- Use environment variables for API keys in production environments

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## License

MIT License - see LICENSE file for details.

## Support

For issues and questions:
- Check the [OpenRouter documentation](https://openrouter.ai/docs)
- Open an issue on GitHub
