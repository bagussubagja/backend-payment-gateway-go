import React from 'react';
import { Link } from 'react-router-dom';

const Home = () => {
  const token = localStorage.getItem('token');

  return (
    <div style={{ textAlign: 'center', maxWidth: '800px', margin: '0 auto' }}>
      <div style={{ 
        backgroundColor: 'white', 
        padding: '3rem 2rem', 
        borderRadius: '8px', 
        boxShadow: '0 2px 10px rgba(0,0,0,0.1)',
        marginBottom: '2rem'
      }}>
        <h1 style={{ fontSize: '3rem', marginBottom: '1rem', color: '#333' }}>Payment Gateway</h1>
        <p style={{ fontSize: '1.25rem', color: '#666', marginBottom: '2rem' }}>
          Simple and secure payment processing with Midtrans integration
        </p>
        
        {token ? (
          <div style={{ display: 'flex', gap: '1rem', justifyContent: 'center' }}>
            <Link 
              to="/dashboard" 
              style={{ 
                padding: '1rem 2rem',
                backgroundColor: '#007bff',
                color: 'white',
                textDecoration: 'none',
                borderRadius: '4px',
                fontSize: '1.1rem'
              }}
            >
              Go to Dashboard
            </Link>
            <Link 
              to="/payment" 
              style={{ 
                padding: '1rem 2rem',
                backgroundColor: '#28a745',
                color: 'white',
                textDecoration: 'none',
                borderRadius: '4px',
                fontSize: '1.1rem'
              }}
            >
              Create Payment
            </Link>
          </div>
        ) : (
          <div style={{ display: 'flex', gap: '1rem', justifyContent: 'center' }}>
            <Link 
              to="/login" 
              style={{ 
                padding: '1rem 2rem',
                backgroundColor: '#007bff',
                color: 'white',
                textDecoration: 'none',
                borderRadius: '4px',
                fontSize: '1.1rem'
              }}
            >
              Login
            </Link>
            <Link 
              to="/register" 
              style={{ 
                padding: '1rem 2rem',
                backgroundColor: '#28a745',
                color: 'white',
                textDecoration: 'none',
                borderRadius: '4px',
                fontSize: '1.1rem'
              }}
            >
              Register
            </Link>
          </div>
        )}
      </div>

      <div style={{ display: 'grid', gridTemplateColumns: 'repeat(auto-fit, minmax(250px, 1fr))', gap: '1rem' }}>
        <div style={{ 
          backgroundColor: 'white', 
          padding: '1.5rem', 
          borderRadius: '8px', 
          boxShadow: '0 2px 10px rgba(0,0,0,0.1)'
        }}>
          <h3 style={{ color: '#007bff', marginBottom: '1rem' }}>Secure Authentication</h3>
          <p style={{ color: '#666' }}>JWT-based authentication system for secure user management</p>
        </div>
        
        <div style={{ 
          backgroundColor: 'white', 
          padding: '1.5rem', 
          borderRadius: '8px', 
          boxShadow: '0 2px 10px rgba(0,0,0,0.1)'
        }}>
          <h3 style={{ color: '#28a745', marginBottom: '1rem' }}>Multiple Payment Methods</h3>
          <p style={{ color: '#666' }}>Support for regular payments and QRIS transactions</p>
        </div>
        
        <div style={{ 
          backgroundColor: 'white', 
          padding: '1.5rem', 
          borderRadius: '8px', 
          boxShadow: '0 2px 10px rgba(0,0,0,0.1)'
        }}>
          <h3 style={{ color: '#ffc107', marginBottom: '1rem' }}>Transaction History</h3>
          <p style={{ color: '#666' }}>Track and monitor all your payment transactions</p>
        </div>
      </div>
    </div>
  );
};

export default Home;