import { createSlice } from "@reduxjs/toolkit";

const processController = createSlice({
  name: "processController",
  initialState: {
    streaming: false,
    streamUrl: '', 
    downloadUrl: '',
    processed: false,
    jobid : null,
    loading: false,
    error: null,
    outputFormat: 'MP4',
    codecFormat: 'H.264',
    videoId: null,
    resolution : 480
  },
  reducers: {
    setStreaming: (state, action) => {
      type = action.payload.mode
      if (type === "transcoding"){
        state.streaming = false;
      } else
      state.streaming = true;
    },
    setError: (state, action) => {
      state.error = action.payload.error;
    },
    changeTranscodeOptions: (state, action) => {
      if (action.payload.codecFormat !== undefined){
        state.codecFormat = action.payload.codecFormat;
      }
      if (action.payload.outputFormat !== undefined) {
        state.outputFormat = action.payload.outputFormat;
      }
      if (action.payload.resolution !== undefined){
        state.resolution = action.payload.resolution;
      }
    },
    setDefaultOptions : (state) => {
      state.codecFormat = 'H.264';
      state.outputFormat = 'MP4';
      state.resolution = 480;
    },
    setJobId : (state, action) => {
      state.jobid = action.payload;
    },
    setLoading: (state, action) => {
      state.loading = action.payload.loading;
    },
    setProcessed : (state, action) => {
      state.processed = action.payload.processed;
      if (action.payload.streamUrl!== undefined){
      state.streamUrl = action.payload.streamUrl;
      }
      if (action.payload.downloadUrl!==undefined){
        state.downloadUrl = action.payload.downloadUrl;
      }
    },
    setVideoId: (state, action) => {
      state.videoId = action.payload.videoId;
    }
  },
});

export const { setStreaming, changeTranscodeOptions, setDefaultOptions, setJobId, setLoading, setProcessed, setVideoId } = processController.actions;
export default processController.reducer;
