<template>
  <TooltipProvider :delay-duration="0">
    <Toaster position="top-right" :duration="3000" />
    <router-view />
  </TooltipProvider>
</template>

<script setup lang="ts">
import { TooltipProvider } from '@/components/ui/tooltip';
import { Toaster } from '@/components/ui/sonner';
import 'vue-sonner/style.css';
import { useTheme } from '@/composables/useTheme';
useTheme();
</script>
