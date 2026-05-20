import { ClearLisSkinsToken, DeleteSavedSkin, GetAppSettings, GetSavedSkins, HasLisSkinsToken, SaveAppSettings, SaveSkin, SearchNewSkins, SetLisSkinsToken, UpdateAllSavedSkinsPrices, UpdateSavedSkinPrice } from "../../../wailsjs/go/main/App";
import type { skins } from "../../../wailsjs/go/models";
import { toApiError } from "../../../shared/api/errors";
import { normalizeAutoRefreshIntervalSeconds, normalizeCurrency, normalizeSavedSkinsViewMode } from "../../../shared/config/settings";
import { logClientEvent } from "../../../shared/lib/logging/logger";
import { toSkinNameColor } from "../../../shared/lib/skinNameColor";
import type {
  BulkPriceUpdateResult,
  NewSkin,
  NewSkinsSearchParams,
  PaginatedResult,
  PriceUpdateResult,
  SaveSkinResult,
  SavedSkin,
  SavedSkinCurrency,
  SavedSkinsViewMode,
} from "../model/types";

const DEFAULT_LIMIT = 20;
const DEFAULT_OFFSET = 0;
const SEARCH_TIMEOUT_MS = 15000;
const inFlightNewSkinsRequests = new Map<string, Promise<PaginatedResult<NewSkin>>>();

const normalizeImageUrl = (imageUrl?: string, pageUrl?: string): string => {
  if (!imageUrl) return "";
  if (imageUrl.startsWith("http://") || imageUrl.startsWith("https://") || imageUrl.startsWith("data:")) return imageUrl;
  if (imageUrl.startsWith("//")) return `https:${imageUrl}`;
  if (imageUrl.startsWith("/")) {
    if (pageUrl?.includes("lis-skins.com")) {
      return `https://lis-skins.com${imageUrl}`;
    }
    return `https://community.akamai.steamstatic.com${imageUrl}`;
  }
  return `https://community.akamai.steamstatic.com/economy/image/${imageUrl}`;
};

const mapSavedSkin = (item: skins.SavedSkinItem): SavedSkin => {
  const steamPageUrl = item.steam_page_url;
  const lisSkinsPageUrl = item.lisskins_page_url;

  return {
    id: item.market_hash_name,
    name: item.market_hash_name,
    title: item.display_name,
    nameColor: toSkinNameColor(item.name_color),
    imageUrl: normalizeImageUrl(item.icon_url, steamPageUrl),
    steamPageUrl,
    lisSkinsPageUrl,
    priceText: item.steam_price_text,
    currency: normalizeCurrency(item.currency),
    updatedAt: item.steam_updated_at,
    steamPriceText: item.steam_price_text,
    steamUpdatedAt: item.steam_updated_at,
    lisSkinsPriceText: item.lisskins_price_text,
    lisSkinsUpdatedAt: item.lisskins_updated_at,
  };
};

const mapNewSkin = (item: skins.NewSkinItem): NewSkin => ({
  id: item.market_hash_name,
  name: item.market_hash_name,
  title: item.display_name,
  nameColor: toSkinNameColor(item.name_color),
  imageUrl: normalizeImageUrl(item.icon_url, item.page_url),
  steamPageUrl: item.page_url,
  lisSkinsPageUrl: "",
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

const buildNewSkinsSearchKey = (
  params: NewSkinsSearchParams,
  limit: number,
  offset: number,
  cursor: string,
): string =>
  JSON.stringify({
    query: params.query.trim(),
    sortColumn: params.sortColumn,
    sortDir: params.sortDir,
    filters: params.filters,
    limit,
    offset,
    cursor,
  });

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
  params: NewSkinsSearchParams,
  limit = DEFAULT_LIMIT,
  offset = DEFAULT_OFFSET,
  cursor = "",
): Promise<PaginatedResult<NewSkin>> => {
  const requestKey = buildNewSkinsSearchKey(params, limit, offset, cursor);
  const activeRequest = inFlightNewSkinsRequests.get(requestKey);
  if (activeRequest) {
    return activeRequest;
  }

  const requestPromise = (async () => {
    try {
      const response = await withTimeout(
        SearchNewSkins({
          market_hash_name: params.query.trim() || undefined,
          sort_column: params.sortColumn,
          sort_dir: params.sortDir,
          price_min: params.filters.priceMin.trim() || undefined,
          price_max: params.filters.priceMax.trim() || undefined,
          search_descriptions: params.filters.searchDescriptions,
          type: params.filters.type,
          weapon: params.filters.weapon,
          rarity: params.filters.rarity,
          exterior: params.filters.exterior,
          item_set: params.filters.itemSet,
          pro_player: params.filters.proPlayer,
          sticker_capsule: params.filters.stickerCapsule,
          tournament_team: params.filters.tournamentTeam,
          limit,
          offset,
          cursor,
        }),
        SEARCH_TIMEOUT_MS,
        "search timeout",
      );
      return {
        items: response.items.map(mapNewSkin),
        total: response.total_count,
        limit: response.limit,
        offset: response.offset,
        nextCursor: response.next_cursor,
      };
    } catch (err) {
      logClientEvent("error", "getNewSkins failed", "skinApi", {
        operation: "getNewSkins",
        query: params.query,
        cursor,
        error: err instanceof Error ? err.message : String(err ?? ""),
      });
      throw toApiError(err);
    } finally {
    inFlightNewSkinsRequests.delete(requestKey);
  }
  })();

  inFlightNewSkinsRequests.set(requestKey, requestPromise);
  return requestPromise;
};

export const updateSkinPrice = async (skinId: string, currency: SavedSkinCurrency): Promise<PriceUpdateResult> => {
  try {
    const response = await UpdateSavedSkinPrice({ market_hash_name: skinId, currency });
    return {
      updated: 1,
      marketHashName: response.market_hash_name,
      steamPageUrl: response.steam_page_url,
      steamPriceText: response.steam_price_text,
      steamUpdatedAt: response.steam_updated_at,
      lisSkinsPageUrl: response.lisskins_page_url,
      lisSkinsPriceText: response.lisskins_price_text,
      lisSkinsUpdatedAt: response.lisskins_updated_at,
      currency: normalizeCurrency(response.currency),
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

export const saveSkin = async (skin: Pick<NewSkin, "id" | "title" | "nameColor" | "imageUrl" | "steamPageUrl">): Promise<SaveSkinResult> => {
  try {
    const response = await SaveSkin({
      market_hash_name: skin.id,
      display_name: skin.title,
      name_color: skin.nameColor?.slice(1) ?? "",
      icon_url: skin.imageUrl,
      page_url: skin.steamPageUrl,
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

export const hasLisSkinsToken = async (): Promise<boolean> => {
  try {
    const response = await HasLisSkinsToken();
    return response.hasToken;
  } catch (err) {
    logClientEvent("error", "hasLisSkinsToken failed", "skinApi", {
      operation: "hasLisSkinsToken",
      error: err instanceof Error ? err.message : String(err ?? ""),
    });
    throw toApiError(err);
  }
};

export const saveLisSkinsToken = async (token: string): Promise<void> => {
  try {
    await SetLisSkinsToken({ token });
  } catch (err) {
    logClientEvent("error", "saveLisSkinsToken failed", "skinApi", {
      operation: "saveLisSkinsToken",
      error: err instanceof Error ? err.message : String(err ?? ""),
    });
    throw toApiError(err);
  }
};

export const clearLisSkinsToken = async (): Promise<void> => {
  try {
    await ClearLisSkinsToken();
  } catch (err) {
    logClientEvent("error", "clearLisSkinsToken failed", "skinApi", {
      operation: "clearLisSkinsToken",
      error: err instanceof Error ? err.message : String(err ?? ""),
    });
    throw toApiError(err);
  }
};

export type AppSettings = {
  currency: SavedSkinCurrency;
  autoRefreshIntervalSeconds: number;
  savedSkinsViewMode: SavedSkinsViewMode;
};

export const getAppSettings = async (): Promise<AppSettings> => {
  try {
    const response = await GetAppSettings();
    return {
      currency: normalizeCurrency(response.currency),
      autoRefreshIntervalSeconds: normalizeAutoRefreshIntervalSeconds(response.auto_refresh_interval_seconds),
      savedSkinsViewMode: normalizeSavedSkinsViewMode(response.saved_skins_view_mode),
    };
  } catch (err) {
    logClientEvent("error", "getAppSettings failed", "skinApi", {
      operation: "getAppSettings",
      error: err instanceof Error ? err.message : String(err ?? ""),
    });
    throw toApiError(err);
  }
};

export const saveAppSettings = async (settings: AppSettings): Promise<AppSettings> => {
  try {
    const normalized = {
      currency: normalizeCurrency(settings.currency),
      autoRefreshIntervalSeconds: normalizeAutoRefreshIntervalSeconds(settings.autoRefreshIntervalSeconds),
      savedSkinsViewMode: normalizeSavedSkinsViewMode(settings.savedSkinsViewMode),
    };
    await SaveAppSettings({
      currency: normalized.currency,
      auto_refresh_interval_seconds: normalized.autoRefreshIntervalSeconds,
      saved_skins_view_mode: normalized.savedSkinsViewMode,
    });
    return normalized;
  } catch (err) {
    logClientEvent("error", "saveAppSettings failed", "skinApi", {
      operation: "saveAppSettings",
      error: err instanceof Error ? err.message : String(err ?? ""),
    });
    throw toApiError(err);
  }
};
