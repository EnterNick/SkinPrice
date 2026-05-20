import React, { useEffect, useRef, useState } from "react";
import { createPortal } from "react-dom";
import { UI_TEXT } from "../../shared/config/uiText";
import { ClipboardSetText } from "../../wailsjs/runtime/runtime";
import { openExternal } from "../../shared/lib/browser/openExternal";
import { formatUpdatedAt } from "../../shared/lib/date/formatUpdatedAt";
import { SkinThumb } from "../../entities/skin/ui/SkinThumb";
import type { SavedSkin, SavedSkinsSortColumn, SavedSkinsSortDirection, SavedSkinsSortState } from "../../entities/skin/model/types";

type SavedSkinsTableProps = {
  items: SavedSkin[];
  isUpdatingAll: boolean;
  updatingSkinIds: Record<string, boolean>;
  deletingSkinIds: Record<string, boolean>;
  sortState: SavedSkinsSortState;
  onSortChange: (column: SavedSkinsSortColumn) => void;
  onRefreshOne: (skinId: string) => Promise<void> | void;
  onDelete: (skinId: string) => Promise<void> | void;
};

const buildPriceTooltip = (source: string, updatedAt?: string) => {
  const formatted = formatUpdatedAt(updatedAt);
  if (formatted === "-") {
    return `${source}: ${UI_TEXT.notUpdatedYet}`;
  }
  return `${source}: ${formatted}`;
};

const buildSkinNameStyle = (nameColor?: string): React.CSSProperties | undefined =>
  nameColor ? { "--skin-name-color": nameColor } as React.CSSProperties : undefined;

const getSortArrow = (direction: SavedSkinsSortDirection) => (direction === "desc" ? "↓" : "↑");
const ROW_ACTIONS_DROPDOWN_WIDTH = 180;

export const SavedSkinsTable: React.FC<SavedSkinsTableProps> = ({
  items,
  isUpdatingAll,
  updatingSkinIds,
  deletingSkinIds,
  sortState,
  onSortChange,
  onRefreshOne,
  onDelete,
}) => {
  const [openMenuId, setOpenMenuId] = useState<string | null>(null);
  const [copiedSkinId, setCopiedSkinId] = useState<string | null>(null);
  const [menuPosition, setMenuPosition] = useState<{ top: number; left: number } | null>(null);
  const copiedResetTimeoutRef = useRef<number | null>(null);
  const toggleRefs = useRef<Record<string, HTMLButtonElement | null>>({});

  useEffect(() => {
    const handlePointerDown = (event: PointerEvent) => {
      const target = event.target;
      if (!(target instanceof Element) || target.closest(".row-actions-menu") || target.closest(".row-actions-dropdown-portal")) return;
      setOpenMenuId(null);
    };

    const handleKeyDown = (event: KeyboardEvent) => {
      if (event.key === "Escape") {
        setOpenMenuId(null);
      }
    };

    window.addEventListener("pointerdown", handlePointerDown);
    window.addEventListener("keydown", handleKeyDown);

    return () => {
      window.removeEventListener("pointerdown", handlePointerDown);
      window.removeEventListener("keydown", handleKeyDown);
    };
  }, []);

  useEffect(() => () => {
    if (copiedResetTimeoutRef.current !== null) {
      window.clearTimeout(copiedResetTimeoutRef.current);
    }
  }, []);

  useEffect(() => {
    if (!openMenuId) {
      setMenuPosition(null);
      return;
    }

    const updateMenuPosition = () => {
      const toggle = toggleRefs.current[openMenuId];
      if (!toggle) {
        setMenuPosition(null);
        return;
      }

      const rect = toggle.getBoundingClientRect();
      const topBelow = rect.bottom + 8;
      const topAbove = rect.top - 8;
      const dropdownHeight = 96;
      const showAbove = topBelow + dropdownHeight > window.innerHeight && rect.top > dropdownHeight + 16;
      const nextTop = showAbove ? Math.max(8, topAbove - dropdownHeight) : topBelow;
      const nextLeft = Math.min(
        Math.max(8, rect.right - ROW_ACTIONS_DROPDOWN_WIDTH),
        Math.max(8, window.innerWidth - ROW_ACTIONS_DROPDOWN_WIDTH - 8),
      );

      setMenuPosition({ top: nextTop, left: nextLeft });
    };

    updateMenuPosition();
    window.addEventListener("resize", updateMenuPosition);
    window.addEventListener("scroll", updateMenuPosition, true);

    return () => {
      window.removeEventListener("resize", updateMenuPosition);
      window.removeEventListener("scroll", updateMenuPosition, true);
    };
  }, [openMenuId]);

  const handleCopyTitle = async (skin: SavedSkin) => {
    const copied = await ClipboardSetText(skin.title);
    if (!copied) return;

    setCopiedSkinId(skin.id);
    if (copiedResetTimeoutRef.current !== null) {
      window.clearTimeout(copiedResetTimeoutRef.current);
    }
    copiedResetTimeoutRef.current = window.setTimeout(() => {
      setCopiedSkinId((current) => (current === skin.id ? null : current));
      copiedResetTimeoutRef.current = null;
    }, 1600);
  };

  return (
    <div className="table-shell">
      <div className="table-scroll">
        <table className="skins-table">
          <thead>
            <tr>
              <th>
                <button className="table-sort-button" type="button" onClick={() => onSortChange("title")}>
                  <span>Скин</span>
                  {sortState?.column === "title" ? <span className="table-sort-arrow">{getSortArrow(sortState.direction)}</span> : null}
                </button>
              </th>
              <th>
                <button className="table-sort-button" type="button" onClick={() => onSortChange("steamPrice")}>
                  <span>{UI_TEXT.steamPriceLabel}</span>
                  {sortState?.column === "steamPrice" ? <span className="table-sort-arrow">{getSortArrow(sortState.direction)}</span> : null}
                </button>
              </th>
              <th>
                <button className="table-sort-button" type="button" onClick={() => onSortChange("lisSkinsPrice")}>
                  <span>{UI_TEXT.lisSkinsPriceLabel}</span>
                  {sortState?.column === "lisSkinsPrice" ? <span className="table-sort-arrow">{getSortArrow(sortState.direction)}</span> : null}
                </button>
              </th>
              <th className="actions-menu-column" aria-label="Действия" />
            </tr>
          </thead>
          <tbody>
            {items.map((skin) => {
              const isUpdating = Boolean(updatingSkinIds[skin.id]);
              const isDeleting = Boolean(deletingSkinIds[skin.id]);
              const isMenuOpen = openMenuId === skin.id;
              return (
                <tr key={skin.id}>
                  <td>
                    <div className="skin-cell">
                      <SkinThumb imageUrl={skin.imageUrl} title={skin.title} />
                      <button
                        className={`skin-title-button${copiedSkinId === skin.id ? " skin-title-button-copied" : ""}`}
                        type="button"
                        title={copiedSkinId === skin.id ? "Скопировано" : "Скопировать название"}
                        style={buildSkinNameStyle(skin.nameColor)}
                        onClick={() => void handleCopyTitle(skin)}
                      >
                        {skin.title}
                      </button>
                    </div>
                  </td>
                  <td>
                    <button
                      className="price-cell-button"
                      type="button"
                      disabled={!skin.steamPageUrl}
                      onClick={() => openExternal(skin.steamPageUrl)}
                      title={buildPriceTooltip(UI_TEXT.sourceSteamShort, skin.steamUpdatedAt)}
                    >
                      <span className="price-cell-value">{skin.steamPriceText || "-"}</span>
                    </button>
                  </td>
                  <td>
                    <button
                      className="price-cell-button"
                      type="button"
                      disabled={!skin.lisSkinsPageUrl}
                      onClick={() => openExternal(skin.lisSkinsPageUrl)}
                      title={buildPriceTooltip(UI_TEXT.sourceLisSkinsShort, skin.lisSkinsUpdatedAt)}
                    >
                      <span className="price-cell-value">{skin.lisSkinsPriceText || "-"}</span>
                    </button>
                  </td>
                  <td className="actions-menu-column">
                    <div className="row-actions-menu" role="group">
                      <button
                        className="row-actions-toggle"
                        type="button"
                        ref={(node) => {
                          toggleRefs.current[skin.id] = node;
                        }}
                        aria-label="Действия"
                        aria-expanded={isMenuOpen}
                        onClick={() => setOpenMenuId((current) => (current === skin.id ? null : skin.id))}
                      >
                        •••
                      </button>
                    </div>
                  </td>
                </tr>
              );
            })}
          </tbody>
        </table>
      </div>
      {openMenuId && menuPosition
        ? createPortal(
            <div className="row-actions-dropdown row-actions-dropdown-portal" style={{ top: `${menuPosition.top}px`, left: `${menuPosition.left}px` }}>
              <button
                className="row-actions-item"
                type="button"
                disabled={Boolean(updatingSkinIds[openMenuId]) || isUpdatingAll || Boolean(deletingSkinIds[openMenuId])}
                onClick={() => {
                  setOpenMenuId(null);
                  void onRefreshOne(openMenuId);
                }}
              >
                {Boolean(updatingSkinIds[openMenuId]) ? UI_TEXT.updateOnePendingCompact : UI_TEXT.updateOne}
              </button>
              <button
                className="row-actions-item toolbar-button-danger-secondary secondary-danger"
                type="button"
                disabled={Boolean(updatingSkinIds[openMenuId]) || isUpdatingAll || Boolean(deletingSkinIds[openMenuId])}
                onClick={() => {
                  setOpenMenuId(null);
                  void onDelete(openMenuId);
                }}
              >
                {Boolean(deletingSkinIds[openMenuId]) ? UI_TEXT.deleteSkinPending : UI_TEXT.deleteSkin}
              </button>
            </div>,
            document.body,
          )
        : null}
    </div>
  );
};
