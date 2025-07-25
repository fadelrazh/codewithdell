'use client';

import Link from 'next/link';
import { ThemeToggle } from '@/components/ui/ThemeToggle';

export function Header() {
  return (
    <header className="bg-white dark:bg-gray-800 shadow-sm border-b border-gray-200 dark:border-gray-700">
      <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
        <div className="flex justify-between items-center h-16">
          {/* Logo */}
          <Link href="/" className="flex items-center space-x-2">
            <div 
              className="w-8 h-8 rounded-lg flex items-center justify-center"
              style={{ backgroundColor: 'var(--primary)' }}
            >
              <span className="text-white font-bold text-sm">C</span>
            </div>
            <span className="text-xl font-bold text-gray-900 dark:text-white">
              CodeWithDell
            </span>
          </Link>

          {/* Navigation */}
          <nav className="hidden md:flex space-x-8">
            <Link 
              href="/" 
              className="text-gray-700 dark:text-gray-300 hover:text-primary dark:hover:text-primary transition-colors"
              style={{
                '--tw-text-opacity': '1',
                color: 'rgb(55 65 81 / var(--tw-text-opacity))'
              } as React.CSSProperties}
            >
              Home
            </Link>
            <Link 
              href="/blog" 
              className="text-gray-700 dark:text-gray-300 hover:text-primary dark:hover:text-primary transition-colors"
              style={{
                '--tw-text-opacity': '1',
                color: 'rgb(55 65 81 / var(--tw-text-opacity))'
              } as React.CSSProperties}
            >
              Blog
            </Link>
          </nav>

          {/* Theme Toggle */}
          <div className="flex items-center space-x-4">
            <ThemeToggle />
          </div>
        </div>
      </div>
    </header>
  );
} 