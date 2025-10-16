/* eslint-disable @typescript-eslint/no-explicit-any */
import request from "@/lib/api";

export async function anchorList(params: any) {
  return request({
    url: "/live/anchor/list",
    method: "get",
    params: params,
  });
}

export async function addAnchor(data: any) {
  return request({
    url: "/live/anchor",
    method: "post",
    data: data,
  });
}

export async function deleteAnchor(id: number) {
  return request({
    url: "/live/anchor/" + id,
    method: "delete",
  });
}

export async function getAnchorStatInfo(id: number) {
  return request({
    url: "/live/anchor/stat/" + id,
    method: "get",
  });
}
