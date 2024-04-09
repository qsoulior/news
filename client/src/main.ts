// import "vfonts/FiraSans.css"
// import "vfonts/IBMPlexSans.css"
import "vfonts/Roboto.css"

import { createApp } from "vue"
import App from "@/App.vue"
import router from "@/router"

const app = createApp(App)

app.use(router)

app.mount("#app")
