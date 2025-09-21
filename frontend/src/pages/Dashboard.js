import React, { useState, useEffect } from 'react';
import { userAPI } from '../services/api';

const Dashboard = () => {
  const [user, setUser] = useState(null);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    const fetchProfile = async () => {
      try {
        const response = await userAPI.getProfile();
        console.log('Profile response:', response.data);
        setUser(response.data);
      } catch (err) {
        console.error('Failed to fetch profile:', err);
      } finally {
        setLoading(false);
      }
    };

    fetchProfile();
  }, []);

  if (loading) {
    return <div style={{ textAlign: 'center', padding: '2rem' }}>Loading...</div>;
  }

  return (
    <div style={{ maxWidth: '800px', margin: '0 auto' }}>
      <h1 style={{ marginBottom: '2rem' }}>Dashboard</h1>
      
      {user && (
        <div style={{ 
          backgroundColor: 'white', 
          padding: '2rem', 
          borderRadius: '8px', 
          boxShadow: '0 2px 10px rgba(0,0,0,0.1)',
          marginBottom: '2rem'
        }}>
          <h2 style={{ marginBottom: '1rem', color: '#333' }}>Profile Information</h2>
          <div style={{ display: 'grid', gridTemplateColumns: '1fr 1fr', gap: '1rem' }}>
            <div>
              <strong>Full Name:</strong> {user.FullName || user.full_name}
            </div>
            <div>
              <strong>Username:</strong> {user.Username || user.username}
            </div>
            <div>
              <strong>Email:</strong> {user.Email || user.email}
            </div>
            <div>
              <strong>Phone:</strong> {user.PhoneNumber || user.phone_number}
            </div>
            <div>
              <strong>Address:</strong> {user.Address || user.address}
            </div>
            <div>
              <strong>City:</strong> {user.City || user.city}, {user.PostalCode || user.postal_code}
            </div>
          </div>
        </div>
      )}
      
      <div style={{ display: 'grid', gridTemplateColumns: 'repeat(auto-fit, minmax(250px, 1fr))', gap: '1rem' }}>
        <div style={{ 
          backgroundColor: 'white', 
          padding: '1.5rem', 
          borderRadius: '8px', 
          boxShadow: '0 2px 10px rgba(0,0,0,0.1)',
          textAlign: 'center'
        }}>
          <h3 style={{ color: '#007bff', marginBottom: '1rem' }}>Create Payment</h3>
          <p style={{ color: '#666', marginBottom: '1rem' }}>Create new payment transactions</p>
          <a href="/payment" style={{ 
            display: 'inline-block',
            padding: '0.5rem 1rem',
            backgroundColor: '#007bff',
            color: 'white',
            textDecoration: 'none',
            borderRadius: '4px'
          }}>
            Go to Payment
          </a>
        </div>
        
        <div style={{ 
          backgroundColor: 'white', 
          padding: '1.5rem', 
          borderRadius: '8px', 
          boxShadow: '0 2px 10px rgba(0,0,0,0.1)',
          textAlign: 'center'
        }}>
          <h3 style={{ color: '#28a745', marginBottom: '1rem' }}>Transaction History</h3>
          <p style={{ color: '#666', marginBottom: '1rem' }}>View your payment history</p>
          <a href="/history" style={{ 
            display: 'inline-block',
            padding: '0.5rem 1rem',
            backgroundColor: '#28a745',
            color: 'white',
            textDecoration: 'none',
            borderRadius: '4px'
          }}>
            View History
          </a>
        </div>
      </div>
    </div>
  );
};

export default Dashboard;