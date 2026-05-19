import { useCallback, useRef, useState } from "react";
import { getNewSkins } from "../../../entities/skin/api/skinApi";
import type { NewSkin } from "../../../entities/skin/model/types";
import { UI_TEXT } from "../../../shared/config/uiText";
import { formatErrorMessage } from "../../../shared/lib/error/formatErrorMessage";

const PAGE_SIZE = 20;

const allowPaint = async () =>
  new Promise<void>((resolve) => {
    window.requestAnimationFrame(() => {
      window.setTimeout(resolve, 0);
    });
  });

export const useNewSkinsSearch = () => {
  const [items, setItems] = useState<NewSkin[]>([]);
  const [loading, setLoading] = useState(false);
  const [loadingMore, setLoadingMore] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const [currentQuery, setCurrentQuery] = useState("");
  const [offset, setOffset] = useState(0);
  const [hasMore, setHasMore] = useState(false);
  const [hasSearched, setHasSearched] = useState(false);
  const requestIdRef = useRef(0);

  const resetSearchState = useCallback(() => {
    requestIdRef.current += 1;
    setItems([]);
    setError(null);
    setHasMore(false);
    setHasSearched(false);
    setCurrentQuery("");
    setOffset(0);
  }, []);

  const loadNewSkins = useCallback(async (searchValue: string, nextOffset = 0, append = false) => {
    const normalizedSearch = searchValue.trim();
    const requestId = requestIdRef.current + 1;
    requestIdRef.current = requestId;

    setHasSearched(true);
    setCurrentQuery(normalizedSearch);
    if (append) {
      setLoadingMore(true);
    } else {
      setLoading(true);
      setItems([]);
      setError(null);
    }

    try {
      await allowPaint();
      const response = await getNewSkins(normalizedSearch, PAGE_SIZE, nextOffset);
      if (requestId !== requestIdRef.current) return;

      const normalizedQuery = normalizedSearch.toLowerCase();
      const filteredItems = normalizedQuery
        ? response.items.filter((item) => item.title.toLowerCase().includes(normalizedQuery) || item.name.toLowerCase().includes(normalizedQuery))
        : response.items;

      setItems((prev) => (append ? [...prev, ...filteredItems] : filteredItems));
      setOffset(nextOffset + response.items.length);
      setHasMore(nextOffset + response.items.length < response.total);
      setLoading(false);
      setLoadingMore(false);
    } catch (err: unknown) {
      if (requestId !== requestIdRef.current) return;
      setError(formatErrorMessage(UI_TEXT.errLoadNew, err));
      if (!append) {
        setItems([]);
      }
      setHasMore(false);
      setLoading(false);
      setLoadingMore(false);
    }
  }, []);

  const loadNextPage = useCallback(async () => {
    if (loading || loadingMore || !hasMore) return;
    await loadNewSkins(currentQuery, offset, true);
  }, [currentQuery, hasMore, loadNewSkins, loading, loadingMore, offset]);

  return {
    items,
    loading,
    loadingMore,
    error,
    hasMore,
    hasSearched,
    loadNewSkins,
    loadNextPage,
    resetSearchState,
  };
};
