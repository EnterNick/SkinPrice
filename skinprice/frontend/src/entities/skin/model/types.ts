export type SavedSkinCurrency = "1" | "3" | "5";
export type SavedSkinsViewMode = "table" | "cards";
export type FontFamilyOptionValue = "inter" | "system" | "nunito" | "roboto" | "ibm-plex-sans" | "manrope" | "monocraft";
export type SavedSkinsSortColumn = "title" | "steamPrice" | "lisSkinsPrice" | "csTmPrice";
export type SavedSkinsSortDirection = "desc" | "asc";
export type SavedSkinsSortState = {
  column: SavedSkinsSortColumn;
  direction: SavedSkinsSortDirection;
} | null;

export type CurrencyOption = {
  value: SavedSkinCurrency;
  label: "USD" | "EUR" | "RUB";
};

export type FontFamilyOption = {
  value: FontFamilyOptionValue;
  label: "Inter" | "System" | "Nunito" | "Roboto" | "IBM Plex Sans" | "Manrope" | "Monocraft";
};

export type NewSkinsSearchSortColumn = "popular" | "price" | "name" | "quantity";
export type NewSkinsSearchSortDir = "asc" | "desc";

export type NewSkinsSearchFilters = {
  priceMin: string;
  priceMax: string;
  searchDescriptions: boolean;
  type: string[];
  weapon: string[];
  rarity: string[];
  exterior: string[];
  itemSet: string[];
  proPlayer: string[];
  stickerCapsule: string[];
  tournamentTeam: string[];
};

export type NewSkinsSearchParams = {
  query: string;
  sortColumn: NewSkinsSearchSortColumn;
  sortDir: NewSkinsSearchSortDir;
  filters: NewSkinsSearchFilters;
};

export type Skin = {
  id: string;
  name: string;
  title: string;
  nameColor?: string;
  imageUrl: string;
  steamPageUrl: string;
  lisSkinsPageUrl: string;
  csTmPageUrl: string;
  priceText: string;
  steamPriceText?: string;
  lisSkinsPriceText?: string;
  csTmPriceText?: string;
  priceCents?: number;
  currency?: SavedSkinCurrency;
  sellListings?: number;
  updatedAt?: string;
  steamUpdatedAt?: string;
  lisSkinsUpdatedAt?: string;
  csTmUpdatedAt?: string;
};

export type SavedSkin = Skin;
export type NewSkin = Skin;

export type PriceUpdateResult = {
  marketHashName?: string;
  steamPageUrl?: string;
  steamPriceText?: string;
  steamUpdatedAt?: string;
  lisSkinsPageUrl?: string;
  lisSkinsPriceText?: string;
  lisSkinsUpdatedAt?: string;
  csTmPageUrl?: string;
  csTmPriceText?: string;
  csTmUpdatedAt?: string;
  currency?: SavedSkinCurrency;
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
  nextCursor?: string;
};
