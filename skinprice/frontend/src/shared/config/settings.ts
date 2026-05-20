import type { CurrencyOption, SavedSkinCurrency, SavedSkinsViewMode } from "../../entities/skin/model/types";

export const DEFAULT_CURRENCY: SavedSkinCurrency = "1";
export const DEFAULT_AUTO_REFRESH_INTERVAL_SECONDS = 30;
export const MIN_AUTO_REFRESH_INTERVAL_SECONDS = 5;
export const DEFAULT_SAVED_SKINS_VIEW_MODE: SavedSkinsViewMode = "table";

export const CURRENCY_OPTIONS: CurrencyOption[] = [
  { value: "1", label: "USD" },
  { value: "5", label: "RUB" },
  { value: "3", label: "EUR" },
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

export const normalizeSavedSkinsViewMode = (value?: string | null): SavedSkinsViewMode => {
  switch (value) {
    case "table":
    case "cards":
      return value;
    default:
      return DEFAULT_SAVED_SKINS_VIEW_MODE;
  }
};
