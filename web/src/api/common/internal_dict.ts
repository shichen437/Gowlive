/* eslint-disable @typescript-eslint/no-explicit-any */
import request from "@/lib/api";

export async function getInternalDict(type: string) {
  const params = {
    dictType: type,
  };
  return request({
    url: "/dict/internal/type",
    method: "get",
    params: params,
  });
}
