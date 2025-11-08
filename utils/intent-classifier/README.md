# Intent Classifier

Semantic intent classification for skills and agents in Claude Code using transformer embeddings.

**Table of Contents**
- [Quick Start](#quick-start)
- [Why Semantic Matching?](#why-semantic-matching)
- [Building from Source](#building-from-source)
- [Usage Examples](#usage)
- [Output Format](#output-format)
- [File Format](#file-format)
- [How It Works](#how-it-works)
- [Model Information](#model-information)

---

## Quick Start

### 1. Prerequisites

Install `libffi` on your system:

```bash
# Ubuntu/Debian
sudo apt install libffi8

# Fedora/RHEL
sudo dnf install libffi

# Arch Linux
sudo pacman -S libffi

# macOS
brew install libffi

# Nix
nix profile install nixpkgs#libffi
```

### 2. Create Your Skills/Agents

Create a directory structure with your skills and agents:

```
my-project/
├── skills/
│   └── python-expert.md
└── agents/
    └── code-reviewer.md
```

### 3. Format Your Files

Each file should have YAML frontmatter:

```markdown
---
name: python-expert
priority: high
---

Expert Python programming assistance including frameworks and best practices.
```

**Frontmatter fields:**
- `name:` - Identifier (defaults to file path if omitted)
- `priority:` - `critical`, `high`, `medium` (default), or `low`
- `type:` - `skill` or `agent` (auto-detected from directory if omitted)

### 4. Run the Classifier

```bash
./intent-classifier \
  --prompt "help me with Python code review" \
  --embed my-project \
  --threshold 0.3
```

**First run**: Downloads llama.cpp (~34MB) and the embedding model (~21MB) automatically.

**Troubleshooting:**
- If you get `libffi` errors, ensure you installed the prerequisite (step 1)
- If you're using Nix, run inside `nix develop` shell for proper library paths
- Lower `--threshold` (e.g., `0.2`) to see more matches; raise it (e.g., `0.5`) for stricter matching

### 5. Output

```
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
🎯 SKILLS & AGENTS ACTIVATION CHECK
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

📚 RECOMMENDED SKILLS:
  → python-expert

💡 SUGGESTED AGENTS:
  → @code-reviewer

ACTION: Use Skill tool and Use @code-reviewer
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
```

### Why Semantic Matching?

The tool understands **meaning**, not just keywords:

| Your Prompt | Keyword Match? | Semantic Match? |
|-------------|----------------|-----------------|
| "create a bodybuilding routine" | ❌ No "fitness" keyword | ✅ Matches `fitness-designer` |
| "help me get in shape" | ❌ No "workout" keyword | ✅ Matches `fitness-designer` |
| "parse JSON file" | ❌ | ❌ Correctly ignored |

No need to maintain exhaustive keyword lists! The embedding model understands that "bodybuilding routine" and "fitness program" are conceptually related.

## Overview

This tool replaces keyword-based detection with intelligent semantic matching. It supports two modes:

1. **Embedding Mode** (fast): Uses sentence transformers to compute embeddings and match via cosine similarity
2. **LLM Mode** (deep): Uses language models to reason about matches and provide confidence scores

All embeddings and LLM responses are cached to optimize performance on subsequent runs.

## Features

- **Semantic matching**: Understands intent beyond exact keywords using embeddings or LLM reasoning
- **Dual-mode operation**:
  - Embedding mode (default): Fast cosine similarity matching with sentence transformers
  - LLM mode (optional): Deep reasoning with language models like SmolLM2
- **Auto-downloading**: Downloads models and llama.cpp automatically on first run
- **Caching**: Intelligently caches embeddings and LLM responses to avoid recomputation
- **Cross-platform**: Supports Linux, macOS, and Windows
- **URL-based models**: Specify any GGUF model via direct URL
- **No Python required**: Pure Go implementation using Yzma/llama.cpp

## Building from Source

If you want to build the tool yourself:

```bash
# Install prerequisites (see Quick Start section above)

# Build
go build -ldflags '-s -w' -o intent-classifier main.go
```

## Usage

### Auto-Detection Mode (Recommended)

The tool automatically detects skills and agents based on directory structure:

```bash
./intent-classifier \
  --prompt "build a Python web app with security checks" \
  --embed testdata \
  --threshold 0.2
```

This searches both `testdata/skills/` and `testdata/agents/` and displays both in a unified output.

### Skills Only

```bash
./intent-classifier \
  --prompt "handle foo operations" \
  --embed testdata/skills
```

### Agents Only

```bash
./intent-classifier \
  --prompt "handle foo operations" \
  --embed testdata/agents
```

### Arguments

**Required:**
- `--prompt`: User prompt to match against
- `--embed`: File or directory to embed and match

**Optional:**
- `--threshold`: Similarity threshold (0.0-1.0, default: `0.4`)
- `--output-type`: Force output type: `auto`, `skills`, or `agents` (default: `auto` - auto-detects from directory structure)
- `--embedding-model`: Embedding model URL or local path (default: all-MiniLM-L6-v2)
- `--lib`: Path to llama.cpp library directory (auto-download if empty)
- `--processor`: Processor type: `cpu`, `cuda`, `vulkan`, `metal` (default: `cpu`)

### First Run

On first run, the tool will:
1. Download llama.cpp binaries (~34MB) from GitHub releases
2. Cache them in `~/.cache/intent-classifier` (Linux/macOS) or `%LOCALAPPDATA%\intent-classifier` (Windows)
3. Load the backend libraries automatically

This is a one-time setup. Subsequent runs use the cached libraries.

## How It Works

### Embedding Mode (Default)

1. **Load Model**: Loads the quantized `all-MiniLM-L6-v2` GGUF model (auto-downloads on first run)
2. **Parse Files**: Reads YAML frontmatter (if present) from all files in the target directory
3. **Preprocess Text**: Smart text compaction before embedding:
   - Strips YAML frontmatter (already parsed, just noise)
   - Normalizes whitespace (collapses multiple spaces/newlines)
   - Removes English stop words ("the", "a", "is", etc.)
   - Reduces token count by ~30-50% while preserving semantic meaning
4. **Compute Embeddings**:
   - Computes 384-dimensional embedding for user prompt
   - Computes embeddings for each file's content (cached for performance)
5. **Match**: Calculates cosine similarity between prompt and each file
6. **Filter**: Returns files above similarity threshold (default: 0.4)
7. **Output**: Renders matches using specified template

### LLM Mode (Optional)

1. **Load Model**: Loads specified LLM model (e.g., SmolLM2-1.7B)
2. **Parse Files**: Reads file content
3. **Reason**: For each file, asks LLM to evaluate match and provide confidence score
4. **Cache**: Stores LLM responses to avoid recomputation
5. **Filter**: Returns files above threshold
6. **Output**: Renders matches using template

## Output Format

Matches are automatically grouped by type and priority level. The tool intelligently detects whether items are skills or agents and displays them in unified output.

### Combined Skills & Agents (Auto-Detected)

```
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
🎯 SKILLS & AGENTS ACTIVATION CHECK
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

⚠️  CRITICAL SKILLS (REQUIRED):
  → security-scanner

📚 RECOMMENDED SKILLS:
  → python-expert

💡 SUGGESTED SKILLS:
  → code-reviewer

📌 OPTIONAL SKILLS:
  → style-guide


⚠️  CRITICAL AGENTS (REQUIRED):
  → @security-expert

📚 RECOMMENDED AGENTS:
  → @python-specialist

💡 SUGGESTED AGENTS:
  → @code-analyzer

📌 OPTIONAL AGENTS:
  → @style-checker

ACTION: Use Skill tool and Use @security-expert, @python-specialist, @code-analyzer, @style-checker
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
```

### Skills Only

```
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
🎯 SKILLS ACTIVATION CHECK
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

📚 RECOMMENDED SKILLS:
  → python-expert

ACTION: Use Skill tool
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
```

### Agents Only

```
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
🤖 AGENTS ACTIVATION CHECK
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

📚 RECOMMENDED AGENTS:
  → @python-specialist

ACTION: Use @python-specialist when responding
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
```

**Type Detection:**
- Items in `/skills/` directories → displayed as skills
- Items in `/agents/` directories → displayed as agents with `@` prefix
- Can be overridden with `type:` field in frontmatter

**Priority levels** are read from the `priority:` field in frontmatter:
- `critical` - Required items (⚠️)
- `high` - Recommended items (📚)
- `medium` - Suggested items (💡) [default]
- `low` - Optional items (📌)

## File Format

Files can optionally have YAML frontmatter:

```markdown
---
name: foo
description: Handles foo operations
priority: high
---

# Content here...
```

**Frontmatter fields:**
- `name:` - Skill identifier (defaults to absolute file path if not specified)
- `priority:` - Priority level: `critical`, `high`, `medium`, `low` (defaults to `medium`)
- `description:` - Skill description (not currently used by matcher)

## Model Information

### Default Embedding Model

- **Model**: `all-MiniLM-L6-v2` (sentence-transformers)
- **Format**: GGUF Q5_K_M quantization
- **Size**: 21MB
- **Dimensions**: 384
- **Source**: https://huggingface.co/second-state/All-MiniLM-L6-v2-Embedding-GGUF

### Supported LLM Models

Any GGUF language model can be specified via URL. Popular choices:

- **SmolLM2-1.7B-Instruct**: Compact reasoning model (~1.2GB)
  - URL: `https://huggingface.co/second-state/SmolLM2-1.7B-Instruct-GGUF/resolve/main/SmolLM2-1.7B-Instruct-Q5_K_M.gguf`
- **SmolLM2-360M-Instruct**: Ultra-lightweight model (~250MB)
  - URL: `https://huggingface.co/bartowski/SmolLM2-360M-Instruct-GGUF/resolve/main/SmolLM2-360M-Instruct-Q5_K_M.gguf`

## Architecture

```
┌─────────────┐
│ User Prompt │
└──────┬──────┘
       │
       v
┌──────────────────┐      ┌───────────────┐
│ llama.cpp/Yzma   │<────>│ libggml-cpu.so│
│ (FFI bindings)   │      │ (backends)    │
└────────┬─────────┘      └───────────────┘
         │
         v
┌────────────────────┐
│ GGUF Model Loader  │
└─────────┬──────────┘
          │
          v
┌─────────────────────┐
│ Embedding Generator │
│ (384-dim vectors)   │
└──────────┬──────────┘
           │
           v
┌──────────────────────┐
│ Cosine Similarity    │
│ Matching (threshold) │
└───────────┬──────────┘
            │
            v
┌───────────────────────┐
│ Priority Grouping     │
│ & Formatted Output    │
└───────────────────────┘
```

## Dependencies

- **Go**: 1.24+
- **Yzma**: v0.8.2+ (llama.cpp Go bindings)
- **libffi**: Runtime dependency for FFI
- **llama.cpp**: Auto-downloaded binaries

## Limitations

- **Not fully static**: Requires libffi shared library at runtime
- **CPU-bound**: Embedding computation takes 1-2s per skill on CPU
- **Threshold tuning**: Default 0.4 similarity threshold may need adjustment for your skills
- **English only**: Sentence transformer model is English-only

## Testing

Run unit tests:

```bash
go test -v
```

Tests include:
- Frontmatter name extraction
- File/directory loading
- Cosine similarity calculations
- Content hashing
- LLM score parsing
- Embedding and LLM response caching

## Future Improvements

- [ ] Support GPU acceleration (CUDA/Vulkan/Metal)
- [ ] Bundle libffi for true portability
- [ ] Add multilingual model support
- [ ] Fine-tune models on specific use cases
- [ ] Support other template formats (JSON, YAML)

## License

MIT
