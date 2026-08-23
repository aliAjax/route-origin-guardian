# Bug Reproduction

## Symptom

Unique capability results share the input Value backing array, and malformed IPv4 prefixes such as `999.1.1.1/24` are accepted.

## Trigger

Call `Unique`, mutate the returned capability Value, and parse the malformed prefix.

## Expected

Returned capability values are independent copies and prefix parsing rejects invalid octets.
