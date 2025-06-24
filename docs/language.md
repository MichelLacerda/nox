# Nox — Language Guide

This guide introduces the core syntax and features of the Nox programming language.

---

## 🧠 Variables

Declare variables using `let`. Reassignment is allowed.

```nox
let x = 10
let name = "Nox"
x = x + 1
```

---

## 🧮 Expressions

Nox supports standard arithmetic and logical operators:

```nox
let total = 5 + 3 * 2
let valid = age > 18 and country == "BR"
```

You can also use:

- `!value` for logical not
- `-number` for negation
- `value ** 2` for exponentiation
- `value % 2` for modulo
- Safe call: `?expr`

---

## 🔁 Control Flow

### If / Else

```nox
if score > 100 {
    print "You win!"
} else {
    print "Try again."
}
```

Nested `if` expressions are allowed.

---

## 🔄 Loops

### Infinite Loop

```nox
for {
    print "Running forever..."
}
```

### For-in over a list

```nox
let names = ["Alan Turing", "Pirulla", "Douglas Adams", "Friedrich Nietzsche"]
for name in names {
    print name
}
```

### For-in over a dictionary

```nox
let user = {"name": "Nicolly Fúlvia", "level": 42}
for key, value in user {
    print key, "=", value
}
```

### Loop control

```nox
for {
    if done {
        break
    }
    if skipping {
        continue
    }
}
```

---

## 📤 Return

Return values from functions using `return`.

```nox
func square(n) {
    return n * n
}
```

Return without expression returns `nil`:

```nox
return
```

---

## 🧱 Blocks and Scope

Blocks `{ ... }` introduce new scopes. Variables declared inside do not escape.

```nox
{
    let x = 100
    print x
}
# x is not accessible here
```

---

## 🧬 Functions

Define functions using `func`:

```nox
func greet(name) {
    print "Hello, " + name
}

greet("Nox")
```

Functions are first-class values:

```nox
let add = func(a, b) {
    return a + b
}

print add(2, 3)  # 5
```

---

## 🧪 Assert

`assert(condition)` throws a runtime error if the condition is `nil` or `false`.

```nox
assert(user != nil)
```

---

## 🤖 Safe Call Operator (`?expr`)

Use `?` to perform a safe access. If any part of the expression fails (e.g. nil access or invalid call), it returns `nil`.

```nox
let name = ?user.profile.name
let first = ?list[0]
let result = ?obj.method()
```

---

## 🧩 Classes and Inheritance

```nox
class Animal {
    func init(name) {
        self.name = name
    }

    func speak() {
        print self.name + " makes a sound."
    }
}

class Dog < Animal {
    func init(name, breed) {
        super.init(name)
        self.breed = breed
    }

    func speak() {
        super.speak()  // chama o método da superclasse
        print self.name + " barks."
    }
}

let dog = Dog("Rex", "Labrador")
dog.speak()
```

### Special keywords:

- `self`: refers to the instance
- `super.foo()`: calls method from superclass

---

## 📦 Modules

### Import

```nox
import "math"
import "foo/bar" as f
```

### Export

```nox
export let gravity = 9.8

export func greet(name) {
    print "Hi " + name
}
```

---

## 📚 Lists

```nox
let nums = [1, 2, 3]
nums.append(4)
print nums[0]  # 1
```

Common methods: `append`, `pop`, `remove`, `insert`, `length`, `contains`, `sort`, `reverse`

---

## 📘 Dictionaries

```nox
let data = {"name": "Ana", "level": 10}
print data["name"]
data.get("age", 21)
```

Methods: `get`, `set`, `remove`, `contains`, `length`, `keys`, `values`

---

## 📞 Method Calls and Indexing

```nox
object.method()
array[0]
dict["key"]
object.method()[1]
```

Chaining and nesting are allowed. Combine with `?` for safe evaluation.

---

## 📁 File I/O

Nox supports basic file operations through the builtin function `open(path, mode)`, which returns a file object. The recommended way to handle files is by using the `with` statement, which ensures proper closing of the file, even in case of runtime errors.

### 🔓 Opening a file

```nox
let file = open("example.txt", "w")
file.write("Hello, Nox!")
file.close()
```

### ✅ Recommended: Using `with` statement

```nox
with open("example.txt", "w") as file {
    file.write("Hello, Nox!\n")
    file.write("Second line.")
}
```

### 📖 Reading content

```nox
with open("example.txt", "r") as file {
    let content = file.read()
    print(content)
}
```

### ➕ Appending to a file

```nox
with open("example.txt", "a") as file {
    file.write("\nAppended line.")
}
```

### 🔧 File Modes

| Mode | Description               |
|------|---------------------------|
| "r"  | Read (default)            |
| "w"  | Write (truncates file)    |
| "a"  | Append                    |
| "r+" | Read/Write                |
| "w+" | Write (truncates) + Read  |
| "a+" | Append + Read             |

### 🧰 File Methods

| Method         | Description                      |
|----------------|----------------------------------|
| `read()`       | Returns full file contents       |
| `readline()`   | Returns a single line            |
| `write(text)`  | Writes a string to the file      |
| `seek(offset)` | Moves file pointer to position   |
| `close()`      | Manually closes the file         |

---

## ✅ Summary

Nox is designed to be concise, expressive, and safe. Use this guide as a reference and explore the examples for deeper understanding.

Check out:

- [Built-in Functions](./builtins.md)
- [Grammar Specification](./grammar.ebnf)