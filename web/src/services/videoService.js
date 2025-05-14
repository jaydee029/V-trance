import apiFetch from './api';

const getUploadUrl = async (videoName, type, resolution, refreshToken, accessToken) => {
  const input = {
    "name":videoName,
    "type":type, 
    "resolution":resolution
  }
  const res = await apiFetch('/tranceapi/upload-Url', {
    method: "POST",
    body: JSON.stringify(input)
  //   headers: {
  //   'Authorization': accessToken,
  //   'X-Refresh-Token': refreshToken
  // }
  });
  // name, videoid, uploadurl
  console.log(res)
  return res; 
}

export const uploadVideo = async (videoName, processType, outputFmt, codecFmt, resolution, videoHeight, videoType, refreshToken, accessToken) => {
  // upload video to the upload url
  urlResponse = await getUploadUrl(videoName, videoType, videoHeight);
  // get videoid
  
  // notify upload

  // TODO update output with extension
  const data = {
    "videoId" : videoId,
    "type" : processType.toUpperCase(),
    "options" : {
      "output" : outputFmt,
      "codec" : codecFmt,
      "resolution" : resolution,
    }
  }
  const res = await api.post('/tranceapi/notifyUpload', data, {
  //   headers: {
  //   'Authorization': accessToken,
  //   'X-Refresh-Token': refreshToken
  // }
  });
  // jobid
  return res.data; 
}

export const getUserVideos = async (accessToken, refreshToken) => {
    const res = await api.get('/tranceapi/getVideos', userId, {
  //   headers: {
  //   'Authorization': accessToken,
  //   'X-Refresh-Token': refreshToken
  // }
  }); 
    // list of [name, url]
    return res.data; 
};

export const getJobStatus = async (jobId, accessToken, refreshToken) => {
  const poll = async () => {
    const res = await api.get(`/tranceapi/jobStatus/${jobId}`, {
  //   headers: {
  //   'Authorization': accessToken,
  //   'X-Refresh-Token': refreshToken
  // }
  });
    const status = res.data?.status;

    if (status === "REJECTED" || status === "COMPLETED") {
      return status;
    } else {
      await new Promise(resolve => setTimeout(resolve, 1000));
      return await poll();
    }
  };

  return await poll();
};


export const fetchStreamVideoUrl = async(videoId, accessToken, refreshToken) => {
  const res = await api.get(`/tranceapi/fetchVideo/${videoId}`, {
  //   headers: {
  //   'Authorization': accessToken,
  //   'X-Refresh-Token': refreshToken
  // }
  })

  return res.data;
  // name
  // stream url
}

export const fetchDownloadVideoUrl = async(videoId, refreshToken, accessToken) => {
  const res = await api.get(`/tranceapi/downloadVideo/${videoId}`, {
  //   headers: {
  //   'Authorization': accessToken,
  //   'X-Refresh-Token': refreshToken
  // }
  })

  return res.data;
  // download url
}
