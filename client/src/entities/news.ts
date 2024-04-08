export interface NewsHead {
  id: string
  title: string
  description: string
  source: string
  publishedAt: Date
}

export interface News extends NewsHead {
  link: string
  authors: string[]
  tags: string[]
  categories: string[]
  content: string
}
