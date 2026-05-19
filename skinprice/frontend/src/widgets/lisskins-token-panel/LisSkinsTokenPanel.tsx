import React from "react";
import { UI_TEXT } from "../../shared/config/uiText";
import { openExternal } from "../../shared/lib/browser/openExternal";

type LisSkinsTokenPanelProps = {
  visible: boolean;
  configured: boolean;
  required: boolean;
  saving?: boolean;
  resetting?: boolean;
  error?: string | null;
  saved?: boolean;
  value: string;
  onChange: (value: string) => void;
  onSave: () => Promise<void> | void;
  onReset?: () => Promise<void> | void;
};

export const LisSkinsTokenPanel: React.FC<LisSkinsTokenPanelProps> = ({
  visible,
  configured,
  required,
  saving,
  resetting,
  error,
  saved,
  value,
  onChange,
  onSave,
  onReset,
}) => {
  if (!visible) return null;

  const statusText = required
    ? UI_TEXT.lisSkinsTokenRequired
    : saved
      ? UI_TEXT.lisSkinsTokenSaved
      : configured
        ? UI_TEXT.lisSkinsTokenConfigured
        : UI_TEXT.lisSkinsTokenHelp;

  const statusTone = error ? "token-panel-status-error" : required ? "token-panel-status-warning" : "token-panel-status-success";

  return (
    <section className="token-panel" aria-labelledby="lisskins-token-title">
      <div className="token-panel-header">
        <div>
          <h2 className="token-panel-title" id="lisskins-token-title">
            {UI_TEXT.lisSkinsTokenTitle}
          </h2>
          <p className="token-panel-description">{UI_TEXT.lisSkinsTokenHelp}</p>
        </div>
        <button className="toolbar-button" type="button" onClick={() => openExternal("https://lis-skins.com/ru/profile/api/")}>
          {UI_TEXT.lisSkinsTokenLink}
        </button>
      </div>
      <div className="token-panel-body">
        <label className="token-panel-field" htmlFor="lisskins-token-input">
          <span className="token-panel-label">{UI_TEXT.lisSkinsTokenLabel}</span>
          <input
            className="search-input"
            id="lisskins-token-input"
            type="password"
            autoComplete="off"
            spellCheck={false}
            placeholder={UI_TEXT.lisSkinsTokenPlaceholder}
            value={value}
            onChange={(event) => onChange(event.target.value)}
            disabled={saving}
            aria-describedby="lisskins-token-status"
          />
        </label>
        <div className="token-panel-actions">
          <button
            className="toolbar-button toolbar-button-primary"
            type="button"
            onClick={() => void onSave()}
            disabled={saving || resetting || value.trim().length === 0}
          >
            {saving ? UI_TEXT.lisSkinsTokenSaving : UI_TEXT.lisSkinsTokenSave}
          </button>
          {onReset && (
            <button
              className="toolbar-button toolbar-button-danger"
              type="button"
              onClick={() => void onReset()}
              disabled={saving || resetting || !configured}
            >
              {resetting ? UI_TEXT.lisSkinsTokenResetting : UI_TEXT.lisSkinsTokenReset}
            </button>
          )}
        </div>
      </div>
      <div className={`token-panel-status ${statusTone}`} id="lisskins-token-status" role={error ? "alert" : "status"}>
        {error || statusText}
      </div>
    </section>
  );
};
