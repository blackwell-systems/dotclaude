# GitHub Pages 404 Troubleshooting

Your Docsify files are committed and pushed. The 404 means GitHub Pages needs configuration.

## Step-by-Step Fix

### 1. Enable GitHub Pages

Go to: **https://github.com/blackwell-systems/dotclaude/settings/pages**

You should see:

```
Build and deployment
├─ Source: Deploy from a branch
│  └─ Branch: [Select branch ▼]  [Select folder ▼]
```

### 2. Configure Settings

Set these values:
- **Source:** Deploy from a branch
- **Branch:** `main`
- **Folder:** `/ (root)`

Click **Save**

### 3. Wait for Build

After clicking Save:
- GitHub will start building (takes 1-2 minutes)
- You'll see a message: "Your site is ready to be published at..."
- Or check: https://github.com/blackwell-systems/dotclaude/actions

### 4. Check Build Status

Go to: **https://github.com/blackwell-systems/dotclaude/actions**

You should see:
- ✅ "pages build and deployment" workflow
- ✅ Green checkmark when done

### 5. Access Your Site

Once the build completes (green checkmark):
**https://blackwell-systems.github.io/dotclaude/**

## Common Issues

### Issue: "There isn't a GitHub Pages site here"

**Fix:** GitHub Pages is not enabled. Follow steps 1-2 above.

### Issue: 404 after enabling

**Causes:**
1. Build is still in progress (wait 2 minutes)
2. Wrong branch selected (use `main`)
3. Wrong folder selected (use `/ (root)`)

**Check:**
- Actions tab shows green checkmark
- Settings → Pages shows green "Your site is live at..."

### Issue: Blank page instead of 404

**Cause:** Files are there but Docsify isn't loading

**Fix:**
- Hard refresh: Ctrl+Shift+R (Windows) or Cmd+Shift+R (Mac)
- Check browser console for errors

## Verification

### Files on GitHub

Check these files exist in your repo:
- https://github.com/blackwell-systems/dotclaude/blob/main/index.html
- https://github.com/blackwell-systems/dotclaude/blob/main/.nojekyll
- https://github.com/blackwell-systems/dotclaude/blob/main/_coverpage.md
- https://github.com/blackwell-systems/dotclaude/blob/main/_sidebar.md

All should show ✅

### Current Status

```bash
# Check locally
cd ~/code/CLAUDE
ls -la | grep -E "index.html|\.nojekyll|_coverpage|_sidebar"

# Should show:
# -rw-r--r-- 1 user user    0 Nov 29 16:11 .nojekyll
# -rw------- 1 user user  612 Nov 29 16:28 _coverpage.md
# -rw------- 1 user user 1621 Nov 29 16:28 _sidebar.md
# -rw------- 1 user user 4176 Nov 29 16:23 index.html
```

## What Should Happen

1. Visit https://github.com/blackwell-systems/dotclaude/settings/pages
2. Enable Pages from main branch, root folder
3. Wait 1-2 minutes
4. See green "Your site is live at..." message
5. Visit https://blackwell-systems.github.io/dotclaude/
6. See beautiful Docsify site with dark cover page

## Still Not Working?

Check these:

**Repository visibility:**
- Repo must be **public** for free GitHub Pages
- Or you need GitHub Pro for private repo pages

**Branch protection:**
- Ensure `main` branch exists
- Ensure you have push access

**Actions enabled:**
- Settings → Actions → General
- Ensure "Allow all actions" is enabled

## Debug Commands

```bash
# Verify files are pushed
git ls-files | grep -E "index.html|\.nojekyll|_coverpage|_sidebar"

# Check remote
git remote -v
# Should show: git@github.com:blackwell-systems/dotclaude.git

# Verify latest commit is pushed
git log origin/main -1 --oneline
# Should show: 8bed353 Remove employer references for public release
```

---

**If you've enabled GitHub Pages and waited 2 minutes, and it's still 404, let me know what you see in:**
1. Settings → Pages (screenshot would help)
2. Actions tab (any failures?)
