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
