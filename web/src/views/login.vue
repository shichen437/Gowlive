<script setup lang="ts">
import { ref, onMounted } from 'vue'
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
import { useI18n } from 'vue-i18n';
import { lang } from '@/api/system/overview';

const { t, locale } = useI18n();

const username = ref('')
const password = ref('')

const userStore = useUserStore()

onMounted(async () => {
    const localLocale = localStorage.getItem('locale');
    if (!localLocale) {
        try {
            const res = await lang();
            if (res && res.data && res.data.lang) {
                localStorage.setItem('locale', res.data.lang);
                locale.value = res.data.lang;
            }
        } catch (error) {
            console.error('Failed to fetch language:', error);
        }
    }
});

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
                <img src="/logo.png" alt="logo" class="h-20 w-20 mb-2" />
                <h2 class="text-3xl font-bold text-primary">Gowlive</h2>
                <p class="text-lg text-muted-foreground mt-2">{{ t('project.login.slogan') }}</p>
            </div>
            <Separator orientation="vertical" />
            <div class="w-1/2 flex items-center justify-center p-8">
                <Card class="w-full max-w-md border-0 shadow-none">
                    <CardHeader class="text-center">
                        <CardTitle class="text-2xl">
                            {{ t('project.login.title') }}
                        </CardTitle>
                        <CardDescription>
                            {{ t('project.login.desc') }}
                        </CardDescription>
                    </CardHeader>
                    <CardContent class="grid gap-4">
                        <div class="grid gap-2">
                            <Label for="username">{{ t('project.login.fields.username') }}</Label>
                            <Input id="username" v-model="username" type="text"
                                :placeholder="t('project.login.placeholders.username')" required />
                        </div>
                        <div class="grid gap-2">
                            <Label for="password">{{ t('project.login.fields.password') }}</Label>
                            <Input id="password" v-model="password" type="password" required
                                :placeholder="t('project.login.placeholders.password')" />
                        </div>
                    </CardContent>
                    <CardFooter class="flex flex-col gap-4">
                        <Button class="w-full" @click="handleLogin">
                            {{ t('project.login.title') }}
                        </Button>
                    </CardFooter>
                </Card>
            </div>
        </div>
    </div>
</template>
