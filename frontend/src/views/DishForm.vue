<template>
  <form @submit.prevent="save">
    <p class="section-kicker">{{ isEdit ? 'Revise' : 'New page' }}</p>
    <h1 class="section-title">{{ isEdit ? '改这道菜' : '记一道新菜' }}</h1>
    <section class="group">
      <h2>基础</h2>
      <label class="field"><span>菜名 *</span><input v-model="form.name" required /></label>
      <label class="field"><span>做法</span><textarea v-model="form.description" rows="6" placeholder="怎么做这道菜" /></label>
      <div class="form-grid">
        <div>
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
      <label class="field">
        <span>来源链接</span>
        <input v-model="form.source_url" placeholder="https://" @blur="onSourceBlur" />
      </label>
      <p v-if="previewing" class="muted">正在识别链接…</p>
      <p v-else-if="previewHint" class="muted">{{ previewHint }}</p>
      <label class="field"><span>来源平台</span><input v-model="form.source_platform" placeholder="小红书 / 抖音 / B站" /></label>
    </section>
    <section class="group">
      <h2>怎么想起它</h2>
      <div class="field">
        <span>主要食材</span>
        <TagSelector
          v-model="ingredients"
          type="ingredient"
          placeholder="输入食材回车添加，也可直接选以前用过的"
        />
      </div>
      <div class="field">
        <span>标签</span>
        <TagSelector
          v-model="form.tags as { name: string; type: string }[]"
          :exclude-types="['ingredient', 'taste', 'scene']"
        />
      </div>
    </section>
    <section class="group">
      <h2>上手难度</h2>
      <div class="field">
        <span>难度</span>
        <RatingInput v-model="form.difficulty as number | null" label="上手难度" />
      </div>
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
import { previewSource } from '../api/source'
import { ApiError } from '../api/client'
import type { Attachment, DishPayload } from '../types'
import { STATUS_LABEL, joinIngredientTags, toIngredientTags } from '../types'
import ImageUploader from '../components/ImageUploader.vue'
import RatingInput from '../components/RatingInput.vue'
import TagSelector from '../components/TagSelector.vue'

const route = useRoute()
const router = useRouter()
const id = computed(() => Number(route.params.id) || 0)
const isEdit = computed(() => route.name === 'dish-edit')
const error = ref('')
const previewing = ref(false)
const previewHint = ref('')
const lastPreviewed = ref('')
const coverId = ref<number | undefined>()
const ingredients = ref<{ name: string; type: string }[]>([])
const form = reactive<DishPayload>({
  name: '',
  status: 'want_to_cook',
  cover_image_url: '',
  source_url: '',
  source_platform: '',
  description: '',
  main_ingredients: '',
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
    description: dish.description,
    main_ingredients: dish.main_ingredients,
    difficulty: dish.difficulty,
    tags: (dish.tags || [])
      .filter((t) => t.type !== 'ingredient' && t.type !== 'taste' && t.type !== 'scene')
      .map((t) => ({ name: t.name, type: t.type })),
  })
  ingredients.value = toIngredientTags(dish.main_ingredients)
})

function onCover(att: Attachment) {
  coverId.value = att.id
  form.cover_image_url = att.file_url
}

function looksLikeURL(raw: string) {
  return /^https?:\/\//i.test(raw)
}

function emptyText(v?: string | null) {
  return !v || !String(v).trim()
}

async function onSourceBlur() {
  const url = (form.source_url || '').trim()
  form.source_url = url
  previewHint.value = ''
  if (!url) return
  if (!looksLikeURL(url)) {
    previewHint.value = '请填写以 http(s) 开头的链接'
    return
  }
  if (url === lastPreviewed.value) return
  previewing.value = true
  try {
    const data = await previewSource(url)
    lastPreviewed.value = url
    if (emptyText(form.name) && data.name) form.name = data.name
    if (emptyText(form.source_platform) && data.source_platform) form.source_platform = data.source_platform
    if (emptyText(form.cover_image_url) && data.cover_image_url) form.cover_image_url = data.cover_image_url
    if (!ingredients.value.length && data.main_ingredients) {
      ingredients.value = toIngredientTags(data.main_ingredients)
    }
    const filled = Boolean(data.name || data.cover_image_url || data.main_ingredients)
    if (data.partial && !filled) {
      previewHint.value = data.source_platform
        ? `已识别为${data.source_platform}，未能从该链接识别更多信息，可继续手填`
        : '未能从该链接识别更多信息，可继续手填'
    } else {
      previewHint.value = '已根据链接预填，请核对'
    }
  } catch (e) {
    lastPreviewed.value = ''
    if (e instanceof ApiError && e.status === 400) {
      previewHint.value = e.message || '无效的链接'
    } else {
      previewHint.value = '未能从该链接识别更多信息，可继续手填'
    }
  } finally {
    previewing.value = false
  }
}

async function save() {
  error.value = ''
  try {
    const payload = {
      ...form,
      main_ingredients: joinIngredientTags(ingredients.value),
      tags: (form.tags || []).filter((t) => t.type !== 'ingredient' && t.type !== 'taste' && t.type !== 'scene'),
      cover_attachment_id: coverId.value,
    }
    const dish = isEdit.value ? await updateDish(id.value, payload) : await createDish(payload)
    router.push(`/dishes/${dish.id}`)
  } catch (e) {
    error.value = e instanceof Error ? e.message : '保存失败'
  }
}
</script>
