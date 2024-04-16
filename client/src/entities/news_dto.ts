export interface NewsHeadDTO {
  id: string
  title: string
  description: string
  source: string
  published_at: string
}

export interface NewsDTO extends NewsHeadDTO {
  link: string
  authors: string[]
  tags: string[]
  categories: string[]
  content: string
}
