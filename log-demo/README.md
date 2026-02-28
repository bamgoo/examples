# log-demo

Run:

```bash
cd examples/log-demo
cp config.toml ../config.toml
cd ..
go run ./log-demo
```

Greptime run:

```bash
cd examples/log-demo
cp config.greptime.toml ../config.toml
cd ..
go run ./log-demo
```

This demo sends a burst of logs from many goroutines to verify async batch write.
It also demonstrates:

- `log.Infof(...)` for format logging
- `log.Noticew(..., Map{...})` for structured logging

Queue overflow strategy:

- `overflow = "drop_oldest"`: queue full时丢最旧日志，优先保留最新日志。
- Optional: `overflow = "drop_newest" | "block"`.

Useful settings:

- `buffer`: queue length
- `batch`: flush batch size
- `timeout`: max flush interval
- `sample`: sampling ratio (`0~1`)

File output:

- `examples/store/log/demo.log`
- rotated files: `examples/store/log/demo.*.log.gz` (when `compress = true`)

Greptime output (when using `config.greptime.toml`):

- endpoint: `127.0.0.1:4001`
- database: `public`
- table: `logs`
