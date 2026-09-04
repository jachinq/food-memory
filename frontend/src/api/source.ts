import { api } from './client'

export interface SourcePreview {
  url: string
  source_platform: string
  name: string
  cover_image_url: string
  main_ingredients: string
  cook_time_minutes: number | null
  partial: boolean
}

export const previewSource = (url: string) =>
  api.post<SourcePreview>('/api/source-preview', { url })
