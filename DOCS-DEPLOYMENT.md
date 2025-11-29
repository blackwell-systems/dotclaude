# Documentation Deployment

The dotclaude documentation site is built with Docsify and hosted on GitHub Pages.

## Viewing Locally

### Option 1: Docsify CLI (Recommended)

```bash
# Install docsify-cli globally
npm install -g docsify-cli

# Serve the docs
cd ~/code/dotclaude
docsify serve .

# Open http://localhost:3000
```

### Option 2: Python HTTP Server

```bash
cd ~/code/dotclaude
python3 -m http.server 3000

# Open http://localhost:3000
```

### Option 3: PHP Server

```bash
cd ~/code/dotclaude
php -S localhost:3000

# Open http://localhost:3000
```

## GitHub Pages Deployment

The site is automatically deployed to GitHub Pages from the `main` branch.

**URL:** https://blackwell-ai.github.io/dotclaude/

### Setup (One-Time)

1. Go to repository Settings → Pages
2. Source: Deploy from a branch
3. Branch: `main` / `(root)`
4. Save

GitHub Pages will automatically build and deploy when you push to main.

### Files

- `index.html` - Docsify configuration and theme
- `_coverpage.md` - Landing page content
- `_sidebar.md` - Navigation structure
- `.nojekyll` - Tells GitHub Pages not to use Jekyll
- `README.md` - Homepage content
- `docs/` - All documentation markdown files

### Theme

Uses Blackwell Systems™ shared Docsify theme:
- Base: vue.css (Docsify default)
- Branding: https://blackwell-systems.github.io/blackwell-docs-theme/docsify.css

Project-specific color overrides in `index.html`:
```css
:root {
  --theme-color: #22c55e;  /* Green for dotclaude */
  --theme-color-dark: #16a34a;
}
```

## Customization

### Theme Colors

Edit `index.html` to change theme colors:

```css
:root {
  --theme-color: #your-color;
  --theme-color-dark: #your-dark-color;
}
```

### Navigation

Edit `_sidebar.md` to modify sidebar navigation.

### Cover Page

Edit `_coverpage.md` to update the landing page.

### Adding New Pages

1. Create markdown file in appropriate location
2. Add entry to `_sidebar.md`
3. Commit and push

GitHub Pages will auto-deploy within ~1 minute.

## Troubleshooting

**Site not updating after push:**
- Wait 1-2 minutes for GitHub Pages build
- Check Actions tab for build status
- Hard refresh browser (Ctrl+Shift+R / Cmd+Shift+R)

**404 errors:**
- Ensure `.nojekyll` file exists
- Check branch settings in GitHub Pages config
- Verify files are in root directory (not in a subdirectory)

**Styles not loading:**
- Check browser console for CSS loading errors
- Verify CDN links are accessible
- Clear browser cache

## Resources

- [Docsify Documentation](https://docsify.js.org/)
- [Blackwell Docs Theme](https://github.com/blackwell-systems/blackwell-docs-theme)
- [GitHub Pages Documentation](https://docs.github.com/en/pages)
