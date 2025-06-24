# Nox Built-in Functions and Constants

This document describes the built-in functions, modules, and constants available in the Nox programming language.

---

## Global Built-ins

These functions are registered directly into the global scope.

### `clock()`
Returns the current time in seconds since the Unix epoch.

### `len(value)`
Returns the length of:
- a string (in runes),
- a list,
- a dictionary,
- a `FileObject` (in bytes).

### `range(...)`
Generates a list of numbers. Accepts:
- `range(end)`
- `range(start, end)`
- `range(start, end, step)`

### `assert(condition, message)`
If condition is false, throws an error with the message. In debug mode, only logs.

### `open(path, mode)`
Opens a file using the specified mode, returning a `FileObject`. Mode is similar to Python ("r", "w", "a", etc).

### `fmt(...)`
String formatter. If first argument is a string, replaces `{}` placeholders with subsequent arguments.

---

## `math` Module

Accessed as `math.function(...)`:

- `abs(value)`
- `sqrt(value)`
- `sin(value)`
- `cos(value)`
- `tan(value)`
- `floor(value)`
- `ceil(value)`
- `round(value)`
- `exp(value)`
- `log(value)` (natural log)
- `pow(base, exponent)`
- `max(a, b)`
- `min(a, b)`

---

## `type` Module

Accessed as `type.function(...)`:

### Type Queries
- `of(value)` → returns the string name of the type.
- `is(value, "typeName")` → true if value matches the type.

### Specific Checkers
- `is_nil(value)`
- `is_bool(value)`
- `is_number(value)`
- `is_string(value)`
- `is_list(value)`
- `is_dict(value)`
- `is_function(value)`
- `is_class(value)`
- `is_instance(value)`

### Advanced Checkers
- `is_iterable(value)` → list, dict or string
- `is_callable(value)` → function or class
- `is_truthy(value)` → evaluates to `true` in logical context
- `is_falsey(value)` → evaluates to `false` in logical context
- `instance_of(instance, class)`

---

## Constants

Globally available:

| Name        | Value      | Description                          |
|-------------|------------|--------------------------------------|
| `PI`        | 3.141592   | π                                    |
| `E`         | 2.718281   | Euler's number                       |
| `PHI`       | 1.618033   | Golden ratio                         |
| `TAU`       | 6.283185   | 2π                                   |
| `sqrt2`     | 1.414213   | √2                                   |
| `sqrtE`     | 1.648721   | √e                                   |
| `sqrtPi`    | 1.772453   | √π                                   |
| `sqrtPhi`   | 1.272019   | √φ                                   |
| `ln2`       | 0.693147   | ln(2)                                |
| `log2E`     | ~1.442695  | log₂(e)                              |
| `ln10`      | 2.302585   | ln(10)                               |
| `log10E`    | ~0.434294  | log₁₀(e)                             |

---

These built-ins are registered automatically when the interpreter starts.
