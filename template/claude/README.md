# Project

A project initialized with Claude Code skills and agent orchestration capabilities.

## What's Included

This template includes three meta-skills that enable collaborative, AI-driven development:

### 🚀 project-init
Detects empty directories and guides you through project setup:
- Interactive stack selection (Go, Node.js, Python, Rust)
- Creates appropriate project structure
- Sets up stack-specific skills and agents
- Configures development workflow

### 🛠️ skill-creator
Create custom Claude Code skills for your workflow:
- Scaffolds skill structure with `init_skill.sh`
- Provides templates and best practices
- Validates skill structure

### 🤖 agent-creator
Create specialized AI agents with personality:
- Template library for common agent types
- Personality customization (Friendly, Professional, Analytical)
- Orchestrator pattern for task delegation
- Stack-specific agent suggestions

## Getting Started

### 1. Enter Development Environment

```bash
# With direnv (recommended)
direnv allow

# Or with nix develop
nix develop
```

### 2. Initialize Your Project

Open Claude Code and let the `project-init` skill guide you:

```
"Let's start a new project"
```

The skill will:
1. Ask about your tech stack
2. Suggest relevant skills and agents
3. Create project structure
4. Set up orchestrator and specialist agents
5. Configure development workflow

### 3. Customize Your Setup

Create additional skills:
```bash
.claude/skills/skill-creator/scripts/init_skill.sh my-skill --path .claude/skills
```

Create custom agents through Claude Code using the agent-creator guidance.

## Project Structure

```
.
├── flake.nix                    # Nix development environment
├── .claude/
│   ├── skills/                  # Claude Code skills
│   │   ├── project-init/        # Project initialization
│   │   ├── skill-creator/       # Skill creation
│   │   └── agent-creator/       # Agent creation
│   └── agents/                  # AI agents (created during init)
│       └── orchestrator.md      # Task coordinator (created during init)
└── README.md
```

After running project-init, you'll have:
- Stack-specific project files (go.mod, package.json, etc.)
- Source and test directories
- Orchestrator agent for task coordination
- Stack-specific specialist agents (test-runner, linter, etc.)

## Agent Orchestration

Every project includes an orchestrator agent that:
- Gathers context by reading files and checking status
- Delegates all execution to specialist agents
- Never makes changes directly
- Coordinates complex workflows

Specialist agents handle:
- **Testing**: Run tests, report failures
- **Linting**: Code quality checks
- **Security**: Vulnerability scanning
- **Documentation**: Generate and update docs
- **Refactoring**: Code improvements

## Philosophy

This template embraces:
- **Collaborative development**: Work with AI as a team
- **Test-first approach**: Quality and testing as first-class citizens
- **Dynamic adaptation**: No static templates, everything customized
- **Personality**: Agents with human-like communication
- **Orchestration**: Coordination over monolithic execution

## Next Steps

1. Run project-init to set up your specific stack
2. Customize your orchestrator's personality
3. Add stack-specific specialist agents
4. Create custom skills for your workflow
5. Start building with your AI team!

## Resources

- [Claude Code Documentation](https://docs.claude.com/en/docs/claude-code)
- Skills: See `.claude/skills/*/SKILL.md` for detailed guides
- Agents: Check `.claude/agents/*.md` for agent configurations

---

Built with ❤️ using Claude Code
