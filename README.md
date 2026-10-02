# claude-pipes

Spawn Claude Code instances programmatically and reuse them over and over.

claude-pipes starts a single persistent headless `claude -p` process and talks
to it over stdin/stdout using `--input-format stream-json` /
`--output-format stream-json`. Each message you send goes to the same instance,
so it keeps its context between turns instead of starting cold every time.

The spawned instance has full access to tools: Read, Write, Edit, Glob, Grep,
WebSearch, WebFetch, Bash and PowerShell. Tool calls are auto-approved, so it
can read and modify files and run shell commands on its own without any
interactive permission prompt.

This is a minimal Go harness meant as a starting point for embedding Claude
Code in your own programs, bots or pipelines. Type a message, the program
forwards it to the running instance and prints the final result. Permission
denials, if any, are printed after the reply.

## Build

```
go build -o claude-pipes.exe .
```

Requires `claude` on PATH.
