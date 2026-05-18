import React from "react";

type DataTableProps<T> = {
  items: T[];
  className?: string;
  renderItem: (item: T) => React.ReactNode;
};

export const DataTable = <T,>({ items, className = "container", renderItem }: DataTableProps<T>) => (
  <div className={className}>{items.map((item) => renderItem(item))}</div>
);
