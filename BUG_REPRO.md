# Logging level normalization bug

## What happened

Deployment log-level configuration did not consistently map to the expected `slog` severity. Values with surrounding whitespace fell back to `INFO`, `WARN` and `ERROR` were swapped, the default was `DEBUG`, filtering used the wrong boundary direction, and level-name helpers returned an incorrect name or panicked for known levels.

## How to trigger it

Use the logging package with values such as `" WARN "`, `"ERROR"`, and an unknown value, then check the default level, severity filtering boundary, and level-name helpers. The buggy environment reports failures for normalization, filtering, default level, and level naming; `MustLevelName` panics even for a known level.

## Observed error

```text
ParseLevel(" WARN ")=INFO, want WARN
level filtering boundary is incorrect
DefaultLevel()=DEBUG, want INFO
LevelName(ERROR)="INFO"
panic: level name unavailable
```
