# PetNexus Backend Architecture Map

Audit date: 2026-07-14  
Source of truth: current repository code.

## Repository folder map

```text
petnexus-backend/
├─ cmd/
│  └─ api/
│     └─ main.go                  # startup, dependency construction, router boot
├─ internal/
│  ├─ config/                     # env + godotenv config loading
│  ├─ database/                   # PostgreSQL connection, ping, startup migrations
│  ├─ dto/                        # request/response API shapes
│  ├─ handlers/                   # HTTP binding and response translation
│  ├─ middleware/                 # JWT auth and role guard middleware
│  ├─ models/                     # GORM models and constants
│  ├─ repositories/               # GORM/raw SQL database access
│  ├─ routes/                     # route registration and dependency bundle
│  ├─ services/                   # validation, ownership, business rules, response mapping
│  └─ utils/                      # response helpers, AppError, JWT, bcrypt, public pet ID
├─ migrations/                    # manual SQL equivalents of startup migrations
├─ docs/                          # implementation docs, progress logs, sprint summaries
├─ docker-compose.yml             # local PostgreSQL
├─ go.mod / go.sum
└─ README.md
```

Future-only placeholders exist for QR, authorization, visits, timeline, notifications, audit logs, and clinic staff. They are not wired into routes and should not be treated as implemented features.

## Startup flow

```mermaid
flowchart TD
    A[main.go] --> B[config.Load]
    B --> C[database.ConnectPostgres]
    C --> D[database.PingPostgres]
    D --> E[database.RunMigrations]
    E --> F[Construct repositories]
    F --> G[Construct services]
    G --> H[Construct handlers]
    H --> I[gin.Default]
    I --> J[routes.Register]
    J --> K[router.Run]
```

Text explanation:

- `config.Load` reads `.env` when present, then environment variables.
- `ConnectPostgres` prefers `DATABASE_URL`; otherwise it builds a local `DB_*` DSN.
- `RunMigrations` executes guarded SQL steps before any route serves traffic.
- Dependency construction is manual and centralized in `cmd/api/main.go`.

## Dependency construction flow

```mermaid
flowchart LR
    DB[(GORM DB)]

    DB --> UserRepo[UserRepository]
    DB --> OwnerRepo[OwnerProfileRepository]
    DB --> BreedRepo[BreedRepository]
    DB --> PetRepo[PetRepository]
    DB --> ClinicRepo[ClinicProfileRepository]
    DB --> AppointmentRepo[AppointmentRepository]
    DB --> PatientRepo[ClinicPatientRepository]
    DB --> MedicalRepo[MedicalRecordRepository]

    UserRepo --> AuthSvc[AuthService]
    OwnerRepo --> OwnerSvc[OwnerProfileService]
    BreedRepo --> BreedSvc[BreedService]
    PetRepo --> PetSvc[PetService]
    ClinicRepo --> ClinicSvc[ClinicProfileService]
    PetRepo --> LookupSvc[ClinicPetLookupService]
    AppointmentRepo --> OwnerApptSvc[OwnerAppointmentService]
    AppointmentRepo --> ClinicApptSvc[ClinicAppointmentService]
    PatientRepo --> PatientSvc[ClinicPatientService]
    MedicalRepo --> MedicalSvc[MedicalRecordService]

    AuthSvc --> AuthHandler
    OwnerSvc --> OwnerHandler
    BreedSvc --> BreedHandler
    PetSvc --> PetHandler
    ClinicSvc --> ClinicHandler
    LookupSvc --> LookupHandler
    OwnerApptSvc --> OwnerAppointmentHandler
    ClinicApptSvc --> ClinicAppointmentHandler
    PatientSvc --> ClinicPatientHandler
    MedicalSvc --> MedicalRecordHandler
```

## Route registration flow

Routes are registered in `internal/routes/routes.go`.

```text
GET /health
GET /health/db

/api/auth
  POST /register
  POST /login

/api/me
  AuthMiddleware

/api/owner
  AuthMiddleware
  RequireRole(owner)
  POST  /profile
  GET   /profile
  PATCH /profile
  POST  /appointments
  GET   /appointments
  GET   /appointments/:id
  PATCH /appointments/:id/cancel

/api/breeds
  GET public

/api/pets
  AuthMiddleware
  RequireRole(owner)
  POST  /
  GET   /
  GET   /:id
  PATCH /:id

/api/clinic
  AuthMiddleware
  RequireRole(clinic, clinic_staff)
  POST  /profile
  GET   /profile
  PATCH /profile
  GET   /pet-lookup
  GET   /patients
  GET   /patients/:petId
  POST  /patients/:petId/medical-records
  POST  /appointments
  GET   /appointments
  GET   /appointments/:id
  PATCH /appointments/:id/status
  PATCH /appointments/:id/cancel
  GET   /medical-records
  GET   /medical-records/:recordId
  PATCH /medical-records/:recordId
```

## Authentication middleware flow

```mermaid
sequenceDiagram
    participant Client
    participant AuthMiddleware
    participant JWT as utils.ParseAccessToken
    participant GinContext
    participant Handler

    Client->>AuthMiddleware: Authorization: Bearer token
    AuthMiddleware->>AuthMiddleware: validate header format
    AuthMiddleware->>JWT: parse HS256 token with exp required
    JWT-->>AuthMiddleware: claims userID, role
    AuthMiddleware->>GinContext: set userID/userRole
    AuthMiddleware->>Handler: c.Next()
```

Role middleware flow:

```mermaid
flowchart TD
    A[RequireRole] --> B[Read userRole from Gin context]
    B --> C{Role present?}
    C -- no --> D[401 UNAUTHORIZED]
    C -- yes --> E{Role allowed?}
    E -- no --> F[403 FORBIDDEN_ROLE]
    E -- yes --> G[Next handler]
```

## Handler -> Service -> Repository -> Database flow

```mermaid
flowchart LR
    Client --> Route
    Route --> Middleware
    Middleware --> Handler
    Handler --> DTO[Request DTO / path / query]
    DTO --> Service
    Service --> Repo[Repository interface]
    Repo --> DB[(PostgreSQL)]
    DB --> Repo
    Repo --> Service
    Service --> ResponseDTO[Response DTO]
    ResponseDTO --> Handler
    Handler --> Envelope[utils.Success / utils.Error]
    Envelope --> Client
```

Responsibility split:

- Handler: parse input and translate errors.
- Service: validate input, resolve JWT-owned profile, enforce business rules.
- Repository: execute database queries and map GORM errors to repository errors.
- Database: enforce schema constraints and indexes.

## Table relationships

```mermaid
erDiagram
    users ||--o| owner_profiles : "user_id"
    users ||--o| clinic_profiles : "user_id"
    owner_profiles ||--o{ pets : "owner_profile_id"
    breeds ||--o{ pets : "breed_id"
    owner_profiles ||--o{ appointments : "owner_profile_id"
    clinic_profiles ||--o{ appointments : "clinic_profile_id"
    pets ||--o{ appointments : "pet_id"
    clinic_profiles ||--o{ medical_records : "clinic_profile_id"
    pets ||--o{ medical_records : "pet_id"
    appointments ||--o| medical_records : "appointment_id"
    users ||--o{ medical_records : "created_by_user_id"
```

Important derived relationship:

```text
Clinic patient =
unique appointments.pet_id
where appointments.clinic_profile_id = current clinic profile
and appointments.status <> 'cancelled'
```

There is no separate `patients` table.

## Full request flow: Login

```mermaid
sequenceDiagram
    participant Client
    participant Routes
    participant AuthHandler
    participant AuthService
    participant UserRepo
    participant DB
    participant JWT

    Client->>Routes: POST /api/auth/login
    Routes->>AuthHandler: Login
    AuthHandler->>AuthHandler: ShouldBindJSON LoginRequest
    AuthHandler->>AuthService: Login(req)
    AuthService->>AuthService: normalize email, validate password present
    AuthService->>UserRepo: FindByEmail(email)
    UserRepo->>DB: SELECT users WHERE email = ?
    DB-->>UserRepo: user or not found
    AuthService->>AuthService: bcrypt compare
    AuthService->>JWT: GenerateAccessToken(user.ID, user.Role)
    AuthService-->>AuthHandler: AuthResponse
    AuthHandler-->>Client: 200 success
```

Flow details:

- Route: `POST /api/auth/login`
- Middleware: none
- Handler: `AuthHandler.Login`
- Request DTO: `dto.LoginRequest`
- Service: `authService.Login`
- Repository: `UserRepository.FindByEmail`
- Models/tables: `models.User`, `users`
- Response DTO: `dto.AuthResponse`
- Error path: invalid JSON 400; missing email/password 422; bad credentials 401; unexpected DB/JWT error 500.

## Full request flow: Create Pet

```mermaid
sequenceDiagram
    participant Client
    participant Middleware
    participant PetHandler
    participant PetService
    participant OwnerRepo
    participant BreedRepo
    participant PetRepo
    participant DB

    Client->>Middleware: POST /api/pets + owner JWT
    Middleware->>Middleware: AuthMiddleware + RequireRole(owner)
    Middleware->>PetHandler: CreatePet
    PetHandler->>PetHandler: bind CreatePetRequest
    PetHandler->>PetService: CreatePet(userID, req)
    PetService->>OwnerRepo: FindByUserID(userID)
    OwnerRepo->>DB: SELECT owner_profiles WHERE user_id = ?
    PetService->>BreedRepo: optional FindByID(breedID)
    PetService->>PetService: validate species/name/breed/gender/date/weight
    PetService->>PetService: GeneratePublicPetID
    PetService->>PetRepo: Create(pet)
    PetRepo->>DB: INSERT pets
    PetService-->>PetHandler: PetResponse
    PetHandler-->>Client: 201 success
```

Flow details:

- Route: `POST /api/pets`
- Middleware: `AuthMiddleware`, `RequireRole(owner)`
- Handler: `PetHandler.CreatePet`
- Request DTO: `dto.CreatePetRequest`
- Service: `petService.CreatePet`
- Repositories: owner profile, breed, pet
- Models/tables: `owner_profiles`, `breeds`, `pets`
- Response DTO: `dto.PetResponse`
- Error path: 400 invalid body/validation; 401 unauthenticated; 403 wrong role; 404 missing owner profile or breed; 500 unexpected DB/randomness.

## Full request flow: Create Owner Appointment

```mermaid
sequenceDiagram
    participant Client
    participant Middleware
    participant Handler as OwnerAppointmentHandler
    participant Service as OwnerAppointmentService
    participant OwnerRepo
    participant ClinicRepo
    participant PetRepo
    participant AppointmentRepo
    participant DB

    Client->>Middleware: POST /api/owner/appointments + owner JWT
    Middleware->>Handler: CreateAppointment
    Handler->>Handler: bind CreateOwnerAppointmentRequest
    Handler->>Service: CreateOwnerAppointment(userID, req)
    Service->>OwnerRepo: FindByUserID(userID)
    Service->>ClinicRepo: FindByID(clinic_profile_id)
    Service->>PetRepo: FindByIDAndOwnerProfileID(pet_id, ownerProfile.ID)
    Service->>Service: validate type/time/duration/note
    Service->>AppointmentRepo: Create(appointment)
    AppointmentRepo->>DB: INSERT appointments
    Service-->>Handler: AppointmentResponse
    Handler-->>Client: 201 success
```

Flow details:

- Route: `POST /api/owner/appointments`
- Middleware: `AuthMiddleware`, `RequireRole(owner)`
- Handler: `OwnerAppointmentHandler.CreateAppointment`
- Request DTO: `dto.CreateOwnerAppointmentRequest`
- Service: `ownerAppointmentService.CreateOwnerAppointment`
- Repositories: owner profile, clinic profile, pet, appointment
- Models/tables: `owner_profiles`, `clinic_profiles`, `pets`, `appointments`
- Response DTO: `dto.AppointmentResponse`
- Error path: invalid UUID/body/time/type/duration 400; missing owner profile/clinic/pet 404; wrong role 403; unexpected DB 500.

## Full request flow: List Clinic Patients

```mermaid
sequenceDiagram
    participant Client
    participant Middleware
    participant Handler as ClinicPatientHandler
    participant Service as ClinicPatientService
    participant ClinicRepo
    participant PatientRepo
    participant DB

    Client->>Middleware: GET /api/clinic/patients + clinic JWT
    Middleware->>Handler: ListPatients
    Handler->>Service: ListClinicPatients(userID, filters)
    Service->>ClinicRepo: FindByUserID(userID)
    Service->>Service: validate q/species/status/limit/offset/sort
    Service->>PatientRepo: FindPatientsByClinicProfileID(clinicProfile.ID, filters)
    PatientRepo->>DB: raw SQL derives pets from non-cancelled appointments
    PatientRepo->>DB: load pets with Breed and OwnerProfile
    Service-->>Handler: []ClinicPatientListItemResponse
    Handler-->>Client: 200 success
```

Flow details:

- Route: `GET /api/clinic/patients`
- Middleware: `AuthMiddleware`, `RequireRole(clinic, clinic_staff)`
- Handler: `ClinicPatientHandler.ListPatients`
- Request DTO/filter: `dto.ClinicPatientFilters`
- Service: `clinicPatientService.ListClinicPatients`
- Repository: `ClinicPatientRepository.FindPatientsByClinicProfileID`
- Models/tables: `appointments`, `pets`, `owner_profiles`, `breeds`, `clinic_profiles`
- Response DTO: `dto.ClinicPatientListItemResponse`
- Error path: invalid query 400; missing auth 401; owner role 403; missing clinic profile 404; DB error 500.

## Full request flow: Create Medical Record

```mermaid
sequenceDiagram
    participant Client
    participant Middleware
    participant Handler as MedicalRecordHandler
    participant Service as MedicalRecordService
    participant ClinicRepo
    participant RecordRepo
    participant DB

    Client->>Middleware: POST /api/clinic/patients/:petId/medical-records
    Middleware->>Handler: CreateMedicalRecord
    Handler->>Handler: parse petId and bind CreateMedicalRecordRequest
    Handler->>Service: CreateMedicalRecord(userID, petID, req)
    Service->>ClinicRepo: FindByUserID(userID)
    Service->>Service: validate visitAt, complaint, follow-up, vitals
    Service->>RecordRepo: PetHasNonCancelledAppointmentWithClinic(clinicID, petID)
    RecordRepo->>DB: COUNT appointments scoped by clinic/pet/status
    alt appointmentId supplied
        Service->>RecordRepo: FindUsableAppointmentForMedicalRecord
        RecordRepo->>DB: SELECT appointment scoped by id/clinic/pet/not cancelled
        Service->>RecordRepo: MedicalRecordExistsByAppointmentID
        RecordRepo->>DB: COUNT medical_records WHERE appointment_id = ?
    end
    Service->>RecordRepo: Create(record)
    RecordRepo->>DB: INSERT medical_records
    Service->>RecordRepo: FindByIDAndClinicProfileID
    RecordRepo->>DB: SELECT with Pet/Breed/OwnerProfile/Appointment/CreatedByUser
    Service-->>Handler: MedicalRecordDetailResponse
    Handler-->>Client: 201 success
```

Flow details:

- Route: `POST /api/clinic/patients/:petId/medical-records`
- Middleware: `AuthMiddleware`, `RequireRole(clinic, clinic_staff)`
- Handler: `MedicalRecordHandler.CreateMedicalRecord`
- Request DTO: `dto.CreateMedicalRecordRequest`
- Service: `medicalRecordService.CreateMedicalRecord`
- Repositories: clinic profile, medical record
- Models/tables: `clinic_profiles`, `appointments`, `medical_records`, `pets`, `owner_profiles`, `users`
- Response DTO: `dto.MedicalRecordDetailResponse`
- Error path: invalid pet ID/body/date/vitals 400; unauthenticated 401; wrong role 403; missing clinic profile/patient/appointment 404; duplicate appointment medical record 409; unexpected DB 500.

## Ownership enforcement map

| Rule | Enforcement location |
| --- | --- |
| Owner can access only own profile | owner route role middleware; `OwnerProfileService` resolves `owner_profiles.user_id` from JWT |
| Owner can access only own pets | owner route role middleware; `PetService.findCurrentOwnerProfile`; `PetRepository.FindByIDAndOwnerProfileID` |
| Owner can use only own pets in appointments | `OwnerAppointmentService.CreateOwnerAppointment`; `PetRepository.FindByIDAndOwnerProfileID` |
| Clinic can access only own profile | clinic route role middleware; `ClinicProfileService` resolves `clinic_profiles.user_id` from JWT |
| Clinic can manage own appointments | `ClinicAppointmentService.currentClinicProfile`; `AppointmentRepository.FindByIDAndClinicProfileID` |
| Clinic patients scoped to current clinic | `ClinicPatientService.currentClinicProfile`; raw SQL filters `appointments.clinic_profile_id` |
| Clinic can read/update own medical records | `MedicalRecordService.currentClinicProfile`; `MedicalRecordRepository.FindByIDAndClinicProfileID` and scoped `Update` |
| Clinic cannot create medical records for unrelated pets | `MedicalRecordService.ensureClinicPatient`; repository counts non-cancelled appointments by clinic/pet |
| Appointment-linked medical records validate pet and clinic | `MedicalRecordService.ensureUsableAppointment`; repository filters by appointment ID, clinic ID, pet ID, and non-cancelled status |
| Immutable medical record ownership fields cannot be patched | `dto.UpdateMedicalRecordRequest` omits ownership fields; `MedicalRecordRepository.Update` uses explicit update map |
