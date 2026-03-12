<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { compareStreams, type DiffEntry } from '@/api/content'

const props = defineProps<{
  owner: string
  loom: string
  sourceStream: string
  targetStream: string
}>()

const diffs = ref<DiffEntry[]>([])
const loading = ref(true)
const error = ref('')
const expanded = ref<Set<string>>(new Set())

onMounted(async () => {
  try {
    diffs.value = await compareStreams(props.owner, props.loom, props.sourceStream, props.targetStream)
  } catch {
    error.value = 'Failed to load diff'
  } finally {
    loading.value = false
  }
})

function toggle(path: string) {
  if (expanded.value.has(path)) {
    expanded.value.delete(path)
  } else {
    expanded.value.add(path)
  }
}

const stats = computed(() => {
  let added = 0, modified = 0, deleted = 0
  for (const d of diffs.value) {
    if (d.status === 'added') added++
    else if (d.status === 'modified') modified++
    else deleted++
  }
  return { added, modified, deleted }
})

function statusColor(status: string) {
  if (status === 'added') return 'text-green-400'
  if (status === 'modified') return 'text-yellow-400'
  return 'text-red-400'
}

function statusBg(status: string) {
  if (status === 'added') return 'bg-green-900/20'
  if (status === 'modified') return 'bg-yellow-900/20'
  return 'bg-red-900/20'
}

function statusIcon(status: string) {
  if (status === 'added') return '+'
  if (status === 'deleted') return '-'
  return '~'
}

interface DiffLine {
  type: 'context' | 'add' | 'remove' | 'header'
  content: string
  oldNum?: number
  newNum?: number
}

function computeUnifiedDiff(entry: DiffEntry): DiffLine[] {
  if (entry.status === 'added') {
    return (entry.source_content || '').split('\n').map((line, i) => ({
      type: 'add' as const,
      content: line,
      newNum: i + 1,
    }))
  }
  if (entry.status === 'deleted') {
    return (entry.target_content || '').split('\n').map((line, i) => ({
      type: 'remove' as const,
      content: line,
      oldNum: i + 1,
    }))
  }

  // Modified — simple line-by-line diff
  const oldLines = (entry.target_content || '').split('\n')
  const newLines = (entry.source_content || '').split('\n')
  const lines: DiffLine[] = []

  // Simple LCS-based diff
  const lcs = computeLCS(oldLines, newLines)
  let oi = 0, ni = 0, li = 0
  while (oi < oldLines.length || ni < newLines.length) {
    if (li < lcs.length && oi < oldLines.length && ni < newLines.length && oldLines[oi] === lcs[li] && newLines[ni] === lcs[li]) {
      lines.push({ type: 'context', content: oldLines[oi], oldNum: oi + 1, newNum: ni + 1 })
      oi++; ni++; li++
    } else if (li < lcs.length && oi < oldLines.length && oldLines[oi] !== lcs[li]) {
      lines.push({ type: 'remove', content: oldLines[oi], oldNum: oi + 1 })
      oi++
    } else if (li < lcs.length && ni < newLines.length && newLines[ni] !== lcs[li]) {
      lines.push({ type: 'add', content: newLines[ni], newNum: ni + 1 })
      ni++
    } else if (oi < oldLines.length) {
      lines.push({ type: 'remove', content: oldLines[oi], oldNum: oi + 1 })
      oi++
    } else if (ni < newLines.length) {
      lines.push({ type: 'add', content: newLines[ni], newNum: ni + 1 })
      ni++
    } else {
      break
    }
  }
  return lines
}

function computeLCS(a: string[], b: string[]): string[] {
  const m = a.length, n = b.length
  // For large files, skip LCS
  if (m > 5000 || n > 5000) return []
  const dp: number[][] = Array.from({ length: m + 1 }, () => Array(n + 1).fill(0))
  for (let i = 1; i <= m; i++) {
    for (let j = 1; j <= n; j++) {
      dp[i][j] = a[i - 1] === b[j - 1] ? dp[i - 1][j - 1] + 1 : Math.max(dp[i - 1][j], dp[i][j - 1])
    }
  }
  const result: string[] = []
  let i = m, j = n
  while (i > 0 && j > 0) {
    if (a[i - 1] === b[j - 1]) {
      result.unshift(a[i - 1])
      i--; j--
    } else if (dp[i - 1][j] > dp[i][j - 1]) {
      i--
    } else {
      j--
    }
  }
  return result
}
</script>

<template>
  <div>
    <div v-if="loading" class="text-gray-500 text-sm">Loading diff...</div>
    <div v-else-if="error" class="text-red-400 text-sm">{{ error }}</div>
    <div v-else-if="diffs.length === 0" class="border border-gray-800 rounded-lg p-6 text-center">
      <p class="text-gray-500">No differences between streams.</p>
    </div>
    <div v-else>
      <!-- Summary -->
      <div class="flex items-center gap-4 mb-4 text-xs text-gray-500">
        <span>{{ diffs.length }} file{{ diffs.length !== 1 ? 's' : '' }} changed</span>
        <span v-if="stats.added" class="text-green-400">{{ stats.added }} added</span>
        <span v-if="stats.modified" class="text-yellow-400">{{ stats.modified }} modified</span>
        <span v-if="stats.deleted" class="text-red-400">{{ stats.deleted }} deleted</span>
      </div>

      <!-- File list -->
      <div class="space-y-2">
        <div v-for="entry in diffs" :key="entry.path" class="border border-gray-800 rounded-lg overflow-hidden">
          <!-- File header -->
          <button
            @click="toggle(entry.path)"
            class="w-full flex items-center gap-3 px-4 py-2 text-left hover:bg-gray-800/50 transition-colors"
            :class="statusBg(entry.status)"
          >
            <span class="text-xs font-bold w-4" :class="statusColor(entry.status)">{{ statusIcon(entry.status) }}</span>
            <span class="text-sm font-mono text-gray-200 flex-1">{{ entry.path }}</span>
            <span class="text-xs capitalize" :class="statusColor(entry.status)">{{ entry.status }}</span>
            <svg
              class="w-4 h-4 text-gray-500 transition-transform"
              :class="expanded.has(entry.path) ? 'rotate-180' : ''"
              fill="none" viewBox="0 0 24 24" stroke="currentColor"
            >
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 9l-7 7-7-7" />
            </svg>
          </button>

          <!-- Diff content -->
          <div v-if="expanded.has(entry.path)" class="overflow-x-auto">
            <table v-if="entry.source_content || entry.target_content" class="w-full text-xs font-mono">
              <tr
                v-for="(line, i) in computeUnifiedDiff(entry)"
                :key="i"
                :class="{
                  'bg-green-900/20': line.type === 'add',
                  'bg-red-900/20': line.type === 'remove',
                }"
              >
                <td class="text-gray-600 text-right px-2 py-0 select-none w-12 align-top border-r border-gray-800">{{ line.oldNum ?? '' }}</td>
                <td class="text-gray-600 text-right px-2 py-0 select-none w-12 align-top border-r border-gray-800">{{ line.newNum ?? '' }}</td>
                <td class="px-1 py-0 select-none w-4 text-center align-top" :class="{
                  'text-green-400': line.type === 'add',
                  'text-red-400': line.type === 'remove',
                }">{{ line.type === 'add' ? '+' : line.type === 'remove' ? '-' : ' ' }}</td>
                <td class="px-2 py-0 whitespace-pre" :class="{
                  'text-green-300': line.type === 'add',
                  'text-red-300': line.type === 'remove',
                  'text-gray-400': line.type === 'context',
                }">{{ line.content }}</td>
              </tr>
            </table>
            <div v-else class="px-4 py-3 text-xs text-gray-500">Binary file or content unavailable</div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>
