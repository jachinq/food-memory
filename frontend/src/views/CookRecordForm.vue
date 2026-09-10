<template>
  <form @submit.prevent="save">
    <p class="section-kicker">Cook log</p>
    <h1 class="section-title">{{ isEdit ? '改这一次' : '记下这一次' }}</h1>
    <p v-if="dish" class="muted">{{ dish.name }}</p>
    <section class="group">
      <label class="field"><span>制作日期</span><input v-model="form.cooked_at" type="date" required /></label>
      <div class="field"><span>成品照片</span>
        <ImageUploader biz-type="cook_record" :biz-id="dish?.id" @uploaded="onPhoto" />
        <div class="row photo-list">
          <div v-for="att in photos" :key="att.id" class="photo-chip">
            <img :src="att.thumbnail_url || att.file_url" alt="" />
            <button type="button" class="photo-chip-remove" @click="removePhoto(att.id)">去掉</button>
          </div>
        </div>
      </div>
      <label class="field">
        <span>本次结果</span>
        <select v-model="form.result">
          <option value="success">成功</option>
          <option value="normal">一般</option>
          <option value="failed">失败</option>
        </select>
      </label>
      <div class="field"><span>本次评分</span><RatingInput v-model="form.rating as number | null" /></div>
      <label v-if="form.result === 'failed'" class="field"><span>翻车原因</span><textarea v-model="form.failure_reason" rows="3" /></label>
      <label class="field"><span>本次改动</span><textarea v-model="form.changes" rows="3" /></label>
      <label class="field"><span>下次注意</span><textarea v-model="form.next_improvement" rows="3" /></label>
      <label class="field"><span>备注</span><textarea v-model="form.notes" rows="3" /></label>
    </section>
    <p v-if="!isEdit && suggest" class="group">这次不错。要把菜品标记为「做成功 / 想再做」吗？
      <span class="row" style="margin-top:10px">
        <button type="button" class="btn btn-primary" @click="form.update_dish_status = 'success'">做成功</button>
        <button type="button" class="btn btn-ghost" @click="form.update_dish_status = 'want_to_recook'">想再做</button>
      </span>
    </p>
    <p v-if="error" class="muted">{{ error }}</p>
    <div class="form-actions">
      <button class="btn btn-primary" type="submit">{{ isEdit ? '保存这次记录' : '写入这次记录' }}</button>
    </div>
  </form>
</template>

<script setup lang="ts">
import { computed, onMounted, reactive, ref } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { getDish } from '../api/dishes'
import { createRecord, updateRecord } from '../api/records'
import type { Attachment, Dish, RecordPayload } from '../types'
import { formatDate } from '../types'
import ImageUploader from '../components/ImageUploader.vue'
import RatingInput from '../components/RatingInput.vue'

const route = useRoute()
const router = useRouter()
const dish = ref<Dish | null>(null)
const photos = ref<Attachment[]>([])
const error = ref('')
const today = new Date().toISOString().slice(0, 10)
const form = reactive<RecordPayload>({
  cooked_at: today,
  result: 'success',
  rating: 4,
  notes: '',
  changes: '',
  failure_reason: '',
  next_improvement: '',
})

const isEdit = computed(() => Boolean(route.params.recordId))
const suggest = computed(() => form.result === 'success' || (form.rating || 0) >= 4)

onMounted(async () => {
  dish.value = await getDish(Number(route.params.id))
  if (!isEdit.value) return
  const rec = dish.value.records?.find((r) => r.id === Number(route.params.recordId))
  if (!rec) {
    error.value = '这条制作记录不存在'
    return
  }
  form.cooked_at = formatDate(rec.cooked_at) || today
  form.result = rec.result || 'normal'
  form.rating = rec.rating
  form.notes = rec.notes || ''
  form.changes = rec.changes || ''
  form.failure_reason = rec.failure_reason || ''
  form.next_improvement = rec.next_improvement || ''
  photos.value = [...(rec.photos || [])]
})

function onPhoto(att: Attachment) {
  photos.value.push(att)
}

function removePhoto(id: number) {
  photos.value = photos.value.filter((p) => p.id !== id)
}

async function save() {
  if (!dish.value) return
  error.value = ''
  const payload: RecordPayload = {
    cooked_at: form.cooked_at,
    result: form.result,
    rating: form.rating,
    notes: form.notes,
    changes: form.changes,
    failure_reason: form.failure_reason,
    next_improvement: form.next_improvement,
    attachment_ids: photos.value.map((p) => p.id),
  }
  try {
    if (isEdit.value) {
      await updateRecord(Number(route.params.recordId), payload)
    } else {
      await createRecord(dish.value.id, {
        ...payload,
        update_dish_status: form.update_dish_status,
      })
    }
    if (!isEdit.value && route.query.from === 'recook') {
      router.push('/recook')
    } else {
      router.push(`/dishes/${dish.value.id}`)
    }
  } catch (e) {
    error.value = e instanceof Error ? e.message : '保存失败'
  }
}
</script>
