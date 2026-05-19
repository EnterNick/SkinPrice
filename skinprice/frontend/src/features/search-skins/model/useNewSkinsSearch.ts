import { useRef, useState } from "react";
import { UI_TEXT } from "../../../shared/config/uiText";
import { getNewSkins, hasLisSkinsToken, setLisSkinsToken } from "../../../entities/skin/api/skinApi";
import type { NewSkin } from "../../../entities/skin/model/types";
import { formatErrorMessage } from "../../../shared/lib/error/formatErrorMessage";
import { STORAGE_KEYS } from "../../../shared/config/storage";

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
  const [tokenRequired, setTokenRequired] = useState(false);
  const [tokenSaving, setTokenSaving] = useState(false);
  const [tokenError, setTokenError] = useState<string | null>(null);
  const [tokenSaved, setTokenSaved] = useState(false);
  const pendingSearchRef = useRef<{ query: string; source: "steam" | "lisskins"; nextOffset: number; append: boolean } | null>(null);
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

    if (source === "lisskins") {
      const hasToken = await hasLisSkinsToken();
      if (!hasToken) {
        pendingSearchRef.current = { query: normalizedSearch, source, nextOffset, append };
        setTokenRequired(true);
        setTokenSaved(false);
        return;
      }
    }

    setTokenRequired(false);
    setTokenError(null);

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
      setError(formatErrorMessage(UI_TEXT.errLoadNew, err));
      if (!append) {
        setItems([]);
      }
      setHasMore(false);
      setLoading(false);
      setLoadingMore(false);
    }
  };

  const saveLisSkinsToken = async (token: string) => {
    setTokenSaving(true);
    setTokenError(null);
    setTokenSaved(false);
    try {
      await setLisSkinsToken(token.trim());
      window.localStorage.setItem(STORAGE_KEYS.lisSkinsTokenSaved, "1");
      setTokenSaved(true);
      setTokenRequired(false);
      const pending = pendingSearchRef.current;
      if (pending) {
        pendingSearchRef.current = null;
        await loadNewSkins(pending.query, pending.source, pending.nextOffset, pending.append);
      }
    } catch (err) {
      setTokenError(formatErrorMessage(UI_TEXT.errLisSkinsTokenSave, err));
    } finally {
      setTokenSaving(false);
    }
  };

  const loadNextPage = async () => {
    if (loading || loadingMore || !hasMore) return;
    await loadNewSkins(currentQuery, currentSource, offset, true);
  };

  return {
    items,
    loading,
    loadingMore,
    error,
    hasMore,
    hasSearched,
    currentQuery,
    currentSource,
    tokenRequired,
    tokenSaving,
    tokenError,
    tokenSaved,
    loadNewSkins,
    loadNextPage,
    saveLisSkinsToken,
  };
};
