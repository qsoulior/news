import type { NewsHead } from "@/entities/news"

const sources = new Map([
  ["ria", "РИА Новости"],
  ["iz", "Известия"],
  ["lenta", "Лента.ру"],
  ["newsdata", "Newsdata"]
])

export function getSourceName(source: string): string {
  return sources.get(source) ?? ""
}

export function getSourceImg(source: string): string {
  return new URL(`/src/assets/icons/icon-${source}.png`, import.meta.url).href
}

export async function getNewsHead(limit: number, skip: number): Promise<NewsHead[]> {
  return Array.from({ length: limit }, (_, i) => ({
    id: (skip + i + 1).toString(),
    title: "Заголовок",
    description: "Описание",
    source: "ria",
    publishedAt: new Date()
  }))
}
