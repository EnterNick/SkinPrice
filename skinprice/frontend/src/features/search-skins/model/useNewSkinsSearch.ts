import { useRef, useState } from "react";
import { UI_TEXT } from "../../../shared/config/uiText";
import { toApiError } from "../../../shared/api/errors";
import { getNewSkins } from "../../../entities/skin/api/skinApi";
import type { NewSkin } from "../../../entities/skin/model/types";

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
  const [currentSource, setCurrentSource] = useState<"steam" | "lisskins">("steam");
  const [offset, setOffset] = useState(0);
  const [hasMore, setHasMore] = useState(false);
  const [hasSearched, setHasSearched] = useState(false);
  const requestIdRef = useRef(0);

  const loadNewSkins = async (searchValue: string, source: "steam" | "lisskins", nextOffset = 0, append = false) => {
    const normalizedSearch = searchValue.trim();
    if (!normalizedSearch) {
      requestIdRef.current += 1;
      setHasSearched(false);
      setItems([]);
      setOffset(0);
      setHasMore(false);
      setError(null);
      setLoading(false);
      setLoadingMore(false);
      return;
    }

    const requestId = requestIdRef.current + 1;
    requestIdRef.current = requestId;
    setHasSearched(true);
    setCurrentQuery(normalizedSearch);
    setCurrentSource(source);
    if (append) {
      setLoadingMore(true);
    } else {
      setLoading(true);
      setItems([]);
      setError(null);
    }

    try {
      await allowPaint();
      const response = await getNewSkins(normalizedSearch, PAGE_SIZE, nextOffset, source);
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
      const apiError = toApiError(err);
      setError(apiError.message || UI_TEXT.errLoadNew);
      if (!append) {
        setItems([]);
      }
      setHasMore(false);
      setLoading(false);
      setLoadingMore(false);
    }
  };

  const loadNextPage = async () => {
    if (loading || loadingMore || !hasMore) return;
    await loadNewSkins(currentQuery, currentSource, offset, true);
  };

  return { items, loading, loadingMore, error, hasMore, hasSearched, currentQuery, currentSource, loadNewSkins, loadNextPage };
};
