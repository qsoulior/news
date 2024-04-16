import type { LocationQuery } from "vue-router"

export function getQueryStr(query: LocationQuery, key: string): string | null {
  const queries = query[key]
  if (queries == undefined) return null

  return Array.isArray(queries) ? queries[0] : queries
}

export function getQueryStrs(query: LocationQuery, key: string): string[] {
  const queries = query[key]
  if (queries == undefined) return []

  if (Array.isArray(queries)) {
    return queries.filter((q) => q != null) as string[]
  }

  return queries != null ? [queries] : []
}

export function getQueryInt(query: LocationQuery, key: string): number | null {
  const queryStr = getQueryStr(query, key)
  if (queryStr == null) return null

  const queryInt = parseInt(queryStr)
  return isNaN(queryInt) ? null : queryInt
}

export function getQueryInts(query: LocationQuery, key: string): number[] {
  const queryStrs = getQueryStrs(query, key)
  return queryStrs.map((queryStr) => parseInt(queryStr)).filter((queryInt) => !isNaN(queryInt))
}
