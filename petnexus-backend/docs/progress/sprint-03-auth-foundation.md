# Sprint 3: Auth Foundation

วันที่อัปเดต: 1 กรกฎาคม 2026

## เป้าหมาย

สร้างระบบ authentication พื้นฐานสำหรับ owner และ clinic staff โดยยังไม่เริ่ม owner profile หรือ feature ของสัตว์เลี้ยง

## สิ่งที่ทำแล้ว

- เพิ่ม migration `user_role` enum และ `users` table
- เพิ่ม User model ที่เก็บเฉพาะ `password_hash`
- เพิ่ม register สำหรับ role `owner` และ `clinic_staff`
- ปิด public registration สำหรับ role `admin`
- normalize email และตรวจ duplicate email
- hash/check password ด้วย bcrypt
- generate/parse JWT access token แบบ HS256 พร้อมวันหมดอายุ
- เพิ่ม JWT auth middleware และ role middleware
- เพิ่ม typed application errors และ response code ที่อ่านง่าย
- เพิ่ม protected endpoint `GET /api/me`
- wire dependencies ผ่าน repository, service, handler และ routes โดยไม่ใช้ global DB

## Files created

- `migrations/001_create_enums.sql`
- `migrations/002_create_users.sql`
- `internal/utils/app_error.go`
- `internal/utils/auth_test.go`

## Main files changed

- `cmd/api/main.go`
- `internal/config/config.go`
- `internal/database/postgres.go`
- `internal/models/user.go`
- `internal/dto/auth_dto.go`
- `internal/repositories/user_repository.go`
- `internal/services/auth_service.go`
- `internal/handlers/auth_handler.go`
- `internal/middleware/auth_middleware.go`
- `internal/middleware/role_middleware.go`
- `internal/routes/routes.go`
- `internal/utils/password.go`
- `internal/utils/jwt.go`
- `README.md`

## Dependencies added

- `github.com/golang-jwt/jwt/v5`
- `github.com/google/uuid`
- `golang.org/x/crypto/bcrypt` (จาก module `golang.org/x/crypto`)

## Endpoints

```text
POST /api/auth/register
POST /api/auth/login
GET  /api/me
```

`GET /health` และ `GET /health/db` ยัง public และ response format เดิม

## Validation และ security rules

- email และ password ต้องมีค่า
- password ขั้นต่ำ 8 ตัวอักษรและไม่เกิน bcrypt limit 72 bytes
- public role ต้องเป็น `owner` หรือ `clinic_staff`
- duplicate email ตอบ `EMAIL_ALREADY_EXISTS`
- login ผิดตอบ `INVALID_CREDENTIALS` โดยไม่บอกว่า email หรือ password ผิด
- JWT รับเฉพาะ HS256 และต้องมี expiration
- missing, malformed, invalid และ expired token ถูกปฏิเสธ
- password, password hash, JWT token และ JWT secret ไม่ถูกเขียนลง log
- API response ไม่มี `passwordHash`

## การตรวจสอบที่ทำแล้ว

- `gofmt` ผ่าน
- `go mod tidy` ผ่าน
- `go test ./...` ผ่าน รวม bcrypt และ JWT expiration tests
- migrations ทั้งสองไฟล์ apply กับ PostgreSQL สำเร็จ
- `GET /health` ตอบ status `ok`
- `GET /health/db` ตอบ status `connected`
- owner registration สำเร็จ
- clinic staff registration สำเร็จ
- duplicate registration ตอบ HTTP 409
- public admin registration ตอบ HTTP 403
- login สำเร็จและคืน access token
- `GET /api/me` คืน user ที่ถูกต้องโดยไม่มี password hash
- missing และ invalid token ตอบ HTTP 401
- query ใน PostgreSQL ยืนยันว่า password ถูกเก็บเป็น hash ไม่ใช่ plaintext

## สิ่งที่ตั้งใจยังไม่ทำ

- owner profile
- pet และ breed
- clinic profile หรือข้อมูล clinic staff เพิ่มเติม
- QR session และ scanner
- access request และ authorization
- visit และ timeline
- notification และ audit log
- refresh token, logout, email verification และ password recovery

## งานถัดไปที่แนะนำ

Sprint 4: Owner Profile โดยเพิ่มเฉพาะ migration, model, repository, service และ routes ของ owner profile บน auth foundation ที่ผ่านการทดสอบแล้ว
