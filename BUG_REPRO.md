# Bug Reproduction

## Symptom

Retryable analysis errors lose their identity after policy and suppression handling, so `errors.Is` cannot classify them.

## Trigger

Return the retryable sentinel from analysis and pass the result through the policy and suppressor layers.

## Expected

The returned error preserves the retryable chain through every layer.
