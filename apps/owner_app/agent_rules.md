# AI Agent Rules
Project: PetNexus Owner App (Flutter)

Version: 1.0

---

# IMPORTANT

This repository is **Frontend only**.

The AI Agent MUST NOT modify, generate, or redesign Backend code.

Backend is maintained separately.

Assume every Backend API already exists.

The AI Agent is responsible ONLY for Flutter code.

---

# Scope

Allowed

✅ Flutter

✅ UI

✅ Provider

✅ Controller

✅ Repository

✅ Models

✅ Responsive Layout

✅ Shared Widgets

✅ Navigation

✅ State Management

✅ JWT Handling (Frontend)

Not Allowed

❌ Backend

❌ Go

❌ Fiber

❌ PostgreSQL

❌ SQL

❌ GORM

❌ API Design

❌ Authentication Logic on Server

❌ Database Schema

---

# Architecture

Always follow

Screen

↓

Controller

↓

Repository

↓

ApiClient

↓

Backend

Never skip a layer.

---

# Backend Rules

The AI Agent MUST assume

Backend is already finished.

Backend endpoints already exist.

Backend request format is fixed.

Backend response format is fixed.

Never redesign API.

Never rename endpoint.

Never suggest changing Backend.

Never modify request body unless explicitly requested.

Never modify response structure.

---

# API Rules

Always use

ApiClient.instance.dio

Never create

Dio()

inside repositories.

Never use

http.post()

http.get()

inside Flutter code.

---

# Repository Rules

Repository is the only layer allowed to call REST API.

Repository responsibilities

• API request

• JSON parsing

• Response model

Repository must NOT

• contain UI

• call notifyListeners()

• navigate

---

# Controller Rules

Controller extends ChangeNotifier.

Responsibilities

Validation

Calling Repository

Managing UI State

notifyListeners()

Controller must NOT

Call REST API directly

Parse JSON

Create Dio

Contain Widgets

---

# Screen Rules

Screen responsibilities

Display UI

Read user input

Call Controller

Navigation

Show Snackbar

Show Dialog

Screen must NOT

Call Repository directly

Call ApiClient

Contain business logic

Parse JSON

---

# Shared Widget Rules

Always reuse

AppButton

AppCard

AppDialog

AppLogo

AppScaffold

AppSectionTitle

AppSocialButton

AppTextField

Never recreate existing widgets.

Never duplicate widgets.

If a widget already exists

Use it.

Do NOT build another version.

---

# Shared Widget Modification Rule

Do NOT add new parameters to Shared Widgets.

Do NOT change Widget API.

If additional behavior is required

Ask first.

---

# Design System

Always use

AppColors

AppTextStyles

AppSpacing

AppRadius

Never use

Color(...)

TextStyle(...)

EdgeInsets.only(...)

directly unless absolutely necessary.

---

# Responsive

Every size must use

context.nw()

context.nh()

context.nf()

Never hardcode

width

height

font size

padding

margin

---

# Provider

Use

context.read()

context.watch()

Never create

Controller()

inside Screen.

---

# Authentication

JWT only.

Store JWT using

FlutterSecureStorage.

Never use

SharedPreferences.

---

# Models

Every API Response must have

Response Model.

Never return

Map<String,dynamic>

to Screen.

---

# Error Handling

Repository

↓

throw Exception

↓

Controller

↓

errorMessage

↓

UI

↓

Dialog / SnackBar

---

# Navigation

Always use

Named Routes

or centralized navigation.

Never duplicate navigation logic.

---

# Code Style

Prefer StatelessWidget.

Split large Widgets.

Avoid duplicate code.

Avoid magic numbers.

Use final whenever possible.

Keep build() readable.

---

# Imports

Prefer

shared/shared.dart

core/core.dart

instead of importing many individual files.

---

# AI Behavior Rules

When generating code

DO

✓ Reuse existing Widgets

✓ Reuse existing Controllers

✓ Reuse existing Repository

✓ Respect existing architecture

✓ Keep naming consistent

✓ Generate compile-ready code

✓ Ask before changing architecture

DON'T

✗ Rewrite project structure

✗ Replace Provider

✗ Replace Dio

✗ Replace JWT system

✗ Modify Backend

✗ Modify API contract

✗ Duplicate Widgets

✗ Invent new design system

✗ Add unnecessary packages

---

# Before Generating Code

The AI Agent must verify

1. Widget already exists?

→ Reuse it.

2. Controller already exists?

→ Reuse it.

3. Repository already exists?

→ Reuse it.

4. Model already exists?

→ Reuse it.

5. Shared Widget API changed?

→ Ask first.

---

# Golden Rule

Frontend only.

Never redesign Backend.

Always respect the existing architecture.

Consistency is more important than creating new code.