import { toApiError } from "../api/errors";

export const formatErrorMessage = (prefix: string, error: unknown): string => {
  const apiError = toApiError(error);
  return `${prefix} ${apiError.message}`;
};
