import { useDispatch } from "react-redux";
import { useNavigate } from "react-router-dom";
import { logout } from '../../store/authSlice';

const CustomHeader = () => {
  const navigate = useNavigate();
  const dispatch = useDispatch();

  const goToLoc = (loc) =>{
    navigate(loc);
  }

  const logoutUser = () => {
    dispatch(logout());                   
  }

    return (
        <div className="dashboard-header">
        { location.pathname === "/dashboard" ?
        <button onClick={() => goToLoc("/profile")}>Profile</button> : 
                <button onClick={() => goToLoc("/dashboard")}>Dashboard</button> }
        <h1 className="sub-heading">V-Trance</h1>
        <button onClick={() => logoutUser()} className="logout-button">
          <img src="./icons/logout.png" alt="" className="logout-icon" />
          Logout
        </button>
      </div>
    );
}

export default CustomHeader;