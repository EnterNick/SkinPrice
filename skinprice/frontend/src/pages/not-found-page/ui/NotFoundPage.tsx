import React from "react";
import { UI_TEXT } from "../../../shared/config/uiText";

export const NotFoundPage: React.FC = () => (
  <div className="container">
    <div className="card">
      <div className="card-body">
        <h2 className="title">404</h2>
        <p className="text">{UI_TEXT.pageNotFound}</p>
      </div>
    </div>
  </div>
);
