import React from "react";
import { UI_TEXT } from "../../shared/config/uiText";

type NewSkinsSearchPanelProps = {
  searchLabel: string;
  fallbackLabel: string;
  onSearch: () => Promise<void> | void;
};

export const NewSkinsSearchPanel: React.FC<NewSkinsSearchPanelProps> = ({ searchLabel, fallbackLabel, onSearch }) => (
  <div className="search-panel">
    <button className="toolbar-button toolbar-button-primary" type="button" onClick={() => void onSearch()}>
      {UI_TEXT.searchAction}
    </button>
    <div className="search-query-label">{searchLabel || fallbackLabel}</div>
  </div>
);
