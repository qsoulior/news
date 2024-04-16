import type { NewsHeadDTO } from "./news_dto"

export class NewsHead {
  id: string
  title: string
  description: string
  source: string
  publishedAt: Date

  constructor() {
    this.id = ""
    this.title = ""
    this.description = ""
    this.source = ""
    this.publishedAt = new Date()
  }

  static from(dto: NewsHeadDTO) {
    const obj = new NewsHead()
    obj.id = dto.id
    obj.title = dto.title
    obj.description = dto.description
    obj.source = dto.source
    obj.publishedAt = new Date(dto.published_at)

    return obj
  }
}

export interface News extends NewsHead {
  link: string
  authors: string[]
  tags: string[]
  categories: string[]
  content: string
}
