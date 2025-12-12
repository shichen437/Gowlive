<script setup lang="ts">
import type { PropType } from 'vue';
import { useRoute } from 'vue-router';
import {
  SidebarMenuItem,
  SidebarMenuButton,
} from '@/components/ui/sidebar';
import {
  Collapsible,
  CollapsibleContent,
  CollapsibleTrigger,
} from '@/components/ui/collapsible';
import { ChevronDown } from 'lucide-vue-next';
import { ref, watch } from 'vue';
import { useI18n } from 'vue-i18n';

const { t } = useI18n();

interface MenuItem {
  title: string;
  icon?: any;
  to?: string;
  children?: MenuItem[];
}

const props = defineProps({
  menu: {
    type: Array as PropType<MenuItem[]>,
    required: true,
  },
});

const route = useRoute();
const openCollapsible = ref<string | null>(null);

const isLinkActive = (to?: string) => to === route.path;

const isChildActive = (children?: MenuItem[]) => {
  if (!children) return false;
  return children.some(child => child.to === route.path);
};

const updateOpenCollapsible = () => {
  for (const item of props.menu) {
    if (item.children && isChildActive(item.children)) {
      openCollapsible.value = item.title;
      return;
    }
  }
};

watch(route, updateOpenCollapsible, { immediate: true });

</script>

<template>
  <template v-for="(item, index) in menu" :key="index">
    <SidebarMenuItem v-if="!item.children">
      <SidebarMenuButton :as-child="true" class="px-4 text-base"
        :variant="isLinkActive(item.to) ? 'outline' : 'default'">
        <router-link :to="item.to || '/'">
          <component :is="item.icon" class="h-5 w-5 shrink-0" />
          <span>{{ t(item.title) }}</span>
        </router-link>
      </SidebarMenuButton>
    </SidebarMenuItem>

    <Collapsible v-else as="li" class="w-full group" :open="openCollapsible === item.title"
      @update:open="openCollapsible = $event ? item.title : null">
      <CollapsibleTrigger as-child>
        <SidebarMenuButton class="w-full px-4 text-base"
          :variant="isChildActive(item.children) ? 'outline' : 'default'">
          <component :is="item.icon" class="h-5 w-5 shrink-0" />
          <span>{{ t(item.title) }}</span>
          <ChevronDown
            class="ml-auto h-4 w-4 shrink-0 transition-transform duration-200 group-data-[state=open]:rotate-180" />
        </SidebarMenuButton>
      </CollapsibleTrigger>
      <CollapsibleContent>
        <div class="py-1 pl-8">
          <template v-for="(child, childIndex) in item.children" :key="childIndex">
            <SidebarMenuItem>
              <SidebarMenuButton :as-child="true" class="text-sm"
                :variant="isLinkActive(child.to) ? 'outline' : 'default'">
                <router-link :to="child.to || '/'">
                  {{ t(child.title) }}
                </router-link>
              </SidebarMenuButton>
            </SidebarMenuItem>
          </template>
        </div>
      </CollapsibleContent>
    </Collapsible>
  </template>
</template>