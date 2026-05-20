import React, { useEffect, useRef, useState } from "react";
import { UI_TEXT } from "../../shared/config/uiText";
import { ClipboardSetText } from "../../wailsjs/runtime/runtime";
import { openExternal } from "../../shared/lib/browser/openExternal";
import { formatUpdatedAt } from "../../shared/lib/date/formatUpdatedAt";
import { CardGrid } from "../../shared/ui/card-grid/CardGrid";
import { SkinThumb } from "../../entities/skin/ui/SkinThumb";
import type { SavedSkin } from "../../entities/skin/model/types";

type SavedSkinsCardsProps = {
  items: SavedSkin[];
  isUpdatingAll: boolean;
  updatingSkinIds: Record<string, boolean>;
  deletingSkinIds: Record<string, boolean>;
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

export const SavedSkinsCards: React.FC<SavedSkinsCardsProps> = ({
  items,
  isUpdatingAll,
  updatingSkinIds,
  deletingSkinIds,
  onRefreshOne,
  onDelete,
}) => {
  const [copiedSkinId, setCopiedSkinId] = useState<string | null>(null);
  const copiedResetTimeoutRef = useRef<number | null>(null);

  useEffect(() => () => {
    if (copiedResetTimeoutRef.current !== null) {
      window.clearTimeout(copiedResetTimeoutRef.current);
    }
  }, []);

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
    <CardGrid
      className="saved-skins-grid"
      items={items}
      renderItem={(skin) => {
        const isUpdating = Boolean(updatingSkinIds[skin.id]);
        const isDeleting = Boolean(deletingSkinIds[skin.id]);
        const isBusy = isUpdating || isDeleting || isUpdatingAll;

        return (
          <article key={skin.id} className="saved-skin-card">
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

            <div className="saved-skin-card-prices">
              <button
                className="saved-skin-price-card"
                type="button"
                disabled={!skin.steamPageUrl}
                title={buildPriceTooltip(UI_TEXT.sourceSteamShort, skin.steamUpdatedAt)}
                onClick={() => openExternal(skin.steamPageUrl)}
              >
                <span className="saved-skin-price-source">{UI_TEXT.steamPriceLabel}</span>
                <span className="saved-skin-price-value">{skin.steamPriceText || "-"}</span>
                <span className="saved-skin-price-updated">{formatUpdatedAt(skin.steamUpdatedAt)}</span>
              </button>
              <button
                className="saved-skin-price-card"
                type="button"
                disabled={!skin.lisSkinsPageUrl}
                title={buildPriceTooltip(UI_TEXT.sourceLisSkinsShort, skin.lisSkinsUpdatedAt)}
                onClick={() => openExternal(skin.lisSkinsPageUrl)}
              >
                <span className="saved-skin-price-source">{UI_TEXT.lisSkinsPriceLabel}</span>
                <span className="saved-skin-price-value">{skin.lisSkinsPriceText || "-"}</span>
                <span className="saved-skin-price-updated">{formatUpdatedAt(skin.lisSkinsUpdatedAt)}</span>
              </button>
            </div>

            <div className="saved-skin-card-actions">
              <button className="toolbar-button toolbar-button-secondary" type="button" disabled={isBusy} onClick={() => void onRefreshOne(skin.id)}>
                {isUpdating ? UI_TEXT.updateOnePendingCompact : UI_TEXT.updateOne}
              </button>
              <button className="toolbar-button toolbar-button-danger-secondary" type="button" disabled={isBusy} onClick={() => void onDelete(skin.id)}>
                {isDeleting ? UI_TEXT.deleteSkinPending : UI_TEXT.deleteSkin}
              </button>
            </div>
          </article>
        );
      }}
    />
  );
};
