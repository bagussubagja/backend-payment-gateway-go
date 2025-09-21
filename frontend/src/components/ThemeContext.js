import React, { createContext, useContext, useState, useEffect } from 'react';

const ThemeContext = createContext();

export const useTheme = () => useContext(ThemeContext);

export const ThemeProvider = ({ children }) => {
  const [isDark, setIsDark] = useState(() => {
    const saved = localStorage.getItem('theme');
    return saved ? saved === 'dark' : false;
  });

  useEffect(() => {
    localStorage.setItem('theme', isDark ? 'dark' : 'light');
  }, [isDark]);

  const toggleTheme = () => setIsDark(!isDark);

  const theme = {
    isDark,
    toggleTheme,
    colors: {
      bg: isDark ? '#1a1a1a' : '#f5f5f5',
      cardBg: isDark ? '#2d2d2d' : '#ffffff',
      text: isDark ? '#ffffff' : '#333333',
      textSecondary: isDark ? '#cccccc' : '#666666',
      border: isDark ? '#404040' : '#dee2e6',
      primary: '#007bff',
      success: '#28a745',
      danger: '#dc3545',
      warning: '#ffc107'
    }
  };

  return (
    <ThemeContext.Provider value={theme}>
      {children}
    </ThemeContext.Provider>
  );
};