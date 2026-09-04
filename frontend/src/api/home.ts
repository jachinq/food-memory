import { api } from './client'
import type { HomeSummary } from '../types'

export const fetchHomeSummary = () => api.get<HomeSummary>('/api/home/summary')
