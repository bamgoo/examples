# ws-demo

最小可跑的 `web + ws` 示例，同时演示：

- `ws` 模块的默认接入
- 自定义 `space`
- `ctx.Upgrade()`
- `ctx.Upgrade("name")`

## 功能

- `web.Context.Upgrade()`
- `space` 隔离
- `ws.Hook`
- `ws.Filter`
- `ws.Message`
- `ws.Command`
- `ctx.Reply()`
- `ctx.Broadcast()`
- `ctx.Groupcast()`
- `ctx.BindUser()`
- `ctx.PushUserResult()`
- 本地投递统计 `BroadcastResult / GroupcastResult`
- 协议导出 `/ws/export`
- 运行指标 `/ws/metrics`
- 首页直接渲染协议文档和指标快照
- 协议文档包含 schema/version 信息

## 运行

```bash
cd examples/ws-demo
go run .
```

默认监听 `http://127.0.0.1:8080/`。

页面里可以直接测试：

- 默认 ws 连接
- 自定义 space 连接
- echo
- join / groupcast
- broadcast
- bind user / push user
- 协议导出和 space 分组视图
- 运行指标实时刷新

其中：

- `/socket` 走 `ctx.Upgrade()`，默认 `space = ctx.Name`
- `/socket/custom` 走 `ctx.Upgrade("custom")`，显式使用 `custom` 空间
