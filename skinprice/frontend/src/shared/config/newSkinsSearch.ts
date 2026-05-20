import type { NewSkinsSearchParams, NewSkinsSearchSortColumn } from "../../entities/skin/model/types";

export const DEFAULT_NEW_SKINS_SEARCH_PARAMS: NewSkinsSearchParams = {
  query: "",
  sortColumn: "popular",
  sortDir: "desc",
  filters: {
    priceMin: "",
    priceMax: "",
    searchDescriptions: false,
    type: [],
    weapon: [],
    rarity: [],
    exterior: [],
    itemSet: [],
    proPlayer: [],
    stickerCapsule: [],
    tournamentTeam: [],
  },
};

export const NEW_SKINS_SORT_OPTIONS: Array<{ value: NewSkinsSearchSortColumn; label: string }> = [
  { value: "popular", label: "Популярность" },
  { value: "price", label: "Цена" },
  { value: "name", label: "Название" },
  { value: "quantity", label: "Количество лотов" },
];

export const NEW_SKINS_SORT_DIR_OPTIONS = [
  { value: "desc", label: "По убыванию" },
  { value: "asc", label: "По возрастанию" },
] as const;

export const NEW_SKINS_FILTER_GROUPS = [
  {
    key: "type",
    label: "Тип предмета",
    options: [
      { value: "tag_CSGO_Type_Rifle", label: "Винтовка" },
      { value: "tag_CSGO_Type_Pistol", label: "Пистолет" },
      { value: "tag_CSGO_Type_SniperRifle", label: "Снайперская" },
      { value: "tag_CSGO_Type_SMG", label: "ПП" },
      { value: "tag_CSGO_Type_Shotgun", label: "Дробовик" },
      { value: "tag_CSGO_Type_Machinegun", label: "Пулемёт" },
      { value: "tag_CSGO_Type_Knife", label: "Нож" },
      { value: "tag_CSGO_Type_WeaponCase", label: "Кейс" },
      { value: "tag_CSGO_Type_WeaponCaseKey", label: "Ключ" },
    ],
  },
  {
    key: "weapon",
    label: "Оружие / предмет",
    options: [
      { value: "tag_weapon_ak47", label: "AK-47" },
      { value: "tag_weapon_awp", label: "AWP" },
      { value: "tag_weapon_m4a1", label: "M4A4" },
      { value: "tag_weapon_m4a1_silencer", label: "M4A1-S" },
      { value: "tag_weapon_deagle", label: "Desert Eagle" },
      { value: "tag_weapon_glock", label: "Glock-18" },
      { value: "tag_weapon_usp_silencer", label: "USP-S" },
      { value: "tag_weapon_knife", label: "Knife" },
    ],
  },
  {
    key: "exterior",
    label: "Износ",
    options: [
      { value: "tag_WearCategory0", label: "Прямо с завода" },
      { value: "tag_WearCategory1", label: "Немного поношенное" },
      { value: "tag_WearCategory2", label: "После полевых" },
      { value: "tag_WearCategory3", label: "Поношенное" },
      { value: "tag_WearCategory4", label: "Закалённое в боях" },
    ],
  },
] as const;
