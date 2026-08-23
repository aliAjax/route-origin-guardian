# Bug Reproduction

## Symptom

Stored route AS paths alias caller memory, and cancelled announce/replace requests still reach storage.

## Trigger

Upsert a route, mutate its AS path, then issue announce and replace with an already-cancelled context.

## Expected

Stored snapshots remain immutable and cancelled service calls return before writing.
