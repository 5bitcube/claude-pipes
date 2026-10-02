# claude-pipes

Minimal Go harness that drives Claude Code headlessly over stdin/stdout using
`--input-format stream-json` / `--output-format stream-json`.

Type a message, the program forwards it to a persistent `claude -p` process and
prints the final result. Permission denials, if any, are printed after the reply.

## Build

```
go build -o claude-pipes.exe .
```

Requires `claude` on PATH.
