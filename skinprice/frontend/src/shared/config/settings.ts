import type {
  CurrencyOption,
  FontFamilyOption,
  FontFamilyOptionValue,
  SavedSkinCurrency,
  SavedSkinsViewMode,
} from "../../entities/skin/model/types";

export const DEFAULT_CURRENCY: SavedSkinCurrency = "1";
export const DEFAULT_AUTO_REFRESH_ENABLED = true;
export const DEFAULT_AUTO_REFRESH_INTERVAL_SECONDS = 30;
export const MIN_AUTO_REFRESH_INTERVAL_SECONDS = 5;
export const DEFAULT_SAVED_SKINS_VIEW_MODE: SavedSkinsViewMode = "table";
export const DEFAULT_FONT_FAMILY: FontFamilyOptionValue = "nunito";
export const DEFAULT_FONT_SIZE_PX = 14;
export const MIN_FONT_SIZE_PX = 10;
export const MAX_FONT_SIZE_PX = 28;

export const CURRENCY_OPTIONS: CurrencyOption[] = [
  { value: "1", label: "USD" },
  { value: "5", label: "RUB" },
  { value: "3", label: "EUR" },
];

export const FONT_FAMILY_OPTIONS: FontFamilyOption[] = [
  { value: "inter", label: "Inter" },
  { value: "system", label: "System" },
  { value: "nunito", label: "Nunito" },
  { value: "roboto", label: "Roboto" },
  { value: "ibm-plex-sans", label: "IBM Plex Sans" },
  { value: "manrope", label: "Manrope" },
  { value: "monocraft", label: "Monocraft" },
];

export const normalizeCurrency = (value?: string | null): SavedSkinCurrency => {
  switch (value) {
    case "1":
    case "USD":
      return "1";
    case "3":
    case "EUR":
      return "3";
    case "5":
    case "RUB":
      return "5";
    default:
      return DEFAULT_CURRENCY;
  }
};

export const normalizeAutoRefreshIntervalSeconds = (value?: string | number | null): number => {
  const numericValue = typeof value === "number" ? value : Number(value);
  if (!Number.isFinite(numericValue)) {
    return DEFAULT_AUTO_REFRESH_INTERVAL_SECONDS;
  }

  if (numericValue < MIN_AUTO_REFRESH_INTERVAL_SECONDS) {
    return DEFAULT_AUTO_REFRESH_INTERVAL_SECONDS;
  }

  return Math.round(numericValue);
};

export const normalizeAutoRefreshEnabled = (value?: boolean | null): boolean => {
  return typeof value === "boolean" ? value : DEFAULT_AUTO_REFRESH_ENABLED;
};

export const normalizeSavedSkinsViewMode = (value?: string | null): SavedSkinsViewMode => {
  switch (value) {
    case "table":
    case "cards":
      return value;
    default:
      return DEFAULT_SAVED_SKINS_VIEW_MODE;
  }
};

export const normalizeFontFamily = (value?: string | null): FontFamilyOptionValue => {
  switch (value) {
    case "inter":
    case "system":
    case "nunito":
    case "roboto":
    case "ibm-plex-sans":
    case "manrope":
    case "monocraft":
      return value;
    default:
      return DEFAULT_FONT_FAMILY;
  }
};

export const normalizeFontSizePx = (value?: string | number | null): number => {
  const numericValue = typeof value === "number" ? value : Number(value);
  if (!Number.isFinite(numericValue)) {
    return DEFAULT_FONT_SIZE_PX;
  }

  const rounded = Math.round(numericValue);
  if (rounded < MIN_FONT_SIZE_PX || rounded > MAX_FONT_SIZE_PX) {
    return DEFAULT_FONT_SIZE_PX;
  }

  return rounded;
};
