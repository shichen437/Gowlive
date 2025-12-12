import request from "@/lib/api";

export async function overview() {
  return request({
    url: "/system/overview",
    method: "get",
  });
}

export async function lang() {
  return request({
    url: "/system/lang",
    method: "get",
    headers: {
      isToken: false,
    },
  });
}
