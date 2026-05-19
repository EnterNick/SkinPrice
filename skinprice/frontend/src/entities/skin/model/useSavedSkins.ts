import { useCallback, useEffect, useState } from "react";
import { UI_TEXT } from "../../../shared/config/uiText";
import { toApiError } from "../../../shared/api/errors";
import { getSavedSkins } from "../api/skinApi";
import type { SavedSkin } from "./types";

export const useSavedSkins = () => {
  const [items, setItems] = useState<SavedSkin[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  const loadSkins = useCallback(async () => {
    const response = await getSavedSkins();
    setItems(response.items);
    setError(null);
    setLoading(false);
    return response;
  }, []);

  useEffect(() => {
    void loadSkins().catch((err: unknown) => {
      const apiError = toApiError(err);
      setItems([]);
      setLoading(false);
      setError(apiError.message || UI_TEXT.errLoadSaved);
    });
  }, [loadSkins]);

  return { items, loading, error, loadSkins };
};
