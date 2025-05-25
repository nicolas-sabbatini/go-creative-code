# Go creative coding

## Compile

To run the projects you must have everything to run a
[Ebiten](https://ebitengine.org/en/documents/install.html) project.

Then you can run the each program running:

```bash
go run cmd/{{project}}
```

Or you can compile all running:

```bash
go run tools/builder.go -target=all -os={{your os}} -arch={{your arch}}
```

For more information on how to use the builder run:

```bash
go run tools/builder.go -h
```
