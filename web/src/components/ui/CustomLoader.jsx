import '../../App.css'

const CustomLoader = ({
  size = "40px",
  borderWidth = "4px",
  borderStyle = "dotted", 
}) => (
  <div style={{ display: "flex", justifyContent: "center", marginTop: "2rem", padding:"20px" }}>
    <div
      className="spinner"
      style={{
        width: size,
        height: size,
        borderWidth: borderWidth,
        borderStyle: borderStyle,
        borderColor: "#f3f3f3",
        borderTopColor: "aquamarine",
        borderRadius: "50%",
        animation: "spin 1s linear infinite",
        boxSizing: "border-box",
      }}
    ></div>
  </div>
);


export default CustomLoader;
