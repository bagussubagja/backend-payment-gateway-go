import React, { useState, useEffect } from 'react';
import { paymentAPI } from '../services/api';

const History = () => {
  const [transactions, setTransactions] = useState([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState('');

  useEffect(() => {
    const fetchHistory = async () => {
      try {
        const response = await paymentAPI.getHistory();
        console.log('History response:', response.data);
        setTransactions(Array.isArray(response.data) ? response.data : []);
      } catch (err) {
        console.error('History error:', err);
        setError(err.response?.data?.error || err.response?.data?.message || 'Failed to fetch history');
      } finally {
        setLoading(false);
      }
    };

    fetchHistory();
  }, []);

  const getStatusColor = (status) => {
    switch (status?.toLowerCase()) {
      case 'settlement':
      case 'capture':
        return '#28a745';
      case 'pending':
        return '#ffc107';
      case 'cancel':
      case 'expire':
      case 'failure':
        return '#dc3545';
      default:
        return '#6c757d';
    }
  };

  if (loading) {
    return <div style={{ textAlign: 'center', padding: '2rem' }}>Loading...</div>;
  }

  return (
    <div style={{ maxWidth: '1000px', margin: '0 auto' }}>
      <h1 style={{ marginBottom: '2rem' }}>Transaction History</h1>
      
      {error && (
        <div style={{ padding: '0.75rem', backgroundColor: '#f8d7da', color: '#721c24', borderRadius: '4px', marginBottom: '1rem' }}>
          {error}
        </div>
      )}

      {transactions.length === 0 ? (
        <div style={{ 
          backgroundColor: 'white', 
          padding: '2rem', 
          borderRadius: '8px', 
          boxShadow: '0 2px 10px rgba(0,0,0,0.1)',
          textAlign: 'center'
        }}>
          <p>No transactions found.</p>
        </div>
      ) : (
        <div style={{ backgroundColor: 'white', borderRadius: '8px', boxShadow: '0 2px 10px rgba(0,0,0,0.1)', overflow: 'hidden' }}>
          <div style={{ overflowX: 'auto' }}>
            <table style={{ width: '100%', borderCollapse: 'collapse' }}>
              <thead>
                <tr style={{ backgroundColor: '#f8f9fa' }}>
                  <th style={{ padding: '1rem', textAlign: 'left', borderBottom: '1px solid #dee2e6' }}>Order ID</th>
                  <th style={{ padding: '1rem', textAlign: 'left', borderBottom: '1px solid #dee2e6' }}>Amount</th>
                  <th style={{ padding: '1rem', textAlign: 'left', borderBottom: '1px solid #dee2e6' }}>Status</th>
                  <th style={{ padding: '1rem', textAlign: 'left', borderBottom: '1px solid #dee2e6' }}>Payment Type</th>
                  <th style={{ padding: '1rem', textAlign: 'left', borderBottom: '1px solid #dee2e6' }}>Date</th>
                  <th style={{ padding: '1rem', textAlign: 'left', borderBottom: '1px solid #dee2e6' }}>Items</th>
                </tr>
              </thead>
              <tbody>
                {transactions.map((transaction, index) => (
                  <tr key={index} style={{ borderBottom: '1px solid #dee2e6' }}>
                    <td style={{ padding: '1rem' }}>{transaction.ID || transaction.id}</td>
                    <td style={{ padding: '1rem' }}>Rp {transaction.Amount?.toLocaleString() || transaction.amount?.toLocaleString()}</td>
                    <td style={{ padding: '1rem' }}>
                      <span style={{ 
                        padding: '0.25rem 0.5rem',
                        borderRadius: '4px',
                        backgroundColor: getStatusColor(transaction.Status || transaction.status),
                        color: 'white',
                        fontSize: '0.875rem'
                      }}>
                        {transaction.Status || transaction.status || 'Unknown'}
                      </span>
                    </td>
                    <td style={{ padding: '1rem' }}>-</td>
                    <td style={{ padding: '1rem' }}>
                      {transaction.CreatedAt ? new Date(transaction.CreatedAt).toLocaleDateString() : 
                       transaction.created_at ? new Date(transaction.created_at).toLocaleDateString() : '-'}
                    </td>
                    <td style={{ padding: '1rem' }}>
                      {transaction.Items && transaction.Items.length > 0 ? (
                        <div>
                          {transaction.Items.map((item, itemIndex) => (
                            <div key={itemIndex} style={{ fontSize: '0.875rem', marginBottom: '0.25rem' }}>
                              {item.Name || item.name} (x{item.Quantity || item.quantity})
                            </div>
                          ))}
                        </div>
                      ) : transaction.items && transaction.items.length > 0 ? (
                        <div>
                          {transaction.items.map((item, itemIndex) => (
                            <div key={itemIndex} style={{ fontSize: '0.875rem', marginBottom: '0.25rem' }}>
                              {item.Name || item.name} (x{item.Quantity || item.quantity})
                            </div>
                          ))}
                        </div>
                      ) : '-'}
                    </td>
                  </tr>
                ))}
              </tbody>
            </table>
          </div>
        </div>
      )}
    </div>
  );
};

export default History;