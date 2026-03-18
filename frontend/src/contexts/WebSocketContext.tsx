import React, { createContext, useContext, useEffect, useMemo, useRef, useState } from 'react';
import { toast } from 'react-toastify';
import { Notification } from '../services/notificationService';

type WSMessage =
  | { type: 'notification'; payload: { notification: Notification } }
  | { type: 'leave_update'; payload: { leave_id: string; status: string; leave?: any } }
  | { type: string; payload?: any };

interface WebSocketContextType {
  connected: boolean;
  lastMessage: WSMessage | null;
}

const WebSocketContext = createContext<WebSocketContextType | undefined>(undefined);

function buildWsUrl(): string {
  const apiUrl = process.env.REACT_APP_API_URL || 'http://localhost:8080/api';
  // Convert http(s)://host[:port]/api => ws(s)://host[:port]/ws
  const wsBase = apiUrl.replace(/^http/, 'ws').replace(/\/api\/?$/, '');
  return `${wsBase}/ws`;
}

export const WebSocketProvider: React.FC<{ children: React.ReactNode }> = ({ children }) => {
  const [connected, setConnected] = useState(false);
  const [lastMessage, setLastMessage] = useState<WSMessage | null>(null);
  const wsRef = useRef<WebSocket | null>(null);
  const reconnectTimer = useRef<number | null>(null);

  const token = localStorage.getItem('access_token');
  const wsUrl = useMemo(() => buildWsUrl(), []);

  useEffect(() => {
    if (!token) {
      // Ensure cleanup on logout
      if (wsRef.current) {
        wsRef.current.close();
        wsRef.current = null;
      }
      setConnected(false);
      return;
    }

    const connect = () => {
      try {
        const url = `${wsUrl}?token=${encodeURIComponent(token)}`;
        const ws = new WebSocket(url);
        wsRef.current = ws;

        ws.onopen = () => {
          setConnected(true);
        };

        ws.onclose = () => {
          setConnected(false);
          // Reconnect with a small backoff
          if (reconnectTimer.current) window.clearTimeout(reconnectTimer.current);
          reconnectTimer.current = window.setTimeout(connect, 1500);
        };

        ws.onerror = () => {
          // Avoid noisy spam; reconnect will happen on close.
        };

        ws.onmessage = (event) => {
          try {
            const msg = JSON.parse(event.data) as WSMessage;
            setLastMessage(msg);

            if (msg.type === 'notification' && msg.payload?.notification) {
              // Optional: small toast for live notifications
              toast.info(`${msg.payload.notification.title}: ${msg.payload.notification.message}`);
            }
          } catch {
            // Ignore malformed messages
          }
        };
      } catch {
        // Retry
        if (reconnectTimer.current) window.clearTimeout(reconnectTimer.current);
        reconnectTimer.current = window.setTimeout(connect, 1500);
      }
    };

    connect();

    return () => {
      if (reconnectTimer.current) window.clearTimeout(reconnectTimer.current);
      if (wsRef.current) wsRef.current.close();
      wsRef.current = null;
      setConnected(false);
    };
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [token, wsUrl]);

  const value = useMemo(() => ({ connected, lastMessage }), [connected, lastMessage]);

  return <WebSocketContext.Provider value={value}>{children}</WebSocketContext.Provider>;
};

export const useWebSocket = () => {
  const ctx = useContext(WebSocketContext);
  if (!ctx) throw new Error('useWebSocket must be used within WebSocketProvider');
  return ctx;
};

