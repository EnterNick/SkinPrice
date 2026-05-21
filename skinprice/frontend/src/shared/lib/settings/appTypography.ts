import type { FontFamilyOptionValue } from "../../../entities/skin/model/types";
import {
  DEFAULT_FONT_FAMILY,
  DEFAULT_FONT_SIZE_PX,
  normalizeFontFamily,
  normalizeFontSizePx,
} from "../../config/settings";

const FONT_FAMILY_MAP: Record<FontFamilyOptionValue, string> = {
  inter: '"Inter", sans-serif',
  system: 'system-ui, -apple-system, BlinkMacSystemFont, "Segoe UI", sans-serif',
  nunito: '"Nunito", sans-serif',
  roboto: '"Roboto", sans-serif',
  "ibm-plex-sans": '"IBM Plex Sans", sans-serif',
  manrope: '"Manrope", sans-serif',
  monocraft: '"Monocraft", monospace',
};

export type TypographySettings = {
  fontFamily?: FontFamilyOptionValue | string | null;
  fontSizePx?: number | string | null;
};

export const applyTypographySettings = (settings?: TypographySettings): void => {
  if (typeof document === "undefined") return;

  const root = document.documentElement;
  const fontFamily = normalizeFontFamily(settings?.fontFamily ?? DEFAULT_FONT_FAMILY);
  const fontSizePx = normalizeFontSizePx(settings?.fontSizePx ?? DEFAULT_FONT_SIZE_PX);

  root.style.setProperty("--app-font-family", FONT_FAMILY_MAP[fontFamily]);
  root.style.setProperty("--app-font-scale", String(fontSizePx / DEFAULT_FONT_SIZE_PX));
  root.style.setProperty("--app-base-font-size", `${fontSizePx}px`);
  root.style.setProperty("--app-ui-font-family", FONT_FAMILY_MAP[fontFamily]);
  root.style.setProperty("--app-ui-font-size", `${fontSizePx}px`);
  root.style.setProperty("--app-ui-font-scale", String(fontSizePx / DEFAULT_FONT_SIZE_PX));
};
