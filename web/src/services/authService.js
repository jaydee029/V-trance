import apiFetch from './api';

export const loginUser = async (credentials) => {
  const res = await apiFetch('/user/login',{
    method:"POST",
    body: JSON.stringify(credentials)
  });
  
  return res; 
};

export const signupUser = async (formData) => {
    const res = await apiFetch('/user/signup',{
    method:"POST",
    body: JSON.stringify(formData)
  }); 
    return res;
};
  