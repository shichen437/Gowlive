/* eslint-disable @typescript-eslint/no-explicit-any */
import { getToken } from "@/store/auth";
import { useUserStore } from "@/store/user";

type SSEConfig = {
  channel: string;
  onMessage: (data: any) => void;
  onError?: (error: any) => void;
  onOpen?: () => void;
  maxRetries?: number;
  retryInterval?: number;
  heartbeatInterval?: number;
  heartbeatTimeout?: number;
  maxRetryInterval?: number;
};

class SSEClient {
  private eventSource: EventSource | null = null;
  private config: SSEConfig;
  private baseURL: string;
  private retryCount = 0;
  private retryTimer: ReturnType<typeof setTimeout> | null = null;
  private heartbeatTimer: ReturnType<typeof setInterval> | null = null;
  private lastMessageTime = Date.now();
  private isReconnecting = false;

  constructor(config: SSEConfig) {
    this.config = {
      maxRetries: 20,
      retryInterval: 3000,
      maxRetryInterval: 60000,
      heartbeatInterval: 12000,
      heartbeatTimeout: 20000,
      ...config,
    };
    this.baseURL = import.meta.env.VITE_APP_BASE_API || "";
  }

  private getEventSourceURL(): string {
    const token = getToken();
    const params = new URLSearchParams({ channel: this.config.channel });
    if (token) params.append("token", token);
    return `${this.baseURL}/sse?${params.toString()}`;
  }

  private startHeartbeatCheck(): void {
    this.stopHeartbeatCheck();
    this.heartbeatTimer = setInterval(() => {
      const now = Date.now();
      const diff = now - this.lastMessageTime;
      if (diff > (this.config.heartbeatTimeout || 30000)) {
        console.warn("SSE: heartbeat timeout, reconnecting...");
        this.scheduleReconnect("heartbeat-timeout");
      }
    }, this.config.heartbeatInterval);
  }

  private stopHeartbeatCheck(): void {
    if (this.heartbeatTimer) {
      clearInterval(this.heartbeatTimer);
      this.heartbeatTimer = null;
    }
  }

  private scheduleReconnect(reason?: string): void {
    if (this.isReconnecting) return;

    if (this.retryCount >= (this.config.maxRetries || 5)) {
      this.config.onError?.(new Error("SSE: reconnect failed"));
      return;
    }

    this.isReconnecting = true;
    this.stopHeartbeatCheck();

    // 计算指数退避延迟
    const base = this.config.retryInterval || 3000;
    const max = this.config.maxRetryInterval || 60000;
    const delay = Math.min(base * Math.pow(2, this.retryCount), max);

    this.retryTimer = setTimeout(() => {
      this.retryCount++;
      console.warn(
        `SSE: try ${this.retryCount} times reconnect after ${
          delay / 1000
        }s... (${reason || "unknown"})`
      );
      this.connect();
      this.isReconnecting = false;
    }, delay);
  }

  /** 建立连接 */
  public connect(): void {
    this.disconnect();
    const url = this.getEventSourceURL();
    this.eventSource = new EventSource(url);
    this.lastMessageTime = Date.now();

    this.eventSource.onopen = () => {
      console.info("SSE: connected");
      this.retryCount = 0;
      this.isReconnecting = false;
      this.config.onOpen?.();
      this.startHeartbeatCheck();
    };

    this.eventSource.onmessage = (event) => {
      this.lastMessageTime = Date.now();
      try {
        const data = JSON.parse(event.data);
        this.config.onMessage(data);
      } catch (error) {
        this.config.onError?.(error);
      }
    };

    this.eventSource.onerror = (error) => {
      const target = this.eventSource;
      this.config.onError?.(error);

      if (!target) return;

      // 登录状态异常
      const status = (target as any).status;
      if (status === 401 || status === -401) {
        useUserStore().logout();
        if (typeof window !== "undefined") window.location.href = "/login";
        return;
      }

      if (target.readyState === EventSource.CLOSED) {
        console.warn("SSE: connection closed, scheduling reconnect...");
        this.scheduleReconnect("eventsource-closed");
      }
    };
  }

  public disconnect(): void {
    if (this.eventSource) {
      this.eventSource.close();
      this.eventSource = null;
    }
    if (this.retryTimer) {
      clearTimeout(this.retryTimer);
      this.retryTimer = null;
    }
    this.stopHeartbeatCheck();
  }
}

export function createSSEConnection(config: SSEConfig): SSEClient {
  const client = new SSEClient(config);
  client.connect();
  return client;
}
