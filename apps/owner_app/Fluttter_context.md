# Flutter Context
Project: PetNexus Owner App

---

# Overview

PetNexus Owner App is a Flutter application for pet owners.

The project follows Clean Architecture with a lightweight structure.

Presentation
↓

Controller (Provider + ChangeNotifier)

↓

Repository

↓

ApiClient (Dio)

↓

Backend (Go Fiber + PostgreSQL)

The application uses JWT Authentication.

JWT is stored securely using FlutterSecureStorage.

---

# Technology

Flutter

Provider

ChangeNotifier

Dio

Flutter Secure Storage

Responsive Layout Extension

REST API

---

# Folder Structure

lib/

app/

core/

shared/

features/

main.dart

---

# Folder Detail

## app

Contains

- routes
- auth gate
- app

---

## core

Contains reusable project-wide utilities.

core/

constants/

network/

services/

theme/

core.dart

---

### constants

Contains

AppColors

AppTextStyles

AppSpacing

AppRadius

---

### network

Contains

ApiClient

ApiConfig

Dio Interceptor

---

### services

Contains

JWTService

Token Helper

---

## shared

Reusable widgets.

Never place business logic here.

shared/

widgets/

shared.dart

---

Widgets

AppButton

AppCard

AppDialog

AppLogo

AppScaffold

AppSectionTitle

AppSocialButton

AppTextField

---

## features

Each feature owns

controller

repository

models

screens

widgets

Example

features/

auth/

home/

pet/

clinic/

appointment/

profile/

---

# Architecture Rules

UI

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

# Screen Responsibilities

A Screen should

Display UI

Read TextEditingController

Call Controller

Navigate

Display Error

A Screen should NOT

Call API

Contain Business Logic

Create Dio

Create Repository manually

---

# Controller Responsibilities

Controller extends ChangeNotifier.

Responsible for

Validation

Calling Repository

Managing State

notifyListeners()

Controller must NOT

Call Dio

Parse JSON

Build Widgets

---

# Repository Responsibilities

Repository handles

REST API

JSON Parsing

DTO

Response Model

Repository communicates with ApiClient only.

Repository must NOT

Build UI

Use notifyListeners()

---

# ApiClient

All requests go through

ApiClient.instance.dio

Never create

Dio()

inside repositories.

---

# Authentication Flow

App

↓

Splash

↓

AuthGate

↓

checkAuthentication()

↓

Token exists?

YES

↓

loadCurrentUser()

↓

Home

NO

↓

First Screen

---

# JWT

Stored using

FlutterSecureStorage

Never use

SharedPreferences

JWT Service provides

saveToken()

getToken()

deleteToken()

isLoggedIn()

---

# Login Flow

Login Screen

↓

AuthController.login()

↓

AuthRepository.login()

↓

ApiClient

↓

Backend

↓

Receive JWT

↓

Save JWT

↓

loadCurrentUser()

↓

Home

---

# Register Flow

Register Screen

↓

AuthController.register()

↓

AuthRepository.register()

↓

Backend

↓

Receive JWT

↓

Save JWT

↓

Complete Profile

---

# Models

Every API response must have a Model.

Never use

Map<String,dynamic>

inside Screen.

Example

LoginResponse

RegisterResponse

UserModel

OwnerProfileModel

---

# Provider

Controllers are injected using Provider.

Correct

context.read<AuthController>()

context.watch<AuthController>()

Wrong

AuthController()

---

# Responsive

Always use

context.nw()

context.nh()

context.nf()

Never use

width: 300

height: 60

fontSize: 18

---

# Shared Widgets

Always reuse

AppButton

AppCard

AppTextField

AppLogo

AppScaffold

AppSectionTitle

AppDialog

AppSocialButton

Never recreate these widgets.

---

# Navigation

Use Named Routes.

Example

/

login

register

complete-profile

home

profile

Never hardcode navigation logic everywhere.

---

# Error Handling

Repository

throws Exception

↓

Controller

updates errorMessage

↓

UI

shows AppDialog or SnackBar

---

# Coding Style

Prefer StatelessWidget.

Use StatefulWidget only when local state is required.

Business logic belongs in Controller.

API belongs in Repository.

UI belongs in Screen.

---

# Design System

Colors

AppColors.primary

AppColors.background

AppColors.surface

AppColors.error

AppColors.success

AppColors.textPrimary

AppColors.textSecondary

Typography

AppTextStyles.heading()

AppTextStyles.title()

AppTextStyles.body()

AppTextStyles.caption()

AppTextStyles.button()

Spacing

AppSpacing.xs

AppSpacing.sm

AppSpacing.md

AppSpacing.lg

AppSpacing.xl

---

# Code Convention

Use final whenever possible.

Avoid magic numbers.

Avoid duplicate widgets.

Split large widgets into reusable widgets.

Never duplicate colors.

Never duplicate TextStyles.

---

# API Convention

Repository returns

Response Models

or throws Exception.

Never return raw JSON to UI.

---

# AI Agent Rules

When generating code

DO

Use existing Shared Widgets.

Use existing Design System.

Use Provider.

Use Repository pattern.

Use Response Models.

Follow folder structure.

Respect Responsive Layout.

DON'T

Create duplicate widgets.

Add parameters to shared widgets unless requested.

Call Dio directly from Screen.

Call REST API from Controller.

Hardcode colors.

Hardcode text styles.

Ignore existing architecture.

---

# Goal

Maintain a scalable, reusable, production-ready Flutter project following a consistent architecture and design system.