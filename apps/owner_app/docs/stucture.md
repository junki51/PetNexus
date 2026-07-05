# Project Structure

## Architecture Philosophy

PetNexus Owner App ใช้แนวคิด

Feature First + Layered Architecture

ทุก Feature จะถูกแยกออกจากกันอย่างชัดเจน
และใช้ Shared Widgets ร่วมกัน

```
lib/
├── app/
├── core/
├── features/
├── layout/
├── shared/
├── repositories/
├── services/
├── models/
└── utils/
```

---

# Folder Responsibilities

## app/

เก็บ

- app.dart
- routes
- providers
- app configuration

ไม่มี Business Logic

---

## core/

เก็บสิ่งที่ใช้ทั้งระบบ

```
core/

constants/

theme/

config/

extensions/

errors/
```

ตัวอย่าง

AppColors

AppTextStyles

AppRadius

AppShadow

AppConfig

---

## features/

แต่ละ Feature แยกออกจากกัน

ตัวอย่าง

```
features/

auth/

home/

pet/

appointment/

notification/

profile/
```

แต่ละ Feature

```
auth/

controllers/

repositories/

screens/

widgets/

models/
```

Feature ต้องไม่พึ่งพา Feature อื่น

---

## shared/

เก็บ Widget ที่ใช้ร่วมกัน

```
shared/

widgets/

dialogs/

animations/
```

ตัวอย่าง

AppButton

AppCard

AppDialog

AppLogo

AppTextField

AppAvatar

AppDivider

AppLoading

AppSpacing

---

## layout/

เก็บ Layout ที่ใช้ซ้ำ

ตัวอย่าง

```
AppScaffold

ResponsiveLayout

BottomNavigation

Drawer

AppHeader
```

---

## repositories/

Repository สำหรับเรียก API

ตัวอย่าง

```
AuthRepository

PetRepository

AppointmentRepository
```

Screen ไม่เรียก API โดยตรง

---

## services/

Service ชั้นล่าง

เช่น

```
ApiService

StorageService

SecureStorageService

NotificationService
```

---

## models/

เก็บ Data Models

ตัวอย่าง

```
UserModel

PetModel

AppointmentModel

MedicalRecordModel
```

---

## utils/

Utility

เช่น

```
Validator

Formatter

DateHelper

StringHelper
```

---

# Layer Architecture

```
Screen

↓

Feature Widget

↓

Controller

↓

Repository

↓

Service

↓

API
```

Screen

ไม่มีหน้าที่

- API
- Database
- Business Logic

---

# Widget Hierarchy

```
Screen

↓

Section

↓

Shared Widget
```

ตัวอย่าง

```
LoginScreen

↓

LoginForm

↓

AppTextField

↓

AppButton
```

---

# State Management

ใช้

Controller

↓

Repository

↓

Service

Controller

จัดการ

- State
- Validation
- Business Logic

Repository

จัดการ

- API

Service

จัดการ

- HTTP
- Local Storage
- Secure Storage

---

# Shared Widget Rule

ก่อนสร้าง Widget ใหม่

ตรวจสอบก่อนว่า

สามารถใช้ Widget เดิมได้หรือไม่

หากใช้ได้

ห้ามสร้างใหม่

Widget กลางทั้งหมด

อยู่ใน

```
shared/widgets/
```

---

# Theme Rule

ทุกหน้าต้องใช้

```
AppColors

AppTextStyles

AppSpacing

AppRadius

AppShadow
```

ห้าม

```
Colors.blue

Colors.red

FontSize 17

BorderRadius.circular(13)
```

โดยตรง

---

# Responsive Rule

ทุก Widget

รองรับ

Phone

Tablet

ใช้

```
context.nw()

context.nh()

context.nf()

context.radius()

context.icon()
```

ห้ามใช้ค่าคงที่

---

# Development Flow

ลำดับการพัฒนา

1. Theme
2. Shared Widgets
3. Layout
4. Screen UI
5. Mock Data
6. Controller
7. Repository
8. API Integration

จนกว่า UI จะเสร็จ

จะยังไม่เชื่อม Backend