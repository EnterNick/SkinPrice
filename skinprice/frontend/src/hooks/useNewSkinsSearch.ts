import { useEffect, useRef, useState } from "react";
import { getNewSkins } from "../api/client";
import { toApiError } from "../api/errors";
import type { NewSkin } from "../api/models";
import { UI_TEXT } from "../constants/uiText";

export const useNewSkinsSearch = () => {
  const [query, setQuery] = useState("");
  const [debouncedQuery, setDebouncedQuery] = useState("");
  const [items, setItems] = useState<NewSkin[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const requestIdRef = useRef(0);

  useEffect(() => {
    const timeout = window.setTimeout(() => setDebouncedQuery(query.trim()), 400);
    return () => window.clearTimeout(timeout);
  }, [query]);

  const loadNewSkins = async (searchValue: string) => {
    const requestId = requestIdRef.current + 1;
    requestIdRef.current = requestId;
    setLoading(true);
    setError(null);

    try {
      const response = await getNewSkins(searchValue);
      if (requestId !== requestIdRef.current) return;

      const normalizedQuery = searchValue.toLowerCase();
      const filteredItems = normalizedQuery
        ? response.items.filter((item) => item.title.toLowerCase().includes(normalizedQuery) || item.name.toLowerCase().includes(normalizedQuery))
        : response.items;

      setItems(filteredItems);
      setLoading(false);
    } catch (err: unknown) {
      if (requestId !== requestIdRef.current) return;
      const apiError = toApiError(err);
      setError(apiError.message || UI_TEXT.errLoadNew);
      setItems([]);
      setLoading(false);
    }
  };

  useEffect(() => {
    void loadNewSkins(debouncedQuery);
  }, [debouncedQuery]);

  return { query, setQuery, debouncedQuery, items, loading, error, loadNewSkins };
};
