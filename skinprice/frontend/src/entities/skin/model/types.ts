export type Skin = {
  id: string;
  name: string;
  title: string;
  imageUrl: string;
  pageUrl: string;
  priceText: string;
  priceCents?: number;
  currency?: string;
  sellListings?: number;
  updatedAt?: string;
};

export type SavedSkin = Skin;
export type NewSkin = Skin;

export type PriceUpdateResult = {
  updated: number;
};

export type ApiErrorCode =
  | "NETWORK_ERROR"
  | "NOT_FOUND"
  | "VALIDATION_ERROR"
  | "RATE_LIMIT"
  | "SERVER_ERROR"
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
