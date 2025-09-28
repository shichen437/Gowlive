import { useColorMode } from "@vueuse/core";
import { useCycleList } from "@vueuse/core";
import { watch } from "vue";

export const useTheme = () => {
  const mode = useColorMode({
    selector: "html",
    attribute: "class",
  });

  const { state, next } = useCycleList(["auto", "light", "dark"], {
    initialValue: mode.value,
  });

  watch(state, (v) => {
    mode.value = v as "auto" | "light" | "dark";
  });

  return {
    mode: state,
    cycle: next,
  };
};
