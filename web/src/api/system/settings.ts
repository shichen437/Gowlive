/* eslint-disable @typescript-eslint/no-explicit-any */
import request from "@/lib/api";

export async function latestVersion() {
  return request({
    url: "/system/latestVersion",
    method: "get",
  });
}