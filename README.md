# corev4-explorer

## Getting private git package

```bash
git config --global url."git@github.com:".insteadOf "https://github.com/"
export GOPRIVATE=github.com/blcvn/*
go mod vendor
```