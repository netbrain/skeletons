# Agent Creator - Claude Integration Tests

These tests verify that Claude correctly uses the agent-creator skill and init script when users request agent creation.

## Test Cases

### Test 1: Basic Agent Creation

**User prompt:**
```
Create a code-reviewer agent
```

**Expected behavior:**
1. agent-creator skill activates
2. Claude runs: `.claude/skills/agent-creator/scripts/init_agent.sh code-reviewer --path .claude/agents`
3. Claude confirms creation and asks user about personality/customization

**Validation:**
```bash
# Run this after Claude creates the agent
test -f .claude/agents/code-reviewer.md || echo "FAIL: Agent file not created"
grep -q "^name: code-reviewer$" .claude/agents/code-reviewer.md || echo "FAIL: name field incorrect"
grep -q "^model:" .claude/agents/code-reviewer.md || echo "FAIL: model field missing"
grep -q "^color:" .claude/agents/code-reviewer.md || echo "FAIL: color field missing"
grep -q "^keywords:" .claude/agents/code-reviewer.md || echo "FAIL: keywords field missing"
grep -q "^patterns:" .claude/agents/code-reviewer.md || echo "FAIL: patterns field missing"
grep -q "## Your Role" .claude/agents/code-reviewer.md || echo "FAIL: System prompt missing"
echo "PASS: Agent created with correct structure"
```

### Test 2: Agent with Personality

**User prompt:**
```
Create a friendly test-runner agent that uses haiku model and green color
```

**Expected behavior:**
1. Claude creates agent with init script
2. Claude customizes frontmatter: `model: haiku`, `color: green`
3. Claude writes friendly personality in system prompt

**Validation:**
```bash
grep -q "^model: haiku$" .claude/agents/test-runner.md || echo "FAIL: Model not set to haiku"
grep -q "^color: green$" .claude/agents/test-runner.md || echo "FAIL: Color not set to green"
# Manual check: verify system prompt has friendly tone
echo "PASS: Agent customized correctly"
```

### Test 3: Agent with Specific Tools

**User prompt:**
```
Create a security-auditor agent that only has Read, Grep, and Glob tools
```

**Expected behavior:**
1. Claude creates agent
2. Claude adds `tools: Read, Grep, Glob` to frontmatter

**Validation:**
```bash
grep -q "^tools: Read, Grep, Glob$" .claude/agents/security-auditor.md || echo "FAIL: Tools not specified"
echo "PASS: Tools restriction added"
```

### Test 4: Orchestrator Agent Creation

**User prompt:**
```
Create an orchestrator agent with a professional, calm personality
```

**Expected behavior:**
1. Claude creates agent with `color: cyan` (per guidelines)
2. System prompt includes orchestrator responsibilities
3. Professional tone in communication style

**Validation:**
```bash
grep -q "^color: cyan$" .claude/agents/orchestrator.md || echo "FAIL: Orchestrator should use cyan"
grep -qi "delegate\|coordinate\|orchestrat" .claude/agents/orchestrator.md || echo "FAIL: Missing orchestrator language"
echo "PASS: Orchestrator created correctly"
```

### Test 5: User-Wide Agent Creation

**User prompt:**
```
Create a global doc-generator agent in my user directory
```

**Expected behavior:**
1. Claude runs: `.claude/skills/agent-creator/scripts/init_agent.sh doc-generator --path ~/.claude/agents`
2. Agent created in user's home directory

**Validation:**
```bash
test -f ~/.claude/agents/doc-generator.md || echo "FAIL: User-wide agent not created"
echo "PASS: User-wide agent created"
```

### Test 6: Error Handling - Invalid Model

**User prompt:**
```
Create an agent with model "gpt4"
```

**Expected behavior:**
1. Claude catches invalid model value
2. Claude suggests valid options: sonnet, opus, haiku, inherit
3. Claude asks user to choose valid model

**Validation:**
```bash
# Manual verification that Claude corrects the mistake
```

## Manual Test Procedure

1. **Start fresh project:**
   ```bash
   cd /tmp && mkdir test-project && cd test-project
   nix flake init -t github:netbrain/skeletons#claude
   nix develop
   ```

2. **Create agents directory:**
   ```bash
   mkdir -p .claude/agents
   ```

3. **Run each test case:**
   - Copy user prompt
   - Send to Claude via `claude "<prompt>"`
   - Observe Claude's actions
   - Run validation commands

4. **Document results:**
   - [ ] Test 1: Basic creation
   - [ ] Test 2: Personality customization
   - [ ] Test 3: Tool restrictions
   - [ ] Test 4: Orchestrator
   - [ ] Test 5: User-wide agent
   - [ ] Test 6: Error handling

## Expected Init Script Usage Patterns

Claude should:
- ✓ Use init script for all agent creation
- ✓ Ask about personality before creating
- ✓ Customize frontmatter based on requirements
- ✓ Write meaningful system prompts (not just templates)
- ✓ Use appropriate color based on agent type
- ✓ Select right model (haiku for simple, sonnet for complex)
- ✓ Remove TODO placeholders
- ✓ Explain how to test the agent

Claude should NOT:
- ✗ Create agent files manually without init script
- ✗ Skip frontmatter fields
- ✗ Leave system prompt as generic template
- ✗ Use invalid field values (model, color, priority)
- ✗ Forget to make agent executable/testable

## Common Patterns by Agent Type

### Test Runners
- Model: haiku (simple task)
- Color: green
- Tools: Bash, Read
- Priority: medium

### Code Reviewers
- Model: sonnet (needs reasoning)
- Color: blue
- Priority: high

### Security Auditors
- Model: sonnet
- Color: red
- Tools: Read, Grep, Glob
- Priority: critical

### Orchestrators
- Model: sonnet
- Color: cyan
- Priority: high
- Should delegate, not execute

## Troubleshooting

**If Claude doesn't use init script:**
- Check agent-creator skill activation keywords/patterns
- Verify init script is executable
- Check that skill frontmatter is complete

**If created agent is missing frontmatter fields:**
- Verify init_agent.sh template includes all fields
- Run: `.claude/skills/agent-creator/scripts/test_init_agent.sh`

**If agent personality is generic:**
- User may need to be more specific in request
- Claude should ask clarifying questions about tone/style
