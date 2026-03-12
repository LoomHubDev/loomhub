<script setup lang="ts">
import { ref, computed, onMounted, watch } from 'vue'
import { useRoute } from 'vue-router'
import { useLoomStore } from '@/stores/loom'
import { getEntityContent, listEntities, listStreams, type Entity } from '@/api/content'
import { createHighlighter, type Highlighter } from 'shiki'

const route = useRoute()
const loomStore = useLoomStore()

const owner = computed(() => route.params.owner as string)
const loomName = computed(() => route.params.loom as string)
const filePath = computed(() => {
  const p = route.params.path
  return Array.isArray(p) ? p.join('/') : (p || '')
})
const streamQuery = computed(() => (route.query.stream as string) || '')
const stream = ref('')

const content = ref('')
const entity = ref<Entity | null>(null)
const loading = ref(true)
const error = ref('')
const highlightedHtml = ref('')
let highlighter: Highlighter | null = null

const lang = computed(() => {
  const ext = filePath.value.split('.').pop()?.toLowerCase() || ''
  const map: Record<string, string> = {
    ts: 'typescript', tsx: 'tsx', js: 'javascript', jsx: 'jsx',
    py: 'python', rb: 'ruby', go: 'go', rs: 'rust',
    java: 'java', kt: 'kotlin', swift: 'swift',
    c: 'c', cpp: 'cpp', h: 'c', hpp: 'cpp',
    cs: 'csharp', css: 'css', scss: 'scss',
    html: 'html', vue: 'vue', svelte: 'svelte',
    json: 'json', yaml: 'yaml', yml: 'yaml', toml: 'toml',
    md: 'markdown', sql: 'sql', sh: 'bash', bash: 'bash',
    xml: 'xml', dockerfile: 'dockerfile',
    dart: 'dart', zig: 'zig', lua: 'lua', php: 'php',
  }
  return map[ext] || 'text'
})

const lineCount = computed(() => content.value.split('\n').length)

async function load() {
  loading.value = true
  error.value = ''
  highlightedHtml.value = ''
  try {
    // Resolve stream
    let streamFilter = streamQuery.value
    if (!streamFilter) {
      const streams = await listStreams(owner.value, loomName.value)
      if (streams.length > 0) streamFilter = streams[0].id
    }
    stream.value = streamFilter
    const entities = await listEntities(owner.value, loomName.value, streamFilter)
    entity.value = entities.find(e => e.path === filePath.value) || null
    if (!entity.value) {
      error.value = 'Entity not found'
      return
    }
    content.value = await getEntityContent(owner.value, loomName.value, entity.value.object_ref)
    await highlight()
  } catch {
    error.value = 'Failed to load entity content'
  } finally {
    loading.value = false
  }
}

async function highlight() {
  if (!content.value) return
  try {
    if (!highlighter) {
      highlighter = await createHighlighter({
        themes: ['github-dark'],
        langs: [lang.value === 'text' ? 'text' : lang.value],
      })
    } else if (lang.value !== 'text') {
      await highlighter.loadLanguage(lang.value as any)
    }
    highlightedHtml.value = highlighter.codeToHtml(content.value, {
      lang: lang.value,
      theme: 'github-dark',
    })
  } catch {
    // Fallback: no highlighting
    highlightedHtml.value = ''
  }
}

onMounted(load)
watch(filePath, load)

function formatSize(bytes: number): string {
  if (bytes < 1024) return `${bytes} B`
  if (bytes < 1024 * 1024) return `${(bytes / 1024).toFixed(1)} KB`
  return `${(bytes / (1024 * 1024)).toFixed(1)} MB`
}
</script>

<template>
  <div>
    <!-- Breadcrumb -->
    <div class="flex items-center gap-2 mb-4 text-sm">
      <RouterLink
        :to="`/${owner}/${loomName}/tree/${stream}`"
        class="text-indigo-400 hover:text-indigo-300"
      >
        {{ loomName }}
      </RouterLink>
      <template v-for="(part, i) in filePath.split('/')" :key="i">
        <span class="text-gray-600">/</span>
        <span :class="i === filePath.split('/').length - 1 ? 'text-white font-medium' : 'text-gray-400'">{{ part }}</span>
      </template>
    </div>

    <div v-if="loading" class="text-gray-500 text-sm">Loading...</div>
    <div v-else-if="error" class="text-red-400 text-sm">{{ error }}</div>
    <div v-else class="border border-gray-800 rounded-lg overflow-hidden">
      <!-- File header -->
      <div class="flex items-center justify-between px-4 py-2 bg-gray-800/50 border-b border-gray-800 text-xs text-gray-400">
        <span>{{ lineCount }} lines</span>
        <span v-if="entity">{{ formatSize(entity.size) }}</span>
      </div>

      <!-- Code content -->
      <div v-if="highlightedHtml" class="code-viewer overflow-x-auto" v-html="highlightedHtml" />
      <pre v-else class="p-4 text-sm font-mono text-gray-300 overflow-x-auto whitespace-pre"><code>{{ content }}</code></pre>
    </div>
  </div>
</template>

<style scoped>
.code-viewer :deep(pre) {
  margin: 0;
  padding: 1rem;
  font-size: 0.8125rem;
  line-height: 1.5;
  background: transparent !important;
}
.code-viewer :deep(code) {
  font-family: ui-monospace, SFMono-Regular, 'SF Mono', Menlo, Consolas, 'Liberation Mono', monospace;
}
</style>
