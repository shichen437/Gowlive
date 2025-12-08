import axios from "axios";
import { toast } from "vue-sonner";
import { getToken } from "@/store/auth";
import { useUserStore } from "@/store/user";
import { pinia } from "@/main";
import errorCode from "@/lib/errorCode";

let isRelogin = { show: false };

// 1. 创建 axios 实例
const service = axios.create({
  baseURL: import.meta.env.VITE_APP_BASE_API,
  timeout: 30000,
  headers: { "Content-Type": "application/json;charset=utf-8" },
});

// 2. 请求拦截器
service.interceptors.request.use(
  (config) => {
    const isToken = (config.headers || {}).isToken !== false;
    if (getToken() && isToken) {
      config.headers.Authorization = `Bearer ${getToken()}`;
    }
    if (config.method === "get" && config.params) {
      let url = config.url + "?";
      for (const propName of Object.keys(config.params)) {
        const value = config.params[propName];
        const part = `${propName}=${encodeURIComponent(value)}&`;
        url += part;
      }
      url = url.slice(0, -1);
      config.params = {};
      config.url = url;
    }
    return config;
  },
  (error) => {
    console.error(error);
    return Promise.reject(error);
  },
);

// 3. 响应拦截器
service.interceptors.response.use(
  (res) => {
    const ct = (res.headers?.["content-type"] || "").toLowerCase();
    const isBinary =
      res.request?.responseType === "blob" ||
      res.request?.responseType === "arraybuffer" ||
      ct.includes("application/octet-stream") ||
      ct.includes(
        "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet",
      ) ||
      ct.includes("text/plain");

    if (isBinary) {
      return res;
    }

    const code = res.data.code || 0;
    const msg = errorCode[code] || res.data.msg || errorCode["default"];

    if (code === 401 || code === -401) {
      if (!isRelogin.show) {
        isRelogin.show = true;
        toast.error("登录状态已过期,请重新登录.");
        const userStore = useUserStore(pinia);
        userStore.logout().then(() => {
          isRelogin.show = false;
          window.location.href = "/login";
        });
      }
      return Promise.reject("无效的会话，或者会话已过期，请重新登录。");
    } else if (code === 500) {
      toast.error(msg);
      return Promise.reject(new Error(msg));
    } else if (code === 601) {
      toast.warning(msg);
      return Promise.reject(new Error(msg));
    } else if (code !== 0) {
      toast.error(msg);
      return Promise.reject("error");
    } else {
      return Promise.resolve(res.data);
    }
  },
  (error) => {
    console.log("err" + error);
    let { message } = error;
    if (message == "Network Error") {
      message = "后端接口连接异常";
    } else if (message.includes("timeout")) {
      message = "系统接口请求超时";
    } else if (message.includes("Request failed with status code")) {
      message = "系统接口" + message.substr(message.length - 3) + "异常";
    }
    toast.error(message);
    return Promise.reject(error);
  },
);

export default service;
