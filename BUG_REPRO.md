# Bug Reproduction

## Bug

Route-origin validation failures lose their machine-readable error identity across the domain, application, persistence, and HTTP boundaries. Wrapped cancellation and validation causes are consequently reported as generic internal failures.

## Trigger

Run the five focused checks in `checks/validation045` individually. They exercise a domain failure, evaluator cancellation, wrapped failure encoding, persistence round-trip, and HTTP status classification.

## Observed errors

All five checks exit with status 1:

```text
domain failure lost sentinel identity: validate origin 203.0.113.0/24 AS64500: no covering ROA
evaluator lost cancellation identity: validation request canceled: context canceled
encode wrapped failure: encode validation failure: validation boundary: validate origin 198.51.100.0/24 AS64501: route origin is not authorized
decoded failure lost sentinel identity: validate origin 192.0.2.0/24 AS64511: RPKI snapshot unavailable
HTTPStatus()=500, want 422
```
