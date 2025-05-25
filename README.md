# OpenRouter CLI

A command-line interface for sending prompts to OpenRouter's API.

## Overview

This CLI tool allows you to send prompts to various language models via OpenRouter's API. It supports multiple input methods, model selection, and flexible API key management.

## Installation

### Prerequisites

- Go 1.16 or later
- An OpenRouter API key (get one from [OpenRouter](https://openrouter.ai/))

### Building from Source

```bash
# Clone the repository
git clone https://github.com/blackopsinc/openrouter-cli.git

# Navigate to the project directory
cd openrouter-cli

# Build the project
go build
```

## Usage

```bash
# Using with an API key and prompt directly as an argument
./openrouter-cli -key "your-api-key-here" "What is the capital of France?"

# Using with an API key and prompt from stdin (when no arguments are provided)
./openrouter-cli -key "your-api-key-here"

# Using with an API key and prompt from a file
./openrouter-cli -key "your-api-key-here" -file "/path/to/file.txt"

# Specify a different model (default is anthropic/claude-3.7-sonnet:thinking)
./openrouter-cli -key "your-api-key-here" -model "openai/gpt-4-turbo"
```

## Available Flags

- `-key`: Your OpenRouter API key (can also be provided via environment variable)
- `-file`: Path to a file whose content will be used as the prompt
- `-model`: The model to use for the request (default is "anthropic/claude-3.7-sonnet:thinking")

## Environment Variables

- `OPENROUTER_API_KEY`: You can set this environment variable instead of using the `-key` flag

## Input Priorities

The CLI prioritizes input sources as follows:

1. If a file is provided with `-file`, it will use the file content as the prompt
2. If no file is provided but a direct prompt is given, it will use the direct prompt
3. If neither a file nor direct prompt is provided, it will read from standard input (stdin)

## Examples

```bash
# Using a direct prompt
./openrouter-cli -key "your-api-key-here" "Explain quantum computing in simple terms"

# Using a file as the prompt source
./openrouter-cli -key "your-api-key-here" -file "code.go"

# Using stdin (interactive mode)
./openrouter-cli -key "your-api-key-here"
# Then type your prompt and press Ctrl+D (Unix/Linux) or Ctrl+Z followed by Enter (Windows) when done

# Using environment variable for API key
export OPENROUTER_API_KEY="your-api-key-here"
./openrouter-cli "Explain quantum computing in simple terms"
```

## Available Models

You can specify various models from OpenRouter's providers, including:

- `anthropic/claude-3.7-sonnet:thinking` (default)
- `openai/gpt-4o`
- `openai/gpt-4-turbo`
- `anthropic/claude-3-opus`
- `anthropic/claude-3-sonnet`
- And many more available through OpenRouter

Check the [OpenRouter documentation](https://openrouter.ai/docs) for the most up-to-date list of available models.

## License

MIT

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.
