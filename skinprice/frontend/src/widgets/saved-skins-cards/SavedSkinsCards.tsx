import React, { useEffect, useRef, useState } from "react";
import { createPortal } from "react-dom";
import { UI_TEXT } from "../../shared/config/uiText";
import { ClipboardSetText } from "../../wailsjs/runtime/runtime";
import { openExternal } from "../../shared/lib/browser/openExternal";
import { formatUpdatedAt } from "../../shared/lib/date/formatUpdatedAt";
import { CardGrid } from "../../shared/ui/card-grid/CardGrid";
import { SkinThumb } from "../../entities/skin/ui/SkinThumb";
import type { PriceSnapshot, SavedSkin } from "../../entities/skin/model/types";

type SavedSkinsCardsProps = {
  items: SavedSkin[];
  isUpdatingAll: boolean;
  updatingSkinIds: Record<string, boolean>;
  deletingSkinIds: Record<string, boolean>;
  sources: Array<{ source: string; label: string }>;
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

const getSourcePrice = (prices: PriceSnapshot[], source: string): PriceSnapshot | undefined => prices.find((price) => price.source === source);

const buildSkinNameStyle = (nameColor?: string): React.CSSProperties | undefined =>
  nameColor ? { "--skin-name-color": nameColor } as React.CSSProperties : undefined;

const ROW_ACTIONS_DROPDOWN_MIN_WIDTH = 0;
const ROW_ACTIONS_DROPDOWN_MAX_WIDTH = 220;

export const SavedSkinsCards: React.FC<SavedSkinsCardsProps> = ({
  items,
  isUpdatingAll,
  updatingSkinIds,
  deletingSkinIds,
  sources,
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
      const estimatedWidth = Math.min(
        ROW_ACTIONS_DROPDOWN_MAX_WIDTH,
        Math.max(ROW_ACTIONS_DROPDOWN_MIN_WIDTH, rect.width + 140),
      );
      const nextLeft = Math.min(
        Math.max(8, rect.right - estimatedWidth),
        Math.max(8, window.innerWidth - estimatedWidth - 8),
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
    <>
      <CardGrid
        className="saved-skins-grid"
        items={items}
        renderItem={(skin) => {
          const isMenuOpen = openMenuId === skin.id;

          return (
            <article key={skin.id} className="saved-skin-card">
              <div className="saved-skin-card-top">
                <div className="saved-skin-card-header">
                  <SkinThumb imageUrl={skin.imageUrl} title={skin.title} />
                  <button
                    className={`skin-title-button saved-skin-card-title${copiedSkinId === skin.id ? " skin-title-button-copied" : ""}`}
                    type="button"
                    title={copiedSkinId === skin.id ? "Скопировано" : "Скопировать название"}
                    style={buildSkinNameStyle(skin.nameColor)}
                    onClick={() => void handleCopyTitle(skin)}
                  >
                    {skin.title}
                  </button>
                </div>
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
              </div>

              <div className="saved-skin-card-prices">
                {sources.map((source) => {
                  const price = getSourcePrice(skin.prices, source.source);
                  return (
                    <button
                      key={source.source}
                      className="saved-skin-price-card"
                      type="button"
                      disabled={!price?.pageUrl}
                      title={buildPriceTooltip(source.label, price?.fetchedAt)}
                      onClick={() => price?.pageUrl && openExternal(price.pageUrl)}
                    >
                      <span className="saved-skin-price-source">{source.label}</span>
                      <span className="saved-skin-price-value">{price?.priceText || "-"}</span>
                    </button>
                  );
                })}
              </div>
            </article>
          );
        }}
      />
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
    </>
  );
};
