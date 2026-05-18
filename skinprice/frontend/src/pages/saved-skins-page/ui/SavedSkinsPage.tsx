import React, { useEffect, useRef, useState } from "react";
import { useNavigate } from "react-router-dom";
import { ROUTES } from "../../../app/router/routes";
import { CURRENCY_OPTIONS, deleteSavedSkin, normalizeCurrency, updateAllSkinPrices, updateSkinPrice } from "../../../entities/skin/api/skinApi";
import type { SavedSkinCurrency } from "../../../entities/skin/model/types";
import { useSavedSkins } from "../../../entities/skin/model/useSavedSkins";
import { UI_TEXT } from "../../../shared/config/uiText";
import { STORAGE_KEYS } from "../../../shared/config/storage";
import { formatErrorMessage } from "../../../shared/lib/error/formatErrorMessage";
import { EmptyState, ErrorState, LoadingState, ToastAlert } from "../../../shared/ui/states/States";
import { PageHeader } from "../../../shared/ui/page-header/PageHeader";
import { SavedSkinsTable } from "../../../widgets/saved-skins-table/SavedSkinsTable";

export const SavedSkinsPage: React.FC = () => {
  const navigate = useNavigate();
  const { items, loading, error, loadSkins } = useSavedSkins();
  const [currency, setCurrency] = useState<SavedSkinCurrency>(() => normalizeCurrency(window.localStorage.getItem(STORAGE_KEYS.currency)));
  const [updatingSkinIds, setUpdatingSkinIds] = useState<Record<string, boolean>>({});
  const [deletingSkinIds, setDeletingSkinIds] = useState<Record<string, boolean>>({});
  const [isUpdatingAll, setIsUpdatingAll] = useState(false);
  const [notice, setNotice] = useState<{ type: "success" | "warning" | "error"; text: string } | null>(null);
  const autoSyncedCurrencyRef = useRef<string | null>(null);

  const refreshOne = async (skinId: string) => {
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
  };

  const refreshAll = async (nextCurrency: SavedSkinCurrency = currency) => {
    setIsUpdatingAll(true);
    setNotice(null);
    try {
      const result = await updateAllSkinPrices(nextCurrency);
      try {
        await loadSkins();
        const text =
          result.failed > 0
            ? UI_TEXT.partialUpdatedAll.replace("{updated}", String(result.updated)).replace("{failed}", String(result.failed))
            : UI_TEXT.successUpdatedAll.replace("{count}", String(result.updated));
        setNotice({ type: result.failed > 0 ? "warning" : "success", text });
      } catch (reloadError) {
        const message = formatErrorMessage("", reloadError).trim();
        setNotice({ type: "warning", text: UI_TEXT.warningPartialUpdated.replace("{message}", message) });
      }
    } catch (updateError) {
      setNotice({ type: "error", text: formatErrorMessage(UI_TEXT.errUpdateAll, updateError) });
    } finally {
      setIsUpdatingAll(false);
    }
  };

  const removeSkin = async (skinId: string) => {
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
  };

  const onCurrencyChange = (nextCurrencyValue: string) => {
    const nextCurrency = normalizeCurrency(nextCurrencyValue);
    setCurrency(nextCurrency);
    window.localStorage.setItem(STORAGE_KEYS.currency, nextCurrency);
    autoSyncedCurrencyRef.current = nextCurrency;
    if (items.length > 0) {
      void refreshAll(nextCurrency);
    }
  };

  useEffect(() => {
    if (loading || isUpdatingAll || items.length === 0) return;

    const needsCurrencySync = items.some((skin) => skin.currency !== currency);
    if (!needsCurrencySync) {
      autoSyncedCurrencyRef.current = currency;
      return;
    }

    if (autoSyncedCurrencyRef.current === currency) return;

    autoSyncedCurrencyRef.current = currency;
    void refreshAll(currency);
  }, [currency, items, loading, isUpdatingAll]);

  if (loading) return <LoadingState text={UI_TEXT.loadingSaved} />;
  if (error) return <ErrorState text={error} />;

  return (
    <div className="saved-skins-page">
      <PageHeader
        eyebrow={UI_TEXT.collectionEyebrow}
        title={UI_TEXT.savedSkinsTitle}
        actions={
          <div className="saved-skins-toolbar">
            <select className="toolbar-select" value={currency} onChange={(e) => onCurrencyChange(e.target.value)}>
              {CURRENCY_OPTIONS.map((option) => (
                <option key={option.value} value={option.value}>
                  {option.label}
                </option>
              ))}
            </select>
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
