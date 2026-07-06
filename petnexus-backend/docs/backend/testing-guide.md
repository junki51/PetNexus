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

Invoke-RestMethod -Method POST `
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
