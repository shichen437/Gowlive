import Cookies from "js-cookie";
import { LAST_FILE_PATH, STREAM_DISPLAY_MODE } from "./consts";

export function getLastFilePath() {
  return Cookies.get(LAST_FILE_PATH);
}

export function setLastFilePath(path: string) {
  return Cookies.set(LAST_FILE_PATH, path, {
    path: "/",
    expires: new Date(Date.now() + 2592000000),
  });
}

export function getStreamDisplayMode() {
  return Cookies.get(STREAM_DISPLAY_MODE);
}

export function setStreamDisplayMode(mode: string) {
  return Cookies.set(STREAM_DISPLAY_MODE, mode, {
    path: "/",
    expires: new Date(Date.now() + 2592000000),
  });
}
