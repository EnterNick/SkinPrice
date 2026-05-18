import { GetSavedSkins, SaveSkin, SearchNewSkins, UpdateAllSavedSkinsPrices, UpdateSavedSkinPrice } from "../../../wailsjs/go/main/App";
import type { skins } from "../../../wailsjs/go/models";
import { toApiError } from "../../../shared/api/errors";
import type { NewSkin, PaginatedResult, PriceUpdateResult, SavedSkin } from "../model/types";

const DEFAULT_LIMIT = 20;
const DEFAULT_OFFSET = 0;

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
  currency: item.currency,
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
    throw toApiError(err);
  }
};

export const getNewSkins = async (query?: string, limit = DEFAULT_LIMIT, offset = DEFAULT_OFFSET): Promise<PaginatedResult<NewSkin>> => {
  try {
    const response = await SearchNewSkins({ market_hash_name: query, limit, offset });
    return {
      items: response.items.map(mapNewSkin),
      total: response.total_count,
      limit: response.limit,
      offset: response.offset,
    };
  } catch (err) {
    throw toApiError(err);
  }
};

export const updateSkinPrice = async (skinId: string, currency: string): Promise<PriceUpdateResult> => {
  try {
    await UpdateSavedSkinPrice({ market_hash_name: skinId, currency });
    return { updated: 1 };
  } catch (err) {
    throw toApiError(err);
  }
};

export const updateAllSkinPrices = async (currency: string): Promise<PriceUpdateResult> => {
  try {
    await UpdateAllSavedSkinsPrices({ currency });
    return { updated: 0 };
  } catch (err) {
    throw toApiError(err);
  }
};

export const saveSkin = async (skin: Pick<NewSkin, "id" | "title" | "imageUrl" | "pageUrl">): Promise<void> => {
  try {
    await SaveSkin({
      market_hash_name: skin.id,
      display_name: skin.title,
      icon_url: skin.imageUrl,
      page_url: skin.pageUrl,
    });
  } catch (err) {
    throw toApiError(err);
  }
};
