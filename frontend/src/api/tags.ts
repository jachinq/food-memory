import { api } from './client'
import type { Tag } from '../types'

export const fetchTags = (type = '', keyword = '') => {
  const q = new URLSearchParams()
  if (type) q.set('type', type)
  if (keyword) q.set('keyword', keyword)
  return api.get<Tag[]>(`/api/tags?${q.toString()}`)
}

export const createTag = (name: string, type = 'custom') => api.post<Tag>('/api/tags', { name, type })
