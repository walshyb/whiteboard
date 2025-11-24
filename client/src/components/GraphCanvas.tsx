import { useEffect, useRef, useState } from "react";
import "./GraphCanvas.css";
import { useWebSocket } from "../hooks";

export default function GraphCanvas() {
  const canvasRef = useRef(null);
  const clientId = useRef(null);
  const [activeClients, setActiveClients] = useState({});

  const wsRef = useWebSocket((e) => {
    const event = JSON.parse(e.data);
    if (event.type === "handshake") {
      clientId.current = event.clientId;
      return;
    }

    if (event.type === "mouseMove") {
      const remoteClientName = event.clientName;
      const { x, y } = event.data;
      setActiveClients((prev) => ({
        ...prev,
        [remoteClientName]: { x, y },
      }));
    }

    // Delete cursors of disconnected clients
    if (event.type === "client_disconnect") {
      const newActiveClients = { ...activeClients };
      delete newActiveClients[event.clientName];
      setActiveClients(newActiveClients);
    }
  });

  const [viewport, setViewport] = useState({
    x: 0, // pan offset x
    y: 0, // pan offset y
    scale: 1,
  });

  const dragging = useRef(false);
  const last = useRef({ x: 0, y: 0 });

  useEffect(() => {
    const canvas = canvasRef.current;
    const ctx = canvas.getContext("2d");

    resizeCanvas(canvas);

    let raf;
    function draw() {
      const width = canvas.clientWidth;
      const height = canvas.clientHeight;
      ctx.clearRect(0, 0, width, height);

      ctx.save();
      ctx.translate(viewport.x, viewport.y);
      ctx.scale(viewport.scale, viewport.scale);

      // Draw dotted grid
      const spacing = 40; // world units between dots
      const dotSize = 2 / viewport.scale; // scales visually with zoom

      // Determine visible world area
      const startX = -viewport.x / viewport.scale - 100;
      const endX = startX + width / viewport.scale + 200;
      const startY = -viewport.y / viewport.scale - 100;
      const endY = startY + height / viewport.scale + 200;

      ctx.fillStyle = "#888";

      for (
        let x = Math.floor(startX / spacing) * spacing;
        x < endX;
        x += spacing
      ) {
        for (
          let y = Math.floor(startY / spacing) * spacing;
          y < endY;
          y += spacing
        ) {
          ctx.beginPath();
          ctx.arc(x, y, dotSize, 0, Math.PI * 2);
          ctx.fill();
        }
      }

      ctx.restore();

      raf = requestAnimationFrame(draw);
    }

    draw();
    return () => cancelAnimationFrame(raf);
  }, [viewport]);

  // Panning
  function onMouseDown(e) {
    dragging.current = true;
    last.current = { x: e.clientX, y: e.clientY };
  }

  const lastSentMouseMovement = useRef(0);
  function onMouseMove(e) {
    const dx = e.clientX - last.current.x;
    const dy = e.clientY - last.current.y;

    if (dragging.current) {
      setViewport((v) => ({
        ...v,
        x: v.x + dx,
        y: v.y + dy,
      }));
      last.current = { x: e.clientX, y: e.clientY };
    }

    // If ws handshake didn't complete,
    // don't even bother sending mouse events
    if (!clientId.current) {
      return;
    }

    // Only send mouse movements to server ~every 30 frames
    const now = Date.now();
    if (now - lastSentMouseMovement.current < 33) return;
    lastSentMouseMovement.current = now;

    // World coordinates
    const wx = (e.clientX - viewport.x) / viewport.scale;
    const wy = (e.clientY - viewport.y) / viewport.scale;

    wsRef.current.send(
      JSON.stringify({
        type: "mouseMove",
        clientId: clientId.current,
        data: { x: wx, y: wy },
      }),
    );
  }

  function onMouseUp() {
    dragging.current = false;
  }

  function onMouseLeave() {
    if (!clientId.current) {
      return;
    }

    wsRef.current.send(
      JSON.stringify({
        type: "mouseMove",
        clientId: clientId.current,
        data: { x: -1, y: -1 },
      }),
    );
  }

  function resizeCanvas(canvas) {
    const ctx = canvas.getContext("2d");
    const dpr = window.devicePixelRatio || 1;
    const rect = canvas.getBoundingClientRect();
    canvas.width = rect.width * dpr;
    canvas.height = rect.height * dpr;
    ctx.scale(dpr, dpr);
  }

  function onWheel(e) {
    e.preventDefault();

    const ZOOM = 1.08;
    const oldScale = viewport.scale;
    const newScale =
      e.deltaY < 0
        ? Math.min(oldScale * ZOOM, 5)
        : Math.max(oldScale / ZOOM, 1);

    const rect = canvasRef.current.getBoundingClientRect();

    // Mouse position
    const mx = e.clientX - rect.left;
    const my = e.clientY - rect.top;

    // Mouse position to world coordinates
    const wx = (mx - viewport.x) / oldScale;
    const wy = (my - viewport.y) / oldScale;

    // Recompute pan so zoom is centered around the mouse
    setViewport({
      scale: newScale,
      x: mx - wx * newScale,
      y: my - wy * newScale,
    });
  }

  return (
    <>
      <canvas
        ref={canvasRef}
        style={{
          width: "100%",
          height: "100%",
          background: "#fff",
          cursor: "grab",
        }}
        onMouseDown={onMouseDown}
        onMouseMove={onMouseMove}
        onMouseLeave={onMouseLeave}
        onMouseUp={onMouseUp}
        onWheel={onWheel}
      />
      {Object.entries(activeClients).map(([name, pos]) => {
        // only show cursor if coordinates are in window
        if (pos.x === -1 && pos.y === -1) return;

        // Convert world coords to canvas coordinates
        const screenX = pos.x * viewport.scale + viewport.x;
        const screenY = pos.y * viewport.scale + viewport.y;

        return (
          <div
            key={name}
            className="remote-cursor"
            style={{
              left: screenX,
              top: screenY,
            }}
          >
            <div className="remote-cursor-label">{name}</div>
            <div className="remote-cursor-dot"></div>
          </div>
        );
      })}
    </>
  );
}
