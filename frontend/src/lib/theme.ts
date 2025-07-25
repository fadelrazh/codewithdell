export type Theme = 'dark' | 'light' | 'system';

class ThemeManager {
  private static instance: ThemeManager;
  private currentTheme: Theme = 'system';
  private isSystemDark: boolean = false;

  private constructor() {
    this.init();
  }

  static getInstance(): ThemeManager {
    if (!ThemeManager.instance) {
      ThemeManager.instance = new ThemeManager();
    }
    return ThemeManager.instance;
  }

  private init() {
    // Only run on client side
    if (typeof window === 'undefined') {
      return;
    }

    // Load theme from localStorage
    const savedTheme = localStorage.getItem('theme') as Theme;
    if (savedTheme && ['dark', 'light', 'system'].includes(savedTheme)) {
      this.currentTheme = savedTheme;
    }

    // Listen for system theme changes
    this.isSystemDark = window.matchMedia('(prefers-color-scheme: dark)').matches;
    
    window.matchMedia('(prefers-color-scheme: dark)').addEventListener('change', (e) => {
      this.isSystemDark = e.matches;
      this.applyTheme();
    });

    this.applyTheme();
  }

  getTheme(): Theme {
    return this.currentTheme;
  }

  setTheme(theme: Theme) {
    this.currentTheme = theme;
    if (typeof window !== 'undefined') {
      localStorage.setItem('theme', theme);
    }
    this.applyTheme();
  }

  private applyTheme() {
    if (typeof window === 'undefined') return;

    const root = document.documentElement;
    const isDark = this.currentTheme === 'dark' || 
                   (this.currentTheme === 'system' && this.isSystemDark);

    // Remove all theme classes first
    root.classList.remove('dark', 'light');
    
    if (isDark) {
      root.classList.add('dark');
    } else {
      root.classList.add('light');
    }
  }

  isDark(): boolean {
    if (typeof window === 'undefined') {
      return false; // Default to light theme on server
    }
    return this.currentTheme === 'dark' || 
           (this.currentTheme === 'system' && this.isSystemDark);
  }

  toggleTheme() {
    const themes: Theme[] = ['light', 'dark', 'system'];
    const currentIndex = themes.indexOf(this.currentTheme);
    const nextIndex = (currentIndex + 1) % themes.length;
    this.setTheme(themes[nextIndex]);
  }
}

export const themeManager = ThemeManager.getInstance();