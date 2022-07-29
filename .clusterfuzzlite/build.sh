#!/bin/bash -eu

go get github.com/AdamKorcz/go-118-fuzz-build/utils
go get github.com/AdamKorcz/go-fuzz-headers
compile_native_go_fuzzer github.com/gocolly/colly/v2 FuzzHtmlElementUnmarshal FuzzHtmlElementUnmarshal

