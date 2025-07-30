import React, { useState, useEffect } from 'react';

function Notification({ message, onClose }) {
  const [isVisible, setIsVisible] = useState(true);

  useEffect(() => {
    const timer = setTimeout(() => {
      setIsVisible(false);
      onClose();
    }, 5000);

    return () => clearTimeout(timer);
  }, [onClose]);

  if (!isVisible) return null;

  return (
    <div className="notification">
      <p>{message}</p>
      <button onClick={() => window.location.reload()}>Reload page</button>
      <button onClick={() => setIsVisible(false)}>Close</button>
    </div>
  );
}

export default Notification;