#!/bin/bash

# Update repository name references
# Usage: ./update-repo-name.sh <org> <repo-name>

ORG="${1:-blackwell-ai}"
REPO="${2:-dotclaude}"

echo "Updating repo references to: $ORG/$REPO"
echo ""

# Files to update
FILES=(
    "_coverpage.md"
    "_sidebar.md"
    "index.html"
    "DOCS-DEPLOYMENT.md"
    "README.md"
)

for file in "${FILES[@]}"; do
    if [ -f "$file" ]; then
        echo "Updating: $file"
        # Update github.com URLs
        sed -i "s|github\.com/blackwell-ai/dotclaude|github.com/$ORG/$REPO|g" "$file"
        sed -i "s|github\.com/blackwell-ai/claude-globals|github.com/$ORG/$REPO|g" "$file"

        # Update github.io URLs
        sed -i "s|blackwell-ai\.github\.io/dotclaude|$ORG.github.io/$REPO|g" "$file"
        sed -i "s|blackwell-ai\.github\.io/claude-globals|$ORG.github.io/$REPO|g" "$file"
    fi
done

echo ""
echo "✓ Updated all references to: $ORG/$REPO"
echo ""
echo "Review changes with: git diff"
echo "Commit with: git add -A && git commit -m 'Update repository name references'"
