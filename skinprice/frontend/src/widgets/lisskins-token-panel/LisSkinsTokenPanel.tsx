import React from "react";
import { UI_TEXT } from "../../shared/config/uiText";
import { openExternal } from "../../shared/lib/browser/openExternal";

type LisSkinsTokenStatus = "missing" | "saved" | "verified" | "authError";

type LisSkinsTokenPanelProps = {
  visible: boolean;
  configured: boolean;
  saving?: boolean;
  resetting?: boolean;
  checking?: boolean;
  error?: string | null;
  value: string;
  canSave: boolean;
  status: LisSkinsTokenStatus;
  revealToken: boolean;
  onToggleReveal: () => void;
  onChange: (value: string) => void;
  onSave: () => Promise<void> | void;
  onCheck: () => Promise<void> | void;
  onReset?: () => Promise<void> | void;
};

export const LisSkinsTokenPanel: React.FC<LisSkinsTokenPanelProps> = ({
  visible,
  configured,
  saving,
  resetting,
  checking,
  error,
  value,
  canSave,
  status,
  revealToken,
  onToggleReveal,
  onChange,
  onSave,
  onCheck,
  onReset,
}) => {
  if (!visible) return null;

  const statusText =
    status === "authError"
      ? UI_TEXT.lisSkinsTokenStatusAuthError
      : status === "verified"
        ? UI_TEXT.lisSkinsTokenStatusVerified
        : status === "saved"
          ? UI_TEXT.lisSkinsTokenStatusSaved
          : UI_TEXT.lisSkinsTokenStatusMissing;

  const statusTone = status === "authError" ? "token-panel-status-error" : status === "missing" ? "token-panel-status-warning" : "token-panel-status-success";

  return (
    <section className="token-panel" aria-labelledby="lisskins-token-title">
      <div className="token-panel-section">
        <h2 className="token-panel-title" id="lisskins-token-title">{UI_TEXT.lisSkinsTokenTitle}</h2>
        <h3 className="token-panel-subtitle">{UI_TEXT.lisSkinsTokenSectionAccess}</h3>
        <p className="token-panel-description">{UI_TEXT.lisSkinsTokenHelp}</p>
      </div>

      <div className="token-panel-divider" />

      <div className="token-panel-section token-panel-body-stack">
        <h3 className="token-panel-subtitle">{UI_TEXT.lisSkinsTokenSectionInput}</h3>
        <label className="token-panel-field" htmlFor="lisskins-token-input">
          <span className="token-panel-label">{UI_TEXT.lisSkinsTokenLabel}</span>
          <input
            className="search-input"
            id="lisskins-token-input"
            type={revealToken ? "text" : "password"}
            autoComplete="off"
            spellCheck={false}
            placeholder={UI_TEXT.lisSkinsTokenPlaceholder}
            value={value}
            onChange={(event) => onChange(event.target.value)}
            disabled={saving || checking}
            aria-describedby="lisskins-token-status"
          />
        </label>
        <div className="token-panel-actions">
          <button className="toolbar-button toolbar-button-secondary" type="button" onClick={onToggleReveal} disabled={saving || checking}>
            {revealToken ? UI_TEXT.lisSkinsTokenHide : UI_TEXT.lisSkinsTokenShow}
          </button>
          <button className="toolbar-button toolbar-button-secondary" type="button" onClick={() => void onCheck()} disabled={checking || saving || value.trim().length === 0}>
            {checking ? UI_TEXT.lisSkinsTokenChecking : UI_TEXT.lisSkinsTokenCheck}
          </button>
          <button className="toolbar-button toolbar-button-primary" type="button" onClick={() => void onSave()} disabled={saving || checking || resetting || !canSave}>
            {saving ? UI_TEXT.lisSkinsTokenSaving : UI_TEXT.lisSkinsTokenSave}
          </button>
          {onReset && (
            <button className="toolbar-button toolbar-button-danger-secondary" type="button" onClick={() => void onReset()} disabled={saving || resetting || checking || !configured}>
              {resetting ? UI_TEXT.lisSkinsTokenResetting : UI_TEXT.lisSkinsTokenReset}
            </button>
          )}
        </div>
      </div>

      <div className="token-panel-divider" />

      <div className="token-panel-section">
        <h3 className="token-panel-subtitle">{UI_TEXT.lisSkinsTokenSectionLink}</h3>
        <button className="toolbar-button" type="button" onClick={() => openExternal("https://lis-skins.com/ru/profile/api/")}>
          {UI_TEXT.lisSkinsTokenLink}
        </button>
      </div>

      <div className={`token-panel-status ${statusTone}`} id="lisskins-token-status" role={error ? "alert" : "status"}>
        {error || statusText}
      </div>
    </section>
  );
};
