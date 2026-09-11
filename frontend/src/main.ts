import { createApp } from 'vue'
import { createPinia } from 'pinia'
import App from './App.vue'
import router from './router'
import { useInstall } from './stores/install'
import './styles/main.css'

const app = createApp(App)
const pinia = createPinia()
app.use(pinia).use(router)

const install = useInstall()
install.hydrate()
window.addEventListener('beforeinstallprompt', install.onBeforeInstall)
window.addEventListener('appinstalled', install.markInstalled)

app.mount('#app')
