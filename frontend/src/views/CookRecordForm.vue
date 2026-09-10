<template>
  <form @submit.prevent="save">
    <p class="section-kicker">Cook log</p>
    <h1 class="section-title">记下这一次</h1>
    <p v-if="dish" class="muted">{{ dish.name }}</p>
    <section class="group">
      <label class="field"><span>制作日期</span><input v-model="form.cooked_at" type="date" required /></label>
      <div class="field"><span>成品照片</span>
        <ImageUploader biz-type="cook_record" :biz-id="dish?.id" @uploaded="onPhoto" />
        <div class="row" style="margin-top:8px">
          <img v-for="url in photos" :key="url" :src="url" alt="" style="width:72px;height:72px;object-fit:cover" />
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
    <p v-if="suggest" class="group">这次不错。要把菜品标记为「做成功 / 想再做」吗？
      <span class="row" style="margin-top:10px">
        <button type="button" class="btn btn-primary" @click="form.update_dish_status = 'success'">做成功</button>
        <button type="button" class="btn btn-ghost" @click="form.update_dish_status = 'want_to_recook'">想再做</button>
      </span>
    </p>
    <p v-if="error" class="muted">{{ error }}</p>
    <div class="form-actions">
      <button class="btn btn-primary" type="submit">写入这次记录</button>
    </div>
  </form>
</template>

<script setup lang="ts">
import { computed, onMounted, reactive, ref } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { getDish } from '../api/dishes'
import { createRecord } from '../api/records'
import type { Attachment, Dish, RecordPayload } from '../types'
import ImageUploader from '../components/ImageUploader.vue'
import RatingInput from '../components/RatingInput.vue'

const route = useRoute()
const router = useRouter()
const dish = ref<Dish | null>(null)
const photos = ref<string[]>([])
const ids = ref<number[]>([])
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

const suggest = computed(() => form.result === 'success' || (form.rating || 0) >= 4)

onMounted(async () => {
  dish.value = await getDish(Number(route.params.id))
})

function onPhoto(att: Attachment) {
  photos.value.push(att.file_url)
  ids.value.push(att.id)
}

async function save() {
  if (!dish.value) return
  error.value = ''
  try {
    await createRecord(dish.value.id, {
      ...form,
      photo_urls: photos.value,
      attachment_ids: ids.value,
    })
    if (route.query.from === 'recook') {
      router.push('/recook')
    } else {
      router.push(`/dishes/${dish.value.id}`)
    }
  } catch (e) {
    error.value = e instanceof Error ? e.message : '保存失败'
  }
}
</script>
