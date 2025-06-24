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

## `random`

### `random.int(min, max)`
Returns a random integer between `min` and `max` (inclusive).

```nox
let n = random.int(1, 10)
print n  # e.g. 7
```

### `random.float()`
Returns a random floating-point number between `0.0` and `1.0`.

```nox
let f = random.float()
print f  # e.g. 0.581
```

---

## `os`

### `os.listdir(path, only_dirs=false)`
Returns a list of files and directories in `path`. If `only_dirs` is `true`, returns only directories.

```nox
print os.listdir(".")
```

### `os.chmod(path, mod)`
Changes the file permissions to mode `mod`.

```nox
os.chmod("script.nox", 0o755)
```

### `os.chdir(path)`
Changes the current working directory to `path`.

```nox
os.chdir("/home/user")
```

### `os.cwd()`
Returns the current working directory.

```nox
print os.cwd()
```

### `os.getenv(key)`
Returns the value of an environment variable.

```nox
print os.getenv("HOME")
```

### `os.setenv(key, value)`
Sets an environment variable.

```nox
os.setenv("MODE", "dev")
```

### `os.mkdir(path)`
Creates a new directory.

```nox
os.mkdir("new_folder")
```

### `os.rmdir(path)`
Removes an empty directory.

```nox
os.rmdir("new_folder")
```

### `os.info(path)`
Returns information about the file: size, type, permissions, etc.

```nox
print os.info("file.txt")
```

### `os.walk(path)`
Recursively walks through directories from `path`. Returns a list of tuples `(dir, subdirs, files)`.

```nox
for (dir, subdirs, files) in os.walk(".") {
    print dir
}
```

### `os.exec(command)`
Executes a system command and returns its output.

```nox
print os.exec("echo Hello")
```

### `os.exit(code=0)`
Terminates the script with the given exit code.

```nox
os.exit(1)
```

---

## `path`

### `path.exists(path)`
Checks if the path exists.

```nox
print path.exists("file.txt")
```

### `path.abs(path)`
Returns the absolute path.

```nox
print path.abs("./file.txt")
```

### `path.join(path, ...)`
Joins multiple path segments.

```nox
print path.join("a", "b", "c.txt")
```

### `path.split(path)`
Splits the path into `(dir, base)`.

```nox
print path.split("a/b/c.txt")  # ("a/b", "c.txt")
```

### `path.splitext(path)`
Splits the path into `(base, ext)`.

```nox
print path.splitext("a/b/c.txt")  # ("a/b/c", ".txt")
```

### `path.basename(path)`
Returns the file name.

```nox
print path.basename("a/b/c.txt")  # "c.txt"
```

### `path.dirname(path)`
Returns the parent directory.

```nox
print path.dirname("a/b/c.txt")  # "a/b"
```

### `path.extname(path)`
Returns the file extension.

```nox
print path.extname("a/b/c.txt")  # ".txt"
```

### `path.relpath(path, start)`
Returns the relative path from the current directory.

```nox
print path.relpath("/usr/local/bin", "usr")
```

### `path.normalize(path)`
Normalizes slashes, `..`, etc.

```nox
print path.normalize("a//b/../c")
```

### `path.size(path)`
Returns the file size (in bytes).

```nox
print path.size("file.txt")
```

### `path.time(path)`
Returns creation, modification, and access timestamps.

```nox
print path.time("file.txt")
```

### `path.isdir(path)`
Returns `true` if the path is a directory.

```nox
print path.isdir("my_folder")
```

### `path.isfile(path)`
Returns `true` if the path is a file.

```nox
print path.isfile("file.txt")
```

### `path.islink(path)`
Returns `true` if the path is a symbolic link.

```nox
print path.islink("shortcut")
```
---

These built-ins are registered automatically when the interpreter starts.
