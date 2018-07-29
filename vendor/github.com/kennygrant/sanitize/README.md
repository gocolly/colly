sanitize [![GoDoc](https://godoc.org/github.com/kennygrant/sanitize?status.svg)](https://godoc.org/github.com/kennygrant/sanitize) [![Go Report Card](https://goreportcard.com/badge/github.com/kennygrant/sanitize)](https://goreportcard.com/report/github.com/kennygrant/sanitize) [![CircleCI](https://circleci.com/gh/kennygrant/sanitize.svg?style=svg)](https://circleci.com/gh/kennygrant/sanitize)
========

Package sanitize provides functions to sanitize html and paths with go (golang).

FUNCTIONS


```go
sanitize.Accents(s string) string
```

Accents replaces a set of accented characters with ascii equivalents.

```go
sanitize.BaseName(s string) string
```

BaseName makes a string safe to use in a file name, producing a sanitized basename replacing . or / with -. Unlike Name no attempt is made to normalise text as a path.

```go
sanitize.HTML(s string) string
```

HTML strips html tags with a very simple parser, replace common entities, and escape < and > in the result. The result is intended to be used as plain text. 

```go
sanitize.HTMLAllowing(s string, args...[]string) (string, error)
```

HTMLAllowing parses html and allow certain tags and attributes from the lists optionally specified by args - args[0] is a list of allowed tags, args[1] is a list of allowed attributes. If either is missing default sets are used. 

```go
sanitize.Name(s string) string
```

Name makes a string safe to use in a file name by first finding the path basename, then replacing non-ascii characters.

```go
sanitize.Path(s string) string
```

Path makes a string safe to use as an url path.


Changes
-------

Version 1.2

Adjusted HTML function to avoid linter warning
Added more tests from https://githubengineering.com/githubs-post-csp-journey/
Chnaged name of license file
Added badges and change log to readme

Version 1.1
Fixed type in comments. 
Merge pull request from Povilas Balzaravicius Pawka 
 - replace br tags with newline even when they contain a space

Version 1.0
First release