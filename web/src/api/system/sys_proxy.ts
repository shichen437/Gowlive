/* eslint-disable @typescript-eslint/no-explicit-any */
import request from "@/lib/api";

export async function allProxies(params: any) {
  return request({
    url: "/system/proxy/list",
    method: "get",
    params: params,
  });
}

export async function addProxy(data: any) {
  return request({
    url: "/system/proxy",
    method: "post",
    data: data,
  });
}

export async function updateProxy(data: any) {
  return request({
    url: "/system/proxy",
    method: "put",
    data: data,
  });
}

export async function deleteProxy(id: number) {
  return request({
    url: "/system/proxy/" + id,
    method: "delete",
  });
}
