/* eslint-disable @typescript-eslint/no-explicit-any */
import request from "@/lib/api";

export async function latestVersion() {
  return request({
    url: "/system/latestVersion",
    method: "get",
  });
}

export async function getSettings(params: any) {
  return request({
    url: "/system/settings",
    method: "get",
    params: params,
  });
}

export async function updateSettings(data: any) {
  return request({
    url: "/system/settings",
    method: "put",
    data: data,
  });
}
