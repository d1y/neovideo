import { createApp } from 'vue'
import App from './App.vue'
import router from './router'
import piniaStore from './store'
import VueLazyload from 'vue-lazyload'

import '@/styles/index.less'
import '@/styles/reset.less'
import Antd from 'ant-design-vue'
import i18n from './locales/index'

const app = createApp(App)

app.use(VueLazyload)
app.use(Antd)
app.use(router)
app.use(piniaStore)
app.use(i18n)
app.mount('#app')
