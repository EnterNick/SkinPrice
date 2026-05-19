import React, { useCallback, useEffect, useMemo, useRef, useState } from "react";
import { useNavigate } from "react-router-dom";
import { ROUTES } from "../../../app/router/routes";
import { deleteSavedSkin, getAppSettings, updateAllSkinPrices, updateSkinPrice } from "../../../entities/skin/api/skinApi";
import type { SavedSkinCurrency } from "../../../entities/skin/model/types";
import { useSavedSkins } from "../../../entities/skin/model/useSavedSkins";
import { DEFAULT_AUTO_REFRESH_INTERVAL_SECONDS, DEFAULT_CURRENCY } from "../../../shared/config/settings";
import { UI_TEXT } from "../../../shared/config/uiText";
import { formatErrorMessage } from "../../../shared/lib/error/formatErrorMessage";
import { EmptyState, ErrorState, LoadingState, ToastAlert } from "../../../shared/ui/states/States";
import { PageHeader } from "../../../shared/ui/page-header/PageHeader";
import { SavedSkinsTable } from "../../../widgets/saved-skins-table/SavedSkinsTable";

export const SavedSkinsPage: React.FC = () => {
  const navigate = useNavigate();
  const { items, loading, error, loadSkins } = useSavedSkins();
  const [currency, setCurrency] = useState<SavedSkinCurrency>(DEFAULT_CURRENCY);
  const [autoRefreshMs, setAutoRefreshMs] = useState(DEFAULT_AUTO_REFRESH_INTERVAL_SECONDS * 1000);
  const [settingsLoading, setSettingsLoading] = useState(true);
  const [settingsError, setSettingsError] = useState<string | null>(null);
  const [updatingSkinIds, setUpdatingSkinIds] = useState<Record<string, boolean>>({});
  const [deletingSkinIds, setDeletingSkinIds] = useState<Record<string, boolean>>({});
  const [isUpdatingAll, setIsUpdatingAll] = useState(false);
  const [notice, setNotice] = useState<{ type: "success" | "warning" | "error"; text: string } | null>(null);
  const autoRefreshInFlightRef = useRef(false);
  const currencyDirty = useMemo(() => items.length > 0 && items.some((skin) => skin.currency !== currency), [currency, items]);

  useEffect(() => {
    let active = true;

    const loadSettings = async () => {
      try {
        const settings = await getAppSettings();
        if (!active) return;
        setCurrency(settings.currency);
        setAutoRefreshMs(settings.autoRefreshIntervalSeconds * 1000);
        setSettingsError(null);
      } catch (err: unknown) {
        if (!active) return;
        setSettingsError(formatErrorMessage("Не удалось загрузить настройки.", err));
      } finally {
        if (active) {
          setSettingsLoading(false);
        }
      }
    };

    void loadSettings();

    return () => {
      active = false;
    };
  }, []);

  const refreshOne = useCallback(async (skinId: string) => {
    setUpdatingSkinIds((prev) => ({ ...prev, [skinId]: true }));
    setNotice(null);
    try {
      await updateSkinPrice(skinId, currency);
      await loadSkins();
    } catch (err: unknown) {
      setNotice({ type: "error", text: formatErrorMessage(UI_TEXT.errUpdateOne, err) });
    } finally {
      setUpdatingSkinIds((prev) => ({ ...prev, [skinId]: false }));
    }
  }, [currency, loadSkins]);

  const refreshAll = useCallback(async (nextCurrency: SavedSkinCurrency = currency, showNotice = true) => {
    setIsUpdatingAll(true);
    if (showNotice) {
      setNotice(null);
    }
    try {
      const result = await updateAllSkinPrices(nextCurrency);
      try {
        await loadSkins();
        if (showNotice) {
          const text =
            result.failed > 0
              ? UI_TEXT.partialUpdatedAll.replace("{updated}", String(result.updated)).replace("{failed}", String(result.failed))
              : UI_TEXT.successUpdatedAll.replace("{count}", String(result.updated));
          setNotice({ type: result.failed > 0 ? "warning" : "success", text });
        }
      } catch (reloadError) {
        if (showNotice) {
          const message = formatErrorMessage("", reloadError).trim();
          setNotice({ type: "warning", text: UI_TEXT.warningPartialUpdated.replace("{message}", message) });
        }
      }
    } catch (updateError) {
      if (showNotice) {
        setNotice({ type: "error", text: formatErrorMessage(UI_TEXT.errUpdateAll, updateError) });
      }
    } finally {
      setIsUpdatingAll(false);
    }
  }, [currency, loadSkins]);

  const removeSkin = useCallback(async (skinId: string) => {
    setDeletingSkinIds((prev) => ({ ...prev, [skinId]: true }));
    setNotice(null);
    try {
      await deleteSavedSkin(skinId);
      await loadSkins();
    } catch (err: unknown) {
      setNotice({ type: "error", text: formatErrorMessage(UI_TEXT.errDelete, err) });
    } finally {
      setDeletingSkinIds((prev) => ({ ...prev, [skinId]: false }));
    }
  }, [loadSkins]);

  useEffect(() => {
    if (loading || settingsLoading || items.length === 0 || !currencyDirty || isUpdatingAll) return;
    void refreshAll(currency, false);
  }, [currency, currencyDirty, items.length, isUpdatingAll, loading, refreshAll, settingsLoading]);

  useEffect(() => {
    if (loading || settingsLoading || items.length === 0 || currencyDirty) return undefined;

    const intervalId = window.setInterval(() => {
      const hasRowOperation = Object.values(updatingSkinIds).some(Boolean) || Object.values(deletingSkinIds).some(Boolean);
      if (autoRefreshInFlightRef.current || isUpdatingAll || hasRowOperation) {
        return;
      }

      autoRefreshInFlightRef.current = true;
      void refreshAll(currency, false).finally(() => {
        autoRefreshInFlightRef.current = false;
      });
    }, autoRefreshMs);

    return () => {
      window.clearInterval(intervalId);
    };
  }, [autoRefreshMs, currency, currencyDirty, deletingSkinIds, isUpdatingAll, items.length, loading, refreshAll, settingsLoading, updatingSkinIds]);

  if (loading || settingsLoading) return <LoadingState text={UI_TEXT.loadingSaved} />;
  if (settingsError) return <ErrorState text={settingsError} />;
  if (error) return <ErrorState text={error} />;

  return (
    <div className="saved-skins-page">
      <PageHeader
        sectionLabel={UI_TEXT.collectionEyebrow}
        title={UI_TEXT.savedSkinsTitle}
        actions={
          <div className="toolbar-group">
            <button className="toolbar-button" type="button" onClick={() => navigate(ROUTES.settings)}>
              {UI_TEXT.ctaToSettings}
            </button>
            <button className="toolbar-button toolbar-button-primary" type="button" disabled={isUpdatingAll} onClick={() => void refreshAll()}>
              {isUpdatingAll ? UI_TEXT.updateAllPending : UI_TEXT.updateAll}
            </button>
            <button className="toolbar-button" type="button" onClick={() => navigate(ROUTES.newSkins)}>
              {UI_TEXT.ctaAddSkins}
            </button>
          </div>
        }
      />
      {notice && <ToastAlert type={notice.type} text={notice.text} />}
      {items.length === 0 ? (
        <div className="empty-shell">
          <EmptyState text={UI_TEXT.notFoundSaved} />
          <p className="empty-hint">{UI_TEXT.emptySavedHint}</p>
          <button className="toolbar-button toolbar-button-primary empty-action" type="button" onClick={() => navigate(ROUTES.newSkins)}>
            {UI_TEXT.ctaAddSkins}
          </button>
        </div>
      ) : (
        <SavedSkinsTable
          items={items}
          isUpdatingAll={isUpdatingAll}
          updatingSkinIds={updatingSkinIds}
          deletingSkinIds={deletingSkinIds}
          onRefreshOne={refreshOne}
          onDelete={removeSkin}
        />
      )}
    </div>
  );
};
