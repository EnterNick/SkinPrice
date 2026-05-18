import React from "react";
import { UI_TEXT } from "../../shared/config/uiText";

type NewSkinsSearchPanelProps = {
  value: string;
  activeSource: "steam" | "lisskins";
  errorText?: string | null;
  helperText: string;
  disabled?: boolean;
  onChange: (value: string) => void;
  onChangeSource: (source: "steam" | "lisskins") => void;
  onSearch: () => Promise<void> | void;
};

export const NewSkinsSearchPanel: React.FC<NewSkinsSearchPanelProps> = ({
  value,
  activeSource,
  errorText,
  helperText,
  disabled,
  onChange,
  onChangeSource,
  onSearch,
}) => (
  <form
    className="search-panel"
    onSubmit={(event) => {
      event.preventDefault();
      void onSearch();
    }}
  >
    <div className="source-tabs" role="tablist" aria-label={UI_TEXT.sourceTabsAriaLabel}>
      <button
        className={`source-tab${activeSource === "steam" ? " source-tab-active" : ""}`}
        type="button"
        role="tab"
        aria-selected={activeSource === "steam"}
        onClick={() => onChangeSource("steam")}
      >
        {UI_TEXT.sourceSteam}
      </button>
      <button
        className={`source-tab${activeSource === "lisskins" ? " source-tab-active" : ""}`}
        type="button"
        role="tab"
        aria-selected={activeSource === "lisskins"}
        onClick={() => onChangeSource("lisskins")}
      >
        {UI_TEXT.sourceLisSkins}
      </button>
    </div>
    <input
      className="search-input"
      type="text"
      value={value}
      placeholder={UI_TEXT.searchPlaceholder}
      onChange={(event) => onChange(event.target.value)}
      disabled={disabled}
    />
    <button className="toolbar-button toolbar-button-primary" type="submit" disabled={disabled}>
      {UI_TEXT.searchAction}
    </button>
    <div className={`search-query-label${errorText ? " search-query-label-error" : ""}`}>{errorText || helperText}</div>
  </form>
);
