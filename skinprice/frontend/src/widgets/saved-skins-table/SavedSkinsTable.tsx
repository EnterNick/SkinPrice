import React from "react";
import { UI_TEXT } from "../../shared/config/uiText";
import { openExternal } from "../../shared/lib/browser/openExternal";
import { formatUpdatedAt } from "../../shared/lib/date/formatUpdatedAt";
import { SkinThumb } from "../../entities/skin/ui/SkinThumb";
import type { SavedSkin } from "../../entities/skin/model/types";

type SavedSkinsTableProps = {
  items: SavedSkin[];
  isUpdatingAll: boolean;
  updatingSkinIds: Record<string, boolean>;
  deletingSkinIds: Record<string, boolean>;
  onRefreshOne: (skinId: string) => Promise<void> | void;
  onDelete: (skinId: string) => Promise<void> | void;
};

export const SavedSkinsTable: React.FC<SavedSkinsTableProps> = ({ items, isUpdatingAll, updatingSkinIds, deletingSkinIds, onRefreshOne, onDelete }) => (
  <div className="table-shell">
    <table className="skins-table">
      <thead>
        <tr>
          <th>Скин</th>
          <th>{UI_TEXT.priceLabel}</th>
          <th>{UI_TEXT.updatedLabel}</th>
          <th>Ссылка</th>
          <th>Действие</th>
        </tr>
      </thead>
      <tbody>
        {items.map((skin) => {
          const isUpdating = Boolean(updatingSkinIds[skin.id]);
          const isDeleting = Boolean(deletingSkinIds[skin.id]);
          return (
            <tr key={skin.id}>
              <td>
                <div className="skin-cell">
                  <SkinThumb imageUrl={skin.imageUrl} title={skin.title} />
                  <div className="skin-title">{skin.title}</div>
                </div>
              </td>
              <td>
                <button
                  className="price-cell-button"
                  type="button"
                  disabled={isUpdating || isUpdatingAll || isDeleting}
                  onClick={() => void onRefreshOne(skin.id)}
                  title={UI_TEXT.updateOne}
                >
                  <span className="price-cell-value">{skin.priceText || "-"}</span>
                </button>
              </td>
              <td>{formatUpdatedAt(skin.updatedAt)}</td>
              <td>
                <button className="table-link-button" type="button" onClick={() => openExternal(skin.pageUrl)}>
                  {UI_TEXT.openLink}
                </button>
              </td>
              <td>
                <div className="table-actions">
                  <button
                    className="toolbar-button"
                    type="button"
                    disabled={isUpdating || isUpdatingAll || isDeleting}
                    onClick={() => void onDelete(skin.id)}
                  >
                    {isDeleting ? UI_TEXT.deleteSkinPending : UI_TEXT.deleteSkin}
                  </button>
                </div>
              </td>
            </tr>
          );
        })}
      </tbody>
    </table>
  </div>
);
