import React, { useRef } from "react";
import { useSelector, useDispatch } from "react-redux";
import {
  setDragging,
  setSelectedFile,
  clearSelectedFile,
} from "../../store/fileSlice";
import { formatFileSize, getVideoMetadata } from "../utils/helperMethods";
import customToast from "./CustomToast";

const MAX_SIZE = 111e6;

const DragDropComponent = () => {
  const dispatch = useDispatch();
  const { isDragging, fileUrl, fileName, fileType, fileSize, fileWidth, fileHeight, fileDuration } = useSelector(
    (state) => state.fileController
  );
  
  const fileInputRef = useRef(null);

  const handleFileUpload = async (e) => {
    const file = e.target.files[0];
    if (!file) return;
      if (file.size > MAX_SIZE) {
        customToast("File size exceeds threshold!");
        return;
      }
    

    try {
      const metadata = await getVideoMetadata(file);
      dispatch(setSelectedFile(metadata));
    } catch (err) {
      console.log(err)
      customToast("Error processing video file.");
    }
  };

  const handleDragEnter = (e) => {
    e.preventDefault();
    dispatch(setDragging(true));
  };

  const handleDragOver = (e) => {
    e.preventDefault();
  };

  const handleDragLeave = () => {
    dispatch(setDragging(false));
  };

  const handleDrop = async (e) => {
    e.preventDefault();
    dispatch(setDragging(false));

    const files = e.dataTransfer.files;
    if (!files || files.length === 0) return;
  
    const file = files[0];
    if (file.size > MAX_SIZE) {
      customToast("File size exceeds threshold!");
      return;
    }
  
    try {
      const metadata = await getVideoMetadata(file);
      dispatch(setSelectedFile(metadata));
    } catch (err) {
      console.log(err)
      customToast("Error processing video file.");
    }
  };

  const removeSelectedFile = () => {
    dispatch(clearSelectedFile());
  };

  return (
    <div>
    <div style={{textAlign: 'center', alignSelf: 'center'}}>
      <h2>Transcode Locally or Stream It Live!</h2>
      <p style={{ maxWidth: '48vw', textAlign: 'center', margin: '0 auto'}}>
        You can either download your media for offline use,
        or receive a streaming link for the transcoded content that can be accessed 
        anytime from your profile
      </p>
    </div>
    <div
      className={`file-drop-zone ${isDragging ? "dragging" : ""}`}
      onDragEnter={handleDragEnter}
      onDragOver={handleDragOver}
      onDragLeave={handleDragLeave}
      onDrop={handleDrop}
    >
      {!fileUrl ? (
        <>
          <img
            src="./icons/upload.png"
            alt=""
            style={{ height: "20vh", width: "auto", margin: "20px" }}
          />
          <br />
          Drag and Drop or{" "}
          <button onClick={() => fileInputRef.current.click()}>
            Upload a Video
          </button>
          <input
            type="file"
            accept="video/*"
            ref={fileInputRef}
            style={{ display: "none" }}
            onChange={handleFileUpload}
          />
        </>
      ) : (
        <div
          style={{
            position: "relative",
            display: "inline-block",
            marginTop: "10px",
          }}
        >
          <video
            src={fileUrl}
            controls
            style={{
              borderRadius: "10px",
              maxWidth: "40vw",
              width: "auto",
              height: "auto",
              objectFit: "contain",
            }}
          />
          <button
            onClick={removeSelectedFile}
            style={{
              position: "absolute",
              top: "-10px",
              right: "-10px",
              background: "red",
              color: "#fff",
              border: "none",
              borderRadius: "50%",
              width: "25px",
              height: "25px",
              cursor: "pointer",
              fontWeight: "bold",
            }}
          >
            ×
          </button>
          <div style={{ marginTop: "10px", fontSize: "14px", color: "#444" }}>
            {`Size: ${formatFileSize(fileSize)}`}
          </div>
        </div>
      )}
    </div>
    </div>
  );
};

export default DragDropComponent;
