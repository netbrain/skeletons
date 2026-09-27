#!/usr/bin/env bash
# Skill activation prompt hook - checks if skills should be activated based on user prompt

set -e

# Read input from stdin
input=$(cat)

# Parse JSON input (for hook input format)
prompt=$(echo "$input" | grep -oP '"prompt"\s*:\s*"\K[^"]+' 2>/dev/null || echo "")

if [[ -z "$prompt" ]]; then
    exit 0
fi

# Get project directory
project_dir="${CLAUDE_PROJECT_DIR}"
if [[ -z "$project_dir" ]]; then
    project_dir=$(echo "$input" | grep -oP '"cwd"\s*:\s*"\K[^"]+' 2>/dev/null || pwd)
fi

# Use keyword/regex matching
prompt=$(echo "$prompt" | tr '[:upper:]' '[:lower:]')

# Arrays to hold matched skills by priority
declare -a critical_skills=()
declare -a high_skills=()
declare -a medium_skills=()
declare -a low_skills=()

# Find and process all SKILL.md files in one operation
while IFS= read -r skill_file; do
    # Extract frontmatter (everything between first --- and second ---)
    frontmatter=$(awk '/^---$/{flag=!flag; next} flag' "$skill_file")

    if [[ -z "$frontmatter" ]]; then
        continue
    fi

    # Parse frontmatter fields
    skill_name=$(echo "$frontmatter" | grep "^name:" | sed 's/^name:[[:space:]]*//' | tr -d '\r')
    priority=$(echo "$frontmatter" | grep "^priority:" | sed 's/^priority:[[:space:]]*//' | tr -d '\r')
    keywords=$(echo "$frontmatter" | grep "^keywords:" | sed 's/^keywords:[[:space:]]*//' | tr -d '\r')
    patterns=$(echo "$frontmatter" | grep "^patterns:" | sed 's/^patterns:[[:space:]]*//' | tr -d '\r')

    # Skip if no activation metadata
    if [[ -z "$priority" && -z "$keywords" && -z "$patterns" ]]; then
        continue
    fi

    # Default priority
    priority="${priority:-medium}"

    # Check keyword triggers (comma-delimited)
    keyword_match=false
    if [[ -n "$keywords" ]]; then
        IFS=',' read -ra keyword_array <<< "$keywords"
        for keyword in "${keyword_array[@]}"; do
            keyword=$(echo "$keyword" | xargs)  # trim whitespace
            keyword_lower=$(echo "$keyword" | tr '[:upper:]' '[:lower:]')
            if [[ "$prompt" == *"$keyword_lower"* ]]; then
                keyword_match=true
                break
            fi
        done
    fi

    # Check intent pattern triggers (comma-delimited)
    intent_match=false
    if [[ -n "$patterns" ]]; then
        IFS=',' read -ra pattern_array <<< "$patterns"
        for pattern in "${pattern_array[@]}"; do
            pattern=$(echo "$pattern" | xargs)  # trim whitespace
            if echo "$prompt" | grep -qiE "$pattern"; then
                intent_match=true
                break
            fi
        done
    fi

    # Add to appropriate priority list if matched
    if [[ "$keyword_match" == true ]] || [[ "$intent_match" == true ]]; then
        case "$priority" in
            critical)
                critical_skills+=("$skill_name")
                ;;
            high)
                high_skills+=("$skill_name")
                ;;
            medium)
                medium_skills+=("$skill_name")
                ;;
            low)
                low_skills+=("$skill_name")
                ;;
        esac
    fi
done < <(find "$project_dir/.claude/skills" -name "SKILL.md" 2>/dev/null || true)

# Generate output if matches found
total_matches=$((${#critical_skills[@]} + ${#high_skills[@]} + ${#medium_skills[@]} + ${#low_skills[@]}))

if [[ $total_matches -gt 0 ]]; then
    echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
    echo "🎯 SKILL ACTIVATION CHECK"
    echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
    echo ""

    if [[ ${#critical_skills[@]} -gt 0 ]]; then
        echo "⚠️  CRITICAL SKILLS (REQUIRED):"
        for skill in "${critical_skills[@]}"; do
            echo "  → $skill"
        done
        echo ""
    fi

    if [[ ${#high_skills[@]} -gt 0 ]]; then
        echo "📚 RECOMMENDED SKILLS:"
        for skill in "${high_skills[@]}"; do
            echo "  → $skill"
        done
        echo ""
    fi

    if [[ ${#medium_skills[@]} -gt 0 ]]; then
        echo "💡 SUGGESTED SKILLS:"
        for skill in "${medium_skills[@]}"; do
            echo "  → $skill"
        done
        echo ""
    fi

    if [[ ${#low_skills[@]} -gt 0 ]]; then
        echo "📌 OPTIONAL SKILLS:"
        for skill in "${low_skills[@]}"; do
            echo "  → $skill"
        done
        echo ""
    fi

    echo "ACTION: Use Skill tool BEFORE responding"
    echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
fi

exit 0
