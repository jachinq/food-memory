import { api } from './client'
import type { Dish, DishPayload, PageResult } from '../types'

export function fetchDishes(params: Record<string, string | number | boolean | undefined>) {
  const q = new URLSearchParams()
  Object.entries(params).forEach(([k, v]) => {
    if (v !== undefined && v !== '' && v !== null) q.set(k, String(v))
  })
  return api.get<PageResult<Dish>>(`/api/dishes?${q.toString()}`)
}

export const getDish = (id: number) => api.get<Dish>(`/api/dishes/${id}`)
export const createDish = (payload: DishPayload) => api.post<Dish>('/api/dishes', payload)
export const updateDish = (id: number, payload: DishPayload) => api.put<Dish>(`/api/dishes/${id}`, payload)
export const deleteDish = (id: number) => api.del(`/api/dishes/${id}`)
