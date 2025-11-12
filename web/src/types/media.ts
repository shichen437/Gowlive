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
