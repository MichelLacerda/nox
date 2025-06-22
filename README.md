# NOX

## Grammar (syntax)

```ebnf
(* EOF is not a Token *)

program     ::= statement* EOF ;

statement   ::= block
                | exprStmt
                | forStmt
                | ifStmt
                | printStmt
                | returnStmt
                | whileStmt
                | importStmt ;

block       ::= "{" declaration* "}" ;

declaration ::= exportDecl
                | funcDecl
                | varDecl
                | statement
                | classDecl ;

exportDecl  ::= "export" (funcDecl | varDecl | classDecl) ;

funcDecl    ::= "func" function ;

function    ::= IDENTIFIER "(" parameters? ")" block;

parameters  ::= IDENTIFIER ( "," IDENTIFIER )* ;

classDecl   ::= "class" IDENTIFIER ( "<" IDENTIFIER )? 
                "{" function* "}" ;

varDecl     ::= "let" IDENTIFIER ( "=" expression )? ";" ;

exprStmt    ::= expression ";" ;

forStmt      ::= "for" ( forSignature | block ) ;

forSignature ::= "(" forClause ")" block
               | identifier "," identifier "in" expression block ;

forClause    ::= (varDecl | exprStmt | ";") expression? ";" expression? ;

ifStmt      ::= "if" "(" expression ")" statement
                ( "else" statement )? ;

printStmt   ::= "print" expression ";" ;

returnStmt  ::= "return" expression ";" ;

whileStmt   ::= "while" "(" expression ")" statement ;

useStmt     ::= "import" STRING ( "as" IDENTIFIER )? ";" ;

expression  ::= assignment ;

assignment  ::= ( call "." )? IDENTIFIER "=" assignment 
                | logic_or ;

logic_or    ::= logic_and ( "or" logic_and )* ;

logic_and   ::= equality ( "and" equality )* ;

equality    ::= comparison ( ( "!=" ) comparison )* ;

comparison  ::= term ( ( ">" | ">=" | "<" | "<=" ) term )* ;

term        ::= factor ( ( "-" | "+" ) factor )* ;

factor      ::= unary ( ( "/" | "*" ) unary )* ;

unary       ::= ( "!" | "-" ) unary | call | primary ;

call        ::= primary ( "(" arguments? ")" | "." IDENTIFIER | "[" expression "]" )* ;

arguments   ::= expression ( "," expression )* ;

primary     ::= NUMBER
                | STRING
                | "true"
                | "false"
                | "nil"
                | "(" expression ")"
                | "self"
                | "super" "." IDENTIFIER 
                | list
                | dict ;

list            ::= "[" ( expression ( "," expression )* ","? )? "]" ;

dict            ::= "{" (dictEntry ( "," dictEntry )* ","? )? "}" ;

dictEntry       ::= STRING ":" expression ;

elements    ::= expression ( "," expression )* ;
```

## Examples:

### ✅ Variable Declarations

```nox
let name = "Nox";
let age = 42;
let active = true;
let empty = nil;
let list = [8, 10, "1.0"];
```

### ✅ Mathematical Operations

```nox
let sum = 1 + 2;
let product = 3 * 4;
let division = 10 / 2;
let difference = 7 - 5;
let result = (1 + 2) * 3;
```

### ✅ Print Statements

```nox
print "Hello, world!";
print name + " is " + age + " years old.";
print list;
print list[0];
```

### ✅ Functions

```nox
func greet(name) {
    print "Hello, " + name + "!";
}

greet("Michel");
```

### ✅ Return Values

```nox
func add(a, b) {
    return a + b;
}

let result = add(10, 20);
print result;
```

### ✅ Conditionals (`if` / `else` / `else if`)

```nox
let age = 20;

if (age >= 21) {
    print "Adult in the another country";
} else if (age >= 18){
    print "Adult in Brazil";
} else {
    print "Minor";
}
```

### ✅ While Loop

```nox
let i = 0;

while (i < 5) {
    print i;
    i = i + 1;
}

while (true) {
    print "This will run forever"; 
}
```

### ✅ For Loop

```nox
let list = [10, 20, 30];
for i, v in list {
    print i, v; // 0 10, 1 20, 2 30
}

let dict = {"a": 1, "b": 2};
for k, v in dict {
    print k, v; // "a" 1, "b" 2
}

for k, _ in dict {
    print  k; // "a", "b"
}

for _, v in dict {
    print  v; // 1, 2
}

for {
    print "infinite loop";
}

for v in range(10, 5, -1) {
    print v; // 10, 9, 8, 7, 6
}

for v in range(0, 5, 1) {
    print v; // 0, 1, 2, 3, 4
}

for v in range(0, 4) {
    print v; // 0, 1, 2, 3
}

for a, _ in 1 {
    print a; // Runtime Error at _: Object is not iterable.
}
```

### ✅ Function

```nox
func prefix(name) {
    return "Mr " + name;
}

func say_my_name(name) {
    func not_my_name() {
        return prefix("White") + ".";
    }

    print not_my_name(); // Output: Mr White.
    print "Sorry, " + prefix(name) + "!";
}

say_my_name("Heisenberg!"); // Output: Hi, Mr White!
```

### ✅ Classes

```nox
class Animal {
    speak() {
        print "The animal makes a sound.";
    }
}

Animal().speak() // Output: The animal makes a sound.
```

### ✅ Inheritance and `super`

```nox
class Dog < Animal {
    speak() {
        print "The dog barks.";
    }
}

class Cat < Animal {
    speak() {
        print "The cat meows.";
    }
}
```

### ✅ Instantiation and Method Access

```nox
let dog = Dog();
let cat = Cat();

dog.speak(); // Output: The dog barks.
cat.speak(); // Output: The cat meows.
```

### ✅ Constructor with `init`

```nox
class Person {
    init(name) {
        self.name = name;
    }

    greet() {
        print "Hello, I am " + self.name;
    }
}

let p = Person("Alan Turing");
p.greet();
```

### ✅ Scoping

```nox
let x = "global";

{
    let x = "local";
    print x; // "local"
}

print x; // "global"
```

### ✅ File

| Modo               | Significado                                       |
| ------------------ | ------------------------------------------------- |
| `"r"`              | Leitura (erro se o arquivo não existir)           |
| `"w"`              | Escrita (cria ou sobrescreve)                     |
| `"a"`              | Escrita no final do arquivo                       |
| `"r+"`             | Leitura e escrita (erro se o arquivo não existir) |
| `"w+"`             | Leitura e escrita (trunca se já existir)          |
| `"a+"`             | Leitura e escrita no final do arquivo             |
| `"rb"`             | Leitura binária                                   |
| `"wb"`             | Escrita binária                                   |
| `"ab"`             | Escrita binária no final                          |
| `"rb+"` ou `"r+b"` | Leitura/escrita binária                           |
| `"wb+"` ou `"w+b"` | Leitura/escrita binária, truncando                |
| `"ab+"` ou `"a+b"` | Leitura/escrita binária no final                  |

Built-in methods: open, read, read_bytes, readline, write, write_bytes, flush, tell, seek, exists, close

```nox
let f = open("example.txt", "w");
f.write("Hello from Nox!");
f.close();

let f2 = open("example.txt", "r");
let content = f2.read();
print content;
f2.close();

with open("README.md", "r") as f {
    print f;
    print f.readline();
}
```
