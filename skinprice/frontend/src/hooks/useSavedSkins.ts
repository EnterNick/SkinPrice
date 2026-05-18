import { useEffect, useState } from "react";
import type { SavedSkin } from "../api/models";
import { getSavedSkins } from "../api/client";
import { toApiError } from "../api/errors";
import { UI_TEXT } from "../constants/uiText";

export const useSavedSkins = () => {
  const [items, setItems] = useState<SavedSkin[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  const loadSkins = async () => {
    const response = await getSavedSkins();
    setItems(response.items);
    setError(null);
    setLoading(false);
    return response;
  };

  useEffect(() => {
    void loadSkins().catch((err: unknown) => {
      const apiError = toApiError(err);
      setItems([]);
      setLoading(false);
      setError(apiError.message || UI_TEXT.errLoadSaved);
    });
  }, []);

  return { items, loading, error, loadSkins };
};
