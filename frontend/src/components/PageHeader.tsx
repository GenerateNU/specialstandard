import React from 'react';
import type { LucideIcon } from 'lucide-react';

interface PageHeaderProps {
  title: string;
  icon: LucideIcon;
  actions?: React.ReactNode;
  description?: string;
}

export function PageHeader({ title, icon: Icon, actions, description }: PageHeaderProps) {
  return (
    <header className="mb-8">
      <div className="flex items-center justify-between mb-4">
        <div className="flex flex-row items-center gap-3">
          <Icon className="w-8 h-8 text-accent" />
          <h1 className="text-4xl font-bold text-primary">{title}</h1>
        </div>
        {actions && <div>{actions}</div>}
      </div>
      {description && (
        <p className="text-secondary text-sm">{description}</p>
      )}
    </header>
  );
}
