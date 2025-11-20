import Cookies from "js-cookie";
import { LAST_FILE_PATH } from "./consts";

export function getLastFilePath() {
  return Cookies.get(LAST_FILE_PATH);
}

export function setLastFilePath(path: string) {
  return Cookies.set(LAST_FILE_PATH, path, {
    path: "/",
    expires: new Date(Date.now() + 2592000000),
  });
}
