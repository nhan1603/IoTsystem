import React from 'react';
import { 
  AppBar, 
  Toolbar, 
  Typography, 
  Button, 
  Badge,
  Box,
  Container,
  MenuItem,
  IconButton
} from '@mui/material';
import { ShoppingCart, Logout, ReceiptLong } from '@mui/icons-material';
import { useNavigate, Link } from 'react-router-dom';
import { useCart } from '../contexts/CartContext';
import RestaurantMenuIcon from '@mui/icons-material/RestaurantMenu';

const Navbar = () => {
  const navigate = useNavigate();
  const { cart, clearCart } = useCart();
  const token = localStorage.getItem('token');
  
  const cartItemCount = cart.reduce((sum, item) => sum + item.quantity, 0);

  const handleLogout = () => {
    // Clear the authentication token
    localStorage.removeItem('token');
    // Clear the cart
    clearCart();
    // Redirect to login page
    navigate('/login');
  };

  // Don't show navbar on login page
  if (!token) {
    return null;
  }

  return (
    <AppBar position="sticky">
      <Container maxWidth="xl">
        <Toolbar disableGutters>
          <Typography 
            variant="h6" 
            component="div" 
            sx={{ flexGrow: 1, cursor: 'pointer' }}
            onClick={() => navigate('/')}
          >
            Campus Food Order
          </Typography>
          <Box sx={{ display: 'flex', gap: 2, alignItems: 'center' }}>
            <Button color="inherit" onClick={() => navigate('/')}>
              <RestaurantMenuIcon sx={{ mr: 1 }} />
              Menu
            </Button>
            <Button 
              color="inherit" 
              onClick={() => navigate('/cart')}
              startIcon={
                <Badge badgeContent={cartItemCount} color="error">
                  <ShoppingCart />
                </Badge>
              }
            >
              Cart
            </Button>
            <Button 
              color="inherit"
              onClick={() => navigate('/orders')}
              startIcon={<ReceiptLong />}
            >
              Orders
            </Button>
            <Button 
              color="inherit"
              onClick={handleLogout}
              startIcon={<Logout />}
            >
              Logout
            </Button>
          </Box>
        </Toolbar>
      </Container>
    </AppBar>
  );
};

export default Navbar;