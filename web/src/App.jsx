import './App.css';
import { BrowserRouter as Router, Route, Routes } from "react-router-dom";
import Dashboard from "./components/Dashboard";
import PrivateRoute from "./router/route";
import Landing from './components/Landing';
import { ToastContainer } from 'react-toastify';
import Profile from './components/Profile';
import Video from './components/Video';

function App() {
  return (
    <Router>
      <Routes>
        <Route path="/" element={<Landing />} />
        {/* <Route exact path='/dashboard' element={<PrivateRoute />}> */}
          <Route exact path='/dashboard' element={<Dashboard />} />
        {/* </Route> */}
        {/* <Route exact path='/profile' element={<PrivateRoute />}> */}
          <Route exact path='/profile' element={<Profile />} />
        {/* </Route> */}
        <Route path="/video" element={<Video />} />
      </Routes>
      <ToastContainer />
    </Router>
  );
}

export default App;
