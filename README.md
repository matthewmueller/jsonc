# jsonc

This package is a fork of [hujson](github.com/tailscale/hujson) that works better with VSCode's built-in [JSON with Comments (jsonc)](https://code.visualstudio.com/docs/languages/json#_json-with-comments) format. This is the same format you'll find on VSCode's Settings and Task files.

This fork adjusts the formatting to match both:

- VSCode's [built-in JSON formatter](https://code.visualstudio.com/docs/languages/json#_formatting)
- The [Prettier](https://prettier.io/) formatter

The specific changes were:

- Does not add a trailing comma `,` on the last field of an object or array
- Switches to using two spaces for indentation instead of tabs
- Removes Go-style field alignment

## Usage in VSCode

VSCode natively recognizes the `.jsonc`.

## Thanks

Many thanks to the [Tailscale](https://github.com/tailscale) team for creating this wonderful package.
