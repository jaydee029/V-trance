import { useState } from "react";
import '../../App.css'
import { useDispatch } from 'react-redux';
import { setCredentials } from '../../store/authSlice';
import { loginUser } from '../../services/authService';
import { setTokens } from "../../store/authSlice";
import { useNavigate } from "react-router-dom";

const Login = () => {
  const dispatch = useDispatch();
  const navigate = useNavigate();

  const handleSubmit = async (e) => {
    e.preventDefault();
    const form = e.target;
    let username = "" 
    let email = ""
    if (form.emailpass.value.includes("@") ){
      email = form.emailpass.value
    } else{
      username = form.emailpass.value
    }
    const data = {
      username: username,
      email: email,
      password: form.password.value,
    };

    try {
      const result = await loginUser(data);
      const tokens = {
        accessToken : result.token,
        refreshToken : result.refresh_token
      }
      dispatch(setTokens(tokens)); 
      navigate("/dashboard")
    } catch (err) {
      console.error('Signup failed', err);
      alert('Signup failed');
    }
  };

  return (
    <form className="app-form" onSubmit={handleSubmit}  >
      <input name="emailpass" type="name" placeholder="Email/Username" required /> <br />
      <input name="password" type="password" placeholder="Password" required /> <br />
      <button type="submit">Login</button>
    </form>
  );
};

export default Login;