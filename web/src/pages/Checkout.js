import React, { useState } from 'react';
import {
  Container,
  Paper,
  Typography,
  Box,
  Divider,
  List,
  ListItem,
  ListItemText,
  CircularProgress
} from '@mui/material';
import { PayPalScriptProvider, PayPalButtons } from '@paypal/react-paypal-js';
import { useCart } from '../contexts/CartContext';
import { useNavigate } from 'react-router-dom';
import { useAuthenticatedFetch } from '../utils/api';

const PAYPAL_CLIENT_ID = "Afbo85wWbwkEpevvCjTbzgVA2ibewJp6tiGL2Cp5gl561j4oTOPLJNf3zyo28Xrq5Q1_uIdmbEO1aMOK";

const Checkout = () => {
  const { cart, total, clearCart } = useCart();
  const navigate = useNavigate();
  const fetchWithAuth = useAuthenticatedFetch();
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState(null);
  const [orderId, setOrderId] = useState(null);
  const [paypalOrderId, setPaypalOrderId] = useState(null);

  // Redirect if cart is empty
  if (cart.length === 0) {
    navigate('/cart');
    return null;
  }

  const createOrder = async () => {
    try {
      const response = await fetchWithAuth('/api/authenticated/v1/paypal/create-order', {
        method: 'POST',
        body: JSON.stringify({
          items: cart.map(item => ({
            menu_item_id: item.id,
            quantity: item.quantity,
          }))
        })
      });
      if (!response.paypal_order_id) {
        throw new Error('Failed to create PayPal order');
      }
      setPaypalOrderId(response.paypal_order_id);
      setOrderId(response.order_id);
      return response.paypal_order_id;
    } catch (err) {
      setError('Failed to create PayPal order. Please try again.');
      throw err;
    }
  };

  const handlePaymentSuccess = async (paypalOrderId) => {
    setLoading(true);
    try {
      // Call backend to capture the PayPal order
      const response = await fetchWithAuth('/api/authenticated/v1/paypal/capture-order', {
        method: 'POST',
        body: JSON.stringify({
          paypal_order_id: paypalOrderId,
          order_id: orderId
        })
      });

      if (!response.success) {
        throw new Error('Failed to capture PayPal payment');
      }

      // Clear cart and redirect to success page
      clearCart();
      navigate('/order-success');
    } catch (err) {
      setError('Failed to process payment. Please try again.');
      console.error('Payment error:', err);
    } finally {
      setLoading(false);
    }
  };

  return (
    <PayPalScriptProvider options={{ 
      "client-id": PAYPAL_CLIENT_ID,
      currency: "GBP"
    }}>
      <Container maxWidth="md" sx={{ py: 4 }}>
        <Paper elevation={3} sx={{ p: 3 }}>
          <Typography variant="h5" gutterBottom>
            Order Summary
          </Typography>
          <List>
            {cart.map((item) => (
              <ListItem key={item.id}>
                <ListItemText
                  primary={item.name}
                  secondary={`Quantity: ${item.quantity}`}
                />
                <Typography>
                  £{(item.price * item.quantity).toFixed(2)}
                </Typography>
              </ListItem>
            ))}
          </List>
          <Divider sx={{ my: 2 }} />
          <Box sx={{ display: 'flex', justifyContent: 'space-between', mb: 3 }}>
            <Typography variant="h6">Total:</Typography>
            <Typography variant="h6">£{total.toFixed(2)}</Typography>
          </Box>

          {error && (
            <Typography color="error" sx={{ mb: 2 }}>
              {error}
            </Typography>
          )}

          {loading ? (
            <Box sx={{ display: 'flex', justifyContent: 'center', my: 2 }}>
              <CircularProgress />
            </Box>
          ) : (
            <PayPalButtons
              createOrder={async () => {
                return await createOrder();
              }}
              onApprove={async (data) => {
                await handlePaymentSuccess(data.orderID);
              }}
              onError={(err) => {
                setError('PayPal payment failed. Please try again.');
                console.error('PayPal Error:', err);
              }}
              style={{
                layout: "vertical",
                color: "gold",
                shape: "rect",
                label: "pay"
              }}
            />
          )}
        </Paper>
      </Container>
    </PayPalScriptProvider>
  );
};

export default Checkout; 