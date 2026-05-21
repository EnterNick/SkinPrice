import React, { useEffect, useMemo, useState } from "react";
import { useNavigate } from "react-router-dom";
import { ROUTES } from "../../../app/router/routes";
import { clearLisSkinsToken, getAppSettings, hasLisSkinsToken, saveAppSettings, saveLisSkinsToken } from "../../../entities/skin/api/skinApi";
import type { SavedSkinCurrency } from "../../../entities/skin/model/types";
import {
  CURRENCY_OPTIONS,
  DEFAULT_AUTO_REFRESH_ENABLED,
  DEFAULT_AUTO_REFRESH_INTERVAL_SECONDS,
  DEFAULT_CURRENCY,
  DEFAULT_FONT_FAMILY,
  DEFAULT_FONT_SIZE_PX,
  DEFAULT_SAVED_SKINS_VIEW_MODE,
  FONT_FAMILY_OPTIONS,
  MAX_FONT_SIZE_PX,
  MIN_FONT_SIZE_PX,
  normalizeAutoRefreshEnabled,
  normalizeAutoRefreshIntervalSeconds,
  normalizeFontFamily,
  normalizeFontSizePx,
} from "../../../shared/config/settings";
import { UI_TEXT } from "../../../shared/config/uiText";
import { formatErrorMessage } from "../../../shared/lib/error/formatErrorMessage";
import { PageHeader } from "../../../shared/ui/page-header/PageHeader";
import { ErrorState, LoadingState, ToastAlert } from "../../../shared/ui/states/States";
import { LisSkinsTokenPanel } from "../../../widgets/lisskins-token-panel/LisSkinsTokenPanel";

type TokenStatus = "missing" | "saved" | "verified" | "authError";

export const SettingsPage: React.FC = () => {
  const navigate = useNavigate();
  const [currency, setCurrency] = useState(DEFAULT_CURRENCY);
  const [autoRefreshEnabled, setAutoRefreshEnabled] = useState(DEFAULT_AUTO_REFRESH_ENABLED);
  const [autoRefreshInput, setAutoRefreshInput] = useState(String(DEFAULT_AUTO_REFRESH_INTERVAL_SECONDS));
  const [savedSkinsViewMode, setSavedSkinsViewMode] = useState(DEFAULT_SAVED_SKINS_VIEW_MODE);
  const [fontFamily, setFontFamily] = useState(DEFAULT_FONT_FAMILY);
  const [fontSizeInput, setFontSizeInput] = useState(String(DEFAULT_FONT_SIZE_PX));
  const [isTokenConfigured, setIsTokenConfigured] = useState(false);
  const [tokenValue, setTokenValue] = useState("");
  const [tokenBaseline, setTokenBaseline] = useState("");
  const [tokenStatus, setTokenStatus] = useState<TokenStatus>("missing");
  const [revealToken, setRevealToken] = useState(false);
  const [loading, setLoading] = useState(true);
  const [savingToken, setSavingToken] = useState(false);
  const [resettingToken, setResettingToken] = useState(false);
  const [checkingToken, setCheckingToken] = useState(false);
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
        setAutoRefreshEnabled(settings.autoRefreshEnabled);
        setAutoRefreshInput(String(settings.autoRefreshIntervalSeconds));
        setSavedSkinsViewMode(settings.savedSkinsViewMode);
        setFontFamily(settings.fontFamily);
        setFontSizeInput(String(settings.fontSizePx));
        setIsTokenConfigured(hasToken);
        setTokenStatus(hasToken ? "saved" : "missing");
        setPageError(null);
      } catch (err: unknown) {
        if (!active) return;
        setPageError(formatErrorMessage("Не удалось загрузить настройки.", err));
      } finally {
        if (active) setLoading(false);
      }
    };

    void loadPageState();
    return () => {
      active = false;
    };
  }, []);

  const canSaveToken = useMemo(() => {
    const trimmed = tokenValue.trim();
    return trimmed.length > 0 && trimmed !== tokenBaseline.trim();
  }, [tokenBaseline, tokenValue]);

  const persistSettings = async (
    nextCurrency: SavedSkinCurrency,
    nextAutoRefreshEnabled: boolean,
    nextInterval: number,
    nextFontFamily: string,
    nextFontSizePx: number,
  ) => {
    const saved = await saveAppSettings({
      currency: nextCurrency,
      autoRefreshEnabled: normalizeAutoRefreshEnabled(nextAutoRefreshEnabled),
      autoRefreshIntervalSeconds: nextInterval,
      savedSkinsViewMode,
      fontFamily: normalizeFontFamily(nextFontFamily),
      fontSizePx: normalizeFontSizePx(nextFontSizePx),
    });
    setCurrency(saved.currency);
    setAutoRefreshEnabled(saved.autoRefreshEnabled);
    setAutoRefreshInput(String(saved.autoRefreshIntervalSeconds));
    setSavedSkinsViewMode(saved.savedSkinsViewMode);
    setFontFamily(saved.fontFamily);
    setFontSizeInput(String(saved.fontSizePx));
    setNotice({ type: "success", text: UI_TEXT.settingsAutoSaved });
  };

  const onCurrencyChange = async (nextCurrency: SavedSkinCurrency) => {
    setNotice(null);
    try {
      await persistSettings(
        nextCurrency,
        autoRefreshEnabled,
        normalizeAutoRefreshIntervalSeconds(autoRefreshInput),
        fontFamily,
        normalizeFontSizePx(fontSizeInput),
      );
    } catch (err: unknown) {
      setNotice({ type: "error", text: formatErrorMessage("Не удалось сохранить настройки.", err) });
    }
  };

  const onAutoRefreshEnabledChange = async (nextValue: boolean) => {
    setNotice(null);
    try {
      await persistSettings(
        currency,
        nextValue,
        normalizeAutoRefreshIntervalSeconds(autoRefreshInput),
        fontFamily,
        normalizeFontSizePx(fontSizeInput),
      );
    } catch (err: unknown) {
      setNotice({ type: "error", text: formatErrorMessage("Не удалось сохранить настройки.", err) });
    }
  };

  const onAutoRefreshBlur = () => {
    const normalized = normalizeAutoRefreshIntervalSeconds(autoRefreshInput);
    void (async () => {
      setNotice(null);
      try {
        await persistSettings(currency, autoRefreshEnabled, normalized, fontFamily, normalizeFontSizePx(fontSizeInput));
      } catch (err: unknown) {
        setNotice({ type: "error", text: formatErrorMessage("Не удалось сохранить настройки.", err) });
      }
    })();
  };

  const onFontFamilyChange = async (nextFontFamily: string) => {
    setNotice(null);
    try {
      await persistSettings(
        currency,
        autoRefreshEnabled,
        normalizeAutoRefreshIntervalSeconds(autoRefreshInput),
        nextFontFamily,
        normalizeFontSizePx(fontSizeInput),
      );
    } catch (err: unknown) {
      setNotice({ type: "error", text: formatErrorMessage("Не удалось сохранить настройки.", err) });
    }
  };

  const onFontSizeBlur = () => {
    const normalized = normalizeFontSizePx(fontSizeInput);
    void (async () => {
      setNotice(null);
      try {
        await persistSettings(
          currency,
          autoRefreshEnabled,
          normalizeAutoRefreshIntervalSeconds(autoRefreshInput),
          fontFamily,
          normalized,
        );
      } catch (err: unknown) {
        setNotice({ type: "error", text: formatErrorMessage("Не удалось сохранить настройки.", err) });
      }
    })();
  };

  const onFontSizeChange = (nextFontSize: string) => {
    setNotice(null);
    setFontSizeInput(nextFontSize);
  };

  const onSaveToken = async () => {
    setSavingToken(true);
    setTokenError(null);
    setNotice(null);
    try {
      const trimmed = tokenValue.trim();
      await saveLisSkinsToken(trimmed);
      setIsTokenConfigured(true);
      setTokenBaseline(trimmed);
      setTokenStatus("saved");
      setNotice({ type: "success", text: UI_TEXT.lisSkinsTokenSaved });
    } catch (err: unknown) {
      setTokenStatus("authError");
      setTokenError(formatErrorMessage(UI_TEXT.errLisSkinsTokenSave, err));
    } finally {
      setSavingToken(false);
    }
  };

  const onCheckToken = async () => {
    setCheckingToken(true);
    setTokenError(null);
    try {
      await saveLisSkinsToken(tokenValue.trim());
      setTokenStatus("verified");
    } catch (err: unknown) {
      setTokenStatus("authError");
      setTokenError(formatErrorMessage(UI_TEXT.errLisSkinsTokenSave, err));
    } finally {
      setCheckingToken(false);
    }
  };

  const onResetToken = async () => {
    setResettingToken(true);
    setTokenError(null);
    setNotice(null);
    try {
      await clearLisSkinsToken();
      setIsTokenConfigured(false);
      setTokenValue("");
      setTokenBaseline("");
      setTokenStatus("missing");
      setNotice({ type: "success", text: UI_TEXT.lisSkinsTokenResetDone });
    } catch (err: unknown) {
      setTokenStatus("authError");
      setTokenError(formatErrorMessage(UI_TEXT.errLisSkinsTokenReset, err));
    } finally {
      setResettingToken(false);
    }
  };

  if (loading) return <LoadingState text="Загрузка настроек..." />;
  if (pageError) return <ErrorState text={pageError} />;

  return (
    <div className="settings-page">
      <PageHeader sectionLabel={UI_TEXT.settingsEyebrow} title={UI_TEXT.settingsTitle} actions={<div className="toolbar-group"><button className="toolbar-button" type="button" onClick={() => navigate(ROUTES.home)}>{UI_TEXT.ctaToHome}</button></div>} />
      {notice && <ToastAlert type={notice.type} text={notice.text} />}
      <LisSkinsTokenPanel
        visible
        configured={isTokenConfigured}
        saving={savingToken}
        resetting={resettingToken}
        checking={checkingToken}
        error={tokenError}
        value={tokenValue}
        canSave={canSaveToken}
        status={tokenStatus}
        revealToken={revealToken}
        onToggleReveal={() => setRevealToken((v) => !v)}
        onChange={(nextValue) => {
          setTokenValue(nextValue);
          setTokenError(null);
          setTokenStatus(nextValue.trim().length === 0 ? "missing" : isTokenConfigured ? "saved" : "missing");
        }}
        onSave={onSaveToken}
        onCheck={onCheckToken}
        onReset={onResetToken}
      />
      <section className="settings-grid">
        <article className="settings-card"><h2 className="settings-card-title">{UI_TEXT.settingsCurrencyTitle}</h2><p className="settings-card-description">{UI_TEXT.settingsCurrencyHelp}</p><label className="token-panel-field" htmlFor="settings-currency-select"><span className="token-panel-label">{UI_TEXT.settingsCurrencyTitle}</span><select id="settings-currency-select" className="toolbar-select settings-select" value={currency} onChange={(event) => void onCurrencyChange(event.target.value as SavedSkinCurrency)}>{CURRENCY_OPTIONS.map((option) => <option key={option.value} value={option.value}>{option.label}</option>)}</select></label></article>
        <article className="settings-card"><h2 className="settings-card-title">{UI_TEXT.settingsRefreshTitle}</h2><p className="settings-card-description">{UI_TEXT.settingsRefreshHelp}</p><label className="search-toggle" htmlFor="settings-refresh-enabled"><input id="settings-refresh-enabled" type="checkbox" checked={autoRefreshEnabled} onChange={(event) => void onAutoRefreshEnabledChange(event.target.checked)} /><span>{UI_TEXT.settingsRefreshEnabledLabel}</span></label><label className="token-panel-field" htmlFor="settings-refresh-interval"><span className="token-panel-label">{UI_TEXT.settingsRefreshLabel}</span><input id="settings-refresh-interval" className="search-input settings-number-input" type="number" min={5} step={5} value={autoRefreshInput} disabled={!autoRefreshEnabled} onChange={(event) => setAutoRefreshInput(event.target.value)} onBlur={onAutoRefreshBlur} /></label><p className="toolbar-hint">{autoRefreshEnabled ? UI_TEXT.settingsRefreshHint : UI_TEXT.settingsRefreshDisabledHint}</p></article>
        <article className="settings-card"><h2 className="settings-card-title">{UI_TEXT.settingsFontFamilyTitle}</h2><p className="settings-card-description">{UI_TEXT.settingsFontFamilyHelp}</p><label className="token-panel-field" htmlFor="settings-font-family-select"><span className="token-panel-label">{UI_TEXT.settingsFontFamilyLabel}</span><select id="settings-font-family-select" className="toolbar-select settings-select" value={fontFamily} onChange={(event) => void onFontFamilyChange(event.target.value)}>{FONT_FAMILY_OPTIONS.map((option) => <option key={option.value} value={option.value}>{option.label}</option>)}</select></label></article>
        <article className="settings-card"><h2 className="settings-card-title">{UI_TEXT.settingsFontSizeTitle}</h2><p className="settings-card-description">{UI_TEXT.settingsFontSizeHelp}</p><label className="token-panel-field" htmlFor="settings-font-size-input"><span className="token-panel-label">{UI_TEXT.settingsFontSizeLabel}</span><input id="settings-font-size-input" className="search-input settings-number-input" type="number" min={MIN_FONT_SIZE_PX} max={MAX_FONT_SIZE_PX} step={1} value={fontSizeInput} onChange={(event) => onFontSizeChange(event.target.value)} onBlur={onFontSizeBlur} /></label><p className="toolbar-hint">{UI_TEXT.settingsFontSizeHint.replace("{min}", String(MIN_FONT_SIZE_PX)).replace("{max}", String(MAX_FONT_SIZE_PX))}</p></article>
      </section>
    </div>
  );
};
