<template>
    <div class="p-4 space-y-4">
        <Card>
            <CardHeader>
                <CardTitle>项目介绍</CardTitle>
            </CardHeader>
            <CardContent>
                <p class="text-base">
                    Gowlive 是一个开源的直播录制工具，旨在为用户提供便捷、高效的直播内容录制解决方案，适用于需要保存、回看或分析直播内容的学习和交流场景。我们鼓励开发者和爱好者共同完善功能，提升用户体验。
                </p>
                <br />
                <p class="text-slate-500 underline">本项目仅供学习和技术交流使用，请勿将其用于任何商业用途或侵犯他人权益的行为。</p>
            </CardContent>
        </Card>

        <Card>
            <CardHeader>
                <CardTitle>项目信息</CardTitle>
            </CardHeader>
            <CardContent>
                <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
                    <p class="flex items-center">
                        <Tag class="w-4 h-4 mr-2" />
                        当前版本: v{{ currentVersion }}
                        <span v-if="hasNewVersion"
                            class="ml-2 px-2 py-1 text-xs font-semibold text-white bg-red-500 rounded-full">
                            发现新版本: {{ latestVersion }}
                        </span>
                    </p>
                    <p class="flex items-center truncate">
                        <Github class="w-4 h-4 mr-2" />
                        访问源码: <a href="https://github.com/shichen437/Gowlive" target="_blank"
                            class="text-slate-500 hover:underline ml-1">Github</a>
                    </p>
                    <p class="flex items-center">
                        <Mail class="w-4 h-4 mr-2" />
                        作者邮箱: shichen437@126.com
                    </p>
                    <p class="flex items-center truncate">
                        <Gitlab class="w-4 h-4 mr-2" />
                        源码镜像: <a href="https://gitee.com/shichen437/Gowlive" target="_blank"
                            class="text-slate-500 hover:underline ml-1">Gitee</a>
                    </p>
                </div>
            </CardContent>
        </Card>

        <div class="flex flex-col items-center justify-center pt-8">
            <h2 class="text-lg font-semibold mb-4">赞助 & 支持</h2>
            <div class="flex space-x-8">
                <div class="flex flex-col items-center">
                    <img src="/support/alipay-qrcode.png" alt="Alipay QR Code" class="w-32 h-32">
                    <p class="text-center mt-2">支付宝</p>
                </div>
                <div class="flex flex-col items-center">
                    <img src="/support/wechat-qrcode.png" alt="WeChat QR Code" class="w-32 h-32">
                    <p class="text-center mt-2">微信</p>
                </div>
            </div>
        </div>
    </div>
</template>

<script setup lang="ts">
import { onMounted, ref } from 'vue'
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card'
import { Github, Gitlab, Tag, Mail } from 'lucide-vue-next'
import { latestVersion as apiLatestVersion } from '@/api/system/settings'

const currentVersion = import.meta.env.VITE_APP_VERSION
const latestVersion = ref('')
const hasNewVersion = ref(false)

onMounted(async () => {
    try {
        const res:any = await apiLatestVersion()
        if (res.code !== 0) {
            return
        }
        const remoteVersion = res.data.latestVersion
        latestVersion.value = remoteVersion
        hasNewVersion.value = compareVersions(remoteVersion, currentVersion) > 0
    } catch (error) {
        console.error('Failed to fetch latest version:', error)
    }
})

function compareVersions(v1: string, v2: string): number {
    const parts1 = v1.replace('v', '').split('.').map(Number)
    const parts2 = v2.replace('v', '').split('.').map(Number)
    const len = Math.max(parts1.length, parts2.length)

    for (let i = 0; i < len; i++) {
        const p1 = parts1[i] || 0
        const p2 = parts2[i] || 0
        if (p1 > p2) return 1
        if (p1 < p2) return -1
    }
    return 0
}
</script>