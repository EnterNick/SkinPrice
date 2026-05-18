import React from "react";
import { UI_TEXT } from "../../shared/config/uiText";

type NewSkinsSearchPanelProps = {
  searchLabel: string;
  fallbackLabel: string;
  activeSource: "steam" | "lisskins";
  onChangeSource: (source: "steam" | "lisskins") => void;
  onSearch: () => Promise<void> | void;
};

export const NewSkinsSearchPanel: React.FC<NewSkinsSearchPanelProps> = ({
  searchLabel,
  fallbackLabel,
  activeSource,
  onChangeSource,
  onSearch,
}) => (
  <div className="search-panel">
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
    <button className="toolbar-button toolbar-button-primary" type="button" onClick={() => void onSearch()}>
      {UI_TEXT.searchAction}
    </button>
    <div className="search-query-label">{searchLabel || fallbackLabel}</div>
  </div>
);
