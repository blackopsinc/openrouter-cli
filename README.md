# OpenRouter CLI

A simple command-line interface for sending prompts to LLM APIs via stdin. Supports OpenRouter, Ollama, and LM Studio. Perfect for piping data from other commands.

## Overview

This CLI tool reads input from stdin and sends it to various LLM providers. It's designed to work seamlessly with Unix pipes, making it easy to process command output through AI models. Supports both cloud (OpenRouter) and local (Ollama, LM Studio) LLM providers.

## Features

- 🚀 **Piped Input**: Read from stdin for seamless integration with other commands
- 🔒 **Environment-Based Configuration**: All settings via environment variables
- 🛡️ **Robust Error Handling**: Clear error messages and proper HTTP handling   
- 📡 **Streaming Support**: Real-time streaming responses (SSE for OpenRouter/LM Studio, newline-delimited JSON for Ollama)
- 🔍 **Verbose/Debug Mode**: Detailed logging for troubleshooting
- 🏠 **Local LLM Support**: Works with Ollama and LM Studio for offline/local AI processing
- ☁️ **Cloud LLM Support**: Works with OpenRouter for cloud-based AI models

## Installation

### Prerequisites

- Go 1.21 or later
- For OpenRouter: An OpenRouter API key (get one from [OpenRouter](https://openrouter.ai/))
- For Ollama: Install and run [Ollama](https://ollama.ai/) locally
- For LM Studio: Install and run [LM Studio](https://lmstudio.ai/) with the local server enabled

### Building from Source

The project includes a Makefile with cross-platform build support. All binaries are built in the `bin/` directory.

#### Quick Build (Current Platform)

```bash
# Navigate to the project directory
cd openrouter-cli

# Build for your current platform
make build

# Or use go directly
go build -o bin/openrouter-cli .
```

#### Cross-Platform Builds

Build for all supported platforms:

```bash
# Build for all platforms (Linux, macOS, Windows)
make all
```

Build for a specific platform:

```bash
# Linux AMD64
make linux-amd64

# Linux ARM64
make linux-arm64

# macOS Intel (AMD64)
make darwin-amd64

# macOS Apple Silicon (ARM64)
make darwin-arm64

# Windows 64-bit
make windows-amd64

# Windows 32-bit
make windows-386
```

All binaries will be placed in the `bin/` directory with platform-specific suffixes:
- Linux: `openrouter-cli-linux-amd64` or `openrouter-cli-linux-arm64`
- macOS: `openrouter-cli-darwin-amd64` or `openrouter-cli-darwin-arm64`
- Windows: `openrouter-cli-windows-amd64.exe` or `openrouter-cli-windows-386.exe`

#### Installation

After building, you can install the binary:

```bash
# Install to system directory (requires sudo)
make install

# Install to user's local bin directory (no sudo required)
make install-user
```

#### Other Make Targets

```bash
# Clean build artifacts
make clean

# Format code
make fmt

# Run tests
make test

# Run tests with coverage
make test-coverage

# Show help
make help
```

## Usage

**Note:** In the examples below, `./openrouter-cli` refers to the binary. If you built it with `make build`, use `bin/openrouter-cli`. If you installed it with `make install` or `make install-user`, you can use `openrouter-cli` directly (assuming it's in your PATH).

### Basic Usage

The tool reads from stdin. Choose your provider:

#### Using OpenRouter (Cloud)

**Linux/macOS:**
```bash
# Set your API key and provider
export OPENROUTER_API_KEY="your-api-key-here"
export LLM_PROVIDER="openrouter"

# Pipe command output to OpenRouter
ps aux | ./openrouter-cli

# Pipe file content
cat file.txt | ./openrouter-cli
```

**Windows (PowerShell):**
```powershell
# Set your API key and provider
$env:OPENROUTER_API_KEY="your-api-key-here"
$env:LLM_PROVIDER="openrouter"

# Pipe command output to OpenRouter
Get-Process | .\openrouter-cli-windows-amd64.exe

# Pipe file content
Get-Content file.txt | .\openrouter-cli-windows-amd64.exe
```

**Windows (Command Prompt):**
```cmd
REM Set your API key and provider
set OPENROUTER_API_KEY=your-api-key-here
set LLM_PROVIDER=openrouter

REM Pipe command output to OpenRouter
tasklist | openrouter-cli-windows-amd64.exe

REM Pipe file content
type file.txt | openrouter-cli-windows-amd64.exe
```

#### Using Ollama (Local)

**Linux/macOS:**
```bash
# Set provider to Ollama (no API key needed)
export LLM_PROVIDER="ollama"
export LLM_MODEL="llama2"  # or any model you have installed

# Make sure Ollama is running, then pipe command output
ps aux | ./openrouter-cli
```

**Windows (PowerShell):**
```powershell
# Set provider to Ollama (no API key needed)
$env:LLM_PROVIDER="ollama"
$env:LLM_MODEL="llama2"

# Make sure Ollama is running, then pipe command output
Get-Process | .\openrouter-cli-windows-amd64.exe
```

#### Using LM Studio (Local)

**Linux/macOS:**
```bash
# Set provider to LM Studio (no API key needed)
export LLM_PROVIDER="lmstudio"
export LLM_MODEL="local-model"  # or the model name in LM Studio

# Make sure LM Studio server is running, then pipe command output
ps aux | ./openrouter-cli
```

**Windows (PowerShell):**
```powershell
# Set provider to LM Studio (no API key needed)
$env:LLM_PROVIDER="lmstudio"
$env:LLM_MODEL="local-model"

# Make sure LM Studio server is running, then pipe command output
Get-Process | .\openrouter-cli-windows-amd64.exe
```

### Environment Variables

The tool uses the following environment variables:

- **`LLM_PROVIDER`** (optional): Provider to use - `openrouter`, `ollama`, or `lmstudio` (default: `openrouter`)
- **`OPENROUTER_API_KEY`** (required for OpenRouter): Your OpenRouter API key
- **`LLM_MODEL`** (optional): Model to use
  - OpenRouter default: `openai/gpt-oss-20b:free`
  - Ollama default: `llama2`
  - LM Studio default: `local-model`
- **`OPENROUTER_PRE_PROMPT`** (optional): Text to prepend to the stdin input    
- **`OPENROUTER_STREAM`** (optional): Enable streaming responses. Set to `1`, `true`, `yes`, or `on` to enable
- **`OPENROUTER_VERBOSE`** (optional): Enable verbose/debug logging. Set to `1`, `true`, `yes`, or `on` to enable
- **`OLLAMA_URL`** (optional): Ollama API URL (default: `http://localhost:11434/api/chat`)
- **`LM_STUDIO_URL`** (optional): LM Studio API URL (default: `http://localhost:1234/v1/chat/completions`)

#### Setting Environment Variables

**Linux/macOS (Bash/Zsh):**
```bash
export OPENROUTER_API_KEY="your-api-key-here"
export LLM_PROVIDER="openrouter"
export LLM_MODEL="openai/gpt-4"
```

**Windows Command Prompt (cmd.exe):**
```cmd
set OPENROUTER_API_KEY=your-api-key-here
set LLM_PROVIDER=openrouter
set LLM_MODEL=openai/gpt-4
```

**Windows PowerShell:**
```powershell
$env:OPENROUTER_API_KEY="your-api-key-here"
$env:LLM_PROVIDER="openrouter"
$env:LLM_MODEL="openai/gpt-4"
```

**Note:** Environment variables set in Command Prompt or PowerShell are session-specific. To make them persistent, use System Properties → Environment Variables, or set them in your shell profile.

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
# OpenRouter
export OPENROUTER_API_KEY="your-api-key"
export LLM_PROVIDER="openrouter"
export LLM_MODEL="openai/gpt-4"
ps aux | ./openrouter-cli

# Ollama
export LLM_PROVIDER="ollama"
export LLM_MODEL="llama3.2"
ps aux | ./openrouter-cli

# LM Studio
export LLM_PROVIDER="lmstudio"
export LLM_MODEL="mistral-7b-instruct"
ps aux | ./openrouter-cli
```

#### With Pre-Prompt

```bash
# Works with any provider
export LLM_PROVIDER="ollama"  # or "openrouter" or "lmstudio"
export OPENROUTER_API_KEY="your-api-key"  # only needed for openrouter
export OPENROUTER_PRE_PROMPT="Analyze the following process list and identify any suspicious processes:"

ps aux | ./openrouter-cli
```

#### With Streaming Responses

```bash
export LLM_PROVIDER="ollama"  # or "openrouter" or "lmstudio"
export OPENROUTER_API_KEY="your-api-key"  # only needed for openrouter
export OPENROUTER_STREAM="true"

echo "Write a short story about a robot" | ./openrouter-cli
```

Streaming mode outputs responses in real-time as they're generated, providing a better user experience for longer responses.

#### With Verbose/Debug Mode

```bash
export LLM_PROVIDER="ollama"  # or "openrouter" or "lmstudio"
export OPENROUTER_API_KEY="your-api-key"  # only needed for openrouter
export OPENROUTER_VERBOSE="true"

echo "Hello world" | ./openrouter-cli
```

Verbose mode provides detailed logging including:
- Request/response details
- HTTP status codes and headers
- Input/output sizes
- Streaming chunk information
- Error details

#### Combined Configuration

```bash
# OpenRouter example
export LLM_PROVIDER="openrouter"
export OPENROUTER_API_KEY="your-api-key"
export LLM_MODEL="anthropic/claude-3-opus"
export OPENROUTER_PRE_PROMPT="Summarize the following:"
export OPENROUTER_STREAM="true"
export OPENROUTER_VERBOSE="true"
cat document.txt | ./openrouter-cli

# Ollama example
export LLM_PROVIDER="ollama"
export LLM_MODEL="llama3.2"
export OPENROUTER_PRE_PROMPT="Summarize the following:"
export OPENROUTER_STREAM="true"
export OPENROUTER_VERBOSE="true"
cat document.txt | ./openrouter-cli
```

## Provider-Specific Notes

### OpenRouter
- Requires an API key
- Supports all models available on OpenRouter
- Uses OpenAI-compatible API format
- Streaming uses Server-Sent Events (SSE)

### Ollama
- No API key required
- Make sure Ollama is running: `ollama serve`
- Install models: `ollama pull llama2` (or any other model)
- Uses Ollama's native API format
- Streaming uses newline-delimited JSON
- Default URL: `http://localhost:11434/api/chat`

### LM Studio
- No API key required
- Make sure LM Studio is running with the local server enabled (Developer tab)
- Load a model in LM Studio before using
- Uses OpenAI-compatible API format
- Streaming uses Server-Sent Events (SSE)
- Default URL: `http://localhost:1234/v1/chat/completions`

## Error Handling

The CLI provides clear error messages for common issues:

- Missing API key (for OpenRouter)
- Invalid provider
- Empty input
- Network errors (connection refused for local providers usually means the service isn't running)
- API errors

## Security

- API keys are only read from environment variables (never from command-line arguments)
- API keys are never logged or stored
- Use environment variables for API keys in production environments
- Local providers (Ollama, LM Studio) don't require API keys and run entirely on your machine

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## License

MIT License - see LICENSE file for details.

## Support

For issues and questions:
- Check the [OpenRouter documentation](https://openrouter.ai/docs)
- Check the [Ollama documentation](https://github.com/ollama/ollama/blob/main/docs/api.md)
- Check the [LM Studio documentation](https://lmstudio.ai/docs)
- Open an issue on GitHub

