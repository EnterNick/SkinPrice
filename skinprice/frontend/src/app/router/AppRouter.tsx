import React from "react";
import { Route, Routes } from "react-router-dom";
import { ROUTES } from "./routes";
import { SavedSkinsPage } from "../../pages/saved-skins-page/ui/SavedSkinsPage";
import { NewSkinsPage } from "../../pages/new-skins-page/ui/NewSkinsPage";
import { NotFoundPage } from "../../pages/not-found-page/ui/NotFoundPage";
import { SettingsPage } from "../../pages/settings-page/ui/SettingsPage";

export const AppRouter: React.FC = () => (
  <Routes>
    <Route path={ROUTES.home} element={<SavedSkinsPage />} />
    <Route path={ROUTES.newSkins} element={<NewSkinsPage />} />
    <Route path={ROUTES.settings} element={<SettingsPage />} />
    <Route path="*" element={<NotFoundPage />} />
  </Routes>
);
