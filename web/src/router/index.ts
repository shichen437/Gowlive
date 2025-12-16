import { createRouter, createWebHistory } from "vue-router";
import { getToken } from "@/store/auth";

import Login from "@/views/login.vue";

// 布局/父路由
import Index from "@/views/index.vue";

// 子页面
import Dashboard from "@/views/dashboard/index.vue";
import Stream from "@/views/stream/index.vue";
import StreamMetrics from "@/views/stream/metrics.vue";
import StreamHistory from "@/views/stream/history.vue";
import StreamCookie from "@/views/stream/cookie.vue";
import StreamAuthor from "@/views/stream/author.vue";
import StreamPreview from "@/views/stream/preview.vue";
import SystemLogs from "@/views/system/logs.vue";
import SystemNotify from "@/views/system/notify.vue";
import SystemChannel from "@/views/system/channel.vue";
import MediaCheck from "@/views/media/check.vue";
import MediaFile from "@/views/media/file.vue";
import MediaPlay from "@/views/media/play.vue";
import Help from "@/views/help/index.vue";
import Users from "@/views/users/index.vue";
import About from "@/views/about/index.vue";
import NotFound from "@/views/NotFound.vue";

const routes = [
  {
    path: "/",
    name: "Index",
    component: Index,
    redirect: "/dashboard",
    children: [
      {
        path: "dashboard",
        name: "Dashboard",
        component: Dashboard,
        meta: { title: "project.router.overview" },
      },
      {
        path: "stream",
        meta: { title: "project.router.liveManagement" },
        children: [
          {
            path: "metrics",
            name: "StreamMetrics",
            component: StreamMetrics,
            meta: { title: "project.router.monitoringMetrics" },
          },
          {
            path: "history",
            name: "StreamHistory",
            component: StreamHistory,
            meta: { title: "project.router.liveHistory" },
          },
          {
            path: "index",
            name: "Stream",
            component: Stream,
            meta: { title: "project.router.roomList" },
          },
          {
            path: "author",
            name: "StreamAuthor",
            component: StreamAuthor,
            meta: { title: "project.router.anchorData" },
          },
          {
            path: "cookie",
            name: "StreamCookie",
            component: StreamCookie,
            meta: { title: "project.router.cookie" },
          },
        ],
      },
      {
        path: "media",
        meta: { title: "project.router.mediaCenter" },
        children: [
          {
            path: "file",
            name: "MediaFile",
            component: MediaFile,
            meta: { title: "project.router.fileManagement" },
          },
          {
            path: "check",
            name: "MediaCheck",
            component: MediaCheck,
            meta: { title: "project.router.fileDetection" },
          },
        ],
      },
      {
        path: "system",
        meta: { title: "project.router.systemManagement" },
        children: [
          {
            path: "channel",
            name: "SystemChannel",
            component: SystemChannel,
            meta: { title: "project.router.pushChannel" },
          },
          {
            path: "notify",
            name: "SystemNotify",
            component: SystemNotify,
            meta: { title: "project.router.notificationCenter" },
          },
          {
            path: "logs",
            name: "SystemLogs",
            component: SystemLogs,
            meta: { title: "project.router.logCenter" },
          },
        ],
      },
      {
        path: "help/index",
        name: "Help",
        component: Help,
        meta: { title: "project.router.commonQuestions" },
      },
      {
        path: "user/index",
        name: "Users",
        component: Users,
        meta: { title: "project.router.personalCenter" },
      },
      {
        path: "about/index",
        name: "About",
        component: About,
        meta: { title: "project.router.about" },
      },
    ],
  },
  {
    path: "/media/play",
    name: "MediaPlay",
    component: MediaPlay,
    meta: { title: "project.router.mediaPlayer" },
  },
  {
    path: "/stream/preview",
    name: "StreamPreview",
    component: StreamPreview,
    meta: { title: "project.router.streamPreview" },
  },
  {
    path: "/login",
    name: "Login",
    component: Login,
  },
  {
    path: "/:pathMatch(.*)*",
    name: "NotFound",
    component: NotFound,
    meta: { title: "404" },
  },
];

const router = createRouter({
  history: createWebHistory(),
  routes,
});

// 全局前置守卫
router.beforeEach((to, _from, next) => {
  const hasToken = getToken();

  if (hasToken) {
    if (to.path === "/login") {
      next({ path: "/" });
    } else {
      next();
    }
  } else {
    if (to.path === "/login") {
      next();
    } else {
      next(`/login?redirect=${to.path}`);
    }
  }
});

export default router;
