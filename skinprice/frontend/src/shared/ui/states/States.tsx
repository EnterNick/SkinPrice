import React from "react";

export const StatusBox: React.FC<{ text: string; className?: string }> = ({ text, className }) => (
  <div className={className ? `status-box ${className}` : "status-box"}>{text}</div>
);

export const LoadingState: React.FC<{ text: string }> = ({ text }) => <StatusBox text={text} />;

export const EmptyState: React.FC<{ text: string }> = ({ text }) => <StatusBox text={text} />;

export const ErrorState: React.FC<{ text: string }> = ({ text }) => <StatusBox text={`Ошибка: ${text}`} />;

export const ToastAlert: React.FC<{ type: "success" | "warning" | "error"; text: string }> = ({ type, text }) => (
  <StatusBox text={text} className={`notice-${type}`} />
);
