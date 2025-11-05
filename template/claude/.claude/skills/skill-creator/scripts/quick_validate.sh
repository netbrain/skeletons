#!/usr/bin/env bash
# Quick validation script for skills - minimal version

set -e

validate_skill() {
    local skill_path="$1"

    # Check SKILL.md exists
    local skill_md="$skill_path/SKILL.md"
    if [[ ! -f "$skill_md" ]]; then
        echo "SKILL.md not found"
        return 1
    fi

    # Read file content
    local content
    content=$(cat "$skill_md")

    # Check if it starts with frontmatter
    if [[ ! "$content" =~ ^--- ]]; then
        echo "No YAML frontmatter found"
        return 1
    fi

    # Extract frontmatter (everything between first --- and second ---)
    local frontmatter
    frontmatter=$(echo "$content" | awk '/^---$/{flag=!flag; next} flag')

    if [[ -z "$frontmatter" ]]; then
        echo "Invalid frontmatter format"
        return 1
    fi

    # Check required fields
    if ! echo "$frontmatter" | grep -q "^name:"; then
        echo "Missing 'name' in frontmatter"
        return 1
    fi

    if ! echo "$frontmatter" | grep -q "^description:"; then
        echo "Missing 'description' in frontmatter"
        return 1
    fi

    # Extract and validate name
    local name
    name=$(echo "$frontmatter" | grep "^name:" | sed 's/^name:[[:space:]]*//' | tr -d '\r')

    if [[ -n "$name" ]]; then
        # Check naming convention (hyphen-case: lowercase with hyphens)
        if [[ ! "$name" =~ ^[a-z0-9-]+$ ]]; then
            echo "Name '$name' should be hyphen-case (lowercase letters, digits, and hyphens only)"
            return 1
        fi

        # Check for invalid hyphen patterns
        if [[ "$name" =~ ^- ]] || [[ "$name" =~ -$ ]] || [[ "$name" =~ -- ]]; then
            echo "Name '$name' cannot start/end with hyphen or contain consecutive hyphens"
            return 1
        fi
    fi

    # Extract and validate description
    local description
    description=$(echo "$frontmatter" | grep "^description:" | sed 's/^description:[[:space:]]*//' | tr -d '\r')

    if [[ -n "$description" ]]; then
        # Check for angle brackets
        if [[ "$description" =~ \< ]] || [[ "$description" =~ \> ]]; then
            echo "Description cannot contain angle brackets (< or >)"
            return 1
        fi
    fi

    echo "Skill is valid!"
    return 0
}

# Main
if [[ $# -ne 1 ]]; then
    echo "Usage: $0 <skill_directory>"
    exit 1
fi

if validate_skill "$1"; then
    exit 0
else
    exit 1
fi
