# 0.3.1 / 2026-02-01

- add an alias

# 0.3.0 / 2026-02-01

- add `UnmarshalExpanded` and `PatchExpanded`. These allow you to preserve dynamic
  variables, while making other changes to the schema (see test for an example).

# 0.2.0 / 2026-02-01

- add higher-level patching functions
- fix some linting and modernize syntax

# 0.1.0 / 2024-11-28

- remove `hujsonfmt` command
- rename hujson => jsonc
- remove object alignment
- switch to using two spaces instead of tabs
- remove trailing comma from format to be compatible with vscode's jsonc
- Fork from https://github.com/tailscale/hujson
