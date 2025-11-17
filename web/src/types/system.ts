export interface SysLogs {
  id: number;
  type: number;
  status: number;
  content: string;
  createdAt: string;
}

export interface SysNotify {
  id: number;
  level: string;
  title: string;
  content: string;
  status: number;
  createdAt: string;
}

export interface PushChannel {
  id: number;
  name: string;
  type: string;
  status: number;
  url: string;
  remark: string;
  email: PushChannelEmail;
  webhook: PushChannelWebhook;
  createdAt: string;
}

export interface PushChannelEmail {
  id: number;
  channelId: number;
  sender: string;
  receiver: string;
  server: string;
  port: number;
  authCode: string;
}

export interface PushChannelWebhook {
  id: number;
  channelId: number;
  webhookUrl: string;
  messageType: number;
  sign: string;
  at: string;
}
