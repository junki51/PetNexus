# PetNexus API Plan

## 1. API Overview

PetNexus uses REST API.

Base URL for local development:

```txt
http://localhost:8080
```

API prefix:

```txt
/api
```

Example:

```txt
GET http://localhost:8080/api/pets
```

Health route does not need `/api` prefix:

```txt
GET http://localhost:8080/health
```

---

## 2. API Design Rules

### 2.1 Use JSON

All request and response bodies should use JSON.

Header:

```txt
Content-Type: application/json
```

---

### 2.2 Use JWT Auth

Protected routes require:

```txt
Authorization: Bearer <access_token>
```

Public routes:

```txt
GET  /health
POST /api/auth/register
POST /api/auth/login
```

All other routes should require authentication unless explicitly stated.

---

### 2.3 Use Consistent Response Format

Success response:

```json
{
  "success": true,
  "message": "Action completed successfully",
  "data": {}
}
```

Error response:

```json
{
  "success": false,
  "message": "Something went wrong",
  "error": {
    "code": "ERROR_CODE",
    "details": "Readable error detail"
  }
}
```

List response:

```json
{
  "success": true,
  "message": "Data fetched successfully",
  "data": [],
  "meta": {
    "page": 1,
    "limit": 20,
    "total": 100
  }
}
```

---

## 3. HTTP Status Code Rules

Use these status codes:

```txt
200 OK                  = successful GET / PATCH
201 Created             = successful create
400 Bad Request         = invalid request body or query
401 Unauthorized        = no token or invalid token
403 Forbidden           = logged in but not allowed
404 Not Found           = resource not found
409 Conflict            = duplicate or conflicting action
422 Unprocessable Entity = validation error
500 Internal Server Error = unexpected server error
```

---

## 4. Core Roles

```txt
owner
clinic_staff
admin
```

---

## 5. Core Permissions

Authorization permissions:

```txt
view_profile
view_history
create_visit
```

Clinic must have approved authorization before accessing full pet data.

---

## 6. Public APIs

## 6.1 Health Check

```txt
GET /health
```

Auth required:

```txt
No
```

Response:

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

---

## 7. Auth APIs

## 7.1 Register

```txt
POST /api/auth/register
```

Auth required:

```txt
No
```

Request body:

```json
{
  "email": "owner@example.com",
  "phone": "0812345678",
  "password": "password123",
  "role": "owner"
}
```

Allowed roles for MVP:

```txt
owner
clinic_staff
```

Response:

```json
{
  "success": true,
  "message": "Registered successfully",
  "data": {
    "user": {
      "id": "user_id",
      "email": "owner@example.com",
      "phone": "0812345678",
      "role": "owner"
    },
    "accessToken": "jwt_token"
  }
}
```

Possible errors:

```txt
400 INVALID_REQUEST
409 EMAIL_ALREADY_EXISTS
422 VALIDATION_ERROR
```

---

## 7.2 Login

```txt
POST /api/auth/login
```

Auth required:

```txt
No
```

Request body:

```json
{
  "email": "owner@example.com",
  "password": "password123"
}
```

Response:

```json
{
  "success": true,
  "message": "Logged in successfully",
  "data": {
    "user": {
      "id": "user_id",
      "email": "owner@example.com",
      "role": "owner"
    },
    "accessToken": "jwt_token"
  }
}
```

Possible errors:

```txt
401 INVALID_CREDENTIALS
422 VALIDATION_ERROR
```

---

## 7.3 Get Current User

```txt
GET /api/me
```

Auth required:

```txt
Yes
```

Response:

```json
{
  "success": true,
  "message": "Current user fetched successfully",
  "data": {
    "user": {
      "id": "user_id",
      "email": "owner@example.com",
      "phone": "0812345678",
      "role": "owner"
    }
  }
}
```

---

## 8. Owner Profile APIs

## 8.1 Create Owner Profile

```txt
POST /api/owner/profile
```

Auth required:

```txt
Yes
```

Allowed role:

```txt
owner
```

Request body:

```json
{
  "fullName": "Sunny Test",
  "nickname": "Sunny",
  "phone": "0812345678",
  "email": "owner@example.com",
  "address": "Bangkok, Thailand",
  "emergencyContactName": "Parent",
  "emergencyContactPhone": "0899999999"
}
```

Response:

```json
{
  "success": true,
  "message": "Owner profile created successfully",
  "data": {
    "ownerProfile": {
      "id": "owner_profile_id",
      "userId": "user_id",
      "fullName": "Sunny Test",
      "nickname": "Sunny",
      "phone": "0812345678",
      "email": "owner@example.com",
      "address": "Bangkok, Thailand",
      "emergencyContactName": "Parent",
      "emergencyContactPhone": "0899999999"
    }
  }
}
```

Possible errors:

```txt
403 FORBIDDEN_ROLE
409 OWNER_PROFILE_ALREADY_EXISTS
422 VALIDATION_ERROR
```

---

## 8.2 Get Owner Profile

```txt
GET /api/owner/profile
```

Auth required:

```txt
Yes
```

Allowed role:

```txt
owner
```

Response:

```json
{
  "success": true,
  "message": "Owner profile fetched successfully",
  "data": {
    "ownerProfile": {
      "id": "owner_profile_id",
      "fullName": "Sunny Test",
      "nickname": "Sunny",
      "phone": "0812345678",
      "email": "owner@example.com",
      "address": "Bangkok, Thailand",
      "emergencyContactName": "Parent",
      "emergencyContactPhone": "0899999999"
    }
  }
}
```

---

## 8.3 Update Owner Profile

```txt
PATCH /api/owner/profile
```

Auth required:

```txt
Yes
```

Allowed role:

```txt
owner
```

Request body:

```json
{
  "nickname": "Sunny",
  "address": "Bangkok, Thailand"
}
```

Response:

```json
{
  "success": true,
  "message": "Owner profile updated successfully",
  "data": {
    "ownerProfile": {
      "id": "owner_profile_id",
      "nickname": "Sunny",
      "address": "Bangkok, Thailand"
    }
  }
}
```

---

## 9. Breed APIs

## 9.1 Get Breeds

```txt
GET /api/breeds
GET /api/breeds?species=dog
GET /api/breeds?species=cat
```

Auth required:

```txt
Yes
```

Response:

```json
{
  "success": true,
  "message": "Breeds fetched successfully",
  "data": {
    "breeds": [
      {
        "id": "breed_id",
        "species": "dog",
        "nameTh": "โกลเด้น รีทรีฟเวอร์",
        "nameEn": "Golden Retriever",
        "imageUrl": "https://example.com/golden.png"
      }
    ]
  }
}
```

Initial dog breeds:

```txt
Golden Retriever
Pomeranian
French Bulldog
Poodle
```

Initial cat breeds:

```txt
British Shorthair
Persian
Scottish Fold
Siamese
```

Also include:

```txt
Not sure / Select later
```

---

## 10. Pet APIs

## 10.1 Create Pet

```txt
POST /api/pets
```

Auth required:

```txt
Yes
```

Allowed role:

```txt
owner
```

Request body:

```json
{
  "name": "Milo",
  "species": "dog",
  "breedId": "breed_id",
  "gender": "male",
  "birthDate": "2021-05-12",
  "weightKg": 12.5,
  "colorNote": "Brown and white",
  "allergyNote": "Chicken allergy",
  "chronicDiseaseNote": "",
  "photoUrl": "https://example.com/milo.png"
}
```

Response:

```json
{
  "success": true,
  "message": "Pet created successfully",
  "data": {
    "pet": {
      "id": "pet_id",
      "ownerId": "owner_profile_id",
      "petNexusId": "PNX-8F3K2A",
      "name": "Milo",
      "species": "dog",
      "breed": {
        "id": "breed_id",
        "nameEn": "Golden Retriever"
      },
      "gender": "male",
      "birthDate": "2021-05-12",
      "weightKg": 12.5,
      "photoUrl": "https://example.com/milo.png"
    }
  }
}
```

Possible errors:

```txt
403 FORBIDDEN_ROLE
404 OWNER_PROFILE_NOT_FOUND
404 BREED_NOT_FOUND
422 VALIDATION_ERROR
```

---

## 10.2 Get Pet List

```txt
GET /api/pets
```

Auth required:

```txt
Yes
```

Allowed role:

```txt
owner
```

Response:

```json
{
  "success": true,
  "message": "Pets fetched successfully",
  "data": {
    "pets": [
      {
        "id": "pet_id",
        "petNexusId": "PNX-8F3K2A",
        "name": "Milo",
        "species": "dog",
        "breedName": "Golden Retriever",
        "gender": "male",
        "weightKg": 12.5,
        "photoUrl": "https://example.com/milo.png"
      }
    ]
  }
}
```

---

## 10.3 Get Pet Detail

```txt
GET /api/pets/:petId
```

Auth required:

```txt
Yes
```

Allowed role:

```txt
owner
```

Rule:

```txt
Owner can only view own pet.
```

Response:

```json
{
  "success": true,
  "message": "Pet fetched successfully",
  "data": {
    "pet": {
      "id": "pet_id",
      "petNexusId": "PNX-8F3K2A",
      "name": "Milo",
      "species": "dog",
      "breedName": "Golden Retriever",
      "gender": "male",
      "birthDate": "2021-05-12",
      "weightKg": 12.5,
      "colorNote": "Brown and white",
      "allergyNote": "Chicken allergy",
      "chronicDiseaseNote": "",
      "photoUrl": "https://example.com/milo.png"
    }
  }
}
```

---

## 10.4 Update Pet

```txt
PATCH /api/pets/:petId
```

Auth required:

```txt
Yes
```

Allowed role:

```txt
owner
```

Rule:

```txt
Owner can only update own pet basic info.
Owner cannot edit clinic_verified visits through this endpoint.
```

Request body:

```json
{
  "weightKg": 13.0,
  "allergyNote": "Chicken allergy"
}
```

Response:

```json
{
  "success": true,
  "message": "Pet updated successfully",
  "data": {
    "pet": {
      "id": "pet_id",
      "weightKg": 13.0,
      "allergyNote": "Chicken allergy"
    }
  }
}
```

---

## 10.5 Get Pet Passport

```txt
GET /api/pets/:petId/passport
```

Auth required:

```txt
Yes
```

Allowed role:

```txt
owner
```

Response:

```json
{
  "success": true,
  "message": "Pet passport fetched successfully",
  "data": {
    "passport": {
      "pet": {
        "id": "pet_id",
        "petNexusId": "PNX-8F3K2A",
        "name": "Milo",
        "species": "dog",
        "breedName": "Golden Retriever",
        "gender": "male",
        "birthDate": "2021-05-12",
        "ageText": "3 years old",
        "weightKg": 13.0,
        "photoUrl": "https://example.com/milo.png"
      },
      "alerts": {
        "allergyNote": "Chicken allergy",
        "chronicDiseaseNote": ""
      },
      "owner": {
        "nickname": "Sunny",
        "emergencyContactName": "Parent",
        "emergencyContactPhone": "0899999999"
      },
      "verification": {
        "hasClinicVerifiedRecords": true,
        "latestVerifiedVisitDate": "2026-06-10"
      }
    }
  }
}
```

---

## 11. QR Session APIs

## 11.1 Create QR Session

```txt
POST /api/pets/:petId/qr-session
```

Auth required:

```txt
Yes
```

Allowed role:

```txt
owner
```

Request body:

```json
{
  "purpose": "clinic_checkin"
}
```

Response:

```json
{
  "success": true,
  "message": "QR session created successfully",
  "data": {
    "qrSession": {
      "id": "qr_session_id",
      "petId": "pet_id",
      "token": "secure_random_token",
      "purpose": "clinic_checkin",
      "expiresAt": "2026-06-30T15:30:00Z"
    },
    "qrPayload": {
      "type": "petnexus_qr_session",
      "token": "secure_random_token"
    }
  }
}
```

Important:

```txt
QR payload must contain token only.
Do not include full pet data in QR.
```

Possible errors:

```txt
403 NOT_PET_OWNER
404 PET_NOT_FOUND
422 VALIDATION_ERROR
```

---

## 11.2 Clinic Scan QR

```txt
POST /api/clinic/scan-qr
```

Auth required:

```txt
Yes
```

Allowed role:

```txt
clinic_staff
```

Request body:

```json
{
  "token": "secure_random_token"
}
```

Response:

```json
{
  "success": true,
  "message": "QR token validated successfully",
  "data": {
    "petPreview": {
      "petId": "pet_id",
      "petNexusId": "PNX-8F3K2A",
      "name": "Milo",
      "species": "dog",
      "breedName": "Golden Retriever",
      "photoUrl": "https://example.com/milo.png"
    },
    "access": {
      "hasApprovedAccess": false,
      "hasPendingRequest": false
    }
  }
}
```

Do not return:

```txt
allergyNote
chronicDiseaseNote
timeline
diagnosis
treatment
medication
owner address
```

Possible errors:

```txt
400 INVALID_TOKEN
401 UNAUTHORIZED
403 FORBIDDEN_ROLE
404 QR_SESSION_NOT_FOUND
410 QR_SESSION_EXPIRED
```

---

## 12. Authorization APIs

## 12.1 Clinic Request Access

```txt
POST /api/clinic/access-request
```

Auth required:

```txt
Yes
```

Allowed role:

```txt
clinic_staff
```

Request body:

```json
{
  "petId": "pet_id",
  "permissions": ["view_profile", "view_history", "create_visit"],
  "reason": "Clinic check-in visit"
}
```

Response:

```json
{
  "success": true,
  "message": "Access request created successfully",
  "data": {
    "authorization": {
      "id": "authorization_id",
      "petId": "pet_id",
      "clinicId": "clinic_id",
      "status": "pending",
      "permissions": ["view_profile", "view_history", "create_visit"],
      "createdAt": "2026-06-30T14:00:00Z"
    }
  }
}
```

Possible errors:

```txt
403 FORBIDDEN_ROLE
404 PET_NOT_FOUND
409 ACCESS_REQUEST_ALREADY_PENDING
409 ACCESS_ALREADY_APPROVED
422 VALIDATION_ERROR
```

---

## 12.2 Owner Get Access Requests

```txt
GET /api/owner/access-requests
```

Auth required:

```txt
Yes
```

Allowed role:

```txt
owner
```

Response:

```json
{
  "success": true,
  "message": "Access requests fetched successfully",
  "data": {
    "requests": [
      {
        "id": "authorization_id",
        "status": "pending",
        "permissions": ["view_profile", "view_history", "create_visit"],
        "pet": {
          "id": "pet_id",
          "name": "Milo",
          "photoUrl": "https://example.com/milo.png"
        },
        "clinic": {
          "id": "clinic_id",
          "name": "Happy Pet Clinic",
          "phone": "021234567"
        },
        "createdAt": "2026-06-30T14:00:00Z"
      }
    ]
  }
}
```

---

## 12.3 Owner Approve Access

```txt
POST /api/authorizations/:id/approve
```

Auth required:

```txt
Yes
```

Allowed role:

```txt
owner
```

Request body:

```json
{
  "expiresAt": "2026-07-30T23:59:59Z"
}
```

Response:

```json
{
  "success": true,
  "message": "Access approved successfully",
  "data": {
    "authorization": {
      "id": "authorization_id",
      "status": "approved",
      "permissions": ["view_profile", "view_history", "create_visit"],
      "expiresAt": "2026-07-30T23:59:59Z"
    }
  }
}
```

Possible errors:

```txt
403 NOT_PET_OWNER
404 AUTHORIZATION_NOT_FOUND
409 AUTHORIZATION_ALREADY_DECIDED
```

---

## 12.4 Owner Reject Access

```txt
POST /api/authorizations/:id/reject
```

Auth required:

```txt
Yes
```

Allowed role:

```txt
owner
```

Response:

```json
{
  "success": true,
  "message": "Access rejected successfully",
  "data": {
    "authorization": {
      "id": "authorization_id",
      "status": "rejected"
    }
  }
}
```

---

## 12.5 Owner Revoke Access

```txt
POST /api/authorizations/:id/revoke
```

Auth required:

```txt
Yes
```

Allowed role:

```txt
owner
```

Response:

```json
{
  "success": true,
  "message": "Access revoked successfully",
  "data": {
    "authorization": {
      "id": "authorization_id",
      "status": "revoked"
    }
  }
}
```

---

## 13. Clinic APIs

## 13.1 Clinic Register

```txt
POST /api/clinic/register
```

Auth required:

```txt
No
```

Request body:

```json
{
  "clinicName": "Happy Pet Clinic",
  "clinicAddress": "Bangkok, Thailand",
  "clinicPhone": "021234567",
  "clinicEmail": "clinic@example.com",
  "staffEmail": "vet@example.com",
  "staffPhone": "0811111111",
  "password": "password123",
  "staffRole": "clinic_admin",
  "licenseNo": "VET-12345"
}
```

Response:

```json
{
  "success": true,
  "message": "Clinic registered successfully",
  "data": {
    "clinic": {
      "id": "clinic_id",
      "name": "Happy Pet Clinic",
      "verifiedStatus": "pending"
    },
    "staff": {
      "id": "clinic_staff_id",
      "role": "clinic_admin",
      "licenseNo": "VET-12345"
    },
    "accessToken": "jwt_token"
  }
}
```

MVP note:

```txt
Clinic verification can stay pending in MVP.
Do not build full admin verification yet.
```

---

## 13.2 Clinic Me

```txt
GET /api/clinic/me
```

Auth required:

```txt
Yes
```

Allowed role:

```txt
clinic_staff
```

Response:

```json
{
  "success": true,
  "message": "Clinic staff fetched successfully",
  "data": {
    "staff": {
      "id": "clinic_staff_id",
      "role": "clinic_admin",
      "licenseNo": "VET-12345"
    },
    "clinic": {
      "id": "clinic_id",
      "name": "Happy Pet Clinic",
      "verifiedStatus": "pending"
    }
  }
}
```

---

## 13.3 Clinic Dashboard

```txt
GET /api/clinic/dashboard
```

Auth required:

```txt
Yes
```

Allowed role:

```txt
clinic_staff
```

Response:

```json
{
  "success": true,
  "message": "Clinic dashboard fetched successfully",
  "data": {
    "stats": {
      "approvedPatients": 12,
      "pendingAccessRequests": 3,
      "visitsToday": 4
    },
    "recentActivity": [
      {
        "type": "clinic_created_visit",
        "message": "Visit created for Milo",
        "createdAt": "2026-06-30T14:00:00Z"
      }
    ]
  }
}
```

MVP note:

```txt
Dashboard can be simple.
Do not build advanced analytics yet.
```

---

## 13.4 Clinic Patients

```txt
GET /api/clinic/patients
```

Auth required:

```txt
Yes
```

Allowed role:

```txt
clinic_staff
```

Rule:

```txt
Return only pets with approved authorization for this clinic.
```

Response:

```json
{
  "success": true,
  "message": "Clinic patients fetched successfully",
  "data": {
    "patients": [
      {
        "petId": "pet_id",
        "petNexusId": "PNX-8F3K2A",
        "name": "Milo",
        "species": "dog",
        "breedName": "Golden Retriever",
        "photoUrl": "https://example.com/milo.png",
        "authorization": {
          "id": "authorization_id",
          "status": "approved",
          "permissions": ["view_profile", "view_history", "create_visit"],
          "expiresAt": "2026-07-30T23:59:59Z"
        }
      }
    ]
  }
}
```

---

## 13.5 Clinic Get Pet Record

```txt
GET /api/clinic/pets/:petId
```

Auth required:

```txt
Yes
```

Allowed role:

```txt
clinic_staff
```

Required permission:

```txt
view_profile
```

Rule:

```txt
Clinic can view full pet record only after owner approval.
```

Response:

```json
{
  "success": true,
  "message": "Pet record fetched successfully",
  "data": {
    "pet": {
      "id": "pet_id",
      "petNexusId": "PNX-8F3K2A",
      "name": "Milo",
      "species": "dog",
      "breedName": "Golden Retriever",
      "gender": "male",
      "birthDate": "2021-05-12",
      "weightKg": 13.0,
      "colorNote": "Brown and white",
      "allergyNote": "Chicken allergy",
      "chronicDiseaseNote": "",
      "photoUrl": "https://example.com/milo.png"
    },
    "owner": {
      "nickname": "Sunny",
      "emergencyContactName": "Parent",
      "emergencyContactPhone": "0899999999"
    }
  }
}
```

Possible errors:

```txt
403 CLINIC_ACCESS_NOT_APPROVED
403 MISSING_PERMISSION
404 PET_NOT_FOUND
```

---

## 14. Visit APIs

## 14.1 Clinic Create Visit

```txt
POST /api/clinic/pets/:petId/visits
```

Auth required:

```txt
Yes
```

Allowed role:

```txt
clinic_staff
```

Required permission:

```txt
create_visit
```

Request body:

```json
{
  "visitDate": "2026-06-30",
  "chiefComplaint": "Annual checkup",
  "diagnosis": "Healthy",
  "treatment": "General physical examination",
  "medication": "",
  "followUpDate": "2026-12-30",
  "note": "Recommend next checkup in 6 months"
}
```

Response:

```json
{
  "success": true,
  "message": "Visit created successfully",
  "data": {
    "visit": {
      "id": "visit_id",
      "petId": "pet_id",
      "clinicId": "clinic_id",
      "vetId": "clinic_staff_id",
      "visitDate": "2026-06-30",
      "chiefComplaint": "Annual checkup",
      "diagnosis": "Healthy",
      "treatment": "General physical examination",
      "medication": "",
      "followUpDate": "2026-12-30",
      "note": "Recommend next checkup in 6 months",
      "verificationStatus": "clinic_verified",
      "createdAt": "2026-06-30T14:00:00Z"
    }
  }
}
```

Possible errors:

```txt
403 CLINIC_ACCESS_NOT_APPROVED
403 MISSING_PERMISSION_CREATE_VISIT
404 PET_NOT_FOUND
422 VALIDATION_ERROR
```

Side effects:

```txt
Create notification for owner.
Create audit log: clinic_created_visit.
Timeline should include this visit.
```

---

## 15. Timeline APIs

## 15.1 Owner Get Pet Timeline

```txt
GET /api/pets/:petId/timeline
```

Auth required:

```txt
Yes
```

Allowed role:

```txt
owner
```

Rule:

```txt
Owner can only view own pet timeline.
```

Response:

```json
{
  "success": true,
  "message": "Timeline fetched successfully",
  "data": {
    "timeline": [
      {
        "id": "visit_id",
        "type": "visit",
        "date": "2026-06-30",
        "title": "Annual checkup",
        "clinicName": "Happy Pet Clinic",
        "diagnosis": "Healthy",
        "treatment": "General physical examination",
        "medication": "",
        "followUpDate": "2026-12-30",
        "verificationStatus": "clinic_verified",
        "createdAt": "2026-06-30T14:00:00Z"
      }
    ]
  }
}
```

---

## 15.2 Clinic Get Pet Timeline

```txt
GET /api/clinic/pets/:petId/timeline
```

Auth required:

```txt
Yes
```

Allowed role:

```txt
clinic_staff
```

Required permission:

```txt
view_history
```

Response:

```json
{
  "success": true,
  "message": "Timeline fetched successfully",
  "data": {
    "timeline": [
      {
        "id": "visit_id",
        "type": "visit",
        "date": "2026-06-30",
        "clinicName": "Happy Pet Clinic",
        "chiefComplaint": "Annual checkup",
        "diagnosis": "Healthy",
        "treatment": "General physical examination",
        "medication": "",
        "verificationStatus": "clinic_verified"
      }
    ]
  }
}
```

Possible errors:

```txt
403 CLINIC_ACCESS_NOT_APPROVED
403 MISSING_PERMISSION_VIEW_HISTORY
404 PET_NOT_FOUND
```

---

## 16. Notification APIs

## 16.1 Get Notifications

```txt
GET /api/notifications
```

Auth required:

```txt
Yes
```

Response:

```json
{
  "success": true,
  "message": "Notifications fetched successfully",
  "data": {
    "notifications": [
      {
        "id": "notification_id",
        "title": "New access request",
        "message": "Happy Pet Clinic requested access to Milo's record.",
        "type": "access_request",
        "isRead": false,
        "createdAt": "2026-06-30T14:00:00Z"
      }
    ]
  }
}
```

---

## 16.2 Mark Notification as Read

```txt
PATCH /api/notifications/:id/read
```

Auth required:

```txt
Yes
```

Response:

```json
{
  "success": true,
  "message": "Notification marked as read",
  "data": {
    "notification": {
      "id": "notification_id",
      "isRead": true
    }
  }
}
```

---

## 17. Audit Log Rules

Audit logs should be created internally by services.

Do not expose audit log API in MVP unless needed for demo.

Important actions:

```txt
clinic_scanned_qr
clinic_requested_access
owner_approved_access
owner_rejected_access
owner_revoked_access
clinic_viewed_pet
clinic_created_visit
```

Audit log fields:

```txt
id
actorUserId
action
targetType
targetId
metadata
createdAt
```

Example metadata:

```json
{
  "clinicId": "clinic_id",
  "petId": "pet_id",
  "authorizationId": "authorization_id"
}
```

---

## 18. Permission Helper Plan

Backend should have helper functions like:

```txt
GetCurrentUser(ctx)
RequireAuth()
RequireRole(role)
CanOwnerAccessPet(ownerUserId, petId)
CanClinicAccessPet(clinicId, petId, permission)
```

Important checks:

```txt
Owner must own pet.
Clinic must have approved authorization.
Authorization must not be revoked.
Authorization must not be expired.
Permission must include required permission.
```

---

## 19. API Implementation Order

Implement APIs in this order:

```txt
1. GET  /health

2. POST /api/auth/register
3. POST /api/auth/login
4. GET  /api/me

5. POST /api/owner/profile
6. GET  /api/owner/profile

7. GET  /api/breeds

8. POST /api/pets
9. GET  /api/pets
10. GET /api/pets/:petId
11. GET /api/pets/:petId/passport

12. POST /api/pets/:petId/qr-session
13. POST /api/clinic/scan-qr

14. POST /api/clinic/access-request
15. GET  /api/owner/access-requests

16. POST /api/authorizations/:id/approve
17. POST /api/authorizations/:id/reject
18. POST /api/authorizations/:id/revoke

19. GET /api/clinic/patients
20. GET /api/clinic/pets/:petId

21. POST /api/clinic/pets/:petId/visits

22. GET /api/pets/:petId/timeline
23. GET /api/clinic/pets/:petId/timeline

24. GET   /api/notifications
25. PATCH /api/notifications/:id/read
```

Do not implement later APIs before this core flow works.

---

## 20. MVP API Test Flow

Use Postman or Thunder Client.

### Step 1: Owner Register

```txt
POST /api/auth/register
```

Save owner token.

---

### Step 2: Owner Create Profile

```txt
POST /api/owner/profile
```

Use owner token.

---

### Step 3: Owner Create Pet

```txt
POST /api/pets
```

Use owner token.

Save petId.

---

### Step 4: Owner Create QR Session

```txt
POST /api/pets/:petId/qr-session
```

Use owner token.

Save QR token.

---

### Step 5: Clinic Register

```txt
POST /api/clinic/register
```

Save clinic token.

---

### Step 6: Clinic Scan QR

```txt
POST /api/clinic/scan-qr
```

Use clinic token.

Expected result:

```txt
Clinic sees pet preview only.
```

---

### Step 7: Clinic Request Access

```txt
POST /api/clinic/access-request
```

Use clinic token.

Save authorizationId.

---

### Step 8: Owner Approves Access

```txt
POST /api/authorizations/:id/approve
```

Use owner token.

---

### Step 9: Clinic Views Full Pet Record

```txt
GET /api/clinic/pets/:petId
```

Use clinic token.

Expected result:

```txt
Clinic sees full approved pet record.
```

---

### Step 10: Clinic Creates Visit

```txt
POST /api/clinic/pets/:petId/visits
```

Use clinic token.

Expected result:

```txt
Visit is created with clinic_verified status.
```

---

### Step 11: Owner Views Timeline

```txt
GET /api/pets/:petId/timeline
```

Use owner token.

Expected result:

```txt
Owner sees clinic verified visit in timeline.
```

---

## 21. API MVP Definition

API MVP is complete when the full test flow works:

```txt
Owner register
Owner create profile
Owner create pet
Owner create QR session
Clinic register
Clinic scan QR
Clinic request access
Owner approve
Clinic view pet
Clinic create visit
Owner view timeline
```

This is enough for the first real demo.
