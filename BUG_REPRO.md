# Migration Transaction Reproduction

## Bug

`ApplyMigrations` executed each migration without one transaction, so a later failure left earlier statements committed. Cancellation also did not reliably roll back the active work, and `DB.Exec` converted driver errors with `%v`, breaking `errors.Is`.

## Reproduce

Run the focused migration checks from the collection contract:

```text
go test -run '^TestMigrationTxnRejectsPlanBeforeBegin$' -count=1 ./internal/platform/storage
go test -run '^TestMigrationTxnRollsBackExecFailure$' -count=1 ./internal/platform/storage
go test -run '^TestMigrationTxnCommitsOnce$' -count=1 ./internal/platform/storage
go test -run '^TestMigrationTxnPropagatesCancel$' -count=1 ./internal/platform/storage
go test -run '^TestMigrationTxnPreservesCause$' -count=1 ./internal/platform/storage
```

Before the fix, the failure path leaves partial writes, cancellation leaves the connection unavailable, or the returned error cannot be matched with `errors.Is`.
