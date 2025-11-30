export const SERVER_URL = import.meta.env.VITE_WS_URL || "localhost:8080";
export const USE_SSL = !!import.meta.env.VITE_WS_URL;
export const WS_URL = USE_SSL ? `wss://${SERVER_URL}` : `ws://${SERVER_URL}`;
export const API_URL = USE_SSL
  ? `https://${SERVER_URL}`
  : `http://${SERVER_URL}`;
