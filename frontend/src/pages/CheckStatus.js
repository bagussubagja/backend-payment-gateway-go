import React, { useState } from 'react';
import { paymentAPI } from '../services/api';
import { useTheme } from '../components/ThemeContext';

const CheckStatus = () => {
  const [orderID, setOrderID] = useState('');
  const [transaction, setTransaction] = useState(null);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState('');
  const { colors } = useTheme();

  const handleSubmit = async (e) => {
    e.preventDefault();
    setLoading(true);
    setError('');
    setTransaction(null);

    try {
      const response = await paymentAPI.getStatus(orderID);
      setTransaction(response.data);
    } catch (err) {
      setError(err.response?.data?.error || 'Failed to fetch transaction status');
    } finally {
      setLoading(false);
    }
  };

  const getStatusColor = (status) => {
    switch (status?.toLowerCase()) {
      case 'settlement':
      case 'capture':
        return colors.success;
      case 'pending':
        return colors.warning;
      case 'cancel':
      case 'expire':
      case 'failure':
        return colors.danger;
      default:
        return colors.textSecondary;
    }
  };

  return (
    <div style={{ maxWidth: '600px', margin: '0 auto' }}>
      <h1 style={{ marginBottom: '2rem', color: colors.text }}>Check Transaction Status</h1>
      
      <div style={{ 
        backgroundColor: colors.cardBg, 
        padding: '2rem', 
        borderRadius: '8px', 
        boxShadow: '0 2px 10px rgba(0,0,0,0.1)',
        marginBottom: '2rem'
      }}>
        <form onSubmit={handleSubmit}>
          <div style={{ marginBottom: '1rem' }}>
            <label style={{ display: 'block', marginBottom: '0.5rem', fontWeight: '500', color: colors.text }}>
              Order ID
            </label>
            <input
              type="text"
              style={{
                width: '100%',
                padding: '0.75rem',
                border: `1px solid ${colors.border}`,
                borderRadius: '4px',
                fontSize: '1rem',
                backgroundColor: colors.cardBg,
                color: colors.text
              }}
              value={orderID}
              onChange={(e) => setOrderID(e.target.value)}
              placeholder="Enter order ID"
              required
            />
          </div>
          
          <button
            type="submit"
            disabled={loading}
            style={{
              width: '100%',
              padding: '0.75rem',
              backgroundColor: loading ? colors.textSecondary : colors.primary,
              color: 'white',
              border: 'none',
              borderRadius: '4px',
              fontSize: '1rem',
              cursor: loading ? 'not-allowed' : 'pointer'
            }}
          >
            {loading ? 'Checking...' : 'Check Status'}
          </button>
        </form>
      </div>

      {error && (
        <div style={{ 
          padding: '0.75rem', 
          backgroundColor: '#f8d7da', 
          color: '#721c24', 
          borderRadius: '4px', 
          marginBottom: '1rem' 
        }}>
          {error}
        </div>
      )}

      {transaction && (
        <div style={{ 
          backgroundColor: colors.cardBg, 
          padding: '2rem', 
          borderRadius: '8px', 
          boxShadow: '0 2px 10px rgba(0,0,0,0.1)'
        }}>
          <h2 style={{ marginBottom: '1rem', color: colors.text }}>Transaction Details</h2>
          <div style={{ display: 'grid', gap: '1rem' }}>
            <div>
              <strong style={{ color: colors.text }}>Order ID:</strong> 
              <span style={{ color: colors.textSecondary, marginLeft: '0.5rem' }}>
                {transaction.ID || transaction.id}
              </span>
            </div>
            <div>
              <strong style={{ color: colors.text }}>Amount:</strong> 
              <span style={{ color: colors.textSecondary, marginLeft: '0.5rem' }}>
                Rp {(transaction.Amount || transaction.amount)?.toLocaleString()}
              </span>
            </div>
            <div>
              <strong style={{ color: colors.text }}>Status:</strong> 
              <span style={{ 
                marginLeft: '0.5rem',
                padding: '0.25rem 0.5rem',
                borderRadius: '4px',
                backgroundColor: getStatusColor(transaction.Status || transaction.status),
                color: 'white',
                fontSize: '0.875rem'
              }}>
                {transaction.Status || transaction.status}
              </span>
            </div>
            <div>
              <strong style={{ color: colors.text }}>Created:</strong> 
              <span style={{ color: colors.textSecondary, marginLeft: '0.5rem' }}>
                {transaction.CreatedAt ? new Date(transaction.CreatedAt).toLocaleString() : 
                 transaction.created_at ? new Date(transaction.created_at).toLocaleString() : '-'}
              </span>
            </div>
            {transaction.Items && transaction.Items.length > 0 && (
              <div>
                <strong style={{ color: colors.text }}>Items:</strong>
                <div style={{ marginTop: '0.5rem' }}>
                  {transaction.Items.map((item, index) => (
                    <div key={index} style={{ 
                      padding: '0.5rem', 
                      backgroundColor: colors.bg, 
                      borderRadius: '4px', 
                      marginBottom: '0.25rem',
                      color: colors.textSecondary
                    }}>
                      {item.Name || item.name} - Qty: {item.Quantity || item.quantity} - 
                      Rp {(item.Price || item.price)?.toLocaleString()}
                    </div>
                  ))}
                </div>
              </div>
            )}
          </div>
        </div>
      )}
    </div>
  );
};

export default CheckStatus;