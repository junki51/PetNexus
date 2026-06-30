# Sprint 1: Backend Foundation

วันที่อัปเดต: 30 มิถุนายน 2026

## เป้าหมาย

สร้างโครง Go backend ที่อ่านง่ายและพร้อมต่อยอด โดยยังไม่ทำ database, authentication หรือ business feature จริง

## สิ่งที่ทำแล้ว

- สร้าง Go module และ layered folder structure
- ตั้งค่า Gin HTTP server
- โหลด environment variables ด้วย godotenv
- เพิ่มค่า default ของ `PORT=8080`
- เพิ่ม response helpers สำหรับ success และ error
- เพิ่ม `GET /health`
- เพิ่ม placeholder สำหรับ models, repositories, services, handlers, DTOs, middleware และ utilities
- เพิ่ม `.env.example`, `.gitignore`, README และเอกสาร migrations

## Endpoint

```text
GET /health
```

Expected result:

```json
{
  "success": true,
  "message": "PetNexus backend is running",
  "data": {
    "status": "ok",
    "service": "petnexus-backend"
  }
}
```

## สิ่งที่ยังไม่ทำใน Sprint นี้

- PostgreSQL/GORM connection
- Tables และ migrations
- Register/login/JWT
- Pet, QR, authorization, visit, timeline และ notification logic

## ผลลัพธ์

Sprint 1 เป็น foundation เท่านั้น ไม่มี fake authentication หรือ business logic ชั่วคราวปะปนอยู่ใน scaffold
