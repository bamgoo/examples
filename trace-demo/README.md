# trace-demo

Run:

```bash
cd examples/trace-demo
cp config.toml ../config.toml
cd ..
go run ./trace-demo
```

This demo writes trace spans/events to three targets:

- console via `trace` default driver
- file via `trace-file` driver
- greptime via `trace-greptime` driver

Trace file output:

- `examples/store/trace/trace.log`

Greptime defaults in this demo:

- endpoint: `127.0.0.1:4001`
- database: `public`
- table: `traces`
