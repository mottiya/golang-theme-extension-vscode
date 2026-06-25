# Installation Manifest

This repository contains a VS Code color theme extension:

```text
Dark Modern Go Extension
```

It is based on VS Code Dark Modern and adds Go-specific semantic highlighting
for `gopls`.

## Requirements

- VS Code
- Official Go extension for VS Code
- `gopls`
- Node.js
- pnpm

## Install From Source

```bash
git clone https://github.com/mottiya/golang-theme-extension-vscode.git
cd golang-theme-extension-vscode
pnpm install
pnpm run package
code --install-extension dark-modern-go-extension-2.9.0.vsix --force
```

Reload VS Code after installation.

## Select Theme

Open VS Code command palette:

```text
Ctrl+Shift+P
```

Run:

```text
Preferences: Color Theme
```

Select:

```text
Dark Modern Go Extension
```

## Recommended Global VS Code Settings

Put these settings in the global VS Code user settings:

```text
~/.config/Code/User/settings.json
```

```json
{
	"workbench.colorTheme": "Dark Modern Go Extension",
	"editor.semanticHighlighting.enabled": true,
	"[go]": {
		"editor.semanticHighlighting.enabled": true
	},
	"go.useLanguageServer": true,
	"gopls": {
		"ui.semanticTokens": true,
		"ui.semanticTokenTypes": {
			"keyword": false
		}
	}
}
```

The `keyword: false` setting lets TextMate keep Dark Modern's keyword colors
for Go control-flow keywords, while semantic tokens still color functions,
methods, types, fields, parameters, labels, package names, and format strings.

If VS Code does not find `gopls`, install it first:

```bash
go install golang.org/x/tools/gopls@latest
```

If it is installed but not available in VS Code's `PATH`, point the Go extension
to it explicitly:

```json
{
	"go.alternateTools": {
		"gopls": "/home/your-user/go/bin/gopls"
	}
}
```

On this machine, the path is:

```text
/home/user22/go/bin/gopls
```

If you do not want global settings on another machine, put the same JSON into
the workspace-only file:

```text
.vscode/settings.json
```

## Development Test

Open this repository in VS Code and press `F5`. In the Extension Development
Host window, select `Dark Modern Go Extension`.

Use this sample file to inspect token coverage:

```text
code/go.go
```

Useful command:

```bash
PATH=/usr/local/go/bin:/home/user22/go/bin:$PATH gopls semtok code/go.go
```

## Included Theme Files

```text
themes/dark-modern-base.json
themes/dark-modern-go-extension.json
```

`dark-modern-base.json` is included intentionally because VS Code themes cannot
reliably inherit from built-in Dark Modern by theme name.
