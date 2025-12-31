export interface FileInfo {
  filename: string;
  size: number;
  isFolder: boolean;
  lastModified: number;
}

export interface FileCheckTask {
  id: number;
  path: string;
  filename: string;
  duration: number;
  progress: number;
  fileStatus: number;
  createdAt: string;
  updatedAt: string;
}

export interface FileSyncTask {
  id: number;
  path: string;
  filename: string;
  syncPath: string;
  duration: number;
  status: number;
  remark: string;
  createdAt: string;
  updatedAt: string;
}

export function canPlay(file: FileInfo) {
  return isVideo(file) || isAudio(file);
}

export function isVideo(file: FileInfo) {
  return (
    !file.isFolder &&
    (file.filename.endsWith(".mp4") ||
      file.filename.endsWith(".flv") ||
      file.filename.endsWith(".mkv") ||
      file.filename.endsWith(".ts"))
  );
}

export function isAudio(file: FileInfo) {
  return !file.isFolder && file.filename.endsWith(".mp3");
}

export function isMp4(file: FileInfo) {
  return !file.isFolder && file.filename.endsWith(".mp4");
}
