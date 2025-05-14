<div align="center">
    <h1>chlog</h1>
    <p>AI-powered changelog generation from Git history</p>
</div>

`chlog` is a CLI tool that uses LLMs to generate clean, structured changelogs based on your Git commits and diffs. It outputs entries in a consistent JSON format. Whether you’re prepping for a release or catching up on recent changes, `chlog` helps you maintain changelogs with minimal manual effort.

## Table of Contents
- [✨ Features](#-features)
- [📦 Installation](#-installation)
- [📄 JSON Format](#-json-format)
- [🚀 Quick Start](#-quick-start)
  * [Example Output](#example-output)
- [🔧 Usage](#-usage)
  * [`chlog`](#chlog)
  * [`chlog generate`](#chlog-generate)
    + [Flags](#flags)
    + [Config File](#config-file)
  * [`chlog models`](#chlog-models)
- [🧠 Design Rationale](#-design-rationale)

## ✨ Features
- Uses commit diffs and messages to generate changelog entries via LLMs
- Outputs structured JSON changelogs
- Supports yaml config files for repeatable changelog generation
- Verbose output (without polluting stdout)
- Configurable model/provider support (currently OpenAI and Gemini)

## 📦 Installation
You can install `chlog` with:
```bash
go install github.com/ammar-ahmed22/chlog@latest
```
> [!IMPORTANT]
> Make sure `$GOPATH/bin` is in your `PATH` to run the `chlog` command.

Alternatively, you can build from source:
```bash
git clone https://github.com/ammar-ahmed22/chlog.git
cd chlog
go install ./...
```
## 📄 JSON Format
The generated changelog entries are in an opinionated JSON format. Here’s an example of what a changelog entry looks like:
```json
  {
    "version": "0.2.0",
    "date": "2025-05-14",
    "from_ref": "v0.1.1",
    "to_ref": "v0.2.0",
    "changes": [
      {
        "id": "support-optional-config-file-for-generate",
        "title": "Support optional config file for generate command flags",
        "description": "...",
        "impact": "...",
        "commits": [
          "c08bd2429546546d992c5364338051d0a5d4edf6"
        ],
        "tags": [
          "feature"
        ]
      },
      {
        "id": "add-version-date-and-git-range-to",
        "title": "Add version, date, and git range to changelog entries manually",
        "description": "...",
        "impact": "...",
        "commits": [
          "804391e4e40b0182605b801b7a9d68120c893410"
        ],
        "tags": [
          "feature",
          "changed"
        ]
      }
    ]
  },
``` 

## 🚀 Quick Start
Generate a prettified changelog entry for the last commit with verbose output:
```bash
chlog generate 0.2.0 --from HEAD~1 --to HEAD --pretty --verbose
```

### Example Output
```bash
> chlog generate 0.2.0 --from HEAD~1 --to HEAD --pretty --verbose

Generating changelog entry "0.2.0" for commits:
e2078de chore: readme

Starting AI changelog generation (provider: openai, model: gpt-4o-mini)
Completed AI changelog generation
tokens used: 1521 (input: 1386, output: 135)
{
  "version": "0.2.0",
  "date": "2025-05-14",
  "from_ref": "HEAD~1",
  "to_ref": "HEAD",
  "changes": [
    {
      "id": "update-readme-description",
      "title": "Update README Description",
      "description": "The README file was updated to clarify that `chlog` is an AI-powered tool for generating changelogs from Git history. This enhancement aims to better inform users about the primary functionality and capabilities of `chlog`.",
      "impact": "Users will have a more accurate understanding of what `chlog` does and how it can help them manage changelogs with minimal effort.",
      "commits": [
        "e2078de294fad79058a6a92d137975b296698682"
      ],
      "tags": [
        "documentation"
      ]
    }
  ]
}
```
> [!TIP]
> You can write the output to a file using `> changelog.json` since verbose output is written to `stderr`.  
> e.g. `chlog generate 0.2.0 --from HEAD~1 --to HEAD --pretty --verbose > changelog.json`

> [!TIP]
> You can use any valid git reference for the `--from` and `--to` flags. So, if you have SemVer tagged commits, you can do something like: `chlog generate 0.2.0 --from v0.1.0 --to v0.2.0` to generate a changelog entry from `v0.1.0` to `v0.2.0`.

## 🔧 Usage
### `chlog` 
```bash
chlog --help
```
```bash
chlog is a command-line tool that uses AI to generate and update structured changelogs from your Git history.

It can automatically summarize changes based on diffs and commit messages and output structured changelogs in JSON formats.
Use it to keep your changelogs clean, consistent, and up-to-date.

Example:
	chlog generate 0.2.0 --from HEAD~10 --to HEAD > changelog.json

Usage:
  chlog [command]

Available Commands:
  completion  Generate the autocompletion script for the specified shell
  generate    Generates the AI-powered changelog entry for the specified version
  help        Help about any command
  models      List supported models and providers

Flags:
  -h, --help   help for chlog

Use "chlog [command] --help" for more information about a command.
```

### `chlog generate`
```bash
chlog generate --help
```
```bash
Generates the AI-powered changelog entry for the specified version

Usage:
  chlog generate <VERSION> (default "2025-05-14") [flags]

Flags:
      --apiKey string     API key for the LLM provider (can also be set via environment variable, see chlog models for details)
  -c, --config string     Path to config file (optional, chlog.yaml will be loaded if present in the current directory)
  -d, --date string       Date for the changelog entry in YYYY-MM-DD format (default "2025-05-14")
      --file string       Path to existing changelog JSON file to update with the new entry (should be an array of changelog entries or empty file)
  -f, --from string       Starting commit reference (e.g. HEAD~3, main, v1.0.0, or abc1234) (default "HEAD~1")
  -h, --help              help for generate
  -m, --model string      LLM model (see chlog models for available options and defaults)
      --pretty            Prettified JSON output
  -p, --provider string   LLM provider (see chlog models for available options) (default "openai")
  -t, --to string         Ending commit reference (e.g. HEAD~3, main, v1.0.0, or abc1234) (default "HEAD")
  -v, --verbose           Enable verbose output
```

#### Flags
| Flag                 | Description                                                                                                     | Set via Config? |
|----------------------|-----------------------------------------------------------------------------------------------------------------|:---------------:|
| `--apiKey`           | API key for the LLM provider. <br>(can also be set via environment variable, use `chlog models` to see details) |        ✅        |
| `--config`<br>`-c`   | Optional path a YAML config file. <br>(`chlog.yaml` is loaded automatically if found in the current directory)  |                 |
| `--date`<br>`-d`     | Date of the entry in `YYYY-MM-DD` format (default: today)                                                       |                 |
| `--file`             | Path to changelog file to update with the generated entry.                                                      |        ✅        |
| `--from`<br>`-f`     | Starting Git reference. <br>Can be any valid git reference including branches, tags, etc. (default: `HEAD~1`)   |                 |
| `--to`<br>`-t`       | Ending Git reference. <br>Can be any valid git reference including branches, tags, etc. (default: `HEAD`)       |                 |
| `--provider`<br>`-p` | LLM provider to use. <br>See `chlog models` to see available providers (default: `openai`)                      |        ✅        |
| `--model`<br>`-m`    | LLM model to use. <br>See `chlog models` to see available models for the selected provider                      |        ✅        |
| `--pretty`           | Format JSON output with indentation                                                                             |        ✅        |
| `--verbose, -v`      | Output verbose output to `stderr`                                                                               |        ✅        |

> [!IMPORTANT]
> When `--file` is specified, `chlog` will:
> - Create an array if the file is empty or does not exist.
> - Prepend the new entry to the existing array.  
>
> The existing JSON must adhere to the format specified in the [JSON Format](#-json-format) section.

#### Config File
You can use a config file (`chlog.yaml` in the current directory) or any other file you specify with the `--config` flag to avoid repeating flags:

Below is an example specifiying all valid options (any of these can be omitted):
```yaml
apiKey: my-api-key
provider: openai
model: gpt-4o-mini
verbose: true
pretty: true
file: ./changelog.json
```

> [!NOTE]
> Any flags specified when calling `chlog generate` will override the config file options.

> [!NOTE]
> The `file` key in the config is relative to the config file.

### `chlog models`
```bash
chlog models
```
```bash
Supported providers and models:

Provider: openai (env var: OPENAI_API_KEY)
 - gpt-4o-mini (default)
 - gpt-4.1-mini

Provider: gemini (env var: GEMINI_API_KEY)
 - gemini-2.0-flash (default)
```
## 🧠 Design Rationale
This section outlines some of the key technical and product decisions made during the development of chlog.

### CLI

Developer tools should meet developers where they work, the terminal. A CLI is intuitive, scriptable, and integrates naturally into developer workflows. By making chlog a CLI-first tool, it can easily be used in CI/CD pipelines (e.g., generating changelogs on tagged releases via GitHub Actions) and avoids the friction of a GUI for users who prefer automation and speed.

### JSON Output

Most modern LLMs have strong native support for JSON, often with explicit structured output modes in their SDKs. This makes JSON the most reliable and developer-friendly format for AI-generated data. The schema used by `chlog` is intentionally opinionated to strike a balance between structure and readability, while still capturing key user-facing details. While custom schemas could be valuable, enforcing a consistent format simplifies validation, enables frontend integrations, and supports potential monetization through standardized rendering.

If I were to continue working on this for longer (which I probably will), I would probably let users define their own schemas optionally. However, that would require some heavier prompt engineering and validation so it was left out for this POC.

### Multiple LLM Providers

Developers often have strong preferences for specific LLM providers, whether due to cost, performance, or available credits (e.g., GCP credits for Gemini). For this reason, `chlog` is configurable with support for both OpenAI and Gemini, making it more flexible and adaptable to different use cases.

As a POC, I only added support for OpenAI and Gemini with a few models. However, the modular design allows for "easy" addition of other providers and models. It would just extend my deadline haha.

### YAML for Configuration

YAML is a well-established standard in the dev ecosystem for config files; readable, writable, and supported by many tools. Supporting it allows developers to define default settings without cluttering the CLI call. That said, the implementation is modular enough to support additional config formats (like JSON or TOML) in the future if needed.

### Logging to stderr

Structured JSON should flow cleanly to stdout so it can be piped into a file or another tool. Verbose logs, like commit info, model usage, and token stats go to stderr to avoid polluting the JSON output. This separation is intentional and aligns with best practices for CLI tools.

### Prepending Entries 

Changelog consumers care most about recent changes. By prepending entries, the latest updates appear at the top, saving users from scrolling through potentially long files to find what’s new.

### Using Go

Go offers an ideal balance of developer ergonomics, type safety, and performance. It’s simple to write and reason about, has first-class support for JSON, and compiles to a single binary; making installation and distribution trivial. Go is particularly well-suited to building CLIs, and its standard library combined with third-party packages made development fast and robust.
