# Sprint 2: Database Foundation

วันที่อัปเดต: 30 มิถุนายน 2026

## เป้าหมาย

ทำให้ backend เชื่อมต่อ PostgreSQL ผ่าน GORM ได้ โดยยังไม่สร้าง tables หรือ feature logic

## สิ่งที่ทำแล้ว

- เพิ่ม `docker-compose.yml` สำหรับ PostgreSQL
- ใช้ database `petnexus`, user `postgres` และ port `5432`
- เพิ่ม named volume และ PostgreSQL healthcheck
- เพิ่ม GORM และ PostgreSQL driver ใน `go.mod`/`go.sum`
- โหลดค่า DB จาก environment config
- สร้าง `ConnectPostgres` สำหรับประกอบ DSN, เปิด connection และ ping PostgreSQL
- บังคับให้ backend หยุดพร้อม error ชัดเจนเมื่อ database connection ล้มเหลว
- เพิ่ม `GET /health/db` ซึ่ง ping database จริง
- อัปเดต README และ migration note

## Endpoint ที่เพิ่ม

```text
GET /health/db
```

Expected result เมื่อ PostgreSQL เชื่อมต่อได้:

```json
{
  "success": true,
  "message": "Database connection is healthy",
  "data": {
    "database": "postgresql",
    "status": "connected"
  }
}
```

## การตรวจสอบที่ทำแล้ว

- `gofmt` ผ่าน
- `go mod tidy` ผ่าน
- `go test ./...` ผ่านทุก package
- ตรวจแล้วว่า `go run ./cmd/api` หยุดพร้อมข้อความ `database connection failed` เมื่อ PostgreSQL ไม่ทำงาน

ข้อจำกัดตอนตรวจ: เครื่องที่ใช้พัฒนาไม่มี Docker/PostgreSQL runtime จึงยังไม่ได้รัน `docker compose up -d` และ smoke test endpoints กับ database จริง

## สิ่งที่ตั้งใจยังไม่ทำ

- SQL migrations และ tables
- GORM `AutoMigrate`
- Models ที่มี GORM fields/tags จริง
- Register/login/JWT/bcrypt
- Pet CRUD
- QR, authorization, visit, timeline และ notification logic

## งานถัดไปที่แนะนำ

1. ติดตั้งหรือเปิด Docker Desktop
2. รัน `docker compose up -d`
3. รัน backend และทดสอบ `/health` กับ `/health/db`
4. เริ่ม Sprint schema/migrations ตาม `docs/database-plan.md`
