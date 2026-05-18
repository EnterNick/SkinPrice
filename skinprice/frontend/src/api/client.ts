import { GetSavedSkins, SaveSkin, SearchNewSkins, UpdateAllSavedSkinsPrices, UpdateSavedSkinPrice } from "../wailsjs/go/main/App";
import type { skins } from "../wailsjs/go/models";
import { toApiError } from "./errors";
import type { NewSkin, PaginatedResult, PriceUpdateResult, SavedSkin } from "./models";

const DEFAULT_LIMIT = 50;
const DEFAULT_OFFSET = 0;
const DEFAULT_CURRENCY = "1";

const mapSavedSkin = (item: skins.SavedSkinItem): SavedSkin => ({
  id: item.market_hash_name,
  name: item.market_hash_name,
  title: item.display_name,
  imageUrl: item.icon_url,
  pageUrl: item.page_url,
  priceText: item.price_text,
  currency: item.currency,
});

const mapNewSkin = (item: skins.NewSkinItem): NewSkin => ({
  id: item.market_hash_name,
  name: item.market_hash_name,
  title: item.display_name,
  imageUrl: item.icon_url,
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

export const getNewSkins = async (query?: string): Promise<PaginatedResult<NewSkin>> => {
  try {
    const response = await SearchNewSkins({ market_hash_name: query, limit: DEFAULT_LIMIT, offset: DEFAULT_OFFSET });
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

export const updateSkinPrice = async (skinId: string): Promise<PriceUpdateResult> => {
  try {
    await UpdateSavedSkinPrice({ market_hash_name: skinId, currency: DEFAULT_CURRENCY });
    return { updated: 1 };
  } catch (err) {
    throw toApiError(err);
  }
};

export const updateAllSkinPrices = async (): Promise<PriceUpdateResult> => {
  try {
    await UpdateAllSavedSkinsPrices({ currency: DEFAULT_CURRENCY });
    const refreshed = await getSavedSkins();
    return { updated: refreshed.items.length };
  } catch (err) {
    throw toApiError(err);
  }
};

export const saveSkin = async (skinId: string): Promise<void> => {
  try {
    const candidates = await getNewSkins(skinId);
    const skin = candidates.items.find((item) => item.id === skinId);

    if (!skin) {
      throw { code: "NOT_FOUND", message: `Скин ${skinId} не найден` };
    }

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
