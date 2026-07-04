// ============================================================
// PetNexus Clinic Platform — Mock Data
// Used during Design First phase (before API integration)
// ============================================================

import type {
  Patient,
  Appointment,
  Activity,
  DashboardStats,
  MedicalRecord,
  Pet,
} from "@/app/types";

// ── Current User (Clinic Staff) ──────────────────────────────

export const MOCK_CURRENT_USER = {
  id: "staff-001",
  name: "Dr. Emily Carter",
  role: "vet" as const,
  avatarUrl: undefined as string | undefined,
  initials: "EC",
};

export const MOCK_CLINIC = {
  id: "clinic-001",
  name: "Happy Paws Veterinary Clinic",
};

// ── Dashboard Stats ───────────────────────────────────────────

export const MOCK_DASHBOARD_STATS: DashboardStats = {
  todayPatients: 18,
  todayPatientsChange: "+12% vs yesterday",
  pendingRequests: 6,
  upcomingAppointments: 7,
  nextAppointmentTime: "Next: 11:30 AM",
  checkedInNow: 5,
};

// ── Today's Schedule ─────────────────────────────────────────

export const MOCK_TODAY_SCHEDULE: Appointment[] = [
  {
    id: "apt-001",
    time: "10:00 AM",
    petName: "Bella",
    petSpecies: "dog",
    ownerName: "Sarah Johnson",
    type: "Vaccination",
    status: "done",
  },
  {
    id: "apt-002",
    time: "10:30 AM",
    petName: "Max",
    petSpecies: "dog",
    ownerName: "Michael Brown",
    type: "Consultation",
    status: "done",
  },
  {
    id: "apt-003",
    time: "11:00 AM",
    petName: "Charlie",
    petSpecies: "cat",
    ownerName: "Lisa Davis",
    type: "Vaccination",
    status: "checked-in",
  },
  {
    id: "apt-004",
    time: "11:30 AM",
    petName: "Luna",
    petSpecies: "cat",
    ownerName: "Emily Chen",
    type: "Follow-up",
    status: "scheduled",
  },
];

// ── Recent Activity ───────────────────────────────────────────

export const MOCK_RECENT_ACTIVITY: Activity[] = [
  {
    id: "act-001",
    type: "check-in",
    description: "Bella (Golden Retriever) checked in by Sarah Johnson",
    time: "10:13 AM",
    petName: "Bella",
    ownerName: "Sarah Johnson",
  },
  {
    id: "act-002",
    type: "new-appointment",
    description: "New appointment booked for Max (Poodle) at 2:30 PM",
    time: "09:30 AM",
    petName: "Max",
    ownerName: "Michael Brown",
  },
  {
    id: "act-003",
    type: "record-updated",
    description: "Medical record updated for Charlie (Persian Cat)",
    time: "09:10 AM",
    petName: "Charlie",
    ownerName: "Lisa Davis",
  },
  {
    id: "act-004",
    type: "visit-completed",
    description: "Luna (Siamese Cat) visit completed",
    time: "08:13 AM",
    petName: "Luna",
    ownerName: "Emily Chen",
  },
];

// ── Patients ─────────────────────────────────────────────────

export const MOCK_PATIENTS: Patient[] = [
  {
    id: "pet-001",
    petNexusId: "PNX-2034-0087",
    name: "Bella",
    species: "dog",
    breed: "Golden Retriever",
    gender: "female",
    birthDate: "2019-03-12",
    weightKg: 28.5,
    photoUrl: undefined,
    ownerId: "owner-001",
    ownerName: "Sarah Johnson",
    ownerPhone: "081-234-5678",
    status: "approved",
    lastVisit: "May 20, 2025",
  },
  {
    id: "pet-002",
    petNexusId: "PNX-2034-0088",
    name: "Billie",
    species: "dog",
    breed: "Golden Retriever",
    gender: "male",
    birthDate: "2020-07-01",
    weightKg: 24.0,
    ownerId: "owner-002",
    ownerName: "Michael Brown",
    ownerPhone: "082-345-6789",
    status: "approved",
    lastVisit: "May 23, 2025",
  },
  {
    id: "pet-003",
    petNexusId: "PNX-2034-0089",
    name: "Neo",
    species: "dog",
    breed: "Poodle",
    gender: "male",
    birthDate: "2021-01-15",
    weightKg: 8.2,
    ownerId: "owner-003",
    ownerName: "Jessica Hall",
    ownerPhone: "083-456-7890",
    status: "approved",
    lastVisit: "May 24, 2025",
  },
  {
    id: "pet-004",
    petNexusId: "PNX-2034-0090",
    name: "Latte",
    species: "cat",
    breed: "Persian Cat",
    gender: "female",
    birthDate: "2020-05-22",
    weightKg: 4.1,
    ownerId: "owner-004",
    ownerName: "Parisa Ean",
    ownerPhone: "084-567-8901",
    status: "approved",
    lastVisit: "May 19, 2025",
  },
  {
    id: "pet-005",
    petNexusId: "PNX-2034-0091",
    name: "Luna",
    species: "cat",
    breed: "Siamese Cat",
    gender: "female",
    birthDate: "2022-09-03",
    weightKg: 3.8,
    ownerId: "owner-005",
    ownerName: "Emily Chen",
    ownerPhone: "085-678-9012",
    status: "approved",
    lastVisit: "May 17, 2025",
  },
  {
    id: "pet-006",
    petNexusId: "PNX-2034-0092",
    name: "Cooper",
    species: "dog",
    breed: "Beagle",
    gender: "male",
    birthDate: "2018-11-30",
    weightKg: 11.3,
    ownerId: "owner-006",
    ownerName: "Robert Wilson",
    ownerPhone: "086-789-0123",
    status: "approved",
    lastVisit: "May 25, 2025",
  },
  {
    id: "pet-007",
    petNexusId: "PNX-2034-0093",
    name: "Mochi",
    species: "cat",
    breed: "British Shorthair",
    gender: "female",
    birthDate: "2023-02-14",
    weightKg: 3.5,
    ownerId: "owner-001",
    ownerName: "Sarah Johnson",
    ownerPhone: "081-234-5678",
    status: "approved",
    lastVisit: "Apr 12, 2025",
  },
];

// ── QR Check-in — Sample Pet ──────────────────────────────────

export const MOCK_QR_PET: Pet = {
  id: "pet-007",
  petNexusId: "PNX-2034-0087",
  name: "Mochi",
  species: "cat",
  breed: "Persian / Female",
  gender: "female",
  birthDate: "2023-02-14",
  weightKg: 3.5,
  ownerId: "owner-001",
  ownerName: "Sarah Johnson",
  ownerPhone: "081-234-5678",
};

// ── Medical Record — Sample ───────────────────────────────────

export const MOCK_MEDICAL_RECORD: MedicalRecord = {
  id: "rec-001",
  petId: "pet-007",
  clinicId: "clinic-001",
  vetId: "staff-001",
  visitDate: "2025-05-20",
  visitType: "Consultation",
  chiefComplaint: "",
  diagnosis: "",
  treatment: "",
  medications: [],
  followUpDate: undefined,
  followUpNote: "",
  note: "",
};

// ── Pet for Medical Record Form ────────────────────────────────

export const MOCK_RECORD_PET: Pet = {
  id: "pet-007",
  petNexusId: "PNX-2034-0087",
  name: "Mochi",
  species: "cat",
  breed: "Persian / Male",
  gender: "male",
  birthDate: "2021-04-10",
  weightKg: 5.0,
  ownerId: "owner-001",
  ownerName: "Sarah Johnson",
  ownerPhone: "081-123-4578",
};
