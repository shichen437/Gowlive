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
    title: "project.router.liveManagement",
    icon: TvMinimal,
    children: [
      { title: "project.router.roomList", to: "/stream/index" },
      { title: "project.router.liveHistory", to: "/stream/history" },
      { title: "project.router.anchorData", to: "/stream/author" },
      { title: "project.router.monitoringMetrics", to: "/stream/metrics" },
      { title: "project.router.cookie", to: "/stream/cookie" },
    ],
  },
  {
    title: "project.router.mediaCenter",
    icon: Disc3,
    children: [
      { title: "project.router.fileManagement", to: "/media/file" },
      { title: "project.router.fileDetection", to: "/media/check" },
    ],
  },
  {
    title: "project.router.systemManagement",
    icon: Infinity,
    children: [
      { title: "project.router.logCenter", to: "/system/logs" },
      { title: "project.router.pushChannel", to: "/system/channel" },
      { title: "project.router.notificationCenter", to: "/system/notify" },
    ],
  },
  {
    title: "project.router.commonQuestions",
    icon: CircleQuestionMark,
    to: "/help/index",
  },
  {
    title: "project.router.personalCenter",
    icon: CircleUserRound,
    to: "/user/index",
  },
  {
    title: "project.router.about",
    icon: Info,
    to: "/about/index",
  },
]);
