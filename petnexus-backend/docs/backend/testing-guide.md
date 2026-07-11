# Backend Testing Guide

## Prerequisites

- Go 1.22+
- Docker Desktop with Docker Compose
- Terminal opened in `petnexus-backend`

Start PostgreSQL:

```powershell
docker compose up -d
docker compose ps
```

## Format and automated tests

Requested formatting command:

```powershell
gofmt -w .
```

If the installed `gofmt` does not recurse into directories, format the changed
Go files explicitly before testing.

Run tests:

```powershell
go test ./...
```

Optional static checks:

```powershell
go vet ./...
```

Run the API:

```powershell
go run ./cmd/api
```

## Base URL

Local:

```powershell
$baseUrl = "http://localhost:8080"
```

Render:

```powershell
$baseUrl = "https://petnexus-api.onrender.com"
```

Run the same functional flow against Render after deployment by changing only
`$baseUrl`. Do not store real tokens in committed files.

## Health checks

```powershell
Invoke-RestMethod -Method GET "$baseUrl/health"
Invoke-RestMethod -Method GET "$baseUrl/health/db"
```

Expected: process status `ok` and database status `connected`.

## Owner flow

Use unique test data:

```powershell
$suffix = [DateTimeOffset]::UtcNow.ToUnixTimeMilliseconds()
$password = "password123"
$ownerEmail = "owner.$suffix@example.com"
```

### Register owner

```powershell
$ownerRegisterBody = @{
  email = $ownerEmail
  phone = "0812345678"
  password = $password
  role = "owner"
} | ConvertTo-Json

Invoke-RestMethod -Method POST `
  -Uri "$baseUrl/api/auth/register" `
  -ContentType "application/json" `
  -Body $ownerRegisterBody
```

Expected: 201 and `data.user.role` equals `owner`.

### Login owner

```powershell
$ownerLoginBody = @{
  email = $ownerEmail
  password = $password
} | ConvertTo-Json

$ownerLogin = Invoke-RestMethod -Method POST `
  -Uri "$baseUrl/api/auth/login" `
  -ContentType "application/json" `
  -Body $ownerLoginBody

$ownerToken = $ownerLogin.data.accessToken
$ownerHeaders = @{ Authorization = "Bearer $ownerToken" }
```

### Current user

```powershell
Invoke-RestMethod -Method GET `
  -Uri "$baseUrl/api/me" `
  -Headers $ownerHeaders
```

### Create owner profile

```powershell
$ownerProfileBody = @{
  first_name = "Sunny"
  last_name = "Example"
  gender = "male"
  date_of_birth = "1995-05-20"
  phone_number = "0812345678"
  province = "Bangkok"
} | ConvertTo-Json

Invoke-RestMethod -Method POST `
  -Uri "$baseUrl/api/owner/profile" `
  -Headers $ownerHeaders `
  -ContentType "application/json" `
  -Body $ownerProfileBody
```

Expected: 201. Repeating POST should return 409.

### List breeds

```powershell
$allBreeds = Invoke-RestMethod -Method GET "$baseUrl/api/breeds"
$dogBreeds = Invoke-RestMethod -Method GET "$baseUrl/api/breeds?species=dog"
$catBreeds = Invoke-RestMethod -Method GET "$baseUrl/api/breeds?species=cat"
$breedId = ($dogBreeds.data | Select-Object -First 1).id
```

Expected seed: 16 total, 8 dog, 8 cat.

### Create pet

```powershell
$petBody = @{
  species = "dog"
  name = "Milo"
  breed_id = $breedId
  gender = "male"
  date_of_birth = "2022-05-10"
  weight_kg = 12.5
  microchip_id = "MC-$suffix"
  color = "Brown"
  distinctive_marks = "White spot on chest"
  is_neutered = $true
} | ConvertTo-Json

$pet = Invoke-RestMethod -Method POST `
  -Uri "$baseUrl/api/pets" `
  -Headers $ownerHeaders `
  -ContentType "application/json" `
  -Body $petBody

$petId = $pet.data.id
$publicPetId = $pet.data.public_pet_id

if ($publicPetId -notmatch '^PNX-PET-[A-Z0-9]{6}$') {
  throw "Unexpected public pet ID: $publicPetId"
}
```

### List my pets

```powershell
Invoke-RestMethod -Method GET `
  -Uri "$baseUrl/api/pets" `
  -Headers $ownerHeaders
```

### Get pet detail

```powershell
Invoke-RestMethod -Method GET `
  -Uri "$baseUrl/api/pets/$petId" `
  -Headers $ownerHeaders
```

### Patch pet

```powershell
$petPatchBody = @{
  name = "Milo Updated"
  weight_kg = 13.2
} | ConvertTo-Json

Invoke-RestMethod -Method PATCH `
  -Uri "$baseUrl/api/pets/$petId" `
  -Headers $ownerHeaders `
  -ContentType "application/json" `
  -Body $petPatchBody
```

## Clinic flow

```powershell
$clinicEmail = "clinic.$suffix@example.com"
```

### Register clinic

```powershell
$clinicRegisterBody = @{
  email = $clinicEmail
  phone = "021234567"
  password = $password
  role = "clinic"
} | ConvertTo-Json

Invoke-RestMethod -Method POST `
  -Uri "$baseUrl/api/auth/register" `
  -ContentType "application/json" `
  -Body $clinicRegisterBody
```

Expected: 201 and returned role `clinic`.

### Login clinic

```powershell
$clinicLoginBody = @{
  email = $clinicEmail
  password = $password
} | ConvertTo-Json

$clinicLogin = Invoke-RestMethod -Method POST `
  -Uri "$baseUrl/api/auth/login" `
  -ContentType "application/json" `
  -Body $clinicLoginBody

$clinicToken = $clinicLogin.data.accessToken
$clinicHeaders = @{ Authorization = "Bearer $clinicToken" }
```

### Create clinic profile

```powershell
$clinicProfileBody = @{
  clinic_name = "Happy Paws Clinic"
  phone_number = "02-123-4567"
  email = "clinic.$suffix@example.com"
  address = "123 Pet Street, Bangkok"
} | ConvertTo-Json

$clinicProfile = Invoke-RestMethod -Method POST `
  -Uri "$baseUrl/api/clinic/profile" `
  -Headers $clinicHeaders `
  -ContentType "application/json" `
  -Body $clinicProfileBody
```

### Get clinic profile

```powershell
Invoke-RestMethod -Method GET `
  -Uri "$baseUrl/api/clinic/profile" `
  -Headers $clinicHeaders
```

### Patch clinic profile

```powershell
$clinicPatchBody = @{ clinic_name = "Happy Paws Bangkok" } | ConvertTo-Json

Invoke-RestMethod -Method PATCH `
  -Uri "$baseUrl/api/clinic/profile" `
  -Headers $clinicHeaders `
  -ContentType "application/json" `
  -Body $clinicPatchBody
```

### Clinic lookup by public pet ID

```powershell
$lookupById = Invoke-RestMethod -Method GET `
  -Uri "$baseUrl/api/clinic/pet-lookup?pet_id=$publicPetId" `
  -Headers $clinicHeaders

$lookupById.data.public_pet_id
```

### Clinic lookup by exact owner phone

```powershell
$lookupByPhone = Invoke-RestMethod -Method GET `
  -Uri "$baseUrl/api/clinic/pet-lookup?owner_phone=0812345678" `
  -Headers $clinicHeaders

$lookupByPhone.data.items
```

Expected: the public ID lookup returns the created pet, and exact phone lookup
contains it. Owner phone is masked. Private owner address and private pet fields
are absent.

## Sprint 8 appointment flow

The commands below reuse `$ownerHeaders`, `$clinicHeaders`, `$petId`,
`$publicPetId`, and `$clinicProfile` created above.

### Owner creates and lists an appointment

```powershell
$scheduledAt = [DateTime]::UtcNow.AddDays(2).ToString("yyyy-MM-ddTHH:mm:ssZ")
$calendarDate = [DateTime]::UtcNow.AddDays(2).ToString("yyyy-MM-dd")

$ownerAppointmentBody = @{
  clinic_profile_id = $clinicProfile.data.id
  pet_id = $petId
  title = "Annual checkup"
  appointment_type = "checkup"
  scheduled_at = $scheduledAt
  duration_minutes = 30
  note = "Bring vaccine card"
} | ConvertTo-Json

$ownerAppointment = Invoke-RestMethod -Method POST `
  -Uri "$baseUrl/api/owner/appointments" `
  -Headers $ownerHeaders `
  -ContentType "application/json" `
  -Body $ownerAppointmentBody

$ownerAppointmentId = $ownerAppointment.data.id

Invoke-RestMethod -Method GET `
  -Uri "$baseUrl/api/owner/appointments?status=requested" `
  -Headers $ownerHeaders

Invoke-RestMethod -Method GET `
  -Uri "$baseUrl/api/owner/appointments/$ownerAppointmentId" `
  -Headers $ownerHeaders
```

Expected: creation returns 201 with status `requested`.

### Clinic creates by public pet ID and uses calendar filters

```powershell
$clinicAppointmentBody = @{
  public_pet_id = $publicPetId
  appointment_type = "vaccination"
  scheduled_at = [DateTime]::UtcNow.AddDays(3).ToString("yyyy-MM-ddTHH:mm:ssZ")
  duration_minutes = 45
} | ConvertTo-Json

$clinicAppointment = Invoke-RestMethod -Method POST `
  -Uri "$baseUrl/api/clinic/appointments" `
  -Headers $clinicHeaders `
  -ContentType "application/json" `
  -Body $clinicAppointmentBody

$clinicAppointmentId = $clinicAppointment.data.id

Invoke-RestMethod -Method GET `
  -Uri "$baseUrl/api/clinic/appointments" `
  -Headers $clinicHeaders

Invoke-RestMethod -Method GET `
  -Uri "$baseUrl/api/clinic/appointments?date=$calendarDate&status=requested&appointment_type=checkup" `
  -Headers $clinicHeaders
```

Expected: clinic creation returns status `scheduled`; the day filter includes
the owner-created appointment.

### Clinic updates status and both roles can cancel scoped appointments

```powershell
$statusBody = @{ status = "checked_in" } | ConvertTo-Json

Invoke-RestMethod -Method PATCH `
  -Uri "$baseUrl/api/clinic/appointments/$ownerAppointmentId/status" `
  -Headers $clinicHeaders `
  -ContentType "application/json" `
  -Body $statusBody

Invoke-RestMethod -Method PATCH `
  -Uri "$baseUrl/api/clinic/appointments/$clinicAppointmentId/cancel" `
  -Headers $clinicHeaders

Invoke-RestMethod -Method PATCH `
  -Uri "$baseUrl/api/owner/appointments/$ownerAppointmentId/cancel" `
  -Headers $ownerHeaders
```

### Sprint 8 negative checks

```powershell
# No token: 401
try {
  Invoke-RestMethod -Method GET "$baseUrl/api/clinic/appointments"
} catch { $_.Exception.Response.StatusCode.value__ }

# Owner on clinic route: 403
try {
  Invoke-RestMethod -Method GET "$baseUrl/api/clinic/appointments" -Headers $ownerHeaders
} catch { $_.Exception.Response.StatusCode.value__ }

# Clinic on owner route: 403
try {
  Invoke-RestMethod -Method GET "$baseUrl/api/owner/appointments" -Headers $clinicHeaders
} catch { $_.Exception.Response.StatusCode.value__ }

# Invalid status: 400
try {
  Invoke-RestMethod -Method PATCH `
    -Uri "$baseUrl/api/clinic/appointments/$ownerAppointmentId/status" `
    -Headers $clinicHeaders `
    -ContentType "application/json" `
    -Body (@{ status = "unknown" } | ConvertTo-Json)
} catch { $_.Exception.Response.StatusCode.value__ }

# Invalid type, duration, or past scheduled_at: 400
$invalidAppointmentBody = @{
  public_pet_id = $publicPetId
  appointment_type = "unsupported"
  scheduled_at = [DateTime]::UtcNow.AddHours(-1).ToString("yyyy-MM-ddTHH:mm:ssZ")
  duration_minutes = 2
} | ConvertTo-Json

try {
  Invoke-RestMethod -Method POST `
    -Uri "$baseUrl/api/clinic/appointments" `
    -Headers $clinicHeaders `
    -ContentType "application/json" `
    -Body $invalidAppointmentBody
} catch { $_.Exception.Response.StatusCode.value__ }
```

## Sprint 9 clinic patient flow

The commands below reuse `$ownerHeaders`, `$clinicHeaders`, `$petId`,
`$publicPetId`, and an existing non-cancelled appointment between the clinic
and pet from the Sprint 8 flow.

### List clinic patients

```powershell
$patients = Invoke-RestMethod -Method GET `
  -Uri "$baseUrl/api/clinic/patients" `
  -Headers $clinicHeaders

$patients.data
```

Expected: the created pet appears once with pet summary, masked owner phone,
and appointment summary.

### Filter clinic patients

```powershell
Invoke-RestMethod -Method GET `
  -Uri "$baseUrl/api/clinic/patients?q=Milo" `
  -Headers $clinicHeaders

Invoke-RestMethod -Method GET `
  -Uri "$baseUrl/api/clinic/patients?species=dog" `
  -Headers $clinicHeaders

Invoke-RestMethod -Method GET `
  -Uri "$baseUrl/api/clinic/patients?limit=20&offset=0&sort=latest_appointment_desc" `
  -Headers $clinicHeaders
```

Expected: filters remain scoped to the authenticated clinic.

### Get clinic patient detail

```powershell
$patientDetail = Invoke-RestMethod -Method GET `
  -Uri "$baseUrl/api/clinic/patients/$petId" `
  -Headers $clinicHeaders

$patientDetail.data
```

Expected: pet detail, owner summary with masked phone, clinic relationship
summary, and recent non-cancelled appointments. The response must not include
`user_id`, `owner_profile_id`, `clinic_profile_id`, owner address, password
data, JWT claims, medical records, visits, reports, payment, or staff schedule
data.

### Sprint 9 negative checks

```powershell
# No token: 401
try {
  Invoke-RestMethod -Method GET "$baseUrl/api/clinic/patients"
} catch { $_.Exception.Response.StatusCode.value__ }

# Owner on clinic patients route: 403
try {
  Invoke-RestMethod -Method GET "$baseUrl/api/clinic/patients" -Headers $ownerHeaders
} catch { $_.Exception.Response.StatusCode.value__ }

# Invalid pet id: 400
try {
  Invoke-RestMethod -Method GET "$baseUrl/api/clinic/patients/not-a-uuid" -Headers $clinicHeaders
} catch { $_.Exception.Response.StatusCode.value__ }

# Invalid species filter: 400
try {
  Invoke-RestMethod -Method GET "$baseUrl/api/clinic/patients?species=bird" -Headers $clinicHeaders
} catch { $_.Exception.Response.StatusCode.value__ }

# Invalid limit: 400
try {
  Invoke-RestMethod -Method GET "$baseUrl/api/clinic/patients?limit=101" -Headers $clinicHeaders
} catch { $_.Exception.Response.StatusCode.value__ }
```

To verify cross-clinic privacy, register/login a second clinic, create its
clinic profile, and call:

```powershell
try {
  Invoke-RestMethod -Method GET "$baseUrl/api/clinic/patients/$petId" -Headers $secondClinicHeaders
} catch { $_.Exception.Response.StatusCode.value__ } # Expected: 404
```

## Sprint 10 medical record flow

The commands below reuse `$ownerHeaders`, `$clinicHeaders`, `$petId`, and
`$clinicAppointmentId` from the owner/clinic/patient/appointment setup flow.
Run this against local PostgreSQL first, then repeat after Render deployment by
changing only `$baseUrl`.

### Create a medical record for a patient

```powershell
$visitAt = [DateTime]::UtcNow.ToString("yyyy-MM-ddTHH:mm:ssZ")
$nextFollowUpAt = [DateTime]::UtcNow.AddDays(7).ToString("yyyy-MM-ddTHH:mm:ssZ")

$medicalRecordBody = @{
  appointmentId = $clinicAppointmentId
  visitAt = $visitAt
  chiefComplaint = "Coughing and reduced appetite"
  clinicalFindings = "Mild fever"
  diagnosis = "Upper respiratory infection"
  treatmentPlan = "Rest and medication"
  medications = "Medication notes as free text"
  followUpInstructions = "Follow up in 7 days"
  nextFollowUpAt = $nextFollowUpAt
  weightKg = 12.5
  temperatureC = 38.2
  notes = "Owner informed"
} | ConvertTo-Json

$medicalRecord = Invoke-RestMethod -Method POST `
  -Uri "$baseUrl/api/clinic/patients/$petId/medical-records" `
  -Headers $clinicHeaders `
  -ContentType "application/json" `
  -Body $medicalRecordBody

$medicalRecordId = $medicalRecord.data.id
```

Expected: 201. Response includes pet, owner, appointment, and minimal
`createdBy` information.

### Create a walk-in or historical record without appointmentId

```powershell
$walkInBody = @{
  visitAt = [DateTime]::UtcNow.AddMinutes(-10).ToString("yyyy-MM-ddTHH:mm:ssZ")
  chiefComplaint = "Walk-in skin irritation"
  diagnosis = "Mild dermatitis"
} | ConvertTo-Json

Invoke-RestMethod -Method POST `
  -Uri "$baseUrl/api/clinic/patients/$petId/medical-records" `
  -Headers $clinicHeaders `
  -ContentType "application/json" `
  -Body $walkInBody
```

Expected: 201 when the pet is already a patient of the clinic.

### List and filter medical records

```powershell
Invoke-RestMethod -Method GET `
  -Uri "$baseUrl/api/clinic/medical-records" `
  -Headers $clinicHeaders

Invoke-RestMethod -Method GET `
  -Uri "$baseUrl/api/clinic/medical-records?pet_id=$petId" `
  -Headers $clinicHeaders

$fromDate = [DateTime]::UtcNow.AddDays(-1).ToString("yyyy-MM-dd")
$toDate = [DateTime]::UtcNow.AddDays(1).ToString("yyyy-MM-dd")
Invoke-RestMethod -Method GET `
  -Uri "$baseUrl/api/clinic/medical-records?from=$fromDate&to=$toDate&page=1&limit=20" `
  -Headers $clinicHeaders
```

Expected: data contains `items` and `pagination`.

### Get and patch a medical record

```powershell
Invoke-RestMethod -Method GET `
  -Uri "$baseUrl/api/clinic/medical-records/$medicalRecordId" `
  -Headers $clinicHeaders

$medicalPatchBody = @{
  chiefComplaint = "Coughing improved"
  diagnosis = "Recovering upper respiratory infection"
  treatmentPlan = "Continue medication"
  weightKg = 12.8
  temperatureC = 37.8
  notes = "Improving"
} | ConvertTo-Json

Invoke-RestMethod -Method PATCH `
  -Uri "$baseUrl/api/clinic/medical-records/$medicalRecordId" `
  -Headers $clinicHeaders `
  -ContentType "application/json" `
  -Body $medicalPatchBody
```

### Sprint 10 negative checks

```powershell
# No token: 401
try {
  Invoke-RestMethod -Method GET "$baseUrl/api/clinic/medical-records"
} catch { $_.Exception.Response.StatusCode.value__ }

# Owner on clinic medical records route: 403
try {
  Invoke-RestMethod -Method GET "$baseUrl/api/clinic/medical-records" -Headers $ownerHeaders
} catch { $_.Exception.Response.StatusCode.value__ }

# Invalid record id: 400
try {
  Invoke-RestMethod -Method GET "$baseUrl/api/clinic/medical-records/not-a-uuid" -Headers $clinicHeaders
} catch { $_.Exception.Response.StatusCode.value__ }

# Invalid pet id on create: 400
try {
  Invoke-RestMethod -Method POST `
    -Uri "$baseUrl/api/clinic/patients/not-a-uuid/medical-records" `
    -Headers $clinicHeaders `
    -ContentType "application/json" `
    -Body $medicalRecordBody
} catch { $_.Exception.Response.StatusCode.value__ }

# Duplicate appointment medical record: 409
try {
  Invoke-RestMethod -Method POST `
    -Uri "$baseUrl/api/clinic/patients/$petId/medical-records" `
    -Headers $clinicHeaders `
    -ContentType "application/json" `
    -Body $medicalRecordBody
} catch { $_.Exception.Response.StatusCode.value__ }

# Invalid vitals/date: 400
$invalidMedicalRecordBody = @{
  visitAt = $visitAt
  chiefComplaint = "Invalid vitals"
  weightKg = 0
  temperatureC = -1
  nextFollowUpAt = [DateTime]::UtcNow.AddDays(-1).ToString("yyyy-MM-ddTHH:mm:ssZ")
} | ConvertTo-Json
try {
  Invoke-RestMethod -Method POST `
    -Uri "$baseUrl/api/clinic/patients/$petId/medical-records" `
    -Headers $clinicHeaders `
    -ContentType "application/json" `
    -Body $invalidMedicalRecordBody
} catch { $_.Exception.Response.StatusCode.value__ }
```

## Negative tests

### No token returns 401

```powershell
try {
  Invoke-RestMethod -Method GET "$baseUrl/api/pets"
} catch {
  $_.Exception.Response.StatusCode.value__ # 401
}
```

### Owner calling clinic profile returns 403

```powershell
try {
  Invoke-RestMethod -Method GET "$baseUrl/api/clinic/profile" -Headers $ownerHeaders
} catch {
  $_.Exception.Response.StatusCode.value__ # 403
}
```

### Clinic pet lookup authorization and validation

```powershell
try {
  Invoke-RestMethod -Method GET `
    -Uri "$baseUrl/api/clinic/pet-lookup?pet_id=$publicPetId"
} catch {
  $_.Exception.Response.StatusCode.value__ # 401
}

try {
  Invoke-RestMethod -Method GET `
    -Uri "$baseUrl/api/clinic/pet-lookup?pet_id=$publicPetId" `
    -Headers $ownerHeaders
} catch {
  $_.Exception.Response.StatusCode.value__ # 403
}

try {
  Invoke-RestMethod -Method GET `
    -Uri "$baseUrl/api/clinic/pet-lookup" `
    -Headers $clinicHeaders
} catch {
  $_.Exception.Response.StatusCode.value__ # 400
}

try {
  Invoke-RestMethod -Method GET `
    -Uri "$baseUrl/api/clinic/pet-lookup?pet_id=PNX-PET-ZZZZZZ" `
    -Headers $clinicHeaders
} catch {
  $_.Exception.Response.StatusCode.value__ # 404
}
```

### Clinic calling owner profile and pet endpoints returns 403

```powershell
try {
  Invoke-RestMethod -Method GET "$baseUrl/api/owner/profile" -Headers $clinicHeaders
} catch {
  $_.Exception.Response.StatusCode.value__ # 403
}

try {
  Invoke-RestMethod -Method GET "$baseUrl/api/pets" -Headers $clinicHeaders
} catch {
  $_.Exception.Response.StatusCode.value__ # 403
}
```

### Invalid species returns 400

```powershell
try {
  Invoke-RestMethod -Method GET "$baseUrl/api/breeds?species=bird"
} catch {
  $_.Exception.Response.StatusCode.value__ # 400
}
```

### Empty PATCH returns 400

```powershell
try {
  Invoke-RestMethod -Method PATCH `
    -Uri "$baseUrl/api/pets/$petId" `
    -Headers $ownerHeaders `
    -ContentType "application/json" `
    -Body '{}'
} catch {
  $_.Exception.Response.StatusCode.value__ # 400
}
```

## Render verification

After deployment:

1. Check `/health` and `/health/db`.
2. Repeat owner and clinic registration/login with unique emails.
3. Repeat profile, pet, and clinic pet lookup flows.
4. Verify wrong-role responses.
5. Inspect Render logs for migration errors without exposing secrets.
