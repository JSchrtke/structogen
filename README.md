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
