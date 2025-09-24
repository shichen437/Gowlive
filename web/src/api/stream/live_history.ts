/* eslint-disable @typescript-eslint/no-explicit-any */
import request from "@/lib/api";

export async function historyList(params: any) {
  return request({
    url: "/live/history/list",
    method: "get",
    params: params,
  });
}

export async function deleteHistory(id: number) {
  return request({
    url: "/live/history/" + id,
    method: "delete",
  });
}