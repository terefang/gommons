# cmdbox

```
Please specify a valid subcommand, choices are:

        bash-completion Generate and output a bash completion-script.
        choose-file     Choose a file, interactively.
        choose-stdin    Choose an item from STDIN, interactively.
        chronic         Run a command quietly, if it succeeds.
        commands        Show all available sub-commands.
        gen-ca          generate a root ca key and certificate.
        gen-key         generate keys in pem format.
        gen-pass        Generate a random password and hashes.
        help            Show usage information.
        template        Populate a template-file.
        version         print version info
        xlua            extended lua environment
```

# gommons

`gommons` is an opinionated collection of common Go libraries 
and utilities for building applications the pragmatic way. 
The name is a play on Go and commons—because every project 
eventually grows the same set of well-tested building blocks.

## Versioning

`gommons` does not follow Semantic Versioning (SemVer). 
Instead, releases use a **date-based versioning scheme**:

```
YYYY.MM.N
```

Where:

* `YYYY` — four-digit release year
* `MM` — two-digit release month
* `N` — incremental release number within that month

For example:

```
2026.07.1
2026.07.2
2026.08.1
```

This scheme reflects the project's release cadence rather 
than attempting to encode compatibility guarantees into 
the version number. Breaking changes, new features, and 
bug fixes are documented in the release notes, making each 
release self-describing without relying on SemVer semantics.
