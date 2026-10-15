---
name: Bug report
about: Something does not work as expected
title: ''
labels: bug
assignees: ''
---

## What happened

<!-- What did you observe? -->

## What you expected

<!-- What should have happened instead? -->

## Steps to reproduce

1.
2.
3.

## Layer

<!--
Which layer is involved? This helps narrow the search:
- domain  — business rules, types, parsing
- storage — JSON persistence
- service — orchestration, checkout
- cmd     — wiring and entry point
-->

## Environment

- Go version (`go version`):
- OS:
- Commit or branch (`git log --oneline -1`):

## Checks

<!-- Paste the output, or note that they pass. -->

```
gofmt -l .
go vet ./...
go test ./...
```