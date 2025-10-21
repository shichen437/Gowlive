import { ref } from "vue";
import {
  TvMinimal,
  Disc3,
  CircleUserRound,
  Info,
  Infinity,
  CircleQuestionMark,
} from "lucide-vue-next";

export const configurableMenu = ref([
  {
    title: "直播管理",
    icon: TvMinimal,
    children: [
      { title: "房间列表", to: "/stream/index" },
      { title: "直播历史", to: "/stream/history" },
      { title: "主播数据", to: "/stream/author" },
      { title: "监控指标", to: "/stream/metrics" },
      { title: "Cookie", to: "/stream/cookie" },
    ],
  },
  {
    title: "媒体中心",
    icon: Disc3,
    children: [
      { title: "文件管理", to: "/media/file" },
      //   { title: "媒体解析", to: "/media/parse" },
      //   { title: "粉丝趋势", to: "/media/followers" },
    ],
  },
  {
    title: "系统管理",
    icon: Infinity,
    children: [
      { title: "日志中心", to: "/system/logs" },
      { title: "推送渠道", to: "/system/channel" },
      { title: "通知中心", to: "/system/notify" },
    ],
  },
  {
    title: "常见问题",
    icon: CircleQuestionMark,
    to: "/help/index",
  },
  {
    title: "个人中心",
    icon: CircleUserRound,
    to: "/user/index",
  },
  {
    title: "关于",
    icon: Info,
    to: "/about/index",
  },
]);
