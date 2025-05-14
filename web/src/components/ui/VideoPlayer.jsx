// import React from "react";
// import ReactPlayer from "react-player";

// const VideoPlayer = () => {
//   return (
//     <div>
//       Video player
//       {/* <ReactPlayer
//         url="http://www.youtube.com/watch?v=ZieW_OSkuiQ"
//         controls
//         width="100%"
//         height="100%"
//       /> */}
//       <br/>
// <video width="560" height="315" part="video" crossOrigin="anonymous" playsInline={true} muted="" src="https://www.youtube.com/watch?v=ZieW_OSkuiQ" preload="metadata"><track kind="metadata" label="cuepoints" data-removeondestroy=""/></video>
//            {/* <iframe width="560" height="315" src="https://videos.pexels.com/video-files/7657449/7657449-hd_1920_1080_25fps.mp4" title="YouTube video player" frameBorder="0" allow="accelerometer; autoplay; clipboard-write; encrypted-media; gyroscope; picture-in-picture; web-share" referrerPolicy="strict-origin-when-cross-origin" allowFullScreen></iframe> */}
//     </div>
//   );
// };

// export default VideoPlayer;

import React, { useRef, useEffect } from 'react';
import Hls from 'hls.js';

const VideoPlayer = ({ src, autoPlay = false, controls = true }) => {
  const videoRef = useRef(null);

  useEffect(() => {
    const video = videoRef.current;

    if (!video) return;
    
    if (src.endsWith('.m3u8')) {
      if (Hls.isSupported()) {
        const hls = new Hls();
        hls.loadSource(src);
        hls.attachMedia(video);
        return () => hls.destroy();
      } else if (video.canPlayType('application/vnd.apple.mpegurl')) {
        video.src = src;
      } else {
        console.error('HLS not supported in this browser.');
      }
    } else {
      video.src = src;
    }
  }, [src]);

  return (
    <video
      ref={videoRef}
      controls={controls}
      autoPlay={autoPlay}
            height={315}
      width={560}
      style={{ borderRadius: '12px' ,padding:"30px"}}
    />
  );
};

export default VideoPlayer;
