import type { News, NewsHead } from "@/entities/news"

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

export async function getNews(id: string): Promise<News> {
  return {
    id: id,
    title: "В Монголии начали производить уникальное удобрение",
    description: "Монгольская компания начала делать удобрение из овечьей шерсти",
    source: "iz",
    publishedAt: new Date(),
    link: "https://lenta.ru/news/2024/04/09/v-mongolii-nachali-proizvodit-unikalnoe-udobrenie/",
    authors: ["John Doe", "Jane Doe"],
    tags: ["Монголия", "Удобрение", "Овцы", "Шерсть", "Monpellets", "Растения"],
    categories: ["Мир", "Среда обитания"],
    content:
      "Монгольская компания Monpellets начала делать полностью органическое удобрение из овечьей шерсти. Об этом сообщает издание Montsame.\nВ компании отметили, что на разработку продукта ушло десять лет. Удобрение без химикатов удается получить благодаря переработке овечьей шерсти в гранулы. Отмечается, что при производстве не используется вода. На сайте Monpellets говорится, что овечья шерсть очень полезна для растений — она содержит азот, фосфор и калий, которые нужны культурам для здорового развития.\nСейчас Monpellets экспортирует свой продукт в некоторые европейские страны, а также обеспечивает им местный рынок. В будущем компания планирует поставлять продукт в Турцию. Кроме того, несколько стран, включая Францию и Британию, попросили прислать образцы удобрения.\nРанее в США нашли новый способ создать безопасное для планеты топливо. Американский стартап Terragia Biofuel привлек шесть миллионов долларов, которые будут направлены на развитие компанией собственного метода переработки биомассы в этанол и другие продукты."
  }
}
