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
import SystemLogs from "@/views/system/logs.vue";
import SystemNotify from "@/views/system/notify.vue";
import SystemChannel from "@/views/system/channel.vue";
import MediaFile from "@/views/media/file.vue";
import Help from "@/views/help/index.vue";
import Users from "@/views/users/index.vue";
import About from "@/views/about/index.vue";

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
        meta: { title: "概览" },
      },
      {
        path: "stream",
        meta: { title: "直播管理" },
        children: [
          {
            path: "metrics",
            name: "StreamMetrics",
            component: StreamMetrics,
            meta: { title: "监控指标" },
          },
          {
            path: "history",
            name: "StreamHistory",
            component: StreamHistory,
            meta: { title: "直播历史" },
          },
          {
            path: "index",
            name: "Stream",
            component: Stream,
            meta: { title: "房间列表" },
          },
          {
            path: "author",
            name: "StreamAuthor",
            component: StreamAuthor,
            meta: { title: "主播数据" },
          },
          {
            path: "cookie",
            name: "StreamCookie",
            component: StreamCookie,
            meta: { title: "Cookie" },
          },
        ],
      },
      {
        path: "media",
        meta: { title: "媒体中心" },
        children: [
          {
            path: "file",
            name: "MediaFile",
            component: MediaFile,
            meta: { title: "文件管理" },
          },
        ],
      },
      {
        path: "system",
        meta: { title: "系统管理" },
        children: [
          {
            path: "channel",
            name: "SystemChannel",
            component: SystemChannel,
            meta: { title: "推送渠道" },
          },
          {
            path: "notify",
            name: "SystemNotify",
            component: SystemNotify,
            meta: { title: "通知中心" },
          },
          {
            path: "logs",
            name: "SystemLogs",
            component: SystemLogs,
            meta: { title: "日志中心" },
          },
        ],
      },
      {
        path: "help/index",
        name: "Help",
        component: Help,
        meta: { title: "常见问题" },
      },
      {
        path: "user/index",
        name: "Users",
        component: Users,
        meta: { title: "个人中心" },
      },
      {
        path: "about/index",
        name: "About",
        component: About,
        meta: { title: "关于" },
      },
    ],
  },
  {
    path: "/login",
    name: "Login",
    component: Login,
  },
  {
    path: "/:pathMatch(.*)*",
    name: "NotFound",
    component: () => import("@/views/NotFound.vue"),
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
