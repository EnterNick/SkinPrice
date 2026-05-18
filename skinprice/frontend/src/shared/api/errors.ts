import type { ApiError, ApiErrorCode } from "../../entities/skin/model/types";

const MESSAGES: Record<ApiErrorCode, string> = {
  invalid_argument: "Некорректные параметры запроса.",
  not_found: "Данные не найдены.",
  already_exists: "Запись уже существует.",
  conflict: "Конфликт данных. Попробуйте снова.",
  unavailable: "Сервис временно недоступен. Попробуйте позже.",
  timeout: "Сервис не ответил вовремя. Попробуйте снова.",
  external: "Не удалось получить данные из внешнего источника.",
  internal: "Внутренняя ошибка приложения.",
  unknown: "Произошла непредвиденная ошибка. Попробуйте снова.",
  UNKNOWN_ERROR: "Произошла непредвиденная ошибка. Попробуйте снова.",
};

const mapCodeFromMessage = (message: string): ApiErrorCode => {
  const normalized = message.toLowerCase();

  if (normalized.includes("not_found") || normalized.includes("not found") || normalized.includes("не найден")) return "not_found";
  if (normalized.includes("invalid_argument") || normalized.includes("invalid") || normalized.includes("validation") || normalized.includes("некоррект")) return "invalid_argument";
  if (normalized.includes("already_exists") || normalized.includes("already exists")) return "already_exists";
  if (normalized.includes("conflict")) return "conflict";
  if (normalized.includes("unavailable") || normalized.includes("bad status")) return "unavailable";
  if (normalized.includes("timeout")) return "timeout";
  if (normalized.includes("external") || normalized.includes("steam")) return "external";
  if (normalized.includes("internal")) return "internal";
  if (normalized.includes("network") || normalized.includes("failed to fetch")) return "external";

  return "UNKNOWN_ERROR";
};

export const toApiError = (err: unknown): ApiError => {
  if (typeof err === "string") {
    try {
      const parsed = JSON.parse(err) as ApiError;
      if (parsed && typeof parsed.code === "string" && typeof parsed.message === "string") {
        return parsed;
      }
    } catch {
      // Ignore invalid JSON and map by message below.
    }
  }

  if (typeof err === "object" && err !== null && "code" in err && "message" in err) {
    const known = err as { code?: ApiErrorCode; message?: string; details?: unknown };
    const code = known.code ?? "UNKNOWN_ERROR";
    return { code, message: known.message ?? MESSAGES.UNKNOWN_ERROR, details: known.details };
  }

  const message = err instanceof Error ? err.message : String(err ?? "");
  const code = mapCodeFromMessage(message);

  return {
    code,
    message: code === "UNKNOWN_ERROR" && message ? message : MESSAGES[code],
    details: err,
  };
};
