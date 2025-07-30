import React from 'react';
import { Container, Typography, Button, Box, Divider } from '@mui/material';
import CartItemList from '../components/CartItemList';
import { useCart } from '../contexts/CartContext';
import { useNavigate } from 'react-router-dom';

const Cart = () => {
  const { cart, total } = useCart();
  const navigate = useNavigate();

  if (cart.length === 0) {
    return (
      <Container maxWidth="md" sx={{ py: 4, textAlign: 'center' }}>
        <Typography variant="h5" gutterBottom>
          Your cart is empty
        </Typography>
        <Button 
          variant="contained" 
          color="primary"
          onClick={() => navigate('/')}
        >
          Browse Menu
        </Button>
      </Container>
    );
  }

  return (
    <Container maxWidth="md" sx={{ py: 4 }}>
      <Typography variant="h4" component="h1" gutterBottom>
        Your Cart
      </Typography>
      <CartItemList />
      <Divider sx={{ my: 3 }} />
      <Box sx={{ 
        display: 'flex', 
        justifyContent: 'space-between', 
        alignItems: 'center',
        mb: 3
      }}>
        <Typography variant="h5">
          Total: ${total.toFixed(2)}
        </Typography>
        <Button
          variant="contained"
          color="primary"
          size="large"
          onClick={() => navigate('/checkout')}
        >
          Proceed to Checkout
        </Button>
      </Box>
    </Container>
  );
};

export default Cart;