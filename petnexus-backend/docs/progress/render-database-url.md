# Render: DATABASE_URL Support

วันที่อัปเดต: 1 กรกฎาคม 2026

## เป้าหมาย

ให้ backend ใช้ PostgreSQL connection string ของ Render ได้ โดยไม่ทำให้ Docker PostgreSQL สำหรับ local development พัง

## สิ่งที่เปลี่ยน

- เพิ่ม `DatabaseURL` ใน application config
- โหลดค่าจาก environment variable `DATABASE_URL` โดย default เป็นค่าว่าง
- ถ้า `DATABASE_URL` มีค่า จะใช้ค่านั้นเป็น GORM PostgreSQL DSN โดยตรง
- ถ้าไม่มีค่า จะ fallback ไปใช้ `DB_HOST`, `DB_PORT`, `DB_USER`, `DB_PASSWORD`, `DB_NAME` และ `DB_SSLMODE` เหมือนเดิม
- เพิ่ม `DATABASE_URL=` ใน `.env.example`
- เพิ่ม unit tests สำหรับ Render URL, local fallback และ whitespace URL

## ผลกระทบ

- Render สามารถตั้ง `DATABASE_URL` จาก managed PostgreSQL ได้
- Local Docker ยังคงใช้ `localhost:5432` และ `sslmode=disable`
- `GET /health/db` ยังตรวจ connection ที่ active อยู่เหมือนเดิม

## การตรวจสอบ

- `gofmt` ผ่าน
- `go test ./...` ผ่าน
- unit tests ยืนยันว่า `DATABASE_URL` มี priority เหนือ DB_* values
- unit tests ยืนยันว่า local DSN และ `sslmode=disable` ยังเหมือนเดิมเมื่อ `DATABASE_URL` ว่าง
