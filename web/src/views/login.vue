<script setup lang="ts">
import { ref } from 'vue'
import { Button } from '@/components/ui/button'
import {
  Card,
  CardContent,
  CardDescription,
  CardFooter,
  CardHeader,
  CardTitle,
} from '@/components/ui/card'
import { Input } from '@/components/ui/input'
import { Label } from '@/components/ui/label'
import { Separator } from '@/components/ui/separator'
import { useUserStore } from '@/store/user'

const username = ref('')
const password = ref('')

const userStore = useUserStore()

async function handleLogin() {
  try {
    await userStore.login(username.value, password.value)
  } catch (error) {
    console.error('登录失败:', error)
  }
}
</script>

<template>
  <div class="flex items-center justify-center h-screen bg-background">
    <div class="flex w-full max-w-4xl h-[600px] rounded-lg border shadow-lg">
      <div class="w-1/2 flex flex-col items-center justify-center p-8 text-center">
        <img src="/logo.png" alt="logo" class="h-14 w-28 mb-2" />
        <h2 class="text-3xl font-bold text-primary">Gowlive</h2>
        <p class="text-lg text-muted-foreground mt-2">一个直播录制平台</p>
      </div>
      <Separator orientation="vertical" />
      <div class="w-1/2 flex items-center justify-center p-8">
        <Card class="w-full max-w-md border-0 shadow-none">
          <CardHeader class="text-center">
            <CardTitle class="text-2xl">
              登录
            </CardTitle>
            <CardDescription>
              输入您的用户名和密码以登录您的帐户
            </CardDescription>
          </CardHeader>
          <CardContent class="grid gap-4">
            <div class="grid gap-2">
              <Label for="username">用户名</Label>
              <Input
                id="username"
                v-model="username"
                type="text"
                placeholder="请输入用户名"
                required
              />
            </div>
            <div class="grid gap-2">
              <Label for="password">密码</Label>
              <Input id="password" v-model="password" type="password" required placeholder="请输入密码" />
            </div>
          </CardContent>
          <CardFooter class="flex flex-col gap-4">
            <Button class="w-full" @click="handleLogin">
              登录
            </Button>
          </CardFooter>
        </Card>
      </div>
    </div>
  </div>
</template>
