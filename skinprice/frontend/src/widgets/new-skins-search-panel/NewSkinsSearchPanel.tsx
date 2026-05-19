import React from "react";
import { UI_TEXT } from "../../shared/config/uiText";

type NewSkinsSearchPanelProps = {
  value: string;
  errorText?: string | null;
  disabled?: boolean;
  onChange: (value: string) => void;
};

export const NewSkinsSearchPanel: React.FC<NewSkinsSearchPanelProps> = ({
  value,
  errorText,
  disabled,
  onChange,
}) => (
  <div className="search-panel">
    <input
      className="search-input"
      id="new-skins-search-input"
      type="text"
      value={value}
      placeholder={UI_TEXT.searchPlaceholder}
      onChange={(event) => onChange(event.target.value)}
      disabled={disabled}
      aria-label={UI_TEXT.searchPlaceholder}
    />
    {errorText && <div className="search-query-label search-query-label-error">{errorText}</div>}
  </div>
);
