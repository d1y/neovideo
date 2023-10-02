import { createApp } from 'vue'
import App from './App.vue'
import router from './router'
import piniaStore from './store'
import VueLazyload from 'vue-lazyload'
import noCover from '@/assets/no_cover.svg'

import '@/styles/index.less'
import '@/styles/reset.less'
import Antd from 'ant-design-vue'
import i18n from './locales/index'

const app = createApp(App)

app.use(VueLazyload, {
  preLoad: 1.0,
  error: noCover,
  attempt: 1
})
app.use(Antd)
app.use(router)
app.use(piniaStore)
app.use(i18n)
app.mount('#app')
