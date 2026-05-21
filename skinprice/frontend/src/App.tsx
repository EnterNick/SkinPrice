import React, { useEffect } from "react";
import { getAppSettings } from "./entities/skin/api/skinApi";
import { RouterProvider } from "./app/providers/router";
import { AppRouter } from "./app/router/AppRouter";

const App: React.FC = () => {
  useEffect(() => {
    void getAppSettings().catch(() => undefined);
  }, []);

  return (
    <RouterProvider>
      <div className="app">
        <main className="layout-content">
          <AppRouter />
        </main>
      </div>
    </RouterProvider>
  );
};

export default App;
