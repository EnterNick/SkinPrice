import React from "react";
import type { NewSkinsSearchParams } from "../../entities/skin/model/types";
import {
  DEFAULT_NEW_SKINS_SEARCH_PARAMS,
  NEW_SKINS_FILTER_GROUPS,
  NEW_SKINS_SORT_DIR_OPTIONS,
  NEW_SKINS_SORT_OPTIONS,
} from "../../shared/config/newSkinsSearch";
import { UI_TEXT } from "../../shared/config/uiText";

type NewSkinsSearchPanelProps = {
  value: NewSkinsSearchParams;
  errorText?: string | null;
  disabled?: boolean;
  onChange: (value: NewSkinsSearchParams) => void;
  onReset: () => void;
};

export const NewSkinsSearchPanel: React.FC<NewSkinsSearchPanelProps> = ({
  value,
  errorText,
  disabled,
  onChange,
  onReset,
}) => {
  const toggleFilterValue = (key: keyof NewSkinsSearchParams["filters"], optionValue: string) => {
    const currentValues = value.filters[key];
    if (!Array.isArray(currentValues)) return;

    const nextValues = currentValues.includes(optionValue)
      ? currentValues.filter((item) => item !== optionValue)
      : [...currentValues, optionValue];

    onChange({
      ...value,
      filters: {
        ...value.filters,
        [key]: nextValues,
      },
    });
  };

  const hasCustomFilters = JSON.stringify(value) !== JSON.stringify(DEFAULT_NEW_SKINS_SEARCH_PARAMS);

  return (
    <div className="search-panel">
      <div className="search-panel-main search-panel-main-extended">
        <label className="token-panel-field" htmlFor="new-skins-search-input">
          <span className="token-panel-label">{UI_TEXT.searchQueryLabel}</span>
          <input
            className="search-input"
            id="new-skins-search-input"
            type="text"
            value={value.query}
            placeholder={UI_TEXT.searchPlaceholder}
            onChange={(event) => onChange({ ...value, query: event.target.value })}
            disabled={disabled}
            aria-label={UI_TEXT.searchPlaceholder}
          />
        </label>

        <label className="token-panel-field" htmlFor="new-skins-sort-column">
          <span className="token-panel-label">{UI_TEXT.searchSortLabel}</span>
          <select
            id="new-skins-sort-column"
            className="toolbar-select"
            value={value.sortColumn}
            disabled={disabled}
            onChange={(event) => onChange({ ...value, sortColumn: event.target.value as NewSkinsSearchParams["sortColumn"] })}
          >
            {NEW_SKINS_SORT_OPTIONS.map((option) => (
              <option key={option.value} value={option.value}>
                {option.label}
              </option>
            ))}
          </select>
        </label>

        <label className="token-panel-field" htmlFor="new-skins-sort-dir">
          <span className="token-panel-label">{UI_TEXT.searchSortDirLabel}</span>
          <select
            id="new-skins-sort-dir"
            className="toolbar-select"
            value={value.sortDir}
            disabled={disabled}
            onChange={(event) => onChange({ ...value, sortDir: event.target.value as NewSkinsSearchParams["sortDir"] })}
          >
            {NEW_SKINS_SORT_DIR_OPTIONS.map((option) => (
              <option key={option.value} value={option.value}>
                {option.label}
              </option>
            ))}
          </select>
        </label>

        <label className="token-panel-field" htmlFor="new-skins-price-min">
          <span className="token-panel-label">{UI_TEXT.searchPriceMinLabel}</span>
          <input
            id="new-skins-price-min"
            className="search-input"
            type="text"
            inputMode="numeric"
            value={value.filters.priceMin}
            placeholder={UI_TEXT.searchPriceMinPlaceholder}
            disabled={disabled}
            onChange={(event) =>
              onChange({
                ...value,
                filters: {
                  ...value.filters,
                  priceMin: event.target.value,
                },
              })
            }
          />
        </label>

        <label className="token-panel-field" htmlFor="new-skins-price-max">
          <span className="token-panel-label">{UI_TEXT.searchPriceMaxLabel}</span>
          <input
            id="new-skins-price-max"
            className="search-input"
            type="text"
            inputMode="numeric"
            value={value.filters.priceMax}
            placeholder={UI_TEXT.searchPriceMaxPlaceholder}
            disabled={disabled}
            onChange={(event) =>
              onChange({
                ...value,
                filters: {
                  ...value.filters,
                  priceMax: event.target.value,
                },
              })
            }
          />
        </label>
      </div>

      <div className="search-panel-toolbar">
        <label className="search-toggle">
          <input
            type="checkbox"
            checked={value.filters.searchDescriptions}
            disabled={disabled}
            onChange={(event) =>
              onChange({
                ...value,
                filters: {
                  ...value.filters,
                  searchDescriptions: event.target.checked,
                },
              })
            }
          />
          <span>{UI_TEXT.searchDescriptionsLabel}</span>
        </label>
        {hasCustomFilters ? (
          <button className="toolbar-button toolbar-button-secondary search-reset-button" type="button" disabled={disabled} onClick={onReset}>
            {UI_TEXT.searchReset}
          </button>
        ) : null}
      </div>

      <div className="search-filter-groups">
        {NEW_SKINS_FILTER_GROUPS.map((group) => (
          <section key={group.key} className="search-filter-group" aria-label={group.label}>
            <h3 className="search-filter-title">{group.label}</h3>
            <div className="search-filter-options">
              {group.options.map((option) => {
                const selected = value.filters[group.key].includes(option.value);

                return (
                  <button
                    key={option.value}
                    className={`search-filter-chip ${selected ? "search-filter-chip-selected" : ""}`}
                    type="button"
                    disabled={disabled}
                    onClick={() => toggleFilterValue(group.key, option.value)}
                  >
                    {option.label}
                  </button>
                );
              })}
            </div>
          </section>
        ))}
      </div>

      {errorText && <div className="search-query-label search-query-label-error">{errorText}</div>}
    </div>
  );
};
