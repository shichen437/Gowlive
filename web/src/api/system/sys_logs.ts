import request from "@/lib/api";

export async function listLogs(params: any) {
  if (params.status === '') {
    delete params.status;
  }
  if (params.type === '') {
    delete params.type;
  }
  return request({
    url: "/system/logs/list",
    method: "get",
    params: params,
  });
}

export async function deleteLogs(id: number) {
  return request({
    url: "/system/logs/" + id,
    method: "delete",
  });
}

export async function deleteAllLogs(data: any) {
  return request({
    url: "/system/logs/all",
    method: "delete",
    data: data,
  });
}
