# Snapshot

[![godev](https://img.shields.io/static/v1?label=godev&message=reference&color=00add8)](https://pkg.go.dev/github.com/NateScarlet/snapshot/pkg/snapshot)
[![build status](https://github.com/NateScarlet/snapshot/workflows/Go/badge.svg)](https://github.com/NateScarlet/snapshot/actions)

Snapshot test for golang.

Store snapshot files under `__snapshots__` folder relative to caller file.

- [x] json snapshot for any go object with type info
- [x] transform value as schema to compare type only
- [x] clean data by regexp before compare

## Update snapshot

set `SNAPSHOT_UPDATE` env var to `true` to update existed snapshot file.


## Usage

```go

import (
    "time"

    "github.com/NateScarlet/snapshot/pkg/snapshot"
)

type Object struct {
    A string
    B int
    C bool
    D time.Time
}

func TestSomeThing(t *testing.T) {
    // Match as text
    snapshot.Match(t, "text")
    // snapshot:
    // text

    // Match as json
    snapshot.MatchJSON(t, Object{})
    // snapshot:
    // {
    //   "$Object": {
    //     "A": "",
    //     "B": 0,
    //     "C": false,
    //     "D": {
    //       "$Time": "0001-01-01 00:00:00 +0000 UTC"
    //     }
    //   }
    // }

    // Match schema as json
    snapshot.MatchJSON(t, Object{}, snapshot.OptionTransform(snapshot.TransformSchema))
    // snapshot:
    // {
    //   "$Object": {
    //     "A": "$string",
    //     "B": "$int",
    //     "C": "$bool",
    //     "D": "$Time"
    //   }
    // }

    // Clean dynamic data to make result deterministic
    snapshot.MatchJSON(t, Object{
        A: "a",
        B: 1,
        C: true,
        D: time.Now(),
    }, OptionCleanRegex(CleanAs("*time*"), `"D": {\s+"\$Time": "(.+)"\s+}`))
    // snapshot:
    // {
    //   "$Object": {
    //     "A": "a",
    //     "B": 1,
    //     "C": true,
    //     "D": {
    //       "$Time": "*time*"
    //     }
    //   }
    // }
}

```
