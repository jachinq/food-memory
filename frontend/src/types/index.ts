export type DishStatus =
  | 'want_to_cook'
  | 'cooked'
  | 'success'
  | 'failed'
  | 'want_to_recook'
  | 'paused'

export type CookResult = 'success' | 'normal' | 'failed'
export type TagType = 'ingredient' | 'taste' | 'scene' | 'cuisine' | 'tool' | 'custom'
export type RecookStatus = 'active' | 'completed' | 'cancelled'

export interface Tag {
  id: number
  name: string
  type: TagType | string
  usage_count: number
}

export interface Attachment {
  id: number
  biz_type: string
  biz_id: number
  file_name: string
  file_url: string
  thumbnail_url: string
  mime_type: string
  file_size: number
  width: number
  height: number
}

export interface CookRecord {
  id: number
  dish_id: number
  cooked_at: string
  result: CookResult | string
  rating: number | null
  notes: string
  changes: string
  failure_reason: string
  next_improvement: string
  photos?: Attachment[]
}

export interface Dish {
  id: number
  name: string
  cover_image_url: string
  cover_thumbnail_url?: string
  source_url: string
  source_platform: string
  description: string
  status: DishStatus | string
  rating: number | null
  difficulty: number | null
  cook_time_minutes: number | null
  main_ingredients: string
  taste: string
  scene: string
  note: string
  last_cooked_at: string | null
  cook_count: number
  is_favorite: boolean
  tags?: Tag[]
  records?: CookRecord[]
}

export interface PageResult<T> {
  items: T[]
  page: number
  pageSize: number
  total: number
}

export interface HomeSummary {
  stats: {
    total_dishes: number
    cooked_count: number
    success_count: number
    month_cook_count: number
  }
  recent: Dish[]
  overdue_high_rating: Dish[]
  random_old: Dish | null
  want_to_cook: Dish[]
}

export interface RecookPlan {
  id: number
  dish_id: number
  planned_date: string | null
  status: RecookStatus | string
  reason: string
  completed_at: string | null
  dish?: Dish
}

export interface DishPayload {
  name: string
  status: string
  cover_image_url?: string
  cover_attachment_id?: number
  source_url?: string
  source_platform?: string
  description?: string
  difficulty?: number | null
  cook_time_minutes?: number | null
  main_ingredients?: string
  taste?: string
  scene?: string
  note?: string
  rating?: number | null
  tags?: { name: string; type: string }[]
}

export interface RecordPayload {
  cooked_at: string
  result: string
  rating?: number | null
  notes?: string
  changes?: string
  failure_reason?: string
  next_improvement?: string
  photo_urls?: string[]
  attachment_ids?: number[]
  update_dish_status?: string
}

export const STATUS_LABEL: Record<string, string> = {
  want_to_cook: '想做',
  cooked: '已做',
  success: '做成功',
  failed: '翻车',
  want_to_recook: '想再做',
  paused: '暂不做',
}

export const RESULT_LABEL: Record<string, string> = {
  success: '成功',
  normal: '一般',
  failed: '失败',
}

export function dishCover(dish: Pick<Dish, 'cover_thumbnail_url' | 'cover_image_url'>): string {
  return dish.cover_thumbnail_url || dish.cover_image_url || ''
}

export function formatDate(value?: string | null, withTime = false): string {
  if (!value) return ''
  const d = new Date(value)
  if (Number.isNaN(d.getTime())) return value.slice(0, 10)
  const y = d.getFullYear()
  const m = String(d.getMonth() + 1).padStart(2, '0')
  const day = String(d.getDate()).padStart(2, '0')
  if (!withTime) return `${y}-${m}-${day}`
  const hh = String(d.getHours()).padStart(2, '0')
  const mm = String(d.getMinutes()).padStart(2, '0')
  return `${y}-${m}-${day} ${hh}:${mm}`
}
