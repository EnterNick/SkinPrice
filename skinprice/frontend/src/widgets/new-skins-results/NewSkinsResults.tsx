import React from "react";
import { CardGrid } from "../../shared/ui/card-grid/CardGrid";
import { LoadingState, StatusBox } from "../../shared/ui/states/States";
import { UI_TEXT } from "../../shared/config/uiText";
import { SkinCard } from "../../entities/skin/ui/SkinCard";
import type { NewSkin } from "../../entities/skin/model/types";

type NewSkinsResultsProps = {
  items: NewSkin[];
  loadingMore: boolean;
  hasMore: boolean;
  savingIds: Record<string, boolean>;
  savedIds: Record<string, boolean>;
  onLoadMore: () => Promise<void> | void;
  onSave: (skin: NewSkin) => Promise<void> | void;
};

export const NewSkinsResults: React.FC<NewSkinsResultsProps> = ({
  items,
  loadingMore,
  hasMore,
  savingIds,
  savedIds,
  onLoadMore,
  onSave,
}) => (
  <>
    <CardGrid
      className="search-results"
      items={items}
      renderItem={(skin) => (
        <SkinCard
          key={skin.id}
          skin={skin}
          saving={Boolean(savingIds[skin.id])}
          saved={Boolean(savedIds[skin.id])}
          onSave={onSave}
        />
      )}
    />
    {hasMore && !loadingMore && (
      <div className="load-more-wrap">
        <button className="toolbar-button toolbar-button-primary" type="button" onClick={() => void onLoadMore()}>
          {UI_TEXT.loadMore}
        </button>
      </div>
    )}
    {loadingMore && <LoadingState text={UI_TEXT.loadingMore} />}
    {!hasMore && <StatusBox text={UI_TEXT.listFullyLoaded} />}
  </>
);
