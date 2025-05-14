import { configureStore } from '@reduxjs/toolkit';
import fileReducer from './fileSlice';
import transcodeReducer from './processSlice';
import authReducer from './authSlice';

export const store = configureStore({
  reducer: {
    fileController: fileReducer,
    processController: transcodeReducer,
    authController: authReducer,
  },
});
