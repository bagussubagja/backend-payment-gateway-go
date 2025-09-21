import React, { useState } from 'react';
import { paymentAPI } from '../services/api';

const Payment = () => {
  const [items, setItems] = useState([{ id: '', name: '', price: '', quantity: 1 }]);
  const [customerDetails, setCustomerDetails] = useState({
    first_name: '',
    last_name: '',
    email: '',
    phone: '',
    address: '',
    city: '',
    postal_code: ''
  });
  const [loading, setLoading] = useState(false);
  const [result, setResult] = useState(null);
  const [error, setError] = useState('');
  const [paymentType, setPaymentType] = useState('regular');

  const addItem = () => {
    setItems([...items, { id: '', name: '', price: '', quantity: 1 }]);
  };

  const removeItem = (index) => {
    setItems(items.filter((_, i) => i !== index));
  };

  const updateItem = (index, field, value) => {
    const newItems = [...items];
    newItems[index][field] = value;
    setItems(newItems);
  };

  const handleSubmit = async (e) => {
    e.preventDefault();
    setLoading(true);
    setError('');
    setResult(null);

    try {
      const formattedItems = items.map(item => ({
        ...item,
        price: parseInt(item.price),
        quantity: parseInt(item.quantity)
      }));

      let response;
      if (paymentType === 'qris') {
        response = await paymentAPI.createQris({ items: formattedItems });
      } else {
        response = await paymentAPI.create({
          items: formattedItems,
          customer_details: customerDetails
        });
      }
      
      setResult(response.data);
    } catch (err) {
      setError(err.response?.data?.message || 'Payment creation failed');
    } finally {
      setLoading(false);
    }
  };

  const inputStyle = {
    width: '100%',
    padding: '0.5rem',
    border: '1px solid #ddd',
    borderRadius: '4px'
  };

  return (
    <div style={{ maxWidth: '800px', margin: '0 auto' }}>
      <h1 style={{ marginBottom: '2rem' }}>Create Payment</h1>
      
      {error && (
        <div style={{ padding: '0.75rem', backgroundColor: '#f8d7da', color: '#721c24', borderRadius: '4px', marginBottom: '1rem' }}>
          {error}
        </div>
      )}

      {result && (
        <div style={{ padding: '1rem', backgroundColor: '#d4edda', color: '#155724', borderRadius: '4px', marginBottom: '1rem' }}>
          <h3>Payment Created Successfully!</h3>
          <p><strong>Order ID:</strong> {result.order_id}</p>
          {result.redirect_url && (
            <p><strong>Payment URL:</strong> <a href={result.redirect_url} target="_blank" rel="noopener noreferrer">Pay Now</a></p>
          )}
          {result.qr_code_url && (
            <div>
              <p><strong>QR Code:</strong></p>
              <img src={result.qr_code_url} alt="QR Code" style={{ maxWidth: '200px' }} />
              <p><strong>Expires:</strong> {result.expiry_time}</p>
            </div>
          )}
        </div>
      )}

      <form onSubmit={handleSubmit}>
        <div style={{ backgroundColor: 'white', padding: '1.5rem', borderRadius: '8px', boxShadow: '0 2px 10px rgba(0,0,0,0.1)', marginBottom: '1rem' }}>
          <h3 style={{ marginBottom: '1rem' }}>Payment Type</h3>
          <div style={{ display: 'flex', gap: '1rem', marginBottom: '1rem' }}>
            <label style={{ display: 'flex', alignItems: 'center', gap: '0.5rem' }}>
              <input
                type="radio"
                value="regular"
                checked={paymentType === 'regular'}
                onChange={(e) => setPaymentType(e.target.value)}
              />
              Regular Payment
            </label>
            <label style={{ display: 'flex', alignItems: 'center', gap: '0.5rem' }}>
              <input
                type="radio"
                value="qris"
                checked={paymentType === 'qris'}
                onChange={(e) => setPaymentType(e.target.value)}
              />
              QRIS Payment
            </label>
          </div>
        </div>

        <div style={{ backgroundColor: 'white', padding: '1.5rem', borderRadius: '8px', boxShadow: '0 2px 10px rgba(0,0,0,0.1)', marginBottom: '1rem' }}>
          <h3 style={{ marginBottom: '1rem' }}>Items</h3>
          {items.map((item, index) => (
            <div key={index} style={{ display: 'grid', gridTemplateColumns: '1fr 2fr 1fr 1fr auto', gap: '0.5rem', marginBottom: '0.5rem', alignItems: 'end' }}>
              <input
                type="text"
                placeholder="Item ID"
                style={inputStyle}
                value={item.id}
                onChange={(e) => updateItem(index, 'id', e.target.value)}
                required
              />
              <input
                type="text"
                placeholder="Item Name"
                style={inputStyle}
                value={item.name}
                onChange={(e) => updateItem(index, 'name', e.target.value)}
                required
              />
              <input
                type="number"
                placeholder="Price"
                style={inputStyle}
                value={item.price}
                onChange={(e) => updateItem(index, 'price', e.target.value)}
                required
              />
              <input
                type="number"
                placeholder="Qty"
                style={inputStyle}
                value={item.quantity}
                onChange={(e) => updateItem(index, 'quantity', e.target.value)}
                min="1"
                required
              />
              {items.length > 1 && (
                <button
                  type="button"
                  onClick={() => removeItem(index)}
                  style={{ padding: '0.5rem', backgroundColor: '#dc3545', color: 'white', border: 'none', borderRadius: '4px' }}
                >
                  ×
                </button>
              )}
            </div>
          ))}
          <button
            type="button"
            onClick={addItem}
            style={{ padding: '0.5rem 1rem', backgroundColor: '#28a745', color: 'white', border: 'none', borderRadius: '4px', marginTop: '0.5rem' }}
          >
            Add Item
          </button>
        </div>

        {paymentType === 'regular' && (
          <div style={{ backgroundColor: 'white', padding: '1.5rem', borderRadius: '8px', boxShadow: '0 2px 10px rgba(0,0,0,0.1)', marginBottom: '1rem' }}>
            <h3 style={{ marginBottom: '1rem' }}>Customer Details</h3>
            <div style={{ display: 'grid', gridTemplateColumns: '1fr 1fr', gap: '1rem', marginBottom: '1rem' }}>
              <input
                type="text"
                placeholder="First Name"
                style={inputStyle}
                value={customerDetails.first_name}
                onChange={(e) => setCustomerDetails({ ...customerDetails, first_name: e.target.value })}
                required
              />
              <input
                type="text"
                placeholder="Last Name"
                style={inputStyle}
                value={customerDetails.last_name}
                onChange={(e) => setCustomerDetails({ ...customerDetails, last_name: e.target.value })}
              />
            </div>
            <div style={{ display: 'grid', gridTemplateColumns: '1fr 1fr', gap: '1rem', marginBottom: '1rem' }}>
              <input
                type="email"
                placeholder="Email"
                style={inputStyle}
                value={customerDetails.email}
                onChange={(e) => setCustomerDetails({ ...customerDetails, email: e.target.value })}
                required
              />
              <input
                type="tel"
                placeholder="Phone"
                style={inputStyle}
                value={customerDetails.phone}
                onChange={(e) => setCustomerDetails({ ...customerDetails, phone: e.target.value })}
                required
              />
            </div>
            <input
              type="text"
              placeholder="Address"
              style={{ ...inputStyle, marginBottom: '1rem' }}
              value={customerDetails.address}
              onChange={(e) => setCustomerDetails({ ...customerDetails, address: e.target.value })}
              required
            />
            <div style={{ display: 'grid', gridTemplateColumns: '1fr 1fr', gap: '1rem' }}>
              <input
                type="text"
                placeholder="City"
                style={inputStyle}
                value={customerDetails.city}
                onChange={(e) => setCustomerDetails({ ...customerDetails, city: e.target.value })}
                required
              />
              <input
                type="text"
                placeholder="Postal Code"
                style={inputStyle}
                value={customerDetails.postal_code}
                onChange={(e) => setCustomerDetails({ ...customerDetails, postal_code: e.target.value })}
                required
              />
            </div>
          </div>
        )}

        <button
          type="submit"
          disabled={loading}
          style={{
            width: '100%',
            padding: '1rem',
            backgroundColor: loading ? '#ccc' : '#007bff',
            color: 'white',
            border: 'none',
            borderRadius: '4px',
            fontSize: '1rem',
            cursor: loading ? 'not-allowed' : 'pointer'
          }}
        >
          {loading ? 'Creating Payment...' : 'Create Payment'}
        </button>
      </form>
    </div>
  );
};

export default Payment;