import React, { useState } from "react";

const CustomEmbedWidget = ({ link }) => {
  const embedCode = `<iframe width="560" height="315" src="${link}" frameborder="0" allowfullscreen></iframe>`;
  const [copiedEmbed, setCopiedEmbed] = useState(false);
  const [copiedLink, setCopiedLink] = useState(false);

  const copyToClipboard = (text, type) => {
    navigator.clipboard.writeText(text).then(() => {
      if (type === "embed") {
        setCopiedEmbed(true);
        setTimeout(() => setCopiedEmbed(false), 2000);
      } else {
        setCopiedLink(true);
        setTimeout(() => setCopiedLink(false), 2000);
      }
    });
  };

  return (
    <div 
    // style={{
    //   backgroundColor: "#1e1e1e",
    //   color: "#e0e0e0",
    //   borderRadius: "12px",
    //   padding: "20px",
    //   maxWidth: "100%",
    //   width: "600px",
    //   fontFamily: "Arial, sans-serif",
    //   margin: "0 auto"
    // }}
    >
      <h2 style={{ marginBottom: "20px", fontSize: "1.5em", color: "#fff" }}>Embed this Stream</h2>

      <label style={{ fontWeight: "bold", display: "block", marginBottom: "8px" }}>Embed Code:</label>
      <div style={{ position: "relative", marginBottom: "24px" }}>
        <textarea
          readOnly
          value={embedCode}
          rows="4"
          style={{
            width: "90%",
            padding: "12px",
            borderRadius: "8px",
            border: "1px solid #444",
            fontFamily: "monospace",
            fontSize: "14px",
            color: "#e0e0e0",
            backgroundColor: "#2a2a2a",
            resize: "none"
          }}
        />
        <button
          onClick={() => copyToClipboard(embedCode, "embed")}
          style={{
            position: "absolute",
            top: "10px",
            right: "10px",
            padding: "6px 12px",
            fontSize: "14px",
            borderRadius: "6px",
            border: "none",
            backgroundColor: "#007BFF",
            color: "#fff",
            cursor: "pointer"
          }}
        >
          {copiedEmbed ? "Copied!" : "Copy"}
        </button>
      </div>

      <label style={{ fontWeight: "bold", display: "block", marginBottom: "8px" }}>Stream Link:</label>
      <div style={{
        display: "flex",
        flexWrap: "wrap",
        gap: "10px",
        alignItems: "center"
      }}>
        <input
          type="text"
          readOnly
          value={link}
          style={{
            flex: "1 1 auto",
            minWidth: "200px",
            padding: "12px",
            borderRadius: "8px",
            border: "1px solid #444",
            backgroundColor: "#2a2a2a",
            color: "#e0e0e0",
            fontSize: "14px"
          }}
        />
        <button
          onClick={() => copyToClipboard(link, "link")}
          style={{
            padding: "10px 16px",
            fontSize: "14px",
            borderRadius: "6px",
            border: "none",
            backgroundColor: "#28A745",
            color: "#fff",
            cursor: "pointer",
            flexShrink: 0
          }}
        >
          {copiedLink ? "Copied!" : "Copy"}
        </button>
      </div>
    </div>
  );
};

export default CustomEmbedWidget;
