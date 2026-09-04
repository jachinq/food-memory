<template>
  <form @submit.prevent="save">
    <p class="section-kicker">{{ isEdit ? 'Revise' : 'New page' }}</p>
    <h1 class="section-title">{{ isEdit ? '改这道菜' : '记一道新菜' }}</h1>
    <section class="group">
      <h2>基础</h2>
      <div class="form-grid">
        <div>
          <label class="field"><span>菜名 *</span><input v-model="form.name" required /></label>
          <label class="field">
            <span>状态 *</span>
            <select v-model="form.status">
              <option v-for="(label, key) in STATUS_LABEL" :key="key" :value="key">{{ label }}</option>
            </select>
          </label>
        </div>
        <ImageUploader v-model="form.cover_image_url" biz-type="dish_cover" :biz-id="id" @uploaded="onCover" />
      </div>
    </section>
    <section class="group">
      <h2>从哪看来的</h2>
      <label class="field"><span>来源链接</span><input v-model="form.source_url" /></label>
      <label class="field"><span>来源平台</span><input v-model="form.source_platform" placeholder="小红书 / 抖音 / B站" /></label>
    </section>
    <section class="group">
      <h2>怎么想起它</h2>
      <label class="field"><span>主要食材</span><input v-model="form.main_ingredients" placeholder="鸡肉,豆腐" /></label>
      <label class="field"><span>口味</span><input v-model="form.taste" placeholder="甜辣" /></label>
      <label class="field"><span>场景</span><input v-model="form.scene" placeholder="下饭" /></label>
      <div class="field"><span>标签</span><TagSelector v-model="form.tags as { name: string; type: string }[]" /></div>
    </section>
    <section class="group">
      <h2>上手难度</h2>
      <div class="form-grid">
        <div class="field">
          <span>难度</span>
          <RatingInput v-model="form.difficulty as number | null" label="上手难度" />
        </div>
        <label class="field"><span>耗时（分钟）</span><input v-model.number="form.cook_time_minutes" type="number" min="0" /></label>
      </div>
    </section>
    <section class="group">
      <h2>随手记</h2>
      <label class="field"><span>备注</span><textarea v-model="form.note" rows="4" /></label>
    </section>
    <p v-if="error" class="muted">{{ error }}</p>
    <div class="form-actions">
      <button class="btn btn-primary" type="submit">记录</button>
    </div>
  </form>
</template>

<script setup lang="ts">
import { computed, onMounted, reactive, ref } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { createDish, getDish, updateDish } from '../api/dishes'
import type { Attachment, DishPayload } from '../types'
import { STATUS_LABEL } from '../types'
import ImageUploader from '../components/ImageUploader.vue'
import RatingInput from '../components/RatingInput.vue'
import TagSelector from '../components/TagSelector.vue'

const route = useRoute()
const router = useRouter()
const id = computed(() => Number(route.params.id) || 0)
const isEdit = computed(() => route.name === 'dish-edit')
const error = ref('')
const coverId = ref<number | undefined>()
const form = reactive<DishPayload>({
  name: '',
  status: 'want_to_cook',
  cover_image_url: '',
  source_url: '',
  source_platform: '',
  main_ingredients: '',
  taste: '',
  scene: '',
  note: '',
  tags: [],
})

onMounted(async () => {
  if (!isEdit.value) return
  const dish = await getDish(id.value)
  Object.assign(form, {
    name: dish.name,
    status: dish.status,
    cover_image_url: dish.cover_image_url,
    source_url: dish.source_url,
    source_platform: dish.source_platform,
    main_ingredients: dish.main_ingredients,
    taste: dish.taste,
    scene: dish.scene,
    note: dish.note,
    difficulty: dish.difficulty,
    cook_time_minutes: dish.cook_time_minutes,
    tags: (dish.tags || []).map((t) => ({ name: t.name, type: t.type })),
  })
})

function onCover(att: Attachment) {
  coverId.value = att.id
  form.cover_image_url = att.file_url
}

async function save() {
  error.value = ''
  try {
    const payload = { ...form, cover_attachment_id: coverId.value }
    const dish = isEdit.value ? await updateDish(id.value, payload) : await createDish(payload)
    router.push(`/dishes/${dish.id}`)
  } catch (e) {
    error.value = e instanceof Error ? e.message : '保存失败'
  }
}
</script>
