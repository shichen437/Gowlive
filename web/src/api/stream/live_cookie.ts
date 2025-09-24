/* eslint-disable @typescript-eslint/no-explicit-any */
import request from "@/lib/api";

export async function allCookies(params: any) {
  return request({
    url: "/live/cookie/list",
    method: "get",
    params: params,
  });
}

export async function addCookie(data: any) {
  return request({
    url: "/live/cookie",
    method: "post",
    data: data,
  });
}

export async function updateCookie(data: any) {
  return request({
    url: "/live/cookie",
    method: "put",
    data: data,
  });
}

export async function deleteCookie(id: number) {
  return request({
    url: "/live/cookie/" + id,
    method: "delete",
  });
}