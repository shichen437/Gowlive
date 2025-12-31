import request from "@/lib/api";

export async function listTasks(params: any) {
  return request({
    url: "/media/sync/list",
    method: "get",
    params: params,
  });
}

export async function resyncTask(id: number) {
  return request({
    url: "/media/sync/resync/" + id,
    method: "put",
  });
}

export async function deleteTask(id: number) {
  return request({
    url: "/media/sync/" + id,
    method: "delete",
  });
}

export async function deleteAll(params: any) {
  return request({
    url: "/media/sync/all",
    method: "delete",
    params: params,
  });
}
