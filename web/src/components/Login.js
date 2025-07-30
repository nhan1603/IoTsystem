import React, { useState } from 'react';
import axios from '../axiosConfig';

function Login({ onLogin }) {
  const [isLogin, setIsLogin] = useState(true);
  const [email, setEmail] = useState('');
  const [password, setPassword] = useState('');
  const [error, setError] = useState('');

  const handleSubmit = async (e) => {
    e.preventDefault();
    setError('');
    try {
      const endpoint = isLogin ? '/api/public/v1/login' : '/api/public/v1/users';
      const response = await axios.post(endpoint, { email, password });
      
      if (isLogin) {
        const { token, user } = response.data;
        onLogin({ token, ...user });
      } else {
        setIsLogin(true);
        setError('Account created successfully. Please log in.');
      }
    } catch (err) {
      console.error(isLogin ? 'Login error:' : 'Signup error:', err);
      setError(isLogin ? 'Invalid credentials' : 'Failed to create account');
    }
  };

  return (
    <div>
      <h2>{isLogin ? 'Login' : 'Create Account'}</h2>
      <form onSubmit={handleSubmit}>
        <input
          type="email"
          value={email}
          onChange={(e) => setEmail(e.target.value)}
          placeholder="Email"
          required
        />
        <input
          type="password"
          value={password}
          onChange={(e) => setPassword(e.target.value)}
          placeholder="Password"
          required
        />
        {error && <p className="error">{error}</p>}
        <button type="submit">{isLogin ? 'Login' : 'Create Account'}</button>
      </form>
      <button onClick={() => setIsLogin(!isLogin)}>
        {isLogin ? 'Need to create an account?' : 'Already have an account?'}
      </button>
    </div>
  );
}

export default Login;