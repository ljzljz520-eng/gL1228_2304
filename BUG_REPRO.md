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
?   	stagebeam/cmd/stagebeam	[no test files]
ok  	stagebeam/integration	0.021s
?   	stagebeam/internal/audit	[no test files]
ok  	stagebeam/internal/config	0.010s
ok  	stagebeam/internal/geometry	0.010s
ok  	stagebeam/internal/model	0.001s
ok  	stagebeam/internal/persist	0.017s
ok  	stagebeam/internal/render	0.001s
ok  	stagebeam/internal/service	0.012s
--- FAIL: TestGestureKeepsBeamLayers (0.00s)
    gesture_test.go:16: finished fist should retain 5 layers, got 4
FAIL
FAIL	stagebeam/internal/stage	0.002s
?   	stagebeam/internal/transport	[no test files]
FAIL
```

## Architecture reproduction

### linux/amd64
- Go toolchain version: exit `0`
- Go build (.): exit `0`
- Go test (.): exit `1`
- Go run smoke (cmd/stagebeam): exit `0`
### linux/arm64
- Go toolchain version: exit `0`
- Go build (.): exit `0`
- Go test (.): exit `1`
- Go run smoke (cmd/stagebeam): exit `0`
