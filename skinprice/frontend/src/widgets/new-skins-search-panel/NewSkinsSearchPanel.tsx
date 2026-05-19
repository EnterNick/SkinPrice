import React from "react";
import { UI_TEXT } from "../../shared/config/uiText";
import { openExternal } from "../../shared/lib/browser/openExternal";

type NewSkinsSearchPanelProps = {
  value: string;
  activeSource: "steam" | "lisskins";
  errorText?: string | null;
  helperText: string;
  disabled?: boolean;
  tokenRequired?: boolean;
  tokenSaving?: boolean;
  tokenError?: string | null;
  tokenSaved?: boolean;
  tokenValue?: string;
  onChange: (value: string) => void;
  onTokenChange: (value: string) => void;
  onSaveToken: () => Promise<void> | void;
  onChangeSource: (source: "steam" | "lisskins") => void;
  onSearch: () => Promise<void> | void;
};

export const NewSkinsSearchPanel: React.FC<NewSkinsSearchPanelProps> = ({
  value,
  activeSource,
  errorText,
  helperText,
  disabled,
  tokenRequired,
  tokenSaving,
  tokenError,
  tokenSaved,
  tokenValue,
  onChange,
  onTokenChange,
  onSaveToken,
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
    {tokenRequired && (
      <div className="search-query-label">
        <div>{UI_TEXT.lisSkinsTokenRequired}</div>
        <div>{UI_TEXT.lisSkinsTokenHelp}</div>
        <button className="toolbar-button" type="button" onClick={() => openExternal("https://lis-skins.com/ru/profile/api/") }>
          {UI_TEXT.lisSkinsTokenLink}
        </button>
        <input
          className="search-input"
          type="text"
          placeholder={UI_TEXT.lisSkinsTokenPlaceholder}
          value={tokenValue}
          onChange={(event) => onTokenChange(event.target.value)}
          disabled={tokenSaving}
        />
        <button className="toolbar-button toolbar-button-primary" type="button" onClick={() => void onSaveToken()} disabled={tokenSaving}>
          {tokenSaving ? UI_TEXT.lisSkinsTokenSaving : UI_TEXT.lisSkinsTokenSave}
        </button>
        {tokenError && <div className="search-query-label search-query-label-error">{tokenError}</div>}
        {tokenSaved && <div className="search-query-label">{UI_TEXT.lisSkinsTokenSaved}</div>}
      </div>
    )}
    <div className={`search-query-label${errorText ? " search-query-label-error" : ""}`}>{errorText || helperText}</div>
  </form>
);
