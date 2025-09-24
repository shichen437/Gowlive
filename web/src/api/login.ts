/* eslint-disable @typescript-eslint/no-explicit-any */
import request from "@/lib/api";

export async function login(username: string, password: string) {
  const data = {
    username,
    password,
  };
  return request({
    url: "/login",
    method: "post",
    headers: {
      isToken: false,
    },
    data: data,
  });
}

export async function logout() {
  return request({
    url: "/logout",
    method: "post",
  });
}

export async function getInfo() {
  return request({
    url: "/user/getInfo",
    method: "get",
  });
}
