"use client";

import React, { createContext, useContext, useState, useEffect } from "react";
import { translations, Locale, TranslationKeys } from "@/app/lib/translations";

interface LanguageContextProps {
  lang: Locale;
  setLang: (lang: Locale) => void;
  t: (key: TranslationKeys, replacements?: Record<string, string | number>) => string;
}

const LanguageContext = createContext<LanguageContextProps | null>(null);

export function LanguageProvider({ children }: { children: React.ReactNode }) {
  const [lang, setLangState] = useState<Locale>("th");

  useEffect(() => {
    const saved = localStorage.getItem("petnexus-lang") as Locale | null;
    if (saved === "th" || saved === "en") {
      Promise.resolve().then(() => {
        setLangState(saved);
      });
    }
  }, []);

  const setLang = (newLang: Locale) => {
    setLangState(newLang);
    localStorage.setItem("petnexus-lang", newLang);
  };

  const t = (key: TranslationKeys, replacements?: Record<string, string | number>): string => {
    const dictionary = translations[lang];
    const text: string = (dictionary as Record<string, string>)[key] || translations["th"][key] || String(key);

    if (!replacements) return text;

    // Perform replacements if variables are provided e.g. {total}
    let formattedText = text;
    Object.entries(replacements).forEach(([k, v]) => {
      formattedText = formattedText.replace(new RegExp(`{${k}}`, "g"), String(v));
    });

    return formattedText;
  };

  return (
    <LanguageContext.Provider value={{ lang, setLang, t }}>
      {children}
    </LanguageContext.Provider>
  );
}

export function useLanguage() {
  const context = useContext(LanguageContext);
  if (!context) {
    throw new Error("useLanguage must be used within a LanguageProvider");
  }
  return context;
}
