import { createSlice } from '@reduxjs/toolkit';

const initialState = {
  user: null,
  accessToken: null,
  doesUserExist: false,
  refreshToken: null,
  isAuthenticated: false,
};

const authController = createSlice({
  name: 'authController',
  initialState: initialState,
  reducers: {
    setCredentials: (state, action) => {
      const { id, name, email, username} = action.payload;

      state.user = { id, name, email, username }; 
      state.doesUserExist = true;
    },
    setTokens: (state, action) => {
      const {accessToken, refreshToken} = action.payload;
      state.accessToken = accessToken;
      state.refreshToken = refreshToken;
      state.isAuthenticated = true;
    },
    logout: (state) => {
      state.user = null;
      state.accessToken = null;
      state.isAuthenticated = false;
    },
    updateAccessToken: (state, action) => {
      state.accessToken = action.payload;
    },
  },
});

export const { setCredentials, logout, updateAccessToken, setTokens } = authController.actions;
export default authController.reducer;
