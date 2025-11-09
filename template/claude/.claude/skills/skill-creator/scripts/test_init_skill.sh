#!/usr/bin/env bash
# Test script for init_skill.sh
# Validates that skills are created with correct structure and frontmatter

set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
INIT_SCRIPT="$SCRIPT_DIR/init_skill.sh"
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

# Test 1: Script creates skill directory
echo "Testing skill directory creation..."
"$INIT_SCRIPT" test-skill --path "$TEST_DIR" > /dev/null 2>&1
[[ -d "$TEST_DIR/test-skill" ]] || fail "Skill directory not created"
pass "Skill directory created"

# Test 2: SKILL.md exists
echo "Testing SKILL.md creation..."
[[ -f "$TEST_DIR/test-skill/SKILL.md" ]] || fail "SKILL.md not created"
pass "SKILL.md created"

# Test 3: Frontmatter exists
echo "Testing frontmatter presence..."
grep -q "^---$" "$TEST_DIR/test-skill/SKILL.md" || fail "Frontmatter delimiters missing"
pass "Frontmatter delimiters present"

# Test 4: Required frontmatter fields
echo "Testing required frontmatter fields..."
SKILL_MD="$TEST_DIR/test-skill/SKILL.md"

grep -q "^name: test-skill$" "$SKILL_MD" || fail "name field missing or incorrect"
pass "name field correct"

grep -q "^description:" "$SKILL_MD" || fail "description field missing"
pass "description field present"

grep -q "^type:" "$SKILL_MD" || fail "type field missing"
pass "type field present"

grep -q "^enforcement:" "$SKILL_MD" || fail "enforcement field missing"
pass "enforcement field present"

grep -q "^priority:" "$SKILL_MD" || fail "priority field missing"
pass "priority field present"

grep -q "^keywords:" "$SKILL_MD" || fail "keywords field missing"
pass "keywords field present"

grep -q "^patterns:" "$SKILL_MD" || fail "patterns field missing"
pass "patterns field present"

# Test 5: Directory structure
echo "Testing directory structure..."
[[ -d "$TEST_DIR/test-skill/scripts" ]] || fail "scripts/ directory missing"
pass "scripts/ directory created"

[[ -d "$TEST_DIR/test-skill/references" ]] || fail "references/ directory missing"
pass "references/ directory created"

[[ -d "$TEST_DIR/test-skill/assets" ]] || fail "assets/ directory missing"
pass "assets/ directory created"

# Test 6: Resource files
echo "Testing resource files..."
[[ -f "$TEST_DIR/test-skill/scripts/example.sh" ]] || fail "scripts/example.sh missing"
pass "scripts/example.sh created"

[[ -x "$TEST_DIR/test-skill/scripts/example.sh" ]] || fail "scripts/example.sh not executable"
pass "scripts/example.sh is executable"

[[ -f "$TEST_DIR/test-skill/references/api_reference.md" ]] || fail "references/api_reference.md missing"
pass "references/api_reference.md created"

[[ -f "$TEST_DIR/test-skill/assets/example_asset.txt" ]] || fail "assets/example_asset.txt missing"
pass "assets/example_asset.txt created"

# Test 7: Duplicate creation prevention
echo "Testing duplicate prevention..."
if "$INIT_SCRIPT" test-skill --path "$TEST_DIR" > /dev/null 2>&1; then
    fail "Script allowed duplicate skill creation"
fi
pass "Duplicate creation prevented"

echo ""
echo -e "${GREEN}All tests passed!${NC}"
