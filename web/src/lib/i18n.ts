import { createI18n } from "vue-i18n";

const i18n = createI18n({
  legacy: false,
  locale: localStorage.getItem("locale") || "zh-CN",
  fallbackLocale: "en",
  messages: loadLocaleMessages(),
});

function loadLocaleMessages() {
  const messages: Record<string, Record<string, any>> = {};
  const modules = import.meta.glob("../locales/**/**/*.json", { eager: true });

  for (const path in modules) {
    const match = path.match(/\..\/locales\/([^\/]+)\/([^\/]+)\.json$/);
    if (!match) continue;
    const [, lang, ns] = match;
    const mod = modules[path] as { default: any } | any;
    const data = (mod?.default ?? mod) as Record<string, any>;
    messages[lang] ??= {};
    messages[lang][ns] = data;
  }
  return messages;
}

export default i18n;
