import React from "react";

type PageHeaderProps = {
  eyebrow: string;
  title: string;
  actions?: React.ReactNode;
};

export const PageHeader: React.FC<PageHeaderProps> = ({ eyebrow, title, actions }) => (
  <div className="page-header">
    <div>
      <p className="eyebrow">{eyebrow}</p>
      <h1 className="page-title">{title}</h1>
    </div>
    {actions}
  </div>
);
