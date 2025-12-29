import { createApp } from 'vue'
import './style.css'
import App from './App.vue'
import router from './router'
import pinia from './stores'
import { setupErrorHandler } from './plugins/errorHandler'

const app = createApp(App)

app.use(pinia)
app.use(router)

// Setup global error handling
setupErrorHandler(app)

app.mount('#app')
