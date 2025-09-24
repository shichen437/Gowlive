import Cookies from "js-cookie";
import { TOKEN_KEY, USERNAME_KEY, LOGGED_KEY } from "./consts";

export function getToken() {
  return Cookies.get(TOKEN_KEY);
}

export function setToken(token: string) {
  return Cookies.set(TOKEN_KEY, token, { path: "/", expires: new Date(Date.now() + 2592000000) });
}

export function removeToken() {
  return Cookies.remove(TOKEN_KEY, { path: "/" });
}

export function setUsername(username: string) {
  return Cookies.set(USERNAME_KEY, encodeURIComponent(username), { path: "/", expires: new Date(Date.now() + 2592000000) });
}

export function getUsername() {
  return Cookies.get(USERNAME_KEY);
}

export function setIsLoggedIn(isLoggedIn: boolean) {
  return Cookies.set(LOGGED_KEY, isLoggedIn.toString(), { path: "/", expires: new Date(Date.now() + 2592000000) });
}

export function getIsLoggedIn() {
  return Cookies.get(LOGGED_KEY) === "true";
}

export function removeIsLoggedIn() {
  return Cookies.remove(LOGGED_KEY, { path: "/" });
}
