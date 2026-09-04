import { api } from './client'
import type { CookRecord, RecordPayload } from '../types'

export interface RecordSaveResult {
  record: CookRecord
  suggest_status?: string
}

export const fetchRecords = (dishId: number) => api.get<CookRecord[]>(`/api/dishes/${dishId}/records`)
export const createRecord = (dishId: number, payload: RecordPayload) =>
  api.post<RecordSaveResult>(`/api/dishes/${dishId}/records`, payload)
export const updateRecord = (id: number, payload: RecordPayload) =>
  api.put<RecordSaveResult>(`/api/records/${id}`, payload)
export const deleteRecord = (id: number) => api.del(`/api/records/${id}`)
