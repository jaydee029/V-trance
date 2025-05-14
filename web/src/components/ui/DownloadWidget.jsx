const CustomDownloadWidget = ({downloadLink = "download link here..."}) => {
  return (
    <div>
      <h2>Download your Video</h2>

      <div>
        {downloadLink}
        <br />
        <br />
      </div>
      
    </div>
  );
}

export default CustomDownloadWidget;