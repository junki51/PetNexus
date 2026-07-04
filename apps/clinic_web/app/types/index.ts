// ============================================================
// PetNexus Clinic Platform — Shared TypeScript Types
// ============================================================

// ── User & Auth ──────────────────────────────────────────────

export type UserRole = "owner" | "clinic_staff" | "admin";

export interface User {
  id: string;
  email: string;
  phone?: string;
  role: UserRole;
  createdAt: string;
}

export interface ClinicStaff {
  id: string;
  userId: string;
  clinicId: string;
  name: string;
  role: "vet" | "assistant" | "clinic_admin";
  licenseNo?: string;
  avatarUrl?: string;
}

export interface Clinic {
  id: string;
  name: string;
  address?: string;
  phone?: string;
  email?: string;
}

// ── Pet & Patient ─────────────────────────────────────────────

export type PetSpecies = "dog" | "cat";
export type PetGender = "male" | "female";
export type AuthorizationStatus =
  | "pending"
  | "approved"
  | "rejected"
  | "revoked";

export interface Pet {
  id: string;
  petNexusId: string;
  name: string;
  species: PetSpecies;
  breed: string;
  gender: PetGender;
  birthDate?: string;
  weightKg?: number;
  allergyNote?: string;
  chronicDiseaseNote?: string;
  photoUrl?: string;
  ownerId: string;
  ownerName: string;
  ownerPhone: string;
}

export interface Patient extends Pet {
  status: AuthorizationStatus;
  lastVisit?: string;
}

// ── Appointment / Schedule ────────────────────────────────────

export type VisitType = "Vaccination" | "Consultation" | "Follow-up" | "Emergency" | "Grooming";

export interface Appointment {
  id: string;
  time: string;
  petName: string;
  petSpecies: PetSpecies;
  ownerName: string;
  type: VisitType;
  status: "scheduled" | "checked-in" | "in-progress" | "done" | "cancelled";
  petPhotoUrl?: string;
}

// ── Medical Record ────────────────────────────────────────────

export interface Medication {
  name: string;
  dosage: string;
  instructions: string;
}

export interface MedicalRecord {
  id: string;
  petId: string;
  clinicId: string;
  vetId: string;
  visitDate: string;
  visitType: VisitType;
  chiefComplaint?: string;
  diagnosis?: string;
  treatment?: string;
  medications: Medication[];
  followUpDate?: string;
  followUpNote?: string;
  note?: string;
}

// ── Activity / Notification ───────────────────────────────────

export type ActivityType =
  | "check-in"
  | "new-appointment"
  | "record-updated"
  | "visit-completed"
  | "access-request";

export interface Activity {
  id: string;
  type: ActivityType;
  description: string;
  time: string;
  petName?: string;
  ownerName?: string;
}

// ── Dashboard Stats ───────────────────────────────────────────

export interface DashboardStats {
  todayPatients: number;
  todayPatientsChange: string;
  pendingRequests: number;
  upcomingAppointments: number;
  nextAppointmentTime?: string;
  checkedInNow: number;
}

// ── QR Session ────────────────────────────────────────────────

export interface QRScanResult {
  pet: Pet;
  authorizationStatus: AuthorizationStatus;
  permissions: string[];
}

// ── Table / Pagination ────────────────────────────────────────

export interface PaginationMeta {
  page: number;
  limit: number;
  total: number;
  totalPages: number;
}

export interface TableColumn<T> {
  key: keyof T | string;
  label: string;
  width?: string;
  align?: "left" | "center" | "right";
  render?: (value: unknown, row: T) => React.ReactNode;
}

// ── Filter Options ────────────────────────────────────────────

export interface SelectOption {
  value: string;
  label: string;
}
