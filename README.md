# Structogen
A generator for structograms (Nassi-Shneiderman-Diagrams), written in Go.

## Prerequisites
The [Go programming language >= v1.16](https://go.dev/dl/)

## Usage
Currently, the parser generates a tree-structure from the parsed structogram. This is meant to
eventually be rendered into a nice looking pdf, but that is TBD.

As such, the only functionality provided at the moment is the generation of said tree.

To run the example program, which will read the included `template.str` and then display the tree as
json, run

```go
go run .
```

## Syntax
Structogen can parse .str files. The entire syntax is documented in `template.str`

```
name("template name")

instruction("counter = 0")

for ("counter != 10") {
    instruction("print counter")

    if ("counter % 2 == 0") {
        call("printEven()")
    } else {
        call("printOdd()")
    }

    instruction("counter++")

    dowhile("counter < 5") {
        switch("counter") {
            case("1") {
                instruction("printOne")
            }
            case("two") {
                call("printTwo")
            }
            case("3") {
                instruction("")
            }
            default {
                instruction("printDefault")
            }
        }

        instruction("counter++")
    }
}
```
