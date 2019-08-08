# <img src="alpaca.svg" width="150" /> Alpaca

Alpaca is a command line utility for building Alfred workflow bundles.

An alpaca project is an `alpaca.yml` file that defines the workflow, alongside any supporting files, such as scripts or images.

## Installation

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
    keyword: say
    with-space: true
    argument: required
    then:
      - object: say
  say:
    type: script
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
- [`object`](#object-schema) An map of objects in the Alfred workflow. Each key is an object name.
- `then` A list representing other objects to connect to, each item has this schema:
  - `object` The name of the object to connect to

### Object Schema

- `icon` A project-relative path to an icon for the object
- `type` The type of object this is. Currently partial support exists for:
  - [`keyword`](#keyword)
  - [`script`](#script)
  - [`script-filter`](#script-filter)

An object can have more properties depending on its type.

#### `keyword`

- `keyword` The keyword that triggers this object
- `with-space` Whether a space is required with this object
- `argument` A string determining whether an argument is required:
  - `required` The argument is required
  - `optional` The argument is optional
  - `none` No argument is accepted

#### `script`

- [`script`](#script-schema) A script configuration object

#### `script-filter`

- [`script`](#script-schema) A script configuration object

### Script Schema

There are a few types of script schemas possible.

#### Executable Script

This version executes the script at the given path (it must be executable).

- `path` The path to the script

#### Inline Script

This version executes an inline script.

- `content` The content of the script
- `type` The type of script, one of:
  - `bash`
  - `php`
  - `ruby`
  - `python`
  - `perl`
  - `zsh`
  - `osascript-as`
  - `osascript-js`

#### Inline Script from Path

This version executes an inline script, but pulls its contents in from a path in the project.

- `path` The path to the script
- `inline` Must be `true`
- `type` The type of script, one of:
  - `bash`
  - `php`
  - `ruby`
  - `python`
  - `perl`
  - `zsh`
  - `osascript-as`
  - `osascript-js`
