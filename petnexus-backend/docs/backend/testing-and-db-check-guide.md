PetNexus Backend Testing & Database Check Guide

คู่มือนี้ใช้สำหรับเช็กว่า backend, database และ API หลักของ PetNexus ทำงานถูกต้องหรือไม่
เหมาะสำหรับเวลาทีมต้องทดสอบระบบ, debug ปัญหา, หรือเช็กว่า server ยังรันอยู่ไหม

1. สิ่งที่ต้องรู้ก่อน

PetNexus backend ใช้ stack หลัก:

Go + Gin
PostgreSQL
GORM
JWT
bcrypt
Docker Compose
Render

Architecture หลัก:

Client
→ Route
→ Middleware
→ Handler
→ Service
→ Repository
→ Database
→ Response

ความหมายแบบสั้น:

Route = บอกว่า endpoint นี้ไป handler ไหน
Middleware = ตรวจ token / role
Handler = รับ request
Service = ตรวจ logic และ permission
Repository = คุย database
Database = เก็บข้อมูลจริง
2. เช็กว่า Backend Local รันอยู่ไหม
2.1 เปิด PostgreSQL local ก่อน

ที่ root project ให้รัน:

docker compose up -d

เช็กว่า container รันอยู่:

docker compose ps

ถ้า PostgreSQL รันถูกต้อง ควรเห็น container สถานะประมาณ:

Up
2.2 รัน Backend

ที่ root project ให้รัน:

go run ./cmd/api

ถ้าสำเร็จควรเห็น log ประมาณ:

database connected successfully
database migration completed successfully
PetNexus backend listening on http://localhost:8080
2.3 เช็ก Health API

เปิด PowerShell อีกหน้าหนึ่ง แล้วรัน:

Invoke-RestMethod http://localhost:8080/health

ถ้าสำเร็จควรได้ response ประมาณ:

{
  "success": true,
  "message": "PetNexus backend is running"
}

ถ้า /health ผ่าน แปลว่า:

Backend server รันอยู่
Route /health ใช้งานได้
Handler ทำงานได้
3. เช็กว่า Backend ต่อ Database ได้ไหม

รัน:

Invoke-RestMethod http://localhost:8080/health/db

ถ้าสำเร็จ แปลว่า:

Backend ต่อ PostgreSQL ได้
Database connection ใช้งานได้

ถ้า fail ให้เช็กตามลำดับ:

1. Docker PostgreSQL รันอยู่ไหม
2. .env ตั้งค่า DB ถูกไหม
3. DB_HOST / DB_PORT / DB_USER / DB_PASSWORD / DB_NAME ถูกไหม
4. Backend log มี error อะไร

คำสั่งเช็ก container:

docker compose ps

คำสั่งดู log ของ container:

docker compose logs
4. เช็ก Database ด้วย psql ผ่าน Docker

ถ้าต้องเข้าไปดู database จริง ให้ใช้คำสั่งนี้:

docker exec -it petnexus-postgres psql -U postgres -d petnexus

ถ้าชื่อ container ไม่ตรง ให้เช็กก่อนด้วย:

docker ps

หลังเข้า psql แล้ว ใช้คำสั่งเหล่านี้:

ดูตารางทั้งหมด
\dt

ควรเห็นตารางหลัก เช่น:

users
owner_profiles
clinic_profiles
breeds
pets
appointments
medical_records
ดูโครงสร้าง table
\d users
\d pets
\d appointments
\d medical_records
ดูข้อมูลใน table
SELECT id, email, role, created_at FROM users;
SELECT id, user_id, full_name, phone FROM owner_profiles;
SELECT id, owner_profile_id, name, public_pet_id FROM pets;
SELECT id, owner_profile_id, clinic_profile_id, pet_id, status FROM appointments;
SELECT id, clinic_profile_id, pet_id, appointment_id, created_by_user_id FROM medical_records;

ออกจาก psql:

\q
5. Smoke Test Flow หลัก

Smoke Test คือการทดสอบแบบเร็ว ๆ ว่าระบบสำคัญยังทำงานอยู่ไหม

ลำดับที่ควรเทส:

1. /health
2. /health/db
3. Register Owner
4. Login Owner
5. GET /api/me
6. Create Owner Profile
7. Create Pet
8. Register/Login Clinic
9. Create Clinic Profile
10. Create Appointment
11. Check Clinic Patients
12. Create Medical Record
6. Auth Test: Register / Login / Token
6.1 Register Owner
$body = @{
  email = "owner_test@example.com"
  password = "password123"
  role = "owner"
} | ConvertTo-Json

Invoke-RestMethod `
  -Method Post `
  -Uri http://localhost:8080/api/auth/register `
  -ContentType "application/json" `
  -Body $body

หลัง register สำเร็จ ควรมี user ใหม่ใน database

เช็กใน psql:

SELECT id, email, role FROM users WHERE email = 'owner_test@example.com';

ข้อสำคัญ:

Database ต้องเก็บ password_hash
ไม่ควรเก็บ password จริง
6.2 Login Owner
$body = @{
  email = "owner_test@example.com"
  password = "password123"
} | ConvertTo-Json

$response = Invoke-RestMethod `
  -Method Post `
  -Uri http://localhost:8080/api/auth/login `
  -ContentType "application/json" `
  -Body $body

$response

ถ้าสำเร็จ response ควรมี token

เก็บ token ไว้ใช้ต่อ:

$token = $response.data.access_token

ถ้า field ชื่อไม่ตรง ให้ดูจาก response จริง แล้วปรับตามนั้น

6.3 GET /api/me
Invoke-RestMethod `
  -Method Get `
  -Uri http://localhost:8080/api/me `
  -Headers @{ Authorization = "Bearer $token" }

ถ้าสำเร็จ แปลว่า:

Token ใช้ได้
Auth Middleware ตรวจ token ผ่าน
Backend รู้ว่า user นี้คือใคร

ถ้าได้ 401:

Token ไม่มี / token ผิด / token หมดอายุ / header ไม่ได้ใช้ Bearer
7. Owner Profile Test
Create Owner Profile
$body = @{
  full_name = "Owner Test"
  phone = "0800000000"
  address = "Bangkok"
} | ConvertTo-Json

Invoke-RestMethod `
  -Method Post `
  -Uri http://localhost:8080/api/owner/profile `
  -Headers @{ Authorization = "Bearer $token" } `
  -ContentType "application/json" `
  -Body $body

เช็ก database:

SELECT id, user_id, full_name, phone FROM owner_profiles;

ข้อสำคัญ:

client ไม่ควรส่ง user_id
backend ต้องเอา user_id จาก JWT เอง
8. Pet Test
Create Pet
$body = @{
  name = "Milo"
  species = "dog"
  breed_id = $null
  birth_date = "2021-01-01"
  sex = "male"
} | ConvertTo-Json

Invoke-RestMethod `
  -Method Post `
  -Uri http://localhost:8080/api/pets `
  -Headers @{ Authorization = "Bearer $token" } `
  -ContentType "application/json" `
  -Body $body

เช็ก database:

SELECT id, owner_profile_id, name, public_pet_id FROM pets;

สิ่งที่ต้องดู:

1. pet ถูกสร้างจริงไหม
2. owner_profile_id ถูกผูกกับ owner ที่ login อยู่ไหม
3. public_pet_id ถูก generate ไหม

ข้อสำคัญ:

client ห้ามส่ง owner_profile_id เอง
backend ต้องหา owner_profile จาก JWT user_id
9. Appointment Test แบบภาพรวม

Appointment ต้องเชื่อม 3 อย่าง:

Owner Profile
Pet
Clinic Profile

เวลาทดสอบ appointment ให้เช็กว่า:

1. owner สร้าง appointment ให้ pet ของตัวเองได้
2. owner สร้าง appointment ให้ pet ของคนอื่นไม่ได้
3. clinic เห็นเฉพาะ appointment ของ clinic ตัวเอง
4. appointment ที่ cancelled ไม่ควรถูกนับเป็น active patient

เช็ก database:

SELECT id, owner_profile_id, clinic_profile_id, pet_id, status, scheduled_at
FROM appointments;
10. Clinic Patient Test

ระบบนี้ไม่มี patients table

Patient ถูก derived จาก appointment

กฎคือ:

ถ้า pet มี appointment กับ clinic
และ appointment status ไม่ใช่ cancelled
pet ตัวนั้นถือเป็น patient ของ clinic

เช็กด้วย SQL:

SELECT DISTINCT pet_id, clinic_profile_id
FROM appointments
WHERE status != 'cancelled';

ถ้า clinic A ไม่ควรเห็น patient ของ clinic B ให้เช็กว่า query ในระบบ scope ด้วย:

clinic_profile_id ของ clinic ที่ login อยู่
11. Medical Record Test

Medical record ต้องผูกกับ:

clinic_profile_id
pet_id
created_by_user_id
optional appointment_id

เช็ก database:

SELECT id, clinic_profile_id, pet_id, appointment_id, created_by_user_id
FROM medical_records;

กฎสำคัญ:

1. clinic สร้าง medical record ได้เฉพาะ pet ที่เป็น patient ของ clinic ตัวเอง
2. clinic อื่นอ่าน/แก้ record นี้ไม่ได้
3. client ห้ามส่ง clinic_profile_id เอง
4. client ห้ามส่ง created_by_user_id เอง
5. PATCH แก้ได้เฉพาะ clinical fields ไม่ใช่ ownership fields

เวลาค้น record ควรมี scope:

record_id + clinic_profile_id

ไม่ควรค้นด้วย:

record_id อย่างเดียว
12. Error Code ที่ต้องจำ
200 OK
request สำเร็จ
201 Created
สร้างข้อมูลใหม่สำเร็จ
400 Bad Request
ข้อมูลที่ส่งมาไม่ถูกต้อง เช่น body ผิด, field หาย, format ผิด
401 Unauthorized
ยังไม่ได้ login
token หาย
token ผิด
token หมดอายุ
403 Forbidden
login แล้ว แต่ role ไม่มีสิทธิ์เข้า endpoint นี้
เช่น owner ไปยิง clinic endpoint
404 Not Found
ไม่เจอข้อมูลใน scope ที่ user มีสิทธิ์เห็น
อาจเป็น id ไม่มีจริง หรือเป็นข้อมูลของคนอื่น
500 Internal Server Error
server/database มีปัญหาภายใน
ต้องดู backend logs
13. Troubleshooting Checklist
กรณี 1: เพื่อนบอกว่า “API เข้าไม่ได้”

เช็กตามนี้:

1. Backend รันอยู่ไหม
2. /health ผ่านไหม
3. URL ถูกไหม
4. Method ถูกไหม เช่น GET/POST/PATCH
5. Port ถูกไหม
6. Backend log มี error ไหม

คำสั่ง:

Invoke-RestMethod http://localhost:8080/health
กรณี 2: เพื่อนบอกว่า “Database ไม่ทำงาน”

เช็กตามนี้:

1. docker compose ps
2. docker compose logs
3. /health/db ผ่านไหม
4. .env ตั้งค่า DB ถูกไหม
5. DB name/user/password ตรงกับ docker-compose.yml ไหม

คำสั่ง:

docker compose ps
docker compose logs
Invoke-RestMethod http://localhost:8080/health/db
กรณี 3: ได้ 401 Unauthorized

สาเหตุที่พบบ่อย:

1. ไม่ได้ส่ง Authorization header
2. ใช้ token ผิด
3. ลืมคำว่า Bearer
4. token หมดอายุ
5. login กับ environment คนละตัว เช่น login local แต่เอา token ไปยิง Render

header ที่ถูก:

Authorization: Bearer <token>
กรณี 4: ได้ 403 Forbidden

แปลว่า:

login แล้ว แต่ role ไม่ถูก

ตัวอย่าง:

owner ไปยิง /api/clinic/...
clinic ไปยิง /api/owner/...

วิธีเช็ก:

1. login ใหม่
2. GET /api/me
3. ดู role ที่ backend ตอบกลับ
กรณี 5: ได้ 404 Not Found ทั้งที่คิดว่าข้อมูลมี

อาจเกิดจาก:

1. id ผิด
2. ข้อมูลอยู่ใน database คนละ environment
3. ข้อมูลเป็นของ owner/clinic คนอื่น
4. query ถูก scope ด้วย owner_profile_id หรือ clinic_profile_id

ในระบบนี้ 404 อาจไม่ได้แปลว่าไม่มีข้อมูลจริงเสมอไป
แต่อาจแปลว่า:

ไม่พบข้อมูลในขอบเขตสิทธิ์ของ user คนนี้
กรณี 6: Backend start ไม่ได้

เช็ก log ว่าติดตรงไหน:

database connection failed
database migration failed
failed to start server

ถ้าเป็น database connection failed:

เช็ก PostgreSQL / .env / DATABASE_URL

ถ้าเป็น migration failed:

เช็ก SQL migration หรือ schema conflict

ถ้าเป็น failed to start server:

เช็กว่า port ถูกใช้ไปแล้วไหม
14. Render Production Check

ถ้าเช็กระบบบน Render ให้ใช้ลำดับนี้:

1. เปิด Render Web Service
2. ดู Logs
3. เช็กว่า deploy ล่าสุดสำเร็จไหม
4. เช็ก environment variables
5. ยิง /health
6. ยิง /health/db
7. ทดสอบ login
8. ทดสอบ endpoint สำคัญ

สิ่งที่ต้องระวัง:

Local database กับ Render database เป็นคนละตัวกัน
ข้อมูลใน local จะไม่ไปอยู่บน Render อัตโนมัติ

ดังนั้นถ้า local มีข้อมูล แต่ Render หาไม่เจอ ไม่ได้แปลว่าระบบพังเสมอไป
อาจเป็นเพราะเป็น database คนละ environment

15. คำสั่งที่ใช้บ่อย
Start PostgreSQL local
docker compose up -d
Stop PostgreSQL local
docker compose down
Check containers
docker compose ps
View container logs
docker compose logs
Run backend
go run ./cmd/api
Run tests
go test ./...
Check health
Invoke-RestMethod http://localhost:8080/health
Check database health
Invoke-RestMethod http://localhost:8080/health/db
Enter PostgreSQL
docker exec -it petnexus-postgres psql -U postgres -d petnexus
16. Quick Debug Order

เวลาเจอปัญหา ให้ไล่แบบนี้:

1. Server เปิดไหม?
   → GET /health

2. Database ต่อได้ไหม?
   → GET /health/db

3. Login ได้ไหม?
   → POST /api/auth/login

4. Token ใช้ได้ไหม?
   → GET /api/me

5. Role ถูกไหม?
   → ดู response จาก /api/me

6. ข้อมูลอยู่ใน DB ไหม?
   → psql SELECT ...

7. ข้อมูลอยู่ใน scope ของ user นี้ไหม?
   → เช็ก owner_profile_id หรือ clinic_profile_id

8. Backend log บอกอะไร?
   → ดู terminal หรือ Render Logs
17. หลักจำสั้น ๆ
401 = ยังไม่รู้ว่าคุณคือใคร
403 = รู้ว่าคุณคือใคร แต่คุณไม่มีสิทธิ์
404 = ไม่เจอข้อมูลใน scope ที่คุณมีสิทธิ์เห็น
500 = backend/database มีปัญหา
Owner scope = owner_profile_id
Clinic scope = clinic_profile_id
User identity = user_id จาก JWT
อย่าเชื่อ ownership field จาก client
ให้ derive จาก JWT เสมอ
