# PetNexus Clinic Platform (Frontend)

## Project Goal

สร้างระบบ Web สำหรับคลินิกสัตวแพทย์ (Clinic Management System)

Frontend จะพัฒนาแยกจาก Backend อย่างสมบูรณ์

Backend ถือว่าเป็น API Provider เท่านั้น

Frontend จะไม่แก้ไข Backend หรือ Database


---

# Development Strategy

เราจะพัฒนาแบบ

Design First

ไม่ใช่

API First

ลำดับการพัฒนา

1. Theme
2. Design System
3. Shared Components
4. Layout
5. Responsive
6. Pages
7. State Management
8. API Integration

จนกว่า UI ทั้งระบบจะเสร็จ จะยังไม่เชื่อม API


---

# Design Direction

Theme ของระบบ

- Modern SaaS Dashboard
- Veterinary Clinic
- Clean UI
- Professional
- Spacious
- Soft Shadow
- Large Radius

Reference

- Linear
- Stripe Dashboard
- Shadcn UI
- Notion
- Vercel Dashboard


---

# Tech Stack

Next.js

React

TypeScript

TailwindCSS

Shadcn/UI

Axios

React Hook Form

Zod


---

# Frontend Architecture

ใช้แนวคิด

Feature First

src/

app/

features/

shared/

layouts/

hooks/

services/

constants/

types/

utils/

assets/


---

# Shared First

ทุกอย่างที่สามารถใช้ซ้ำได้

ต้องสร้างใน

shared/

ก่อน

ห้ามสร้าง Component ซ้ำใน Feature

ตัวอย่าง

Button

Input

Card

Dialog

Avatar

Badge

Table

Pagination

Sidebar

Topbar

Loading

Empty

Search

StatCard


---

# Layout First

ก่อนสร้างหน้า

ต้องสร้าง Layout ก่อน

เช่น

AuthLayout

DashboardLayout

MainLayout

SettingsLayout

ทุก Page จะอยู่ภายใต้ Layout


---

# Responsive Strategy

Desktop First

Breakpoints

1440

1280

1024

768

640

ทุก Component ต้อง Responsive ตั้งแต่สร้าง

ห้ามแก้ Responsive ทีหลัง


---

# Design System

ก่อนทำ Feature

ต้องมี

Theme

Colors

Typography

Spacing

Radius

Shadow

Animation

Icons

Button

Input

Card

Dialog

Badge

Avatar

Table

Navbar

Sidebar

Dropdown

Modal


---

# Component Philosophy

Component ต้อง

Reusable

Composable

Stateless ถ้าเป็นไปได้

Typed

แยกหน้าที่ชัดเจน

ไม่ทำหลายหน้าที่ใน Component เดียว


---

# Feature Development Order

Dashboard

↓

Authentication

↓

Patient

↓

QR Check-in

↓

Medical Records

↓

Appointments

↓

Reports

↓

Settings


---

# State Management

State จะทำหลังจาก UI เสร็จ

ก่อนหน้านั้น

Component ต้องเป็น Mock UI

ใช้ Dummy Data

เมื่อ UI เสร็จ

ค่อยเชื่อม

Context

React Query

หรือ Zustand


---

# API Integration

API เป็นขั้นตอนสุดท้าย

Frontend

ไม่เปลี่ยน API

ไม่เปลี่ยน Backend

ไม่เปลี่ยน Database

Frontend มีหน้าที่

เรียก API

Map Data

แสดงผล


---

# Coding Principle

Small Components

Reuse Everything

Single Responsibility

No Duplicate UI

No Hardcode Theme

No Inline Styles

Feature Isolation


---

# Current Planning

ตอนนี้

กำลังสร้าง

Theme

↓

Shared Components

↓

Layouts

↓

Responsive System

หลังจากนั้น

จึงเริ่มทำแต่ละหน้า


---

# UI Reference

Theme หลักอ้างอิงจาก

PetNexus Clinic Platform

ประกอบด้วย

Login

Dashboard

Patient List

QR Check-in

Medical Record

ใช้สี

- Teal
- Navy
- White
- Light Gray

สไตล์

Modern SaaS Dashboard

Clean

Minimal

Professional