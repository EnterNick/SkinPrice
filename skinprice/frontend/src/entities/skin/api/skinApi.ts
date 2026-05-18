import { DeleteSavedSkin, GetSavedSkins, SaveSkin, SearchNewSkins, UpdateAllSavedSkinsPrices, UpdateSavedSkinPrice } from "../../../wailsjs/go/main/App";
import type { skins } from "../../../wailsjs/go/models";
import { toApiError } from "../../../shared/api/errors";
import { logClientEvent } from "../../../shared/lib/logging/logger";
import type { BulkPriceUpdateResult, CurrencyOption, NewSkin, PaginatedResult, PriceUpdateResult, SaveSkinResult, SavedSkin, SavedSkinCurrency } from "../model/types";

const DEFAULT_LIMIT = 20;
const DEFAULT_OFFSET = 0;
const SEARCH_TIMEOUT_MS = 8000;

export const DEFAULT_CURRENCY: SavedSkinCurrency = "1";
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

const normalizeImageUrl = (imageUrl?: string): string => {
  if (!imageUrl) return "";
  if (imageUrl.startsWith("http://") || imageUrl.startsWith("https://") || imageUrl.startsWith("data:")) return imageUrl;
  if (imageUrl.startsWith("//")) return `https:${imageUrl}`;
  if (imageUrl.startsWith("/")) return `https://community.akamai.steamstatic.com${imageUrl}`;
  return `https://community.akamai.steamstatic.com/economy/image/${imageUrl}`;
};

const mapSavedSkin = (item: skins.SavedSkinItem): SavedSkin => ({
  id: item.market_hash_name,
  name: item.market_hash_name,
  title: item.display_name,
  imageUrl: normalizeImageUrl(item.icon_url),
  pageUrl: item.page_url,
  priceText: item.price_text,
  currency: normalizeCurrency(item.currency),
  updatedAt: item.updated_at,
});

const mapNewSkin = (item: skins.NewSkinItem): NewSkin => ({
  id: item.market_hash_name,
  name: item.market_hash_name,
  title: item.display_name,
  imageUrl: normalizeImageUrl(item.icon_url),
  pageUrl: item.page_url,
  priceText: item.price_text,
  priceCents: item.price_cents,
  sellListings: item.sell_listings,
});

const withTimeout = async <T>(promise: Promise<T>, timeoutMs: number, message: string): Promise<T> => {
  let timeoutId: number | undefined;

  try {
    return await Promise.race([
      promise,
      new Promise<T>((_, reject) => {
        timeoutId = window.setTimeout(() => reject(new Error(message)), timeoutMs);
      }),
    ]);
  } finally {
    if (timeoutId !== undefined) {
      window.clearTimeout(timeoutId);
    }
  }
};

export const getSavedSkins = async (): Promise<PaginatedResult<SavedSkin>> => {
  try {
    const response = await GetSavedSkins({ limit: DEFAULT_LIMIT, offset: DEFAULT_OFFSET });
    return {
      items: response.items.map(mapSavedSkin),
      total: response.total_count,
      limit: response.limit,
      offset: response.offset,
    };
  } catch (err) {
    logClientEvent("error", "getSavedSkins failed", "skinApi", {
      operation: "getSavedSkins",
      error: err instanceof Error ? err.message : String(err ?? ""),
    });
    throw toApiError(err);
  }
};

export const getNewSkins = async (
  query?: string,
  limit = DEFAULT_LIMIT,
  offset = DEFAULT_OFFSET,
  source: "steam" | "lisskins" = "steam",
): Promise<PaginatedResult<NewSkin>> => {
  try {
    const response = await withTimeout(
      SearchNewSkins({ market_hash_name: query, source, limit, offset }),
      SEARCH_TIMEOUT_MS,
      "search timeout",
    );
    return {
      items: response.items.map(mapNewSkin),
      total: response.total_count,
      limit: response.limit,
      offset: response.offset,
    };
  } catch (err) {
    logClientEvent("error", "getNewSkins failed", "skinApi", {
      operation: "getNewSkins",
      source,
      query: query ?? "",
      error: err instanceof Error ? err.message : String(err ?? ""),
    });
    throw toApiError(err);
  }
};

export const updateSkinPrice = async (skinId: string, currency: SavedSkinCurrency): Promise<PriceUpdateResult> => {
  try {
    const response = await UpdateSavedSkinPrice({ market_hash_name: skinId, currency });
    return {
      updated: 1,
      marketHashName: response.market_hash_name,
      priceText: response.price_text,
      currency: normalizeCurrency(response.currency),
      updatedAt: response.updated_at,
    };
  } catch (err) {
    logClientEvent("error", "updateSkinPrice failed", "skinApi", {
      operation: "updateSkinPrice",
      skinId,
      currency,
      error: err instanceof Error ? err.message : String(err ?? ""),
    });
    throw toApiError(err);
  }
};

export const updateAllSkinPrices = async (currency: SavedSkinCurrency): Promise<BulkPriceUpdateResult> => {
  try {
    const response = await UpdateAllSavedSkinsPrices({ currency });
    return {
      updated: response.updated_count,
      failed: response.failed_count,
      failures: response.failures.map((failure) => ({
        marketHashName: failure.market_hash_name,
        message: failure.message,
      })),
    };
  } catch (err) {
    logClientEvent("error", "updateAllSkinPrices failed", "skinApi", {
      operation: "updateAllSkinPrices",
      currency,
      error: err instanceof Error ? err.message : String(err ?? ""),
    });
    throw toApiError(err);
  }
};

export const saveSkin = async (skin: Pick<NewSkin, "id" | "title" | "imageUrl" | "pageUrl">): Promise<SaveSkinResult> => {
  try {
    const response = await SaveSkin({
      market_hash_name: skin.id,
      display_name: skin.title,
      icon_url: skin.imageUrl,
      page_url: skin.pageUrl,
    });
    return { created: response.created };
  } catch (err) {
    logClientEvent("error", "saveSkin failed", "skinApi", {
      operation: "saveSkin",
      skinId: skin.id,
      error: err instanceof Error ? err.message : String(err ?? ""),
    });
    throw toApiError(err);
  }
};

export const deleteSavedSkin = async (skinId: string): Promise<void> => {
  try {
    await DeleteSavedSkin({ market_hash_name: skinId });
  } catch (err) {
    logClientEvent("error", "deleteSavedSkin failed", "skinApi", {
      operation: "deleteSavedSkin",
      skinId,
      error: err instanceof Error ? err.message : String(err ?? ""),
    });
    throw toApiError(err);
  }
};
