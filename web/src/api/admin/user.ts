/* eslint-disable @typescript-eslint/no-explicit-any */
import request from "@/lib/api";

export async function putProfile(nickname: string, sex: number) {
  const data = {
    nickname,
    sex,
  };
  return request({
    url: "/user/profile",
    method: "put",
    data: data,
  });
}

export async function putPassword(oldPwd: string, newPwd: string) {
  const data = {
    oldPwd,
    newPwd,
  };
  return request({
    url: "/user/password",
    method: "put",
    data: data,
  });
}