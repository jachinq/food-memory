import { api } from './client'
import type { RecookPlan } from '../types'

export const fetchRecookPlans = (status = '') => {
  const q = status ? `?status=${encodeURIComponent(status)}` : ''
  return api.get<RecookPlan[]>(`/api/recook-plans${q}`)
}

export const createRecookPlan = (payload: { dish_id: number; planned_date?: string; reason?: string }) =>
  api.post<RecookPlan>('/api/recook-plans', payload)

export const updateRecookPlan = (id: number, payload: { planned_date?: string; status?: string; reason?: string }) =>
  api.put<RecookPlan>(`/api/recook-plans/${id}`, payload)

export const completeRecookPlan = (id: number) => api.post<RecookPlan>(`/api/recook-plans/${id}/complete`)
export const cancelRecookPlan = (id: number) => api.post(`/api/recook-plans/${id}/cancel`)
export const deleteRecookPlan = (id: number) => api.del(`/api/recook-plans/${id}`)
