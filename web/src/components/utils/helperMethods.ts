export function formatFileSize(bytes) : number {
    if (bytes === 0) return 0;
    const k = 1024;
    // const sizes = ['Bytes', 'KB', 'MB', 'GB', 'TB'];
    const size = 'MB'
    const i = Math.floor(Math.log(bytes) / Math.log(k));
    return parseFloat((bytes / Math.pow(k, i)).toFixed(2));
  }

export const getVideoMetadata = (file: File): Promise<{
    fileUrl: string;
    fileName: string;
    fileSize: number; // MB
    fileType: string; 
    fileDuration: number; // seconds
    fileHeight: number; // px
    fileWidth: number; // px
  }> => {
    return new Promise((resolve, reject) => {
      const video = document.createElement("video");
      const fileUrl = URL.createObjectURL(file);
      video.preload = "metadata";
  
      video.onloadedmetadata = () => {
        resolve({
          fileUrl,
          fileName: file.name,
          fileSize: formatFileSize(file.size),
          fileType: file.type,
          fileDuration: video.duration,
          fileHeight: video.videoHeight,
          fileWidth: video.videoWidth,
        });
      };
  
      video.onerror = () => {
        reject("Failed to load video metadata.");
      };
  
      video.src = fileUrl;
    });
  };

export const getValidResolutions = (res : number) : number[] => {
  let resolutions = [360, 480, 720, 1080];
  const valid = resolutions.filter(r=> r>=res);
  return valid;
}

export const getValidCodecFormats = (videoFormat : string) : string[] => {
  const codecMap: {[key: string] : string[]} = {
    "MP4" : ["H.264", "H.265", "AV1"],
    "WebM" : ["VP9", "AV1"],
    "MKV" : ["H.264", "H.265", "VP9", "AV1"],
    "DASH" : ["H.264", "H.265", "VP9", "AV1"],
  }
  if (!(videoFormat in codecMap)) return [];
  return codecMap[videoFormat]
}