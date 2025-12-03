# dotclaude Onboarding Audit

**Date**: 2025-12-02
**Focus**: New user experience, installation clarity, first-run workflow

## Critical Issues

### 1. README Install Method Doesn't Work ⛔

**Problem**: README shows this:

```bash
curl -fsSL https://raw.githubusercontent.com/blackwell-systems/dotclaude/main/install.sh | bash
```

**Why it fails**:
- `install.sh` uses `REPO_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"`
- When piped to bash, there's no repo cloned
- Script can't find `base/`, `profiles/`, `examples/`
- Installation fails

**What users see**: Errors about missing directories

**Expected flow**: Git clone, then run install.sh from within repo

**Fix needed**: Either fix the script to handle curl | bash, or change README to show git clone method only

---

### 2. Conflicting Install Instructions 🔴

**README says**:
```bash
curl -fsSL ... | bash
Then create your first profile:
cp -r examples/sample-profile profiles/my-project
```

**GETTING-STARTED.md says**:
```bash
git clone ... ~/code/dotclaude
cd ~/code/dotclaude
./install.sh
```

**Problems**:
- Two different methods with different outcomes
- README method (curl) doesn't give you examples/
- No explanation of which to use or why
- User confusion: "Do I git clone or curl?"

**Fix needed**: Pick ONE install method and use it consistently

---

### 3. Missing "What's Next" After Install 🟡

**Current**: Install completes with:

```
✓ Installation Complete
🍃 Tip: Run 'dotclaude help' to see all commands
```

**Problems**:
- Doesn't tell user what to do next
- Shows 4 commands but no guidance on which to run first
- No clear path to "create your first profile"

**Better**:
```
✓ Installation Complete

Next steps:
  1. Create your first profile:
     dotclaude create my-project

  2. Activate it:
     dotclaude activate my-project

  3. Verify:
     dotclaude show

Run 'dotclaude help' for all commands.
```

---

### 4. Profile Creation Unclear 🟡

**README says**: `cp -r examples/sample-profile profiles/my-project`

**Problems**:
- `dotclaude create` command exists but isn't mentioned
- Manual cp is error-prone (where is examples/?)
- No guidance on what to edit after copying

**Fix needed**:
- Document `dotclaude create <name>` as primary method
- Mention cp -r as alternative for advanced users
- Show what to do after creation

---

### 5. DOTCLAUDE_REPO_DIR Confusing 🟠

**Install script ends with**:
```
Set DOTCLAUDE_REPO_DIR in your shell (optional):
  export DOTCLAUDE_REPO_DIR="$REPO_DIR"
```

**Problems**:
- Not clear when this is needed
- Not clear what happens if you don't set it
- "Optional" but when would you need it?

**Fix needed**: Explain the use case or remove if truly optional

---

## Medium Priority Issues

### 6. First Profile Creation Has No Template 🟠

**Current**: User must either:
- `cp -r examples/sample-profile` (where is examples?)
- `dotclaude create` (creates empty profile)

**Problem**: New users don't know what a profile should contain

**Fix needed**: `dotclaude create` should scaffold from sample-profile template

---

### 7. No Validation After Install 🟠

**Current**: Install script assumes everything worked

**Better**: End with a validation check:
```bash
# After install
echo "Validating installation..."
if command -v dotclaude >/dev/null; then
  echo "✓ dotclaude CLI installed"
else
  echo "✗ dotclaude not in PATH"
fi

if [ -d "$CLAUDE_DIR/scripts" ]; then
  echo "✓ Management scripts installed"
fi
```

---

### 8. Profile Activation Skipped = Dead End 🟠

**Install asks**: "Which profile would you like to activate? (or 'skip' to skip)"

**If user skips**:
- They see: "To activate a profile later: dotclaude activate <profile-name>"
- But they have NO profiles created yet
- Next step should be: "Create your first profile with: dotclaude create <name>"

---

## Minor Issues

### 9. Help Text Inconsistent 🟡

**dotclaude help** shows all commands but:
- No clear grouping (profile management vs git tools vs status)
- Some commands have descriptions, others don't
- No "most common commands" section for new users

---

### 10. Example Profile Not Self-Documenting 🟡

`examples/sample-profile/CLAUDE.md` should have:
- Comments explaining each section
- Examples of what to customize
- Link to full documentation

---

## Recommendations

### Immediate Fixes (Breaking Issues)

1. **Fix README install method**
   - Change to git clone only (safest)
   - OR make install.sh work with curl | bash

2. **Unify install instructions**
   - Pick one method
   - Use it in README and GETTING-STARTED
   - Document alternative methods in FAQ

3. **Improve post-install messaging**
   - Clear "what's next" steps
   - Numbered flow: create → edit → activate → verify

### High Priority Improvements

4. **Make `dotclaude create` use template**
   - Copy from examples/sample-profile automatically
   - OR bundle template in base/

5. **Add install validation**
   - Check PATH includes ~/.local/bin
   - Verify dotclaude CLI works
   - Confirm base files copied

6. **Better first-run experience**
   - If no profiles exist, guide to `dotclaude create`
   - If base not customized, suggest editing base/CLAUDE.md first

### Nice to Have

7. **Interactive first-run wizard**
   ```bash
   dotclaude init
   → Creates sample profile
   → Guides through customization
   → Activates automatically
   ```

8. **Onboarding checklist in docs**
   ```
   ✓ Installed base configuration
   ✓ Created first profile
   ✓ Customized base standards
   ✓ Activated and verified
   ```

---

## Testing the Current Experience

I simulated a new user install:

```bash
# Step 1: Follow README
curl -fsSL https://raw.githubusercontent.com/blackwell-systems/dotclaude/main/install.sh | bash
# → FAILS: Can't find base/ directory

# Step 2: Try GETTING-STARTED
git clone https://github.com/blackwell-systems/dotclaude.git ~/code/dotclaude
cd ~/code/dotclaude
./install.sh
# → Works! But then what?

# Step 3: Follow README's "create first profile"
cp -r examples/sample-profile profiles/my-project
# → Works, but feels manual

# Step 4: Activate
dotclaude activate my-project
# → Works!

# Conclusion: Git clone method works, curl method doesn't
```

---

## Proposed Install Flow (Ideal)

```bash
# 1. Install
git clone https://github.com/blackwell-systems/dotclaude.git ~/.dotclaude
cd ~/.dotclaude
./install.sh

# 2. Create first profile (with template)
dotclaude create my-project
# → Copies from examples/sample-profile automatically
# → Opens in $EDITOR

# 3. Activate
dotclaude activate my-project
# → Merges base + profile to ~/.claude/CLAUDE.md

# 4. Verify
dotclaude show
# → Shows active profile and merged content preview
```

---

## Summary

**Blockers**:
- README install method doesn't work (curl | bash fails)
- Conflicting install instructions across docs

**High Priority**:
- No clear "what's next" after install
- Profile creation workflow unclear
- `dotclaude create` should use template

**Fix Priority**:
1. Fix install method (immediate)
2. Unify documentation (immediate)
3. Improve post-install messaging (high)
4. Make create use template (high)
5. Add validation checks (medium)
