import { LogClientEvent } from "../../../wailsjs/go/main/App";

export type ClientLogLevel = "debug" | "info" | "warn" | "error";

type ClientLogContext = Record<string, unknown>;

const writeConsole = (level: ClientLogLevel, message: string, context?: ClientLogContext): void => {
  const method =
    level === "debug" ? console.debug :
    level === "warn" ? console.warn :
    level === "error" ? console.error :
    console.info;

  if (context && Object.keys(context).length > 0) {
    method(message, context);
    return;
  }

  method(message);
};

const sendToBackend = async (
  level: ClientLogLevel,
  message: string,
  component?: string,
  context?: ClientLogContext,
): Promise<void> => {
  try {
    await LogClientEvent({
      level,
      message,
      component: component ?? "",
      context: context ?? {},
    });
  } catch (error) {
    console.error("failed to forward client log event", error);
  }
};

export const logClientEvent = (
  level: ClientLogLevel,
  message: string,
  component?: string,
  context?: ClientLogContext,
): void => {
  writeConsole(level, message, context);
  void sendToBackend(level, message, component, context);
};

export const installGlobalErrorLogging = (): void => {
  window.addEventListener("error", (event) => {
    logClientEvent("error", "window error", "window", {
      message: event.message,
      filename: event.filename,
      lineno: event.lineno,
      colno: event.colno,
    });
  });

  window.addEventListener("unhandledrejection", (event) => {
    const reason = event.reason instanceof Error ? event.reason.message : String(event.reason ?? "");
    logClientEvent("error", "unhandled promise rejection", "window", { reason });
  });
};
