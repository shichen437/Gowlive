<template>
    <Toaster position="top-right" :duration="3000" />
    <router-view />
</template>

<script setup lang="ts">
import { Toaster } from '@/components/ui/sonner'
import 'vue-sonner/style.css'
import { useTheme } from "@/composables/useTheme";
useTheme();
</script>
