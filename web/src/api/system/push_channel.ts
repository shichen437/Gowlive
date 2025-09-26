import request from "@/lib/api";

export async function listPushChannel(params: any) {
  return request({
    url: "/system/push/channel/list",
    method: "get",
    params: params,
  });
}

export async function addPushChannel(data: any) {
  return request({
    url: "/system/push/channel",
    method: "post",
    data: data,
  });
}

export async function updatePushChannel(data: any) {
  return request({
    url: "/system/push/channel",
    method: "put",
    data: data,
  });
}

export async function getPushChannel(id: number) {
  return request({
    url: "/system/push/channel/" + id,
    method: "get",
  });
}

export async function deletePushChannel(id: number) {
  return request({
    url: "/system/push/channel/" + id,
    method: "delete",
  });
}
