import { getInternalDict as getInternalDictFromApi } from "@/api/common/internal_dict";

export function getInternalDict(type: string) {
  return getInternalDictFromApi(type);
}
