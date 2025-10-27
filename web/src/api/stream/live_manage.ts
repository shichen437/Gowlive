/* eslint-disable @typescript-eslint/no-explicit-any */
import request from "@/lib/api";

export async function roomList(params: any) {
  return request({
    url: "/live/room/list",
    method: "get",
    params: params,
  });
}

export async function addRoom(data: any) {
  return request({
    url: "/live/room/manage",
    method: "post",
    data: data,
  });
}

export async function addBatchRoom(data: any) {
  return request({
    url: "/live/room/manage/batch",
    method: "post",
    data: data,
  });
}

export async function updateRoom(data: any) {
  return request({
    url: "/live/room/manage",
    method: "put",
    data: data,
  });
}

export async function deleteRoom(liveId: number) {
  return request({
    url: "/live/room/manage/" + liveId,
    method: "delete",
  });
}

export async function roomDetail(liveId: number) {
  return request({
    url: "/live/room/manage/" + liveId,
    method: "get",
  });
}

export async function startRoom(id: number) {
  return request({
    url: "/live/room/manage/start/" + id,
    method: "put",
  });
}

export async function stopRoom(id: number) {
  return request({
    url: "/live/room/manage/stop/" + id,
    method: "put",
  });
}
