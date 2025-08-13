const THEME_STORAGE_KEY = 'theme';

export const themeService = {
  getStoredTheme() {
    try {
      return localStorage.getItem(THEME_STORAGE_KEY);
    } catch (_) {
      return null;
    }
  },

  getPreferredTheme() {
    if (typeof window !== 'undefined' && window.matchMedia) {
      return window.matchMedia('(prefers-color-scheme: dark)').matches ? 'dark' : 'light';
    }
    return 'light';
  },

  getActiveTheme() {
    return this.getStoredTheme() || this.getPreferredTheme();
  },

  applyTheme(theme) {
    const root = document.documentElement;
    if (theme === 'dark') {
      root.classList.add('dark');
    } else {
      root.classList.remove('dark');
    }
  },

  setTheme(theme) {
    try {
      localStorage.setItem(THEME_STORAGE_KEY, theme);
    } catch (_) {
      // ignore storage errors
    }
    this.applyTheme(theme);
  },

  initTheme() {
    this.applyTheme(this.getActiveTheme());
  },

  isDark() {
    return this.getActiveTheme() === 'dark';
  },
};

export default themeService;

