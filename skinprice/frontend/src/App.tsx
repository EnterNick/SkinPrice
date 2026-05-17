import React, { useEffect, useState } from "react";
import { GetSavedSkins, type SavedSkinItem } from "./wailsjs/go/main/App";
import "./styles.css";

type SavedSkinsState = {
  items: SavedSkinItem[];
  loading: boolean;
  error: string | null;
};

const App: React.FC = () => {
  const [state, setState] = useState<SavedSkinsState>({
    items: [],
    loading: true,
    error: null,
  });

  useEffect(() => {
    void GetSavedSkins({ limit: 50, offset: 0 })
      .then((response) => {
        setState({ items: response.items, loading: false, error: null });
      })
      .catch((err: unknown) => {
        setState({
          items: [],
          loading: false,
          error: err instanceof Error ? err.message : "Не удалось загрузить сохранённые скины",
        });
      });
  }, []);

  if (state.loading) {
    return <div className="app"><div className="container">Загрузка...</div></div>;
  }

  if (state.error) {
    return <div className="app"><div className="container">Ошибка: {state.error}</div></div>;
  }

  return (
    <div className="app">
      <div className="container">
        {state.items.map((skin) => (
          <a key={skin.market_hash_name} href={skin.page_url} className="card" target="_blank" rel="noreferrer">
            <div className="image-wrapper">
              <img src={skin.icon_url} alt={skin.display_name} className="card-image" />
            </div>
            <div className="card-body">
              <h2 className="title">{skin.display_name}</h2>
              <p className="text">{skin.market_hash_name}</p>
            </div>
          </a>
        ))}
      </div>
    </div>
  );
};

export default App;
