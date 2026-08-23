# Bug Reproduction

## Symptom

Cancelled RPKI refreshes continued reading and could advance the index with a partial snapshot.

## Trigger

Abort the refresh context while the source fetch/import chain is active, then inspect the loader's last error and index version.

## Expected

Cancellation and deadlines propagate through source, loader, and import; the index remains unchanged.
