# xotp

![xotpbanner](./misc/banner.png)

self-hostable lightweight otp service with generation, verification, expiry and
configuration endpoints.

## Usage

```powershell
go run ./cmd/xotp
```

The API listens on `http://localhost:8080`. Use `PORT` to configure another
port. Start Redis with `docker compose up -d redis`; use `REDIS_ADDR` to
configure its address (default: `localhost:6379`).

## Endpoints

| Method | Path | Request body | Response body |
| --- | --- | --- | --- |
| `GET` | `/` | — | `{"alive":true}` |
| `GET` | `/generateOtp` | `{"userId":"user-123"}` | `{"otp":123456}` |
| `POST` | `/verifyOtp` | `{"userId":"user-123","otp":123456}` | `{"match":true}` |
| `DELETE` | `/expiry` | `{"userId":"user-123"}` | `204 No Content` |
