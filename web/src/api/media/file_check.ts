import request from "@/lib/api";

export async function listTasks(params: any) {
  return request({
    url: "/media/check/list",
    method: "get",
    params: params,
  });
}

export async function postTask(data: any) {
  return request({
    url: "/media/check",
    method: "post",
    data: data,
  });
}

export async function deleteTask(id: number) {
  return request({
    url: "/media/check/" + id,
    method: "delete",
  });
}
