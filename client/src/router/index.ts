import { createRouter, createWebHistory } from "vue-router"
import { getQueryInt } from "@/router/query"

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: [
    {
      path: "/:page?",
      name: "list",
      props: (route) => ({ page: getQueryInt(route.query, "page") ?? undefined }),
      component: () => import("@/views/ListView.vue")
    },
    {
      path: "/news/:id",
      name: "item",
      props: true,
      component: () => import("@/views/ItemView.vue")
    }
  ]
})

export default router
