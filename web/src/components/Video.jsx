import React, { useEffect } from "react";
import DragDropComponent from "./ui/DragDrop";
import { useSelector, useDispatch } from "react-redux";
import "../App.css";
import CustomHeader from "./ui/CustomHeader";
import Options from "./ui/Options";
import { useNavigate } from "react-router-dom";
import CustomLoader from "./ui/CustomLoader";
import CustomEmbedWidget from "./ui/EmbedWidget";
import CustomDownloadWidget from "./ui/DownloadWidget";
import VideoPlayer from "./ui/VideoPlayer";
import { getJobStatus, fetchStreamVideoUrl, fetchDownloadVideoUrl } from "../services/videoService";
import customToast from "./ui/CustomToast";
import { setProcessed, setLoading } from "../store/processSlice";

const Video = () => {
  const navigate = useNavigate();
  const dispatch = useDispatch();
  const loading = useSelector((state) => state.processController.loading);
  const refreshToken = useSelector((state) => state.authController.refreshToken);
  const accessToken = useSelector((state) => state.authController.accessToken);
  const processed = useSelector((state) => state.processController.processed);
  const jobId = useSelector((state) => state.processController.jobid);
  const videoId = useSelector((state) => state.processController.videoId);
  const streamUrl = useSelector((state) => state.processController.streamUrl);
  const downloadUrl = useSelector((state) => state.processController.downloadUrl);
  const streaming = useSelector((state) => state.processController.streaming);

  const goToProfile = () => {
    // customToast("Go to Profile!")
    navigate("/profile");
  };


  return (
    <div className="app-container">
      <CustomHeader />
      <h2>Stream your Video Here!</h2>
      <p>or use the streaming link to run the stream on any platform!</p>
      <div className="dashboard-body">
        
        <VideoPlayer src="https://videos.pexels.com/video-files/4614907/4614907-uhd_2560_1440_30fps.mp4" />
        <CustomEmbedWidget link="https://videos.pexels.com/video-files/4614907/4614907-uhd_2560_1440_30fps.mp4"/> 
      </div>
    </div>
  );
};

export default Video;
