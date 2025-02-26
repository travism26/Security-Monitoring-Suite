import React from "react";

export function SystemHealthSkeleton() {
  const skeletonCard = (
    <div className="bg-white dark:bg-gray-800 rounded-lg shadow-sm p-4 animate-pulse">
      <div className="flex items-center justify-between mb-2">
        <div className="h-4 bg-gray-200 dark:bg-gray-700 rounded w-24" />
        <div className="h-3 w-3 bg-gray-200 dark:bg-gray-700 rounded-full" />
      </div>
      <div className="flex items-baseline">
        <div className="h-8 bg-gray-200 dark:bg-gray-700 rounded w-16" />
      </div>
      <div className="mt-1 h-3 bg-gray-200 dark:bg-gray-700 rounded w-20" />
    </div>
  );

  return (
    <div className="space-y-4">
      <div className="flex items-center justify-between mb-4">
        <div className="h-6 bg-gray-200 dark:bg-gray-700 rounded w-32" />
        <div className="h-4 bg-gray-200 dark:bg-gray-700 rounded w-24" />
      </div>
      <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-4">
        {[...Array(4)].map((_, i) => (
          <div key={i}>{skeletonCard}</div>
        ))}
      </div>
    </div>
  );
}
