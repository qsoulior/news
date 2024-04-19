import type { NewsDTO, NewsHeadDTO } from "./news_dto"

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

export class News extends NewsHead {
  link: string
  authors: string[]
  tags: string[]
  categories: string[]
  content: string

  constructor() {
    super()
    this.link = ""
    this.authors = []
    this.tags = []
    this.categories = []
    this.content = ""
  }

  static from(dto: NewsDTO) {
    const obj = new News()
    Object.assign(obj, { ...NewsHead.from(dto) })

    obj.link = dto.link
    obj.authors = dto.authors
    obj.tags = dto.tags
    obj.categories = dto.categories
    obj.content = dto.content

    return obj
  }
}
