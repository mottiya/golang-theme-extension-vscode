# Dark Modern Go Extension

VS Code Dark Modern theme with additional Go semantic highlighting.

The theme keeps the built-in Dark Modern palette as a base and adds Go-focused
rules for semantic tokens from `gopls`: functions, methods, structs,
interfaces, fields, parameters, labels, format strings, imported packages, and
function values.

## Usage

Recommended VS Code settings:

```json
{
	"workbench.colorTheme": "Dark Modern Go Extension",
	"editor.semanticHighlighting.enabled": true,
	"gopls": {
		"semanticTokens": true,
		"semanticTokenTypes": {
			"keyword": false
		}
	}
}
```

Disabling only `keyword` semantic tokens lets VS Code keep Dark Modern's
TextMate keyword split for Go, so `func`, `type`, `struct`, `interface`, and
`map` stay visually separate from control-flow keywords such as `return`, `if`,
`defer`, and `go`.

## Local Development

Open this repository in VS Code and press `F5`. In the Extension Development
Host window, run `Preferences: Color Theme` or press `Ctrl+K Ctrl+T`, then
select `Dark Modern Go Extension`.

For Go-focused changes, edit:

```text
themes/dark-modern-go-extension.json
```

Use `Developer: Inspect Editor Tokens and Scopes` on `code/go.go` to see which
TextMate scopes and semantic tokens affect a specific token.

You can also inspect `gopls` tokens from the terminal:

```bash
PATH=/usr/local/go/bin:/home/user22/go/bin:$PATH gopls semtok code/go.go
```

## Build And Install

Requires Node.js and pnpm.

```bash
pnpm install
pnpm run package
code --install-extension dark-modern-go-extension-2.9.0.vsix --force
```

After installing, reload VS Code and select `Dark Modern Go Extension`.

## Theme Files

- `themes/dark-modern-base.json` is the resolved Dark Modern base.
- `themes/dark-modern-go-extension.json` includes the base and adds Go-specific rules.
- `code/go.go` is a local sample file for checking Go token coverage.

## License

MIT.
