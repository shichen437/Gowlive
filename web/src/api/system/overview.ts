import request from "@/lib/api";

export async function overview() {
  return request({
    url: "/system/overview",
    method: "get",
  });
}