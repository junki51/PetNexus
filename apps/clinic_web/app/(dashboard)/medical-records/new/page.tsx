"use client";

import React, { useState } from "react";
import { ArrowLeft, Plus, Trash2, Save, FileEdit, CheckCircle2 } from "lucide-react";
import { Card, CardHeader, CardTitle, CardBody } from "@/app/components/ui/Card";
import { Input } from "@/app/components/ui/Input";
import { Select } from "@/app/components/ui/Select";
import { Button } from "@/app/components/ui/Button";
import { Avatar } from "@/app/components/ui/Avatar";
import { MOCK_RECORD_PET } from "@/app/lib/mock-data";
import type { Medication, SelectOption } from "@/app/types";
import Link from "next/link";
import { useRouter } from "next/navigation";
import { useLanguage } from "@/app/components/LanguageContext";

export default function NewMedicalRecordPage() {
  const router = useRouter();
  const { t } = useLanguage();
  const pet = MOCK_RECORD_PET;

  // Form states
  const [visitDate, setVisitDate] = useState("2025-05-20");
  const [visitTime, setVisitTime] = useState("10:30");
  const [visitType, setVisitType] = useState("Consultation");
  const [symptoms, setSymptoms] = useState("");
  const [diagnosis, setDiagnosis] = useState("");
  const [treatment, setTreatment] = useState("");
  const [medications, setMedications] = useState<Medication[]>([
    { name: "", dosage: "", instructions: "" },
  ]);
  const [followUpDate, setFollowUpDate] = useState("");
  const [followUpNote, setFollowUpNote] = useState("");
  const [vet, setVet] = useState("Dr. Emily Carter");
  const [recordStatus, setRecordStatus] = useState("In-Progress");

  const [loading, setLoading] = useState(false);
  const [success, setSuccess] = useState(false);

  const VISIT_TYPE_OPTIONS: SelectOption[] = React.useMemo(() => [
    { value: "Consultation", label: t("visit_type") + ": Consultation" },
    { value: "Vaccination", label: t("visit_type") + ": Vaccination" },
    { value: "Follow-up", label: t("visit_type") + ": Follow-up" },
    { value: "Emergency", label: t("visit_type") + ": Emergency" },
    { value: "Grooming", label: t("visit_type") + ": Grooming" },
  ], [t]);

  const VET_OPTIONS: SelectOption[] = [
    { value: "Dr. Emily Carter", label: "Dr. Emily Carter" },
    { value: "Dr. James Wilson", label: "Dr. James Wilson" },
  ];

  const STATUS_OPTIONS: SelectOption[] = [
    { value: "In-Progress", label: "In-Progress" },
    { value: "Completed", label: "Completed" },
  ];

  // Medications handlers
  const handleAddMedication = () => {
    setMedications([...medications, { name: "", dosage: "", instructions: "" }]);
  };

  const handleRemoveMedication = (index: number) => {
    setMedications(medications.filter((_, i) => i !== index));
  };

  const handleMedicationChange = (
    index: number,
    field: keyof Medication,
    value: string
  ) => {
    const updated = medications.map((med, i) => {
      if (i === index) {
        return { ...med, [field]: value };
      }
      return med;
    });
    setMedications(updated);
  };

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault();
    setLoading(true);

    // Simulate saving medical record
    setTimeout(() => {
      setLoading(false);
      setSuccess(true);
      setTimeout(() => {
        router.push("/patients");
      }, 1500);
    }, 1500);
  };

  return (
    <div className="flex flex-col gap-6 max-w-7xl mx-auto">
      {/* Header with Back button */}
      <div className="flex items-center gap-3">
        <Link href="/patients">
          <button className="p-2 bg-white border border-navy-200 text-navy-600 hover:text-navy-800 hover:bg-navy-50 rounded-lg transition-colors cursor-pointer shadow-sm">
            <ArrowLeft size={16} />
          </button>
        </Link>
        <div>
          <div className="flex items-center gap-2 text-xs text-navy-500 font-medium">
            <span>{t("patients")}</span>
            <span>/</span>
            <span>{pet.name}</span>
            <span>/</span>
            <span className="text-teal-600 font-semibold">{t("new_record_title")}</span>
          </div>
          <h1 className="text-2xl font-bold text-navy-900 mt-1">{t("new_record_title")}</h1>
        </div>
      </div>

      {success && (
        <div className="bg-emerald-50 border border-emerald-200 text-emerald-800 rounded-xl p-4 flex items-center gap-3 animate-[slide-up_0.2s_ease-out]">
          <CheckCircle2 className="text-emerald-500 w-5 h-5 shrink-0" />
          <div className="text-sm font-medium">
            {t("success_save_record")}
          </div>
        </div>
      )}

      <form onSubmit={handleSubmit} className="grid grid-cols-1 lg:grid-cols-12 gap-6 items-start">
        {/* Left Column — Form Fields (lg:col-span-8) */}
        <div className="lg:col-span-8 flex flex-col gap-6">
          <Card className="bg-white">
            <CardBody className="flex flex-col gap-6">
              {/* Visit Information */}
              <div>
                <h3 className="text-sm font-bold text-navy-800 uppercase tracking-wider mb-4">
                  {t("visit_info")}
                </h3>
                <div className="grid grid-cols-1 sm:grid-cols-3 gap-4">
                  <Input
                    label={t("visit_date")}
                    type="date"
                    value={visitDate}
                    onChange={(e) => setVisitDate(e.target.value)}
                    required
                  />
                  <Input
                    label={t("visit_time")}
                    type="time"
                    value={visitTime}
                    onChange={(e) => setVisitTime(e.target.value)}
                    required
                  />
                  <Select
                    label={t("visit_type")}
                    options={VISIT_TYPE_OPTIONS}
                    value={visitType}
                    onChange={setVisitType}
                    required
                  />
                </div>
              </div>

              <hr className="border-navy-100" />

              {/* Symptoms / Complaint */}
              <div>
                <label className="block text-sm font-bold text-navy-800 uppercase tracking-wider mb-2">
                  {t("symptoms_complaint")}
                </label>
                <textarea
                  className="w-full min-h-[100px] p-3 rounded-lg border border-navy-200 bg-white text-sm text-navy-800 placeholder:text-navy-400 focus:outline-none focus:ring-2 focus:ring-teal-500/30 focus:border-teal-500 transition-colors"
                  placeholder={t("describe_symptoms")}
                  value={symptoms}
                  onChange={(e) => setSymptoms(e.target.value)}
                  rows={4}
                  required
                />
              </div>

              <hr className="border-navy-100" />

              {/* Diagnosis */}
              <div>
                <label className="block text-sm font-bold text-navy-800 uppercase tracking-wider mb-2">
                  {t("diagnosis")}
                </label>
                <textarea
                  className="w-full min-h-[80px] p-3 rounded-lg border border-navy-200 bg-white text-sm text-navy-800 placeholder:text-navy-400 focus:outline-none focus:ring-2 focus:ring-teal-500/30 focus:border-teal-500 transition-colors"
                  placeholder={t("describe_diagnosis")}
                  value={diagnosis}
                  onChange={(e) => setDiagnosis(e.target.value)}
                  rows={3}
                  required
                />
              </div>

              <hr className="border-navy-100" />

              {/* Treatment / Procedures */}
              <div>
                <label className="block text-sm font-bold text-navy-800 uppercase tracking-wider mb-2">
                  {t("treatment_procedures")}
                </label>
                <textarea
                  className="w-full min-h-[100px] p-3 rounded-lg border border-navy-200 bg-white text-sm text-navy-800 placeholder:text-navy-400 focus:outline-none focus:ring-2 focus:ring-teal-500/30 focus:border-teal-500 transition-colors"
                  placeholder={t("describe_treatment")}
                  value={treatment}
                  onChange={(e) => setTreatment(e.target.value)}
                  rows={4}
                />
              </div>

              <hr className="border-navy-100" />

              {/* Medications Prescribed */}
              <div>
                <div className="flex items-center justify-between mb-4">
                  <h3 className="text-sm font-bold text-navy-800 uppercase tracking-wider">
                    {t("medications_prescribed")}
                  </h3>
                  <Button
                    type="button"
                    variant="outline"
                    size="sm"
                    onClick={handleAddMedication}
                    icon={<Plus size={14} />}
                    className="cursor-pointer"
                  >
                    {t("add_medication")}
                  </Button>
                </div>

                <div className="flex flex-col gap-3">
                  {medications.map((med, index) => (
                    <div key={index} className="flex gap-3 items-end flex-wrap sm:flex-nowrap">
                      <div className="flex-1 min-w-[200px]">
                        <Input
                          placeholder={t("med_name")}
                          value={med.name}
                          onChange={(e) =>
                            handleMedicationChange(index, "name", e.target.value)
                          }
                        />
                      </div>
                      <div className="w-full sm:w-1/3 min-w-[120px]">
                        <Input
                          placeholder={t("dosage")}
                          value={med.dosage}
                          onChange={(e) =>
                            handleMedicationChange(index, "dosage", e.target.value)
                          }
                        />
                      </div>
                      <div className="flex-1 min-w-[200px]">
                        <Input
                          placeholder={t("instructions")}
                          value={med.instructions}
                          onChange={(e) =>
                            handleMedicationChange(index, "instructions", e.target.value)
                          }
                        />
                      </div>
                      {medications.length > 1 && (
                        <button
                          type="button"
                          onClick={() => handleRemoveMedication(index)}
                          className="p-2 text-red-500 hover:text-red-700 hover:bg-red-50 rounded-lg transition-colors cursor-pointer shrink-0 h-10 border border-transparent"
                        >
                          <Trash2 size={16} />
                        </button>
                      )}
                    </div>
                  ))}
                </div>
              </div>

              <hr className="border-navy-100" />

              {/* Follow-up Plan */}
              <div>
                <h3 className="text-sm font-bold text-navy-800 uppercase tracking-wider mb-4">
                  {t("follow_up_plan")}
                </h3>
                <div className="grid grid-cols-1 sm:grid-cols-2 gap-4">
                  <Input
                    label={t("follow_up_date")}
                    type="date"
                    value={followUpDate}
                    onChange={(e) => setFollowUpDate(e.target.value)}
                  />
                  <div className="flex flex-col gap-1.5">
                    <label className="text-sm font-medium text-navy-700">
                      {t("follow_up_notes")}
                    </label>
                    <input
                      type="text"
                      placeholder={t("notes_placeholder")}
                      className="h-10 px-3 rounded-lg border border-navy-200 bg-white text-sm text-navy-900 focus:outline-none focus:ring-2 focus:ring-teal-500/30 focus:border-teal-500 transition-colors"
                      value={followUpNote}
                      onChange={(e) => setFollowUpNote(e.target.value)}
                    />
                  </div>
                </div>
              </div>
            </CardBody>
          </Card>

          {/* Action buttons footer */}
          <div className="flex justify-end gap-3">
            <Link href="/patients">
              <Button type="button" variant="ghost">
                {t("cancel")}
              </Button>
            </Link>
            <Button type="button" variant="outline" icon={<FileEdit size={16} />} className="cursor-pointer">
              {t("save_draft")}
            </Button>
            <Button type="submit" loading={loading} icon={<Save size={16} />} className="cursor-pointer">
              {t("save_record")}
            </Button>
          </div>
        </div>

        {/* Right Column — Patient Details Sidebar (lg:col-span-4) */}
        <div className="lg:col-span-4 flex flex-col gap-6 sticky top-20">
          {/* Pet Info Card */}
          <Card className="bg-white">
            <CardBody>
              <div className="flex flex-col items-center text-center pb-6 border-b border-navy-100">
                <Avatar name={pet.name} size="xl" className="mb-4" />
                <h3 className="text-lg font-bold text-navy-900">{pet.name}</h3>
                <p className="text-xs font-semibold text-navy-500 mt-1 uppercase tracking-wide">
                  {pet.breed}
                </p>
                <span className="text-[10px] font-mono bg-teal-50 text-teal-600 border border-teal-200 rounded px-2 py-0.5 mt-2 font-bold">
                  ID: {pet.petNexusId}
                </span>
              </div>

              {/* Pet metadata list */}
              <div className="py-6 border-b border-navy-100 space-y-4 text-sm">
                <div className="flex justify-between">
                  <span className="text-navy-500 font-medium">{t("pet_owner")}</span>
                  <span className="text-navy-800 font-semibold">{pet.ownerName}</span>
                </div>
                <div className="flex justify-between">
                  <span className="text-navy-500 font-medium">{t("pet_phone")}</span>
                  <span className="text-navy-800 font-semibold">{pet.ownerPhone}</span>
                </div>
                <div className="flex justify-between">
                  <span className="text-navy-500 font-medium">{t("pet_weight")}</span>
                  <span className="text-navy-800 font-semibold">{pet.weightKg} kg</span>
                </div>
                <div className="flex justify-between">
                  <span className="text-navy-500 font-medium">{t("pet_age")}</span>
                  <span className="text-navy-800 font-semibold">{t("pet_age") === "อายุ" ? "2 ปี" : "2 years"}</span>
                </div>
                <div className="flex justify-between">
                  <span className="text-navy-500 font-medium">{t("last_visit")}</span>
                  <span className="text-navy-800 font-semibold">Apr 12, 2025</span>
                </div>
              </div>

              <div className="pt-4 text-center">
                <button
                  type="button"
                  onClick={() => {}}
                  className="text-xs font-semibold text-teal-600 hover:text-teal-700 cursor-pointer"
                >
                  {t("view_full_history")}
                </button>
              </div>
            </CardBody>
          </Card>

          {/* Admin Info Card */}
          <Card className="bg-white">
            <CardHeader>
              <CardTitle subtitle={t("admin_info")}>
                {t("admin_info")}
              </CardTitle>
            </CardHeader>
            <CardBody className="flex flex-col gap-4">
              <Select
                label={t("veterinarian")}
                options={VET_OPTIONS}
                value={vet}
                onChange={setVet}
              />
              <Select
                label={t("record_status")}
                options={STATUS_OPTIONS}
                value={recordStatus}
                onChange={setRecordStatus}
              />
            </CardBody>
          </Card>
        </div>
      </form>
    </div>
  );
}
