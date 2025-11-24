import { useEffect, useRef } from "react";

export function useWebSocket(handleOnMessage) {
  const wsRef = useRef(null);
  const reconnectTimeout = useRef(null);
  const initialized = useRef(false);
  const reconnectDelay = 2000; // start delay in ms

  const connect = () => {
    const socket = new WebSocket("ws://localhost:8080/ws");
    wsRef.current = socket;

    socket.onmessage = function (e) {
      handleOnMessage(e);
    };

    socket.onerror = (err) => {
      console.error("WS error", err);
      socket.close(); // trigger reconnect via onclose
    };

    socket.onclose = () => {
      console.log("WS closed");
      reconnectTimeout.current = setTimeout(connect, reconnectDelay);
    };
  };

  useEffect(() => {
    if (initialized.current) return; // skip second mount
    initialized.current = true;
    connect();

    return () => {
      if (wsRef.current) wsRef.current.close();
      if (reconnectTimeout.current) clearTimeout(reconnectTimeout.current);
    };
  }, []);

  return wsRef;
}
