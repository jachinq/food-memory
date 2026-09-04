import { api } from './client'
import type { Attachment } from '../types'

export function uploadImage(file: File, bizType = 'dish_cover', bizId?: number) {
  const form = new FormData()
  form.append('file', file)
  form.append('biz_type', bizType)
  if (bizId) form.append('biz_id', String(bizId))
  return api.post<Attachment>('/api/upload', form)
}
