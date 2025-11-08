#!/usr/bin/env bash
# Agent Initializer - Creates a new agent from template
#
# Usage:
#     init_agent.sh <agent-name> --path <path>
#
# Examples:
#     init_agent.sh code-reviewer --path .claude/agents
#     init_agent.sh test-runner --path ~/.claude/agents

set -e

title_case_agent_name() {
    local agent_name="$1"
    echo "$agent_name" | sed 's/-/ /g; s/\b\(.\)/\u\1/g'
}

init_agent() {
    local agent_name="$1"
    local path="$2"

    # Determine agent file path
    local agent_file="$path/$agent_name.md"
    agent_file=$(realpath "$agent_file" 2>/dev/null || readlink -f "$agent_file" || echo "$agent_file")

    # Create directory if needed
    local agent_dir=$(dirname "$agent_file")
    if [[ ! -d "$agent_dir" ]]; then
        if ! mkdir -p "$agent_dir"; then
            echo "❌ Error creating directory: $agent_dir"
            return 1
        fi
        echo "✅ Created directory: $agent_dir"
    fi

    # Check if file already exists
    if [[ -f "$agent_file" ]]; then
        echo "❌ Error: Agent file already exists: $agent_file"
        return 1
    fi

    # Create agent file from template
    local agent_title
    agent_title=$(title_case_agent_name "$agent_name")

    cat > "$agent_file" <<EOF
---
name: $agent_name
description: [TODO: Brief description of what this agent does and when to use it. Include "PROACTIVELY" if it should auto-trigger.]
model: sonnet
color: blue
type: agent
enforcement: suggest
priority: medium
keywords: [TODO: Add comma-separated keywords for fallback matching, e.g., test, testing, review]
patterns: [TODO: Add comma-separated regex patterns for fallback matching, e.g., run.*test, review.*code]
---

You are $agent_title, a specialized agent.

## Your Role

[TODO: Define the agent's core responsibility and purpose]

## Your Approach

[TODO: Describe how the agent thinks and works]

## Communication Style

[TODO: Define the agent's personality and tone - friendly, professional, analytical, quirky]

## Key Tasks

[TODO: List the main tasks this agent handles]

## Integration

[TODO: Describe how this agent works with other agents or the orchestrator]
EOF

    echo "✅ Created agent file: $agent_file"

    # Print next steps
    echo ""
    echo "✅ Agent '$agent_name' initialized successfully at $agent_file"
    echo ""
    echo "Next steps:"
    echo "1. Edit the agent file to complete all TODO items"
    echo "2. Customize the system prompt to define behavior"
    echo "3. Test by explicitly invoking: 'Use the $agent_name agent'"
    echo ""
    echo "Frontmatter fields to customize:"
    echo "  - description: Explain what the agent does and when to use it"
    echo "  - model: sonnet (complex), haiku (simple), opus (advanced), inherit"
    echo "  - color: cyan, blue, green, red, yellow, purple, orange, pink, gray"
    echo "  - priority: critical, high, medium, low"
    echo "  - keywords: Comma-separated keywords for matching"
    echo "  - patterns: Comma-separated regex patterns for matching"
    echo "  - tools: (optional) Comma-separated tool list, e.g., Read, Write, Bash"

    return 0
}

# Main
if [[ $# -ne 3 ]] || [[ "$2" != "--path" ]]; then
    echo "Usage: $0 <agent-name> --path <path>"
    echo ""
    echo "Agent name requirements:"
    echo "  - Hyphen-case identifier (e.g., 'code-reviewer', 'test-runner')"
    echo "  - Lowercase letters, digits, and hyphens only"
    echo "  - Max 40 characters"
    echo ""
    echo "Examples:"
    echo "  $0 code-reviewer --path .claude/agents"
    echo "  $0 test-runner --path ~/.claude/agents"
    echo "  $0 orchestrator --path .claude/agents"
    exit 1
fi

agent_name="$1"
path="$3"

echo "🚀 Initializing agent: $agent_name"
echo "   Location: $path"
echo ""

if init_agent "$agent_name" "$path"; then
    exit 0
else
    exit 1
fi
