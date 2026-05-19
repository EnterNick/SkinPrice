import React, { useEffect, useState } from "react";
import { useNavigate } from "react-router-dom";
import { ROUTES } from "../../../app/router/routes";
import { clearLisSkinsToken, getAppSettings, hasLisSkinsToken, saveAppSettings, saveLisSkinsToken } from "../../../entities/skin/api/skinApi";
import type { SavedSkinCurrency } from "../../../entities/skin/model/types";
import {
  CURRENCY_OPTIONS,
  DEFAULT_AUTO_REFRESH_INTERVAL_SECONDS,
  DEFAULT_CURRENCY,
  normalizeAutoRefreshIntervalSeconds,
} from "../../../shared/config/settings";
import { UI_TEXT } from "../../../shared/config/uiText";
import { formatErrorMessage } from "../../../shared/lib/error/formatErrorMessage";
import { PageHeader } from "../../../shared/ui/page-header/PageHeader";
import { ErrorState, LoadingState, ToastAlert } from "../../../shared/ui/states/States";
import { LisSkinsTokenPanel } from "../../../widgets/lisskins-token-panel/LisSkinsTokenPanel";

export const SettingsPage: React.FC = () => {
  const navigate = useNavigate();
  const [currency, setCurrency] = useState(DEFAULT_CURRENCY);
  const [autoRefreshInput, setAutoRefreshInput] = useState(String(DEFAULT_AUTO_REFRESH_INTERVAL_SECONDS));
  const [isTokenConfigured, setIsTokenConfigured] = useState(false);
  const [tokenValue, setTokenValue] = useState("");
  const [loading, setLoading] = useState(true);
  const [savingToken, setSavingToken] = useState(false);
  const [resettingToken, setResettingToken] = useState(false);
  const [tokenSaved, setTokenSaved] = useState(false);
  const [tokenError, setTokenError] = useState<string | null>(null);
  const [pageError, setPageError] = useState<string | null>(null);
  const [notice, setNotice] = useState<{ type: "success" | "warning" | "error"; text: string } | null>(null);

  useEffect(() => {
    let active = true;

    const loadPageState = async () => {
      try {
        const [settings, hasToken] = await Promise.all([getAppSettings(), hasLisSkinsToken()]);
        if (!active) return;
        setCurrency(settings.currency);
        setAutoRefreshInput(String(settings.autoRefreshIntervalSeconds));
        setIsTokenConfigured(hasToken);
        setPageError(null);
      } catch (err: unknown) {
        if (!active) return;
        setPageError(formatErrorMessage("Не удалось загрузить настройки.", err));
      } finally {
        if (active) {
          setLoading(false);
        }
      }
    };

    void loadPageState();

    return () => {
      active = false;
    };
  }, []);

  const persistSettings = async (nextCurrency: SavedSkinCurrency, nextInterval: number) => {
    const saved = await saveAppSettings({
      currency: nextCurrency,
      autoRefreshIntervalSeconds: nextInterval,
    });
    setCurrency(saved.currency);
    setAutoRefreshInput(String(saved.autoRefreshIntervalSeconds));
    setNotice({ type: "success", text: UI_TEXT.settingsAutoSaved });
  };

  const onCurrencyChange = async (nextCurrency: SavedSkinCurrency) => {
    setNotice(null);
    try {
      await persistSettings(nextCurrency, normalizeAutoRefreshIntervalSeconds(autoRefreshInput));
    } catch (err: unknown) {
      setNotice({ type: "error", text: formatErrorMessage("Не удалось сохранить настройки.", err) });
    }
  };

  const onAutoRefreshBlur = () => {
    const normalized = normalizeAutoRefreshIntervalSeconds(autoRefreshInput);
    void (async () => {
      setNotice(null);
      try {
        await persistSettings(currency, normalized);
      } catch (err: unknown) {
        setNotice({ type: "error", text: formatErrorMessage("Не удалось сохранить настройки.", err) });
      }
    })();
  };

  const onSaveToken = async () => {
    setSavingToken(true);
    setTokenError(null);
    setNotice(null);
    try {
      await saveLisSkinsToken(tokenValue);
      setIsTokenConfigured(true);
      setTokenSaved(true);
      setTokenValue("");
      setNotice({ type: "success", text: UI_TEXT.lisSkinsTokenSaved });
    } catch (err: unknown) {
      setTokenSaved(false);
      setTokenError(formatErrorMessage(UI_TEXT.errLisSkinsTokenSave, err));
    } finally {
      setSavingToken(false);
    }
  };

  const onResetToken = async () => {
    setResettingToken(true);
    setTokenError(null);
    setNotice(null);
    try {
      await clearLisSkinsToken();
      setIsTokenConfigured(false);
      setTokenSaved(false);
      setTokenValue("");
      setNotice({ type: "success", text: UI_TEXT.lisSkinsTokenResetDone });
    } catch (err: unknown) {
      setTokenSaved(false);
      setTokenError(formatErrorMessage(UI_TEXT.errLisSkinsTokenReset, err));
    } finally {
      setResettingToken(false);
    }
  };

  if (loading) return <LoadingState text="Загрузка настроек..." />;
  if (pageError) return <ErrorState text={pageError} />;

  return (
    <div className="settings-page">
      <PageHeader
        sectionLabel={UI_TEXT.settingsEyebrow}
        title={UI_TEXT.settingsTitle}
        actions={
          <div className="toolbar-group">
            <button className="toolbar-button" type="button" onClick={() => navigate(ROUTES.home)}>
              {UI_TEXT.ctaToHome}
            </button>
          </div>
        }
      />
      {notice && <ToastAlert type={notice.type} text={notice.text} />}
      <LisSkinsTokenPanel
        visible
        configured={isTokenConfigured}
        required={!isTokenConfigured}
        saving={savingToken}
        resetting={resettingToken}
        error={tokenError}
        saved={tokenSaved}
        value={tokenValue}
        onChange={(nextValue) => {
          setTokenValue(nextValue);
          setTokenSaved(false);
          setTokenError(null);
        }}
        onSave={onSaveToken}
        onReset={onResetToken}
      />
      <section className="settings-grid">
        <article className="settings-card">
          <h2 className="settings-card-title">{UI_TEXT.settingsCurrencyTitle}</h2>
          <p className="settings-card-description">{UI_TEXT.settingsCurrencyHelp}</p>
          <label className="token-panel-field" htmlFor="settings-currency-select">
            <span className="token-panel-label">{UI_TEXT.settingsCurrencyTitle}</span>
            <select
              id="settings-currency-select"
              className="toolbar-select settings-select"
              value={currency}
              onChange={(event) => void onCurrencyChange(event.target.value as SavedSkinCurrency)}
            >
              {CURRENCY_OPTIONS.map((option) => (
                <option key={option.value} value={option.value}>
                  {option.label}
                </option>
              ))}
            </select>
          </label>
        </article>
        <article className="settings-card">
          <h2 className="settings-card-title">{UI_TEXT.settingsRefreshTitle}</h2>
          <p className="settings-card-description">{UI_TEXT.settingsRefreshHelp}</p>
          <label className="token-panel-field" htmlFor="settings-refresh-interval">
            <span className="token-panel-label">{UI_TEXT.settingsRefreshLabel}</span>
            <input
              id="settings-refresh-interval"
              className="search-input settings-number-input"
              type="number"
              min={5}
              step={5}
              value={autoRefreshInput}
              onChange={(event) => setAutoRefreshInput(event.target.value)}
              onBlur={onAutoRefreshBlur}
            />
          </label>
          <p className="toolbar-hint">{UI_TEXT.settingsRefreshHint}</p>
        </article>
      </section>
    </div>
  );
};
