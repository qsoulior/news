import { News, NewsHead } from "@/entities/news"
import type { NewsDTO, NewsHeadDTO } from "@/entities/news_dto"

const baseURL = new URL(import.meta.env.VITE_BACKEND_URL)
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

// YYYY-MM-DD
export function toDateString(date: Date) {
  const utc = new Date(date.getTime() - date.getTimezoneOffset() * 60000)
  return utc.toISOString().slice(0, 10)
}

interface GetNewsHeadQuery {
  text?: string
  tags?: string[]
  sources?: string[]
  dateFrom?: Date
  dateTo?: Date
}

interface GetNewsHeadOptions {
  limit: number
  skip?: number
  sort?: number
}

interface GetNewsHeadData {
  results?: NewsHeadDTO[]
  count?: number
  total_count?: number
}

export async function getNewsHead(query: GetNewsHeadQuery, opts: GetNewsHeadOptions) {
  const url = new URL("/news", baseURL)
  const urlParams = url.searchParams

  // options
  urlParams.set("limit", opts.limit.toFixed())
  if (opts.skip) urlParams.set("skip", opts.skip.toFixed())
  if (opts.sort) urlParams.set("sort", opts.sort.toFixed())

  // query
  if (query.text) urlParams.set("text", query.text)
  query.tags?.forEach((tag) => urlParams.append("tags[]", tag))
  query.sources?.forEach((source) => urlParams.append("sources[]", source))
  if (query.dateFrom) urlParams.set("date_from", toDateString(query.dateFrom))
  if (query.dateTo) urlParams.set("date_to", toDateString(query.dateTo))

  const resp = await fetch(url, {
    method: "GET"
  })

  if (resp.status != 200) {
    throw new Error(`incorrect status code: ${resp.status}`)
  }

  const respData = (await resp.json()) as GetNewsHeadData
  const results = (respData.results ?? []).map((dto) => NewsHead.from(dto))
  return {
    results: results,
    count: respData.count ?? results?.length ?? 0,
    totalCount: respData.total_count ?? 0
  }
}

export async function getNews(id: string): Promise<News> {
  const url = new URL("/news/" + id, baseURL)
  const resp = await fetch(url, {
    method: "GET"
  })

  if (resp.status != 200) {
    throw new Error(`incorrect status code: ${resp.status}`)
  }

  const respData = (await resp.json()) as NewsDTO
  return News.from(respData)
}
