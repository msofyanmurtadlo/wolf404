# Wolf404 VS Code Extension Installation

## Quick Install

1. **Copy extension folder**:

   ```powershell
   Copy-Item -Recurse "wolf404-vscode" "$env:USERPROFILE\.vscode\extensions\wolf404-language-1.0.1"
   ```

2. **Restart VS Code**

3. **Enable File Icon Theme**:
   - Press `Ctrl+Shift+P`
   - Type "File Icon Theme"
   - Select "A Wolf404 Icons"

## Manual Install

1. Open VS Code
2. Press `Ctrl+Shift+X` to open Extensions
3. Click "..." menu → "Install from VSIX"
4. Navigate to `wolf404-vscode` folder
5. Select the extension

## Verify Installation

After installation, `.wlf` files should show the Wolf404 logo icon in:

- File Explorer sidebar
- Editor tabs
- File tree

## Troubleshooting

If icons don't appear:

1. Restart VS Code completely
2. Check File Icon Theme is set to "A Wolf404 Icons"
3. Reload window: `Ctrl+Shift+P` → "Reload Window"
