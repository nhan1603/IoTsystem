import React from 'react';
import { Card, CardContent, CardMedia, Typography, Button, Box } from '@mui/material';
import { useCart } from '../contexts/CartContext';

const MenuItemCard = ({ item }) => {
  const { addToCart } = useCart();

  return (
    <Card sx={{ height: '100%', display: 'flex', flexDirection: 'column' }}>
      <CardMedia
        component="img"
        height="200"
        image={item.imageUrl || 'https://via.placeholder.com/300x200?text=Food+Item'}
        alt={item.name}
      />
      <CardContent sx={{ flexGrow: 1 }}>
        <Typography gutterBottom variant="h5" component="h2">
          {item.name}
        </Typography>
        <Typography variant="body2" color="text.secondary">
          {item.description}
        </Typography>
        <Box sx={{ mt: 2, display: 'flex', justifyContent: 'space-between', alignItems: 'center' }}>
          <Typography variant="h6" color="primary">
            £{item.price.toFixed(2)}
          </Typography>
          <Button 
            variant="contained" 
            onClick={() => addToCart(item)}
            disabled={!item.isAvailable}
          >
            Add to Cart
          </Button>
        </Box>
      </CardContent>
    </Card>
  );
};

export default MenuItemCard;