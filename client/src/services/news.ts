import type { NewsHead } from "@/entities/news"

export async function getNewsHead(limit: number, skip: number): Promise<NewsHead[]> {
  return Array.from({ length: limit }, (_, i) => ({
    id: (skip + i + 1).toString(),
    title: "Заголовок",
    description: "Описание",
    source: "РИА Новости",
    publishedAt: new Date()
  }))
}
