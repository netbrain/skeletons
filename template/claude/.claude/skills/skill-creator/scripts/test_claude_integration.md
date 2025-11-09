# Skill Creator - Claude Integration Tests

These tests verify that Claude correctly uses the skill-creator skill and init script when users request skill creation.

## Test Cases

### Test 1: Basic Skill Creation

**User prompt:**
```
Create a new skill called "api-helper" for handling API integrations
```

**Expected behavior:**
1. skill-creator skill activates
2. Claude runs: `.claude/skills/skill-creator/scripts/init_skill.sh api-helper --path .claude/skills`
3. Claude confirms creation and asks user to customize frontmatter

**Validation:**
```bash
# Run this after Claude creates the skill
test -f .claude/skills/api-helper/SKILL.md || echo "FAIL: SKILL.md not created"
grep -q "^name: api-helper$" .claude/skills/api-helper/SKILL.md || echo "FAIL: name field incorrect"
grep -q "^keywords:" .claude/skills/api-helper/SKILL.md || echo "FAIL: keywords field missing"
grep -q "^patterns:" .claude/skills/api-helper/SKILL.md || echo "FAIL: patterns field missing"
test -d .claude/skills/api-helper/scripts || echo "FAIL: scripts/ directory missing"
echo "PASS: Skill created with correct structure"
```

### Test 2: Skill Creation with Custom Path

**User prompt:**
```
Create a database skill in a custom location: skills/database
```

**Expected behavior:**
1. Claude runs: `.claude/skills/skill-creator/scripts/init_skill.sh database --path skills`
2. Skill created at `skills/database/`

**Validation:**
```bash
test -f skills/database/SKILL.md || echo "FAIL: Custom path not respected"
echo "PASS: Custom path used correctly"
```

### Test 3: Skill Customization

**User prompt:**
```
Create a Python testing skill with keywords: pytest, testing, unit tests
and patterns: run.*test, pytest.*
```

**Expected behavior:**
1. Claude creates skill with init script
2. Claude edits SKILL.md to add proper keywords and patterns
3. Claude removes TODO placeholders

**Validation:**
```bash
grep -q "keywords: pytest, testing, unit tests" .claude/skills/python-testing/SKILL.md || echo "FAIL: Keywords not customized"
grep -q "patterns: run.*test, pytest.*" .claude/skills/python-testing/SKILL.md || echo "FAIL: Patterns not customized"
! grep -q "\[TODO" .claude/skills/python-testing/SKILL.md || echo "FAIL: TODOs still present"
echo "PASS: Skill properly customized"
```

### Test 4: Error Handling - Duplicate Skill

**User prompt:**
```
Create a skill-creator skill
```

**Expected behavior:**
1. Claude runs init script
2. Script fails (skill-creator already exists)
3. Claude reports error to user and suggests alternative name

**Validation:**
```bash
# Manually verify Claude reports the error and doesn't create duplicate
```

## Manual Test Procedure

1. **Start fresh project:**
   ```bash
   cd /tmp && mkdir test-project && cd test-project
   nix flake init -t github:netbrain/skeletons#claude
   nix develop
   ```

2. **Run each test case:**
   - Copy user prompt
   - Send to Claude via `claude "<prompt>"`
   - Observe Claude's actions
   - Run validation commands

3. **Document results:**
   - [ ] Test 1: Basic creation
   - [ ] Test 2: Custom path
   - [ ] Test 3: Customization
   - [ ] Test 4: Error handling

## Expected Init Script Usage Patterns

Claude should:
- ✓ Use absolute or relative paths to init script
- ✓ Pass `--path` parameter correctly
- ✓ Capture script output
- ✓ Follow up by editing SKILL.md if user provides specific requirements
- ✓ Explain next steps to user (customize frontmatter, add content)

Claude should NOT:
- ✗ Create SKILL.md manually without using init script
- ✗ Skip frontmatter fields
- ✗ Leave all TODOs in place without guidance
- ✗ Create skills without directory structure

## Troubleshooting

**If Claude doesn't use init script:**
- Check that skill-creator skill is properly activated
- Verify init script is executable: `chmod +x .claude/skills/skill-creator/scripts/init_skill.sh`
- Check SKILL.md frontmatter includes proper keywords/patterns

**If created skill is missing frontmatter:**
- Verify init_skill.sh template includes all fields
- Run: `.claude/skills/skill-creator/scripts/test_init_skill.sh`
