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

const Dashboard = () => {
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

  useEffect(() => {
    if (!loading || !jobId) return;

    const checkStatus = async () => {
      const status = await getJobStatus(jobId, accessToken, refreshToken);
      if (status === "REJECTED") {
        customToast("Sorry, we couldn't process your media :(")
        dispatch(setLoading({ loading: false }))
        dispatch(setProcessed({ processed: false }))
      } else {
        dispatch(setLoading({ loading: true }));
        if (!streaming) {
          const downloadUrl = await fetchDownloadVideoUrl(videoId, accessToken, refreshToken);
          const urlRes = downloadUrl.downloadUrl;
          dispatch(setProcessed({ processed: true, downloadUrl: urlRes }))
          return;
        } else {
          const streamUrl = await fetchStreamVideoUrl(videoId, accessToken, refreshToken);
          const urlRes = streamUrl.streamUrl;
          dispatch(setProcessed({ processed: true, streamUrl: urlRes }))
        }
      }
    };

    checkStatus();
  }, [loading, jobId]);

  return (
    <div className="app-container">
      <CustomHeader />
      <div className="dashboard-body">
        {processed && streaming ?
          <VideoPlayer src="https://youtu.be/hk1ANYx1w64?si=9JqeaK0SOHpPSG5S" autoPlay />

          : <DragDropComponent />}

        {loading ? (
          <div style={{ alignItems: "center" }}>
            <CustomLoader size="10vw" />
            <h2>Hold on while we process your media...</h2>
          </div>
        ) : (
          processed ? streaming ? <CustomEmbedWidget /> : <CustomDownloadWidget downloadLink={downloadUrl} /> :
            <Options />
        )}
      </div>
    </div>
  );
};

export default Dashboard;
