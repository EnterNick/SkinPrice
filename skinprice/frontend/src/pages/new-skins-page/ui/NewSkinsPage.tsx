import React, { useEffect, useState } from "react";
import { useNavigate } from "react-router-dom";
import { ROUTES } from "../../../app/router/routes";
import { saveSkin } from "../../../entities/skin/api/skinApi";
import type { NewSkin } from "../../../entities/skin/model/types";
import { useNewSkinsSearch } from "../../../features/search-skins/model/useNewSkinsSearch";
import { MIN_SEARCH_LENGTH } from "../../../shared/config/search";
import { UI_TEXT } from "../../../shared/config/uiText";
import { formatErrorMessage } from "../../../shared/lib/error/formatErrorMessage";
import { PageHeader } from "../../../shared/ui/page-header/PageHeader";
import { EmptyState, ErrorState, LoadingState, ToastAlert } from "../../../shared/ui/states/States";
import { NewSkinsResults } from "../../../widgets/new-skins-results/NewSkinsResults";
import { NewSkinsSearchPanel } from "../../../widgets/new-skins-search-panel/NewSkinsSearchPanel";

export const NewSkinsPage: React.FC = () => {
  const navigate = useNavigate();
  const {
    items,
    loading,
    loadingMore,
    error,
    hasMore,
    hasSearched,
    loadNewSkins,
    loadNextPage,
    resetSearchState,
  } = useNewSkinsSearch();
  const [savingIds, setSavingIds] = useState<Record<string, boolean>>({});
  const [savedIds, setSavedIds] = useState<Record<string, boolean>>({});
  const [notice, setNotice] = useState<{ type: "success" | "error"; text: string } | null>(null);
  const [query, setQuery] = useState("");

  const onSave = async (skin: NewSkin) => {
    setSavingIds((prev) => ({ ...prev, [skin.id]: true }));
    setNotice(null);
    try {
      const result = await saveSkin(skin);
      setSavedIds((prev) => ({ ...prev, [skin.id]: true }));
      setNotice({ type: "success", text: result.created ? UI_TEXT.saveDoneShort : UI_TEXT.saveSkinAlready });
    } catch (err: unknown) {
      setNotice({ type: "error", text: formatErrorMessage(UI_TEXT.errSave, err) });
    } finally {
      setSavingIds((prev) => ({ ...prev, [skin.id]: false }));
    }
  };

  useEffect(() => {
    const value = query.trim();
    if (value.length === 0) {
      setNotice(null);
      void loadNewSkins("");
      return;
    }

    if (value.length < MIN_SEARCH_LENGTH) {
      setNotice(null);
      resetSearchState();
      return;
    }

    const timeoutId = window.setTimeout(() => {
      setNotice(null);
      void loadNewSkins(value);
    }, 1000);

    return () => {
      window.clearTimeout(timeoutId);
    };
  }, [loadNewSkins, query, resetSearchState]);

  return (
    <div className="new-skins-page">
      <PageHeader
        sectionLabel={UI_TEXT.searchEyebrow}
        title={UI_TEXT.newSkinsTitle}
        actions={
          <div className="toolbar-group">
            <button className="toolbar-button" type="button" onClick={() => navigate(ROUTES.home)}>
              {UI_TEXT.ctaToHome}
            </button>
          </div>
        }
      />
      <NewSkinsSearchPanel
        value={query}
        errorText={notice?.type === "error" ? notice.text : null}
        disabled={false}
        onChange={setQuery}
      />
      {notice && <ToastAlert type={notice.type} text={notice.text} />}
      {loading && <LoadingState text={UI_TEXT.loadingNew} />}
      {error && !loading && <ErrorState text={error} />}
      {!loading && !error && hasSearched && items.length === 0 && <EmptyState text={UI_TEXT.notFoundSearch} />}
      {!loading && !error && items.length > 0 && (
        <NewSkinsResults
          items={items}
          loadingMore={loadingMore}
          hasMore={hasMore}
          savingIds={savingIds}
          savedIds={savedIds}
          onLoadMore={loadNextPage}
          onSave={onSave}
        />
      )}
    </div>
  );
};
