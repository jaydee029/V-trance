import { useEffect, useState } from "react";
import { getUserVideos } from "../services/videoService";
import { useDispatch, useSelector } from "react-redux";
import { useNavigate } from "react-router-dom";
import CustomHeader from "./ui/CustomHeader";
import customToast from "./ui/CustomToast";
import CustomLoader from "./ui/CustomLoader";

const Profile = () => {
  const [videos, setVideos] = useState([]);
  const [loading, setLoading] = useState(false);
  const { user } = useSelector((state) => state.authController);

  const dispatch = useDispatch();
    const refreshToken = useSelector((state) => state.authController.refreshToken);
    const accessToken = useSelector((state) => state.authController.accessToken);
  const navigate = useNavigate();

  const fetchVideos = async () => {
    try {
      setLoading(true);
      res = await getUserVideos(accessToken, refreshToken);
      setVideos(res);
    } catch (error) {
      customToast("Failed to fetch user videos :(");
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    fetchVideos();
  }, []);

  return (
    <div className="app-container">
      <CustomHeader/>
      <h1 style={{ color: "aquamarine" }}> Welcome user.name </h1>
      <div style={{ justifyItems: "start", padding: "30px" , backgroundColor:"black"}}>
        <h3><i> Videos streamed by you </i></h3>
        <div style={{ display: "flex", justifySelf: "center", padding:"40px" , backgroundColor:"black", justifyContent:"center"}}>
          {loading ? (
            <CustomLoader />
          ) : videos.length === 0 ? (
            <p>No videos found</p>
          ) : (
            <ul>
              {videos.map((video) => (
                <li key={video.id}>TITLE: {video.name}
                <br/>
                URL: {video.url}
                </li>
              ))}
            </ul>
          )}
        </div>
      </div>
    </div>
  );
};

export default Profile;
