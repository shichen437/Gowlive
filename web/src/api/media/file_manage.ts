import request from "@/lib/api";

export async function listFiles(params: any) {
  return request({
    url: "/media/file/list",
    method: "get",
    params: params,
  });
}

export async function getRoomFilePath(params: any) {
  return request({
    url: "/media/file/roomPath",
    method: "get",
    params: params,
  });
}

export async function deleteFile(data: any) {
  return request({
    url: "/media/file",
    method: "delete",
    data: data,
  });
}
