# Claude Code Project Templates

Nix flake templates for Claude Code projects with intelligent skills and agent orchestration.

## Quick Start

Initialize a new Claude Code project:

```bash
# Create and enter new project directory
mkdir my-project && cd my-project

# Initialize from template
nix flake init -t github:netbrain/skeletons#claude

# Enter development environment
nix develop
# or with direnv: direnv allow

# Start Claude Code and initialize your project
# Just say: "Let's start a new project"
```

## What You Get

### 🚀 project-init skill
Intelligent project initialization that:
- Detects empty directories automatically
- Asks about your tech stack (Go, Node.js, Python, Rust)
- Creates appropriate project structure
- Sets up stack-specific skills and agents
- Configures orchestrator with your preferred personality

### 🛠️ skill-creator skill
Create custom Claude Code skills:
- Scaffolds skill structure with templates
- Provides best practices and examples
- Validates skill format
- Packages for distribution

### 🤖 agent-creator skill
Create AI agents with personality:
- Template library for common patterns
- Personality options (Friendly, Professional, Analytical)
- Orchestrator pattern for task delegation
- Stack-specific recommendations

## Templates

### `claude` (default)
General-purpose Claude Code project template.

```bash
nix flake init -t github:netbrain/skeletons#claude
```

Includes:
- Three meta-skills (project-init, skill-creator, agent-creator)
- Basic flake.nix with common dev tools
- .gitignore configured for Nix
- Comprehensive README

## Workflow Example

```bash
# 1. Initialize project
mkdir my-api && cd my-api
nix flake init -t github:netbrain/skeletons#claude

# 2. Enter dev environment
nix develop

# 3. Open Claude Code and initialize
# Claude detects empty directory and activates project-init
# You'll be asked:
#   - Tech stack? (Go, Node.js, Python, Rust)
#   - Project type? (API, CLI, library, full-stack)
#   - Agent personality? (Friendly, Professional, Analytical)
#   - Which agents? (test-runner, linter, security-auditor, etc.)

# 4. Project is now set up with:
#   - Stack-specific files (go.mod, package.json, etc.)
#   - Directory structure (src/, tests/, cmd/, etc.)
#   - Orchestrator agent coordinating work
#   - Specialist agents for testing, linting, etc.
#   - Stack-specific TDD skill

# 5. Start coding with your AI team!
```

## Philosophy

This template embraces:

**Collaborative Development**: Work with AI agents as a coordinated team, not a single monolithic assistant.

**Orchestration Pattern**: An orchestrator agent gathers context and delegates execution to specialists, following the SessionProcessManager pattern.

**Dynamic Customization**: No rigid templates - everything is built collaboratively based on your needs and preferences.

**Personality-Driven**: Agents communicate like team members with distinct personalities (friendly, professional, analytical).

**Quality First**: Testing, linting, and code review are first-class citizens with dedicated agents.

**Progressive Disclosure**: Skills and agents are created when needed, not all at once.

## Agent Orchestration

Every project includes an **orchestrator agent** that:
- Gathers information by reading files and checking status
- Analyzes what needs to be done
- Delegates all execution to appropriate specialist agents
- Never makes changes directly
- Explains reasoning before delegation

Example orchestrator personalities:

**Professional (Maestro):**
```
"I see you want to add authentication. Let me check the current
structure... Based on what I found, I'll have security-auditor
review implications first, then code-reviewer can assess the
implementation approach."
```

**Friendly (Buddy):**
```
"Hey! Let's add authentication! 🔐 I'll take a look at what we
have... Okay, I'm going to ask our security expert to check things
out first. Sound good?"
```

**Analytical (Architect):**
```
"Authentication module request detected. Executing codebase scan...
Security audit required. Delegating to security-auditor agent.
Awaiting analysis completion."
```

## Stack-Specific Features

### Go Projects
- `go-tdd` skill for test-driven development
- `go-test-runner` agent (runs `go test`)
- `go-linter` agent (golangci-lint)
- Standard Go project structure (cmd/, pkg/, internal/)

### Node.js/TypeScript
- `node-tdd` skill with Vitest/Jest
- `vitest-runner` agent for testing
- `type-checker` agent (TypeScript)
- Modern Node.js structure (src/, tests/)

### Python
- `python-tdd` skill with pytest
- `pytest-runner` agent
- `ruff-linter` agent for code quality
- Python package structure (src/, tests/)

### Rust
- `rust-tdd` skill with cargo test
- `cargo-test-runner` agent
- `clippy-linter` agent
- Standard Cargo structure

## Creating Custom Skills

```bash
# Initialize new skill
python .claude/skills/skill-creator/scripts/init_skill.py my-skill --path .claude/skills

# Edit the generated SKILL.md
vim .claude/skills/my-skill/SKILL.md

# Validate
python .claude/skills/skill-creator/scripts/quick_validate.py .claude/skills/my-skill

# Package for distribution
python .claude/skills/skill-creator/scripts/package_skill.py .claude/skills/my-skill
```

## Creating Custom Agents

Agents are simple markdown files with YAML frontmatter:

```markdown
---
name: my-agent
description: What this agent does and when to use it
model: sonnet
color: blue
tools: Read, Write, Bash
---

You are [agent name], specialized in [domain].

[Define personality, approach, and behavior]
```

See `.claude/skills/agent-creator/SKILL.md` for templates and examples.

## Advanced Usage

### Adding More Skills
```bash
# Create skill for your domain
python .claude/skills/skill-creator/scripts/init_skill.py api-design --path .claude/skills
# Customize for your API patterns
```

### Creating Agent Teams
```bash
# Create specialized agents for your workflow
# - data-validator (validates data integrity)
# - api-tester (tests API endpoints)
# - performance-monitor (watches for regressions)
```

### Multi-Stack Projects
```
my-monorepo/
├── backend/ (Go API)
│   └── .claude/agents/go-*
├── frontend/ (TypeScript)
│   └── .claude/agents/ts-*
└── .claude/
    ├── skills/ (shared)
    └── agents/
        └── orchestrator.md (coordinates both)
```

## Contributing

Contributions welcome! Areas for improvement:
- Additional language templates
- More agent templates
- Enhanced orchestration patterns
- Documentation improvements

## Resources

- [Claude Code Documentation](https://docs.claude.com/en/docs/claude-code)
- [Nix Flakes](https://nixos.wiki/wiki/Flakes)
- [Agent Orchestration Patterns](https://github.com/netbrain/skeletons/blob/main/.claude/skills/agent-creator/SKILL.md)

## License

MIT

---

Built with ❤️ using Claude Code and Nix
