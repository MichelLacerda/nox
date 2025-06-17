
# NOX

## Grammar (syntax)

```ebnf
program     ::= statement* EOF ; (* EOF is not a Token *)
statement   ::= block
                | exprStmt
                | forStmt
                | ifStmt
                | printStmt
                | returnStmt
                | whileStmt ;
block       ::= "{" declaration* "}" ;
declaration ::= funcDecl
                | varDecl
                | statement
                | classDecl ;
funcDecl    ::= "func" function ;
function    ::= IDENTIFIER "(" parameters? ")" block;
parameters  ::= IDENTIFIER ( "," IDENTIFIER )* ;
classDecl   ::= "class" IDENTIFIER ( "<" IDENTIFIER )? 
                "{" function* "}" ;
varDecl     ::= "let" IDENTIFIER ( "=" expression )? ";" ;
exprStmt    ::= expression ";" ;
forStmt     ::= "for" "(" ( varDecl | exprStmt | ";" )
                expression? ";"
                expression? ")" statement ;
ifStmt      ::= "if" "(" expression ")" statement
                ( "else" statement )? ;
printStmt   ::= "print" expression ";" ;
returnStmt  ::= "return" expression ";" ;
whileStmt   ::= "while" "(" expression ")" statement ;
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
                | "[" elements? "]" ;
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
for (let i = 0; i < len(list); i = i + 1) {
    print list[i];
}

for (;;) {
    print "This will run forever";
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


## Utils

Windows: List all `.go` files separated by `==== FILENAME ====`

```powershell
Get-ChildItem -Recurse -Filter *.go | ForEach-Object {
    Write-Host "`n==== $($_.FullName) ====" -ForegroundColor Cyan
    Get-Content $_
}
```