import { defineStore } from "pinia";
import router from "@/router";
import { login as apiLogin, logout as apiLogout, getInfo } from "@/api/login";
import { type UserInfo } from "@/types/user";
import {
  setToken,
  removeToken,
  setUsername,
  setIsLoggedIn,
  removeIsLoggedIn,
} from "./auth";
import { USER_KEY, LOGGED_KEY } from "./consts";

interface UserState {
  userInfo: UserInfo | null;
}

export const useUserStore = defineStore("user", {
  state: (): UserState => ({
    userInfo: (() => {
      try {
        const data = localStorage.getItem(USER_KEY);
        return data ? (JSON.parse(data) as UserInfo) : null;
      } catch {
        return null;
      }
    })(),
  }),

  actions: {
    async getUserInfo(refresh = false): Promise<UserInfo | null> {
      if (!refresh && this.userInfo) {
        return this.userInfo;
      }
      try {
        const res = await getInfo();
        const info = res.data.userInfo;
        this.userInfo = info;
        localStorage.setItem(USER_KEY, JSON.stringify(info));
        return info;
      } catch (error) {
        this.userInfo = null;
        localStorage.removeItem(USER_KEY);
        return null;
      }
    },

    async login(username: string, password: string): Promise<void> {
      try {
        const res = await apiLogin(username, password);
        setToken(res.data.token);
        setUsername(username);
        setIsLoggedIn(true);
        localStorage.setItem(LOGGED_KEY, "true");
        await this.getUserInfo(true);
        router.push("/");
      } catch (error: any) {
        throw new Error(error.message || "登录失败");
      }
    },

    async logout(): Promise<void> {
      try {
        await apiLogout();
      } finally {
        this.clearUserInfo();
      }
    },

    clearUserInfo() {
      this.userInfo = null;
      removeIsLoggedIn();
      removeToken();
      localStorage.removeItem(USER_KEY);
      localStorage.removeItem(LOGGED_KEY);
    },
  },
});
