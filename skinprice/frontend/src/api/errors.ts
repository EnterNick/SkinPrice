import type { ApiError, ApiErrorCode } from "./models";

const MESSAGES: Record<ApiErrorCode, string> = {
  NETWORK_ERROR: "Проблема с подключением. Проверьте интернет и попробуйте снова.",
  NOT_FOUND: "Данные не найдены.",
  VALIDATION_ERROR: "Некорректные параметры запроса.",
  RATE_LIMIT: "Слишком много запросов. Повторите позже.",
  SERVER_ERROR: "Сервер временно недоступен. Попробуйте позже.",
  UNKNOWN_ERROR: "Произошла непредвиденная ошибка. Попробуйте снова.",
};

const mapCodeFromMessage = (message: string): ApiErrorCode => {
  const normalized = message.toLowerCase();

  if (normalized.includes("network") || normalized.includes("failed to fetch")) return "NETWORK_ERROR";
  if (normalized.includes("not found") || normalized.includes("не найден")) return "NOT_FOUND";
  if (normalized.includes("invalid") || normalized.includes("validation") || normalized.includes("некоррект")) return "VALIDATION_ERROR";
  if (normalized.includes("rate") || normalized.includes("too many")) return "RATE_LIMIT";
  if (normalized.includes("internal") || normalized.includes("server") || normalized.includes("timeout")) return "SERVER_ERROR";

  return "UNKNOWN_ERROR";
};

export const toApiError = (err: unknown): ApiError => {
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
