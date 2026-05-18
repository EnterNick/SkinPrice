import React from "react";

type CardGridProps<T> = {
  items: T[];
  className?: string;
  renderItem: (item: T) => React.ReactNode;
};

export const CardGrid = <T,>({ items, className = "container", renderItem }: CardGridProps<T>) => (
  <div className={className}>{items.map((item) => renderItem(item))}</div>
);
