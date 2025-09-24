import { ref, onMounted, computed } from "vue";
import { getInternalDict as getInternalDictFromApi } from "@/api/common/internal_dict";
import type { InternalDict } from "@/types/dict";

const cache = new Map<string, InternalDict[]>();

export function useDict(type: string) {
  const dict = ref<InternalDict[]>([]);

  const getDict = async () => {
    if (cache.has(type)) {
      dict.value = cache.get(type)!;
      return;
    }
    try {
      const res = await getInternalDictFromApi(type);
      if (res.data.data) {
        cache.set(type, res.data.data);
        dict.value = res.data.data;
      }
    } catch (error) {
      console.error(`Failed to fetch dictionary for type "${type}":`, error);
    }
  };

  onMounted(getDict);

  const getLabel = (value: string | number) => {
    const item = dict.value.find((d) => d.dictValue === value);
    return item ? item.dictLabel : value;
  };

  const options = computed(() =>
    dict.value.map((item) => ({
      label: item.dictLabel,
      value: item.dictValue,
    }))
  );

  return {
    dict,
    getLabel,
    options,
  };
}
