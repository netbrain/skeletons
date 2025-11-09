#!/usr/bin/env bash
# Test script for init_agent.sh
# Validates that agents are created with correct structure and frontmatter

set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
INIT_SCRIPT="$SCRIPT_DIR/init_agent.sh"
TEST_DIR=$(mktemp -d)

cleanup() {
    rm -rf "$TEST_DIR"
}
trap cleanup EXIT

# Colors
RED='\033[0;31m'
GREEN='\033[0;32m'
NC='\033[0m' # No Color

pass() {
    echo -e "${GREEN}✓${NC} $1"
}

fail() {
    echo -e "${RED}✗${NC} $1"
    exit 1
}

# Test 1: Script creates agent file
echo "Testing agent file creation..."
"$INIT_SCRIPT" test-agent --path "$TEST_DIR" > /dev/null 2>&1
[[ -f "$TEST_DIR/test-agent.md" ]] || fail "Agent file not created"
pass "Agent file created"

# Test 2: Frontmatter exists
echo "Testing frontmatter presence..."
grep -q "^---$" "$TEST_DIR/test-agent.md" || fail "Frontmatter delimiters missing"
pass "Frontmatter delimiters present"

# Test 3: Required frontmatter fields
echo "Testing required frontmatter fields..."
AGENT_MD="$TEST_DIR/test-agent.md"

grep -q "^name: test-agent$" "$AGENT_MD" || fail "name field missing or incorrect"
pass "name field correct"

grep -q "^description:" "$AGENT_MD" || fail "description field missing"
pass "description field present"

grep -q "^model:" "$AGENT_MD" || fail "model field missing"
pass "model field present"

grep -q "^color:" "$AGENT_MD" || fail "color field missing"
pass "color field present"

grep -q "^type:" "$AGENT_MD" || fail "type field missing"
pass "type field present"

grep -q "^enforcement:" "$AGENT_MD" || fail "enforcement field missing"
pass "enforcement field present"

grep -q "^priority:" "$AGENT_MD" || fail "priority field missing"
pass "priority field present"

grep -q "^keywords:" "$AGENT_MD" || fail "keywords field missing"
pass "keywords field present"

grep -q "^patterns:" "$AGENT_MD" || fail "patterns field missing"
pass "patterns field present"

# Test 4: System prompt sections exist
echo "Testing system prompt structure..."
grep -q "## Your Role" "$AGENT_MD" || fail "System prompt missing 'Your Role' section"
pass "'Your Role' section present"

grep -q "## Your Approach" "$AGENT_MD" || fail "System prompt missing 'Your Approach' section"
pass "'Your Approach' section present"

grep -q "## Communication Style" "$AGENT_MD" || fail "System prompt missing 'Communication Style' section"
pass "'Communication Style' section present"

# Test 5: Valid field values
echo "Testing field value validity..."
MODEL=$(grep "^model:" "$AGENT_MD" | cut -d: -f2 | xargs)
[[ "$MODEL" =~ ^(sonnet|haiku|opus|inherit)$ ]] || fail "Invalid model value: $MODEL"
pass "model value is valid"

COLOR=$(grep "^color:" "$AGENT_MD" | cut -d: -f2 | xargs)
[[ "$COLOR" =~ ^(cyan|blue|green|red|yellow|purple|orange|pink|gray)$ ]] || fail "Invalid color value: $COLOR"
pass "color value is valid"

TYPE=$(grep "^type:" "$AGENT_MD" | cut -d: -f2 | xargs)
[[ "$TYPE" == "agent" ]] || fail "Invalid type value: $TYPE"
pass "type value is correct"

PRIORITY=$(grep "^priority:" "$AGENT_MD" | cut -d: -f2 | xargs)
[[ "$PRIORITY" =~ ^(critical|high|medium|low)$ ]] || fail "Invalid priority value: $PRIORITY"
pass "priority value is valid"

# Test 6: Directory creation if needed
echo "Testing directory auto-creation..."
SUBDIR="$TEST_DIR/subdir"
"$INIT_SCRIPT" nested-agent --path "$SUBDIR" > /dev/null 2>&1
[[ -f "$SUBDIR/nested-agent.md" ]] || fail "Agent file not created in new directory"
pass "Directory auto-created"

# Test 7: Duplicate creation prevention
echo "Testing duplicate prevention..."
if "$INIT_SCRIPT" test-agent --path "$TEST_DIR" > /dev/null 2>&1; then
    fail "Script allowed duplicate agent creation"
fi
pass "Duplicate creation prevented"

echo ""
echo -e "${GREEN}All tests passed!${NC}"
