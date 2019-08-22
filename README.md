# <img src="alpaca.svg" width="150" /> Alpaca

[![](https://github.com/jclem/alpaca/workflows/.github/workflows/ci.yml/badge.svg)](https://github.com/jclem/alpaca/actions)

Alpaca is a command line utility for building Alfred workflow bundles.

An alpaca project is an `alpaca.yml` file that defines the workflow, alongside any supporting files, such as scripts or images.

**Note:** Alpaca is still in the proof-of-concept phase. Huge portions of Alfred functionality are unimplemented.

<details>
<summary><strong>Contents</strong></summary>

- [Installation](#installation)
- [Usage](#usage)
  - [`alpaca pack`](#alpaca-pack-dir)
- [Schema](#schema)
  - [Example](#example)
  - [Root Schema](#root-schema)
  - [Object Schema](#object-schema)
    - [`applescript`](#applescript)
    - [`clipboard`](#clipboard)
    - [`keyword`](#keyword)
    - [`open-url`](#open-url)
    - [`script`](#script)
    - [`script-filter`](#script-filter)
  - [Script Schema](#script-schema)
    - [Executable Script](#executable-script)
    - [Inline Script](#inline-script)

</details>

## Installation

[Download the latest release](https://github.com/jclem/alpaca/releases/latest), or:

```shell
$ go get github.com/jclem/alpaca
```

## Usage

### `alpaca pack <dir>`

Pack an Alpaca project into an Alfred workflow. The workflow will be output into the current directory.

```shell
$ alpaca pack .
```

## Schema

### Example

This workflow uses a trigger "say" to say words passed as arguments. For example, the user might type "say hello world" into Alfred.

```yaml
name: Say

objects:
  trigger:
    type: keyword
    config:
      keyword: say
      with-space: true
      argument: required
    then: say

  say:
    type: script
    config:
      script:
        content: say {query}
        type: bash
```

### Root Schema

- `name` The name of the project
- `version` The version of the project
- `author` The author of the project (i.e. `Jonathan Clem <jonathan@example.com>`)
- `bundle-id` The Alfred workflow bundle ID
- `description` A short description of the workflow
- `readme` A longer description of the workflow, seen when users import it
- `url` A homepage URL for the workflow
- `icon` A project-relative path to an icon to use for the worflow
- `variables` A map of variable names and their default values
- [`object`](#object-schema) An map of objects in the Alfred workflow. Each key is an object name.

### Object Schema

- `icon` A project-relative path to an icon for the object
- `type` The type of object this is. Currently partial support exists for:
  - [`applescript`](#applescript)
  - [`clipboard`](#clipboard)
  - [`keyword`](#keyword)
  - [`open-url`](#open-url)
  - [`script`](#script)
  - [`script-filter`](#script-filter)
- `config` A type-specific configuration object, see each type schema for details
- `then` A string, list of strings, or a list of objects representing other objects to connect to, each objects having this schema:
  - `object` The name of the object to connect to

#### `applescript`

- `cache` (`bool`, default `true`) Whether to cache the compiled AppleScript
- [`script`](#script-schema) A script configuration object, but only `content` is respected

#### `clipboard`

- `text` (`string`, default `"{query}"`) The text to copy to the clipboard—use `"{query}"` for the exact query

#### `keyword`

- `keyword` (`string`) The keyword that triggers this object
- `with-space` (`bool`) Whether a space is required with this object
- `title` (`string`) The title of the object
- `subtitle` (`string`) The subtitle of the object
- `argument` (`string`) A string determining whether an argument is required. One of:
  - `required` The argument is required
  - `optional` The argument is optional
  - `none` No argument is accepted

#### `open-url`

- `url` (`string`) The URL to open. Use `"{query}"` for the exact query

#### `script`

- [`script`](#script-schema) A script configuration object

#### `script-filter`

- `argument` (`string`) A string determining whether an argument is required. One of:
  - `required` The argument is required
  - `optional` The argument is optional
  - `none` No argument is accepted
- `argument-trim` (`string`) Whether to trim the argument's whitespace. One of:
  - `auto` Argument automatically trimmed
  - `off` Argument not trimmed
- `escaping` (`[]string`) A list of strings to escape in the query text, if the script object's `arg-type` is `query`. Any of:
  - `spaces`
  - `backquotes`
  - `double-quote`
  - `brackets`
  - `semicolons`
  - `dollars`
  - `backslashes`
- `ignore-empty-argument` (`bool`) Whether an empty argument (when `arg-type` is `argv` in script config) is omitted from `argv` (when `false`, it will be an empty string)
- `keyword` (`string`) The keyword that triggers this object
- `running-subtitle` (`string`) A subtitle to display while this filter runs
- `subtitle` (`string`) A subtitle for this object
- `title` (`string`) A title for this object
- `with-space` (`bool`) Whether a space is required with this object
- `alfred-filters-results` An object describing how Alfred filters results (it does not if this is omitted):
  - `mode` (`string`) The mode Alfred uses to filter results. One of:
    - `exact-boundary`
    - `exact-start`
    - `word-match`
- `run-behavior` An object describing behavior of the script run
  - `immediate` (`bool`) Whether to run immediately always for the first character typed
  - `queue-mode` (`string`) A mode for how script runs are queued. One of:
    - `wait` Wait for previous run to complete
    - `terminate` Terminate previous run immediately and start a new one
  - `queue-delay` (`string`) A delay mode for queueing script runs. One of:
    - `immediate` No delay
    - `automatic` Automatic after last character typed
    - `100ms` 100ms after last character typed
    - `200ms` 200ms after last character typed
    - `300ms` 300ms after last character typed
    - `400ms` 400ms after last character typed
    - `500ms` 500ms after last character typed
    - `600ms` 600ms after last character typed
    - `700ms` 700ms after last character typed
    - `800ms` 800ms after last character typed
    - `900ms` 900ms after last character typed
    - `1000ms` 1000ms after last character typed
- [`script`](#script-schema) A script configuration object

### Script Schema

There are a few types of script schemas possible, in addition to these options:

#### Executable Script

This version executes the script at the given path (it must be executable).

- `path` (`string`) The path to the script

#### Inline Script

This version executes an inline script.

- `arg-type` (`string`) The way the argument is passed to the script. One of:
  - `query` Interpolated as (`query`)
  - `argv` Passed into process arguments
- `content` (`string`) The content of the script
- `type` (`string`) The type of script, one of:
  - `bash`
  - `php`
  - `ruby`
  - `python`
  - `perl`
  - `zsh`
  - `osascript-as`
  - `osascript-js`
