import React, { useCallback, useEffect, useMemo, useRef, useState } from "react";
import { useNavigate } from "react-router-dom";
import { ROUTES } from "../../../app/router/routes";
import { deleteSavedSkin, getAppSettings, saveAppSettings, updateAllSkinPrices, updateSkinPrice } from "../../../entities/skin/api/skinApi";
import type { SavedSkinCurrency, SavedSkinsSortColumn, SavedSkinsSortState, SavedSkinsViewMode } from "../../../entities/skin/model/types";
import { useSavedSkins } from "../../../entities/skin/model/useSavedSkins";
import { DEFAULT_AUTO_REFRESH_INTERVAL_SECONDS, DEFAULT_CURRENCY, DEFAULT_SAVED_SKINS_VIEW_MODE } from "../../../shared/config/settings";
import { UI_TEXT } from "../../../shared/config/uiText";
import { formatErrorMessage } from "../../../shared/lib/error/formatErrorMessage";
import { EmptyState, ErrorState, LoadingState, ToastAlert } from "../../../shared/ui/states/States";
import { PageHeader } from "../../../shared/ui/page-header/PageHeader";
import { SavedSkinsCards } from "../../../widgets/saved-skins-cards/SavedSkinsCards";
import { SavedSkinsTable } from "../../../widgets/saved-skins-table/SavedSkinsTable";

const parsePriceValue = (priceText?: string): number => {
  const normalized = (priceText || "").replace(/\s/g, "").replace(",", ".");
  const matched = normalized.match(/\d+(?:\.\d+)?/);
  if (!matched) return Number.POSITIVE_INFINITY;
  const value = Number.parseFloat(matched[0]);
  return Number.isFinite(value) ? value : Number.POSITIVE_INFINITY;
};

export const SavedSkinsPage: React.FC = () => {
  const navigate = useNavigate();
  const { items, loading, error, loadSkins } = useSavedSkins();
  const [currency, setCurrency] = useState<SavedSkinCurrency>(DEFAULT_CURRENCY);
  const [autoRefreshMs, setAutoRefreshMs] = useState(DEFAULT_AUTO_REFRESH_INTERVAL_SECONDS * 1000);
  const [viewMode, setViewMode] = useState<SavedSkinsViewMode>(DEFAULT_SAVED_SKINS_VIEW_MODE);
  const [settingsLoading, setSettingsLoading] = useState(true);
  const [settingsError, setSettingsError] = useState<string | null>(null);
  const [updatingSkinIds, setUpdatingSkinIds] = useState<Record<string, boolean>>({});
  const [deletingSkinIds, setDeletingSkinIds] = useState<Record<string, boolean>>({});
  const [isUpdatingAll, setIsUpdatingAll] = useState(false);
  const [isSavingViewMode, setIsSavingViewMode] = useState(false);
  const [notice, setNotice] = useState<{ type: "success" | "warning" | "error"; text: string } | null>(null);
  const [searchQuery, setSearchQuery] = useState("");
  const [sortState, setSortState] = useState<SavedSkinsSortState>(null);
  const [isSortMenuOpen, setIsSortMenuOpen] = useState(false);
  const autoRefreshInFlightRef = useRef(false);
  const sortMenuRef = useRef<HTMLDivElement | null>(null);
  const currencyDirty = useMemo(() => items.length > 0 && items.some((skin) => skin.currency !== currency), [currency, items]);

  const filteredItems = useMemo(() => {
    const query = searchQuery.trim().toLowerCase();
    const base = query.length === 0 ? items : items.filter((skin) => skin.title.toLowerCase().includes(query));
    if (!sortState) {
      return base;
    }

    const directionFactor = sortState.direction === "asc" ? 1 : -1;
    return [...base].sort((a, b) => {
      let result = 0;
      if (sortState.column === "title") {
        result = a.title.localeCompare(b.title, "ru");
      } else if (sortState.column === "steamPrice") {
        result = parsePriceValue(a.steamPriceText) - parsePriceValue(b.steamPriceText);
      } else {
        result = parsePriceValue(a.lisSkinsPriceText) - parsePriceValue(b.lisSkinsPriceText);
      }

      if (result === 0) {
        result = a.title.localeCompare(b.title, "ru");
      }

      return result * directionFactor;
    });
  }, [items, searchQuery, sortState]);

  const toggleSort = useCallback((column: SavedSkinsSortColumn) => {
    setSortState((current) => {
      if (!current || current.column !== column) {
        return { column, direction: "desc" };
      }
      if (current.direction === "desc") {
        return { column, direction: "asc" };
      }
      return null;
    });
  }, []);

  useEffect(() => {
    if (!isSortMenuOpen) return undefined;

    const handlePointerDown = (event: PointerEvent) => {
      const target = event.target;
      if (!(target instanceof Node)) return;
      if (sortMenuRef.current?.contains(target)) return;
      setIsSortMenuOpen(false);
    };

    const handleKeyDown = (event: KeyboardEvent) => {
      if (event.key === "Escape") {
        setIsSortMenuOpen(false);
      }
    };

    window.addEventListener("pointerdown", handlePointerDown);
    window.addEventListener("keydown", handleKeyDown);

    return () => {
      window.removeEventListener("pointerdown", handlePointerDown);
      window.removeEventListener("keydown", handleKeyDown);
    };
  }, [isSortMenuOpen]);

  useEffect(() => {
    let active = true;

    const loadSettings = async () => {
      try {
        const settings = await getAppSettings();
        if (!active) return;
        setCurrency(settings.currency);
        setAutoRefreshMs(settings.autoRefreshIntervalSeconds * 1000);
        setViewMode(settings.savedSkinsViewMode);
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

  const changeViewMode = useCallback(async (nextViewMode: SavedSkinsViewMode) => {
    if (nextViewMode === viewMode) return;

    setIsSavingViewMode(true);
    setIsSortMenuOpen(false);
    setNotice(null);
    try {
      const saved = await saveAppSettings({
        currency,
        autoRefreshIntervalSeconds: Math.round(autoRefreshMs / 1000),
        savedSkinsViewMode: nextViewMode,
      });
      setViewMode(saved.savedSkinsViewMode);
    } catch (err: unknown) {
      setNotice({ type: "error", text: formatErrorMessage("Не удалось сохранить режим отображения.", err) });
    } finally {
      setIsSavingViewMode(false);
    }
  }, [autoRefreshMs, currency, viewMode]);

  const toggleViewMode = useCallback(async () => {
    await changeViewMode(viewMode === "table" ? "cards" : "table");
  }, [changeViewMode, viewMode]);

  const applySortOption = useCallback((column: SavedSkinsSortColumn) => {
    toggleSort(column);
    setIsSortMenuOpen(false);
  }, [toggleSort]);

  const getSortOptionLabel = useCallback((column: SavedSkinsSortColumn) => {
    if (column === "title") return UI_TEXT.savedSkinsSortTitle;
    if (column === "steamPrice") return UI_TEXT.savedSkinsSortSteam;
    return UI_TEXT.savedSkinsSortLisSkins;
  }, []);

  const getSortButtonLabel = useCallback(() => {
    if (!sortState) return UI_TEXT.savedSkinsSortMenu;

    const suffix = sortState.direction === "desc" ? "↓" : "↑";
    return `${getSortOptionLabel(sortState.column)} ${suffix}`;
  }, [getSortOptionLabel, sortState]);

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
            <button
              className="toolbar-button toolbar-button-secondary view-mode-icon-button"
              type="button"
              disabled={isSavingViewMode}
              title={viewMode === "table" ? UI_TEXT.savedSkinsViewToggleToCards : UI_TEXT.savedSkinsViewToggleToTable}
              aria-label={viewMode === "table" ? UI_TEXT.savedSkinsViewToggleToCards : UI_TEXT.savedSkinsViewToggleToTable}
              onClick={() => void toggleViewMode()}
            >
              <span className={`view-mode-icon${viewMode === "table" ? " view-mode-icon-table" : " view-mode-icon-cards"}`} aria-hidden="true" />
            </button>
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
          <h3 className="empty-title">В коллекции пока нет скинов</h3>
          <p className="empty-hint">Добавьте первый скин, чтобы начать отслеживать цену</p>
          <button className="toolbar-button toolbar-button-primary empty-action" type="button" onClick={() => navigate(ROUTES.newSkins)}>
            Добавить скин
          </button>
        </div>
      ) : (
        <>
          <div className="saved-skins-controls">
            <input
              className="search-input"
              type="text"
              value={searchQuery}
              onChange={(event) => setSearchQuery(event.target.value)}
              placeholder="Поиск по названию"
            />
            {viewMode === "cards" ? (
              <div className="saved-skins-sort-dropdown" ref={sortMenuRef}>
                <button
                  className={`toolbar-button toolbar-button-secondary saved-skins-sort-trigger${isSortMenuOpen ? " view-mode-button-active" : ""}`}
                  type="button"
                  aria-haspopup="menu"
                  aria-expanded={isSortMenuOpen}
                  onClick={() => setIsSortMenuOpen((current) => !current)}
                >
                  <span>{getSortButtonLabel()}</span>
                </button>
                {isSortMenuOpen ? (
                  <div className="saved-skins-sort-menu" role="menu" aria-label={UI_TEXT.savedSkinsSortMenu}>
                    <button className="saved-skins-sort-item" type="button" role="menuitem" onClick={() => applySortOption("title")}>
                      {getSortOptionLabel("title")}
                      {sortState?.column === "title" ? <span className="saved-skins-sort-state">{sortState.direction === "desc" ? "↓" : "↑"}</span> : null}
                    </button>
                    <button className="saved-skins-sort-item" type="button" role="menuitem" onClick={() => applySortOption("steamPrice")}>
                      {getSortOptionLabel("steamPrice")}
                      {sortState?.column === "steamPrice" ? <span className="saved-skins-sort-state">{sortState.direction === "desc" ? "↓" : "↑"}</span> : null}
                    </button>
                    <button className="saved-skins-sort-item" type="button" role="menuitem" onClick={() => applySortOption("lisSkinsPrice")}>
                      {getSortOptionLabel("lisSkinsPrice")}
                      {sortState?.column === "lisSkinsPrice" ? <span className="saved-skins-sort-state">{sortState.direction === "desc" ? "↓" : "↑"}</span> : null}
                    </button>
                    <button className="saved-skins-sort-item" type="button" role="menuitem" onClick={() => { setSortState(null); setIsSortMenuOpen(false); }}>
                      {UI_TEXT.savedSkinsSortReset}
                    </button>
                  </div>
                ) : null}
              </div>
            ) : null}
          </div>
          {viewMode === "table" ? (
            <SavedSkinsTable
              items={filteredItems}
              isUpdatingAll={isUpdatingAll}
              updatingSkinIds={updatingSkinIds}
              deletingSkinIds={deletingSkinIds}
              sortState={sortState}
              onSortChange={toggleSort}
              onRefreshOne={refreshOne}
              onDelete={removeSkin}
            />
          ) : (
            <SavedSkinsCards
              items={filteredItems}
              isUpdatingAll={isUpdatingAll}
              updatingSkinIds={updatingSkinIds}
              deletingSkinIds={deletingSkinIds}
              onRefreshOne={refreshOne}
              onDelete={removeSkin}
            />
          )}
          {filteredItems.length === 0 && <EmptyState text="По текущему поиску ничего не найдено" />}
        </>
      )}
    </div>
  );
};
