import React, { useEffect } from "react";
import { createPortal } from "react-dom";

export const StatusBox: React.FC<{ text: string; className?: string }> = ({ text, className }) => (
  <div className={className ? `status-box ${className}` : "status-box"}>{text}</div>
);

export const LoadingState: React.FC<{ text: string }> = ({ text }) => <StatusBox text={text} />;

export const EmptyState: React.FC<{ text: string }> = ({ text }) => <StatusBox text={text} />;

export const ErrorState: React.FC<{ text: string }> = ({ text }) => <StatusBox text={`Ошибка: ${text}`} />;

export const ToastAlert: React.FC<{ type: "success" | "warning" | "error"; text: string; onClose?: () => void }> = ({
  type,
  text,
  onClose,
}) => {
  useEffect(() => {
    if (!onClose) return undefined;

    const timeoutId = window.setTimeout(() => {
      onClose();
    }, 5000);

    return () => {
      window.clearTimeout(timeoutId);
    };
  }, [onClose, text, type]);

  return createPortal(
    <div className="toast-stack" role="status" aria-live="polite">
      <div className={`toast-alert toast-alert-${type}`}>
        <div className="toast-alert-body">{text}</div>
        <button className="toast-alert-close" type="button" aria-label="Закрыть уведомление" onClick={onClose}>
          ×
        </button>
      </div>
    </div>,
    document.body,
  );
};
