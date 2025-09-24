import { createApp } from 'vue'
import './index.css'
import App from './App.vue'
import router from './router'
import { createPinia } from 'pinia'

export const pinia = createPinia()

createApp(App).use(router).use(pinia).mount('#app')
