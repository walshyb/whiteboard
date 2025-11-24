import React from "react";

export default function RemoteCursor({ x, y, name, color = "#4A90E2" }) {
  const size = 14;

  const cursorStyle = {
    position: "absolute",
    left: x,
    top: y,
    transform: "translate(-50%, -50%)",
    pointerEvents: "none",

    display: "flex",
    flexDirection: "column",
    alignItems: "center",
    gap: "2px",

    transition: "left 50ms linear, top 50ms linear",
  };

  const dotStyle = {
    width: size,
    height: size,
    borderRadius: "50%",
    backgroundColor: color,
    boxShadow: `0 0 6px ${color}55`,
  };

  const labelStyle = {
    background: color,
    color: "white",
    padding: "2px 6px",
    borderRadius: "4px",
    fontSize: "10px",
    whiteSpace: "nowrap",
    userSelect: "none",
  };

  return (
    <div style={cursorStyle}>
      <div style={labelStyle}>{name}</div>
      <div style={dotStyle}></div>
    </div>
  );
}
