const skinNameColorPattern = /^[0-9a-fA-F]{6}$/;

export const toSkinNameColor = (value?: string | null): string | undefined => {
  if (!value || !skinNameColorPattern.test(value)) {
    return undefined;
  }

  return `#${value}`;
};
