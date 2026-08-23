# Bug Reproduction

## Symptom

A zero-value metrics Registry panics on its first counter write with `assignment to entry in nil map`.

## Trigger

Instantiate `Registry{}` and call `Counter` for a new metric name.

## Expected

The registry initializes its counter map safely and returns a usable counter without panic.
