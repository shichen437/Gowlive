<template>
    <div class="space-y-4">
        <Card>
            <CardHeader>
                <CardTitle>{{ t('project.about.intro.title') }}</CardTitle>
            </CardHeader>
            <CardContent>
                <p class="text-base">
                    {{ t('project.about.intro.desc') }}
                </p>
                <br />
                <p class="text-slate-500 underline">{{ t('project.about.intro.security') }}</p>
            </CardContent>
        </Card>

        <Card>
            <CardHeader>
                <CardTitle>{{ t('project.about.info.title') }}</CardTitle>
            </CardHeader>
            <CardContent>
                <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
                    <p class="flex items-center">
                        <Tag class="w-4 h-4 mr-2" />
                        {{ t('project.about.info.currentVersion') }}: v{{ currentVersion }}
                        <span v-if="hasNewVersion"
                            class="ml-2 px-2 py-1 text-xs font-semibold text-white bg-red-500 rounded-full">
                            {{ t('project.about.info.new') }}: {{ latestVersion }}
                        </span>
                    </p>
                    <p class="flex items-center truncate">
                        <Github class="w-4 h-4 mr-2" />
                        {{ t('project.about.info.sourceCode') }}: <a href="https://github.com/shichen437/Gowlive"
                            target="_blank" class="text-slate-500 hover:underline ml-1">Github</a>
                        &nbsp;&nbsp;&nbsp;&nbsp;
                        <a href="https://gitee.com/shichen437/Gowlive" target="_blank"
                            class="text-slate-500 hover:underline ml-1">Gitee</a>
                    </p>
                    <p class="flex items-center">
                        <Mail class="w-4 h-4 mr-2" />
                        {{ t('project.about.info.email') }}: shichen437@126.com
                    </p>
                    <p class="flex items-center truncate">
                        <Link class="w-4 h-4 mr-2" />
                        {{ t('project.about.info.website.title') }}: <a href="https://shichen437.github.io"
                            target="_blank" class="text-slate-500 hover:underline ml-1">{{
                                t('project.about.info.website.main') }}</a> &nbsp;&nbsp;&nbsp;&nbsp;
                        <a href="https://gowlive-6b1dv010m.maozi.io" target="_blank"
                            class="text-slate-500 hover:underline ml-1">{{
                                t('project.about.info.website.backup') }}</a>
                    </p>
                </div>
            </CardContent>
        </Card>

        <div class="flex flex-col items-center justify-center pt-8">
            <h2 class="text-lg font-semibold mb-4">{{ t('project.about.support.title') }}</h2>
            <div class="flex space-x-8">
                <div class="flex flex-col items-center">
                    <img src="/support/alipay-qrcode.png" alt="Alipay QR Code" class="w-32 h-32">
                    <p class="text-center mt-2">{{ t('project.about.support.alipay') }}</p>
                </div>
                <div class="flex flex-col items-center">
                    <img src="/support/wechat-qrcode.png" alt="WeChat QR Code" class="w-32 h-32">
                    <p class="text-center mt-2">{{ t('project.about.support.wechat') }}</p>
                </div>
            </div>
        </div>
    </div>
</template>

<script setup lang="ts">
import { onMounted, ref } from 'vue'
import { useI18n } from 'vue-i18n';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card'
import { Github, Link, Tag, Mail } from 'lucide-vue-next'
import { latestVersion as apiLatestVersion } from '@/api/system/settings'

const { t } = useI18n()
const currentVersion = import.meta.env.VITE_APP_VERSION
const latestVersion = ref('')
const hasNewVersion = ref(false)

onMounted(async () => {
    try {
        const res: any = await apiLatestVersion()
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
