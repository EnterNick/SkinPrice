export type SavedSkinCurrency = "1" | "3" | "5";

export type CurrencyOption = {
  value: SavedSkinCurrency;
  label: "USD" | "EUR" | "RUB";
};

export type Skin = {
  id: string;
  name: string;
  title: string;
  imageUrl: string;
  pageUrl: string;
  priceText: string;
  priceCents?: number;
  currency?: SavedSkinCurrency;
  sellListings?: number;
  updatedAt?: string;
};

export type SavedSkin = Skin;
export type NewSkin = Skin;

export type PriceUpdateResult = {
  marketHashName?: string;
  priceText?: string;
  currency?: SavedSkinCurrency;
  updatedAt?: string;
  updated: number;
};

export type BulkPriceUpdateFailure = {
  marketHashName: string;
  message: string;
};

export type BulkPriceUpdateResult = {
  updated: number;
  failed: number;
  failures: BulkPriceUpdateFailure[];
};

export type SaveSkinResult = {
  created: boolean;
};

export type ApiErrorCode =
  | "invalid_argument"
  | "not_found"
  | "already_exists"
  | "conflict"
  | "unavailable"
  | "timeout"
  | "external"
  | "internal"
  | "unknown"
  | "UNKNOWN_ERROR";

export type ApiError = {
  code: ApiErrorCode;
  message: string;
  details?: unknown;
};

export type PaginatedResult<T> = {
  items: T[];
  total: number;
  limit: number;
  offset: number;
};
