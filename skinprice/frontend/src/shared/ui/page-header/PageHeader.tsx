import React from "react";

type PageHeaderProps = {
  sectionLabel: string;
  title: string;
  actions?: React.ReactNode;
};

export const PageHeader: React.FC<PageHeaderProps> = ({ sectionLabel, title, actions }) => (
  <div className="page-header">
    <div className="page-header-content">
      <p className="section-label">{sectionLabel}</p>
      <h1 className="page-title">{title}</h1>
    </div>
    {actions ? <div className="page-header-actions">{actions}</div> : null}
  </div>
);
