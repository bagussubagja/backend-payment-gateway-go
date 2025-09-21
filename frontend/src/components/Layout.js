import React from 'react';
import { Link, useNavigate } from 'react-router-dom';
import { useTheme } from './ThemeContext';

const Layout = ({ children }) => {
  const navigate = useNavigate();
  const token = localStorage.getItem('token');
  const { colors, isDark, toggleTheme } = useTheme();

  const handleLogout = () => {
    localStorage.removeItem('token');
    navigate('/login');
  };

  return (
    <div style={{ minHeight: '100vh', backgroundColor: colors.bg }}>
      <nav style={{ 
        backgroundColor: colors.cardBg, 
        padding: '1rem 2rem', 
        boxShadow: '0 2px 4px rgba(0,0,0,0.1)',
        display: 'flex',
        justifyContent: 'space-between',
        alignItems: 'center'
      }}>
        <Link to="/" style={{ fontSize: '1.5rem', fontWeight: 'bold', textDecoration: 'none', color: colors.text }}>
          Payment Gateway
        </Link>
        <div style={{ display: 'flex', gap: '1rem', alignItems: 'center' }}>
          <button onClick={toggleTheme} style={{ 
            padding: '0.5rem', 
            backgroundColor: 'transparent', 
            color: colors.text, 
            border: `1px solid ${colors.border}`, 
            borderRadius: '4px',
            cursor: 'pointer'
          }}>
            {isDark ? 'Light' : 'Dark'}
          </button>
          {token ? (
            <>
              <Link to="/dashboard" style={{ textDecoration: 'none', color: colors.primary }}>Dashboard</Link>
              <Link to="/payment" style={{ textDecoration: 'none', color: colors.primary }}>Create Payment</Link>
              <Link to="/history" style={{ textDecoration: 'none', color: colors.primary }}>History</Link>
              <Link to="/status" style={{ textDecoration: 'none', color: colors.primary }}>Check Status</Link>
              <button onClick={handleLogout} style={{ 
                padding: '0.5rem 1rem', 
                backgroundColor: colors.danger, 
                color: 'white', 
                border: 'none', 
                borderRadius: '4px',
                cursor: 'pointer'
              }}>
                Logout
              </button>
            </>
          ) : (
            <>
              <Link to="/login" style={{ textDecoration: 'none', color: colors.primary }}>Login</Link>
              <Link to="/register" style={{ textDecoration: 'none', color: colors.primary }}>Register</Link>
            </>
          )}
        </div>
      </nav>
      <main style={{ padding: '2rem' }}>
        {children}
      </main>
    </div>
  );
};

export default Layout;