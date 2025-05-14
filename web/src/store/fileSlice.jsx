import { createSlice } from "@reduxjs/toolkit";

const initialState = {
  isDragging: false,
  fileUrl: null,
  fileName: null,
  fileSize: 0, // MB
  fileType: null, 
  fileDuration: 0, 
  fileHeight: 0, 
  fileWidth: 0,
};

const fileController = createSlice({
  name: "fileController",
  initialState,
  reducers: {
    setDragging: (state, action) => {
      state.isDragging = action.payload;
    },
    setSelectedFile: (state, action) => {
      state.fileUrl = action.payload.fileUrl;
      state.fileDuration = action.payload.fileDuration;
      state.fileHeight = action.payload.fileHeight;
      state.fileWidth = action.payload.fileWidth;
      state.fileName = action.payload.fileName;
      state.fileSize = action.payload.fileSize;
      state.fileType = action.payload.fileType;
    },
    clearSelectedFile: (state) => {
      // URL.revokeObjectURL(fileUrl);
      state.fileUrl = null;
      state.fileDuration = 0;
      state.fileHeight = 0;
      state.fileWidth = 0;
      state.fileName = null;
      state.fileSize = 0;
      state.fileType = null;
    },
  },
});

export const { setDragging, setSelectedFile, clearSelectedFile } = fileController.actions;
export default fileController.reducer;
