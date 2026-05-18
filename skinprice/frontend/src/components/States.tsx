import React from "react";

export const LoadingState: React.FC<{ text: string }> = ({ text }) => <div className="status-box">{text}</div>;

export const EmptyState: React.FC<{ text: string }> = ({ text }) => <div className="status-box">{text}</div>;

export const ErrorState: React.FC<{ text: string }> = ({ text }) => <div className="status-box">Ошибка: {text}</div>;

export const ToastAlert: React.FC<{ type: "success" | "warning" | "error"; text: string }> = ({ text }) => (
  <div className="status-box">{text}</div>
);
