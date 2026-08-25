这是一个 Web/后端项目：喀斯特泉域观测服务围绕泉点、传感器读数、水文脉冲事件、采样批次、告警和维护调度提供 API。
# Karst Pulse

## 运行

```bash
export GOTOOLCHAIN=local
go build ./...
go run . -addr :8080 -db data/karst-pulse.db
```

启动后访问 `GET /healthz`，静态操作页位于 `/`。

## 主要端点

- `GET/POST /api/v1/springs`
- `GET /api/v1/springs/{id}`
- `GET/POST /api/v1/sensors`
- `GET/POST /api/v1/readings`
- `GET /api/v1/events`
- `GET/POST /api/v1/batches`
- `GET/POST /api/v1/batches/{id}/samples`
- `GET /api/v1/alerts`、`POST /api/v1/alerts/{id}/ack`
- `GET/POST /api/v1/maintenance`、`POST /api/v1/maintenance/{id}/complete`
- `GET /api/v1/analysis?spring_id=1`

读数提交由后台 worker 通过 `channel` 排队，事件检测器根据连续读数推进上升、峰值、衰减和确认阶段。SQLite 数据文件默认为 `data/karst-pulse.db`，进程重启后保留业务数据。
