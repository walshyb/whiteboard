import { useEffect, useRef } from "react";
import { ClientMessage, ServerMessage } from "./proto/generated/events";
import { WS_URL } from "./helpers";

export function useWebSocket(
  handleOnMessage: Function,
  dependencies: Array<any>,
): [React.RefObject<WebSocket | null>, Function] {
  const wsRef: React.RefObject<WebSocket | null> = useRef(null);
  const reconnectTimeout: React.RefObject<number | null> = useRef(null);
  const initialized = useRef(false);
  const reconnectDelay = 2000; // start delay in ms

  function sendWsMessage(message: ClientMessage) {
    const encodedMessage = ClientMessage.encode(message).finish();
    wsRef.current?.send(encodedMessage);
  }

  const connect = () => {
    const socket = new WebSocket(WS_URL + "/ws");
    socket.binaryType = "arraybuffer";
    wsRef.current = socket;

    socket.onmessage = function (e: MessageEvent) {
      try {
        const uint8Array = new Uint8Array(e.data as ArrayBuffer);
        const serverMessage: ServerMessage = ServerMessage.decode(uint8Array);

        handleOnMessage(serverMessage);
      } catch (error: Error | unknown) {
        console.error("Failed to decode Protobuf message", error);
      }
    };

    socket.onerror = (_: Error | unknown) => {
      //console.error("WS error", err);
      socket.close(); // trigger reconnect via onclose
    };

    socket.onclose = () => {
      //console.log("WS closed");
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
  }, dependencies);

  return [wsRef, sendWsMessage];
}
