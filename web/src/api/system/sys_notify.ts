import request from "@/lib/api";

export async function listNotify(params: any) {
  return request({
    url: "/system/notify/list",
    method: "get",
    params: params,
  });
}

export async function markNotify(id: number) {
  return request({
    url: "/system/notify/" + id,
    method: "put",
  });
}

export async function markAllNotify() {
  return request({
    url: "/system/notify/all",
    method: "put",
  });
}

export async function deleteNotify(id: number) {
  return request({
    url: "/system/notify/" + id,
    method: "delete",
  });
}

export async function deleteAllNotify() {
  return request({
    url: "/system/notify/all",
    method: "delete",
  });
}
