import { useColorMode, useStorage } from "@vueuse/core";
import { useCycleList } from "@vueuse/core";
import { watch } from "vue";

export const themes = [
  { value: "default", name: "project.topbar.colorMode.default" },
  { value: "glacier", name: "project.topbar.colorMode.glacier" },
  { value: "zhuqing-ink", name: "project.topbar.colorMode.zhuqing-ink" },
  {
    value: "midnightviolet",
    name: "project.topbar.colorMode.midnightviolet",
  },
  { value: "neonpro", name: "project.topbar.colorMode.neonpro" },
];

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

  const theme = useStorage("theme", "default");

  watch(
    theme,
    (newTheme, oldTheme) => {
      if (typeof document === "undefined") return;
      const html = document.documentElement;
      if (oldTheme && oldTheme !== "default") {
        html.classList.remove(`theme-${oldTheme}`);
      }
      if (newTheme && newTheme !== "default") {
        html.classList.add(`theme-${newTheme}`);
      }
    },
    { immediate: true },
  );

  const setTheme = (newTheme: string) => {
    theme.value = newTheme;
  };

  return {
    mode: state,
    cycle: next,
    theme,
    setTheme,
    themes,
  };
};
