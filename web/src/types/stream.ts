export interface LiveManage {
  id: number;
  roomUrl: string;
  interval: number;
  format: 'flv' | 'mp4' | 'mp3';
  monitorType: number;
  monitorStartAt: string;
  monitorStopAt: string;
  remark: string;
}

export interface RoomInfo {
  id: number;
  liveId: number;
  roomName: string;
  anchor: string;
  platform: string;
  status: number;
  isRecording: boolean;
}

export interface LiveHistory {
  id: number;
  liveId: number;
  anchor: string;
  startedAt: string;
  endedAt: string;
  duration: string;
}

export interface LiveCookie {
  id: number;
  platform: string;
  cookie: string;
  remark: string;
}

export interface AnchorInfo {
  id: number;
  platform: string;
  anchorName: string;
  signature: string;
  followingCount: number;
  followerCount: number;
  likeCount: number;
  videoCount: number;
  createdAt: string;
}
