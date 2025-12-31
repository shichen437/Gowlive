import request from "@/lib/api";

export async function openlistStatus() {
  return request({
    url: "/third/openlist/status",
    method: "get",
  });
}
