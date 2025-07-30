import React, { useState } from 'react';
import axios from '../axiosConfig';

function VideoShare({ user , onVideoShared }) {
  const [videoUrl, setVideoUrl] = useState('');
  const [videoTitle, setVideoTitle] = useState('');

  const handleSubmit = async (e) => {
    e.preventDefault();
    try {
      const token = localStorage.getItem('authToken');
      await axios.post('/api/authenticated/v1/share-video', {
        videoUrl,
        videoTitle,
      },
      {
        headers: {
          'Authorization': `${token}`
        }
      });
      setVideoUrl('');
      setVideoTitle('');
      alert('Video shared successfully!');
      onVideoShared(); // Call this function to trigger a refresh of the shared videos list
    } catch (error) {
      console.error('Error sharing video:', error);
      alert('Failed to share video. Please try again.');
    }
  };

  return (
    <form onSubmit={handleSubmit}>
      <input
        type="text"
        value={videoUrl}
        onChange={(e) => setVideoUrl(e.target.value)}
        placeholder="YouTube Video URL"
        required
      />
      <input
        type="text"
        value={videoTitle}
        onChange={(e) => setVideoTitle(e.target.value)}
        placeholder="Say something about the video"
        required
      />
      <button type="submit">Share Video</button>
    </form>
  );
}

export default VideoShare;