# (Idiotic) SimpleScript
A stupidly minimal and simplistic interpreted language written in go.

## Instalation

Prerequisites: Go 1.21 or higher


do this for installing the interpreter to your system
```sh
make install
```

do this if you want to build it locally
```sh
make build
```

do this to uninstall it from system
```sh
make uninstall
```

## Usage

### REPL Mode

```sh
idiot
```

### Run a Script

```sh
idiot run script.simp
```

## Language Feature

### Datatype

* Integers : 1, 67, -69
* Booleans : true, false
* Strings  : "Hello World!"
* Arrays   : [1,2,3]
* Hashmaps : {6:"seven"}

### Variables

```sh
let x = 5;
let six = "seven";
let eliteBallKnowledge = ["apple", "grapes", "durian"];
```

### Functions

```sh
let add = ft(x, y) {
    return x + y;
}

let brainrot = add(6,7);
```

### Control Flow

```sh
if (x == 6) {
    return 7;
} else {
    return 89;
}

for () {
    print("six sevennnn");
    if (x == 67){
        break;
    }
    x = x + 1;
}
```

### Operators

* Arithmetic : +,-,\*,/
* Comparison : ==,!=,\<,\>,\<=,\>=
* Logical    : !
* Assignment : =

### Builtin Functions

- len(x) - Returns length of string or array
- head(arr) - Returns first element of array
- tail(arr) - Returns last element of array
- killHead(arr) - Returns array without first element
- push(arr, x) - Returns array with element appended
- print(x) - Prints values to stdout

## Example

```sh
let fibonacci = ft(n) {
   if (n < 2) {
       return n;
   }
   return fibonacci(n - 1) + fibonacci(n - 2);
};

let result = fibonacci(10);
print(result);
```

## Development

### Run Test

```sh
make test
```

### Run Directy

```sh
make run
```

## Project Structure

* lexer/      - Tokenization 
* token/      - Token definitions
* ast/        - Abstract syntax tree
* parser/     - Parser implementation
* object/     - Object system and environment
* evaluator/  - Tree-walking evaluator
* repl/       - Read-eval-print loop



