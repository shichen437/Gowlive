<template>
  <div class="flex min-h-screen items-center justify-center bg-background">
    <div class="text-center">
      <h1 class="text-9xl font-bold text-primary">404</h1>
      <p class="mt-4 text-2xl font-semibold text-foreground">页面不存在</p>
      <p class="mt-2 text-muted-foreground">请检查你输入的地址是否正确。</p>
      <div class="mt-6">
        <Button @click="goHome">返回主页</Button>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { Button } from "@/components/ui/button";
import { useRouter } from "vue-router";

const router = useRouter();

function goHome() {
  router.push("/");
}
</script>
