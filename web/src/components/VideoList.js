import React from 'react';

function VideoList({ videos }) {
  return (
    <div>
      <h2>Shared Videos</h2>
      <ul>
        {videos.map((video) => (
          <li key={video.id}>
            <h3>{video.title}</h3>
            <iframe
              width="560"
              height="315"
              src={`https://www.youtube.com/embed/${video.youtube_video_id}`}
              title={video.title}
              frameBorder="0"
              allow="accelerometer; autoplay; clipboard-write; encrypted-media; gyroscope; picture-in-picture"
              allowFullScreen
            ></iframe>
          </li>
        ))}
      </ul>
    </div>
  );
}

export default VideoList;
