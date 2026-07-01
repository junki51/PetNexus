# PetNexus Backend Progress

โฟลเดอร์นี้เป็น development log สำหรับส่งต่องานให้เพื่อนหรือผู้พัฒนาคนถัดไปอ่าน โดยสรุปเฉพาะสิ่งที่ทำจริง ผลการตรวจสอบ ข้อจำกัด และงานถัดไป

## Current status

| Sprint | สถานะ | ผลลัพธ์หลัก |
| --- | --- | --- |
| [Sprint 1: Backend Foundation](./sprint-01-backend-foundation.md) | เสร็จแล้ว | Go/Gin scaffold และ `GET /health` |
| [Sprint 2: Database Foundation](./sprint-02-database-foundation.md) | โค้ดเสร็จแล้ว | Docker Compose, GORM, PostgreSQL connection และ `GET /health/db` |
| [Sprint 3: Auth Foundation](./sprint-03-auth-foundation.md) | เสร็จและทดสอบแล้ว | Register, login, JWT middleware และ `GET /api/me` |

## วิธีใช้งาน

หลังจบงานแต่ละ Sprint:

1. สร้างไฟล์ใหม่จาก [template](./update-template.md)
2. ระบุวันที่และขอบเขตของงาน
3. บันทึกไฟล์ที่สร้างหรือแก้ไข
4. บันทึกคำสั่งทดสอบและผลจริง
5. ระบุสิ่งที่ตั้งใจยังไม่ทำ
6. อัปเดตตาราง Current status ด้านบน

Development log ต้องไม่เก็บ password, token, secret, ข้อมูลส่วนบุคคล หรือข้อความสนทนาดิบ
