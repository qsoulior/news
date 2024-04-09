import { createRouter, createWebHistory } from "vue-router"

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: [
    {
      path: "/:page?",
      name: "list",
      props: (route) => {
        const params = route.params.page
        const param = parseInt(Array.isArray(params) ? params[0] : params)
        return { page: isNaN(param) ? undefined : param }
      },
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
