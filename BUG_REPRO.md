# BUG_REPRO

The following failures were observed while validating the initial project state.
Each section records what failed, how to reproduce it, and the complete command output.
They are preserved intentionally; only failing build gates are omitted from the generated Dockerfile.

## Failure 1: Go test (.)

- Observed problem: `Go test (.)` failed in the initial project state.
- Working directory: `.`
- Command: `cd /app && GOTOOLCHAIN=local GOPROXY=off GOSUMDB=off go test -count=1 ./...`
- Exit status: `1`

```text
ok  	drivingmaterials/catalog	0.009s
?   	drivingmaterials/cmd/drivehub	[no test files]
ok  	drivingmaterials/domain	0.001s
?   	drivingmaterials/presenter	[no test files]
--- FAIL: TestCourseSearchResultsStayIndependent (0.01s)
    regression_test.go:38: first result labels changed: []string{"light", "three", "south", "2026-06-01"}
FAIL
FAIL	drivingmaterials/integration	0.018s
ok  	drivingmaterials/persistence	0.011s
ok  	drivingmaterials/query	0.005s
ok  	drivingmaterials/transport	0.005s
ok  	drivingmaterials/workflow	0.010s
FAIL
```

## Architecture reproduction

### linux/amd64
- Go toolchain version: exit `0`
- Go build (.): exit `0`
- Go test (.): exit `1`
- Go run smoke (cmd/drivehub): exit `0`
### linux/arm64
- Go toolchain version: exit `0`
- Go build (.): exit `0`
- Go test (.): exit `1`
- Go run smoke (cmd/drivehub): exit `0`
