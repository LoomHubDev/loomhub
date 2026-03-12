<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { useRoute } from 'vue-router'
import { useAuthStore } from '@/stores/auth'
import {
  getWeaveRequest, weaveWR, closeWR,
  listComments, createComment,
  listReviews, createReview,
  type WeaveRequest, type Comment, type Review,
} from '@/api/weave-requests'
import { timeAgo } from '@/utils/format'
import DiffViewer from '@/components/DiffViewer.vue'

const route = useRoute()
const auth = useAuthStore()
const owner = computed(() => route.params.owner as string)
const loomName = computed(() => route.params.loom as string)
const number = computed(() => Number(route.params.number))

const wr = ref<WeaveRequest | null>(null)
const comments = ref<Comment[]>([])
const reviews = ref<Review[]>([])
const loading = ref(true)

const newComment = ref('')
const submittingComment = ref(false)
const reviewStatus = ref<'approved' | 'changes_requested' | 'commented'>('approved')
const reviewBody = ref('')
const submittingReview = ref(false)
const showReviewForm = ref(false)

async function load() {
  loading.value = true
  try {
    const [wrData, commentsData, reviewsData] = await Promise.all([
      getWeaveRequest(owner.value, loomName.value, number.value),
      listComments(owner.value, loomName.value, number.value),
      listReviews(owner.value, loomName.value, number.value),
    ])
    wr.value = wrData
    comments.value = commentsData
    reviews.value = reviewsData
  } finally {
    loading.value = false
  }
}

onMounted(load)

async function addComment() {
  if (!newComment.value.trim()) return
  submittingComment.value = true
  try {
    const c = await createComment(owner.value, loomName.value, number.value, newComment.value)
    comments.value.push(c)
    newComment.value = ''
  } finally {
    submittingComment.value = false
  }
}

async function submitReview() {
  submittingReview.value = true
  try {
    const r = await createReview(owner.value, loomName.value, number.value, {
      status: reviewStatus.value,
      body: reviewBody.value,
    })
    reviews.value.push(r)
    reviewBody.value = ''
    showReviewForm.value = false
  } finally {
    submittingReview.value = false
  }
}

async function doWeave() {
  if (!confirm('Weave this request?')) return
  const updated = await weaveWR(owner.value, loomName.value, number.value)
  wr.value = updated
}

async function doClose() {
  if (!confirm('Close this weave request?')) return
  const updated = await closeWR(owner.value, loomName.value, number.value)
  wr.value = updated
}

function statusBadge(status: string) {
  if (status === 'open') return 'bg-green-900/50 text-green-400 border-green-800'
  if (status === 'woven') return 'bg-purple-900/50 text-purple-400 border-purple-800'
  return 'bg-gray-800 text-gray-400 border-gray-700'
}

function reviewBadge(status: string) {
  if (status === 'approved') return 'text-green-400'
  if (status === 'changes_requested') return 'text-yellow-400'
  return 'text-gray-400'
}
</script>

<template>
  <div v-if="loading" class="text-gray-500 text-sm">Loading...</div>
  <div v-else-if="wr" class="max-w-3xl">
    <!-- Header -->
    <div class="mb-6">
      <h2 class="text-xl font-semibold text-white">
        {{ wr.title }}
        <span class="text-gray-600 font-normal">#{{ wr.number }}</span>
      </h2>
      <div class="flex items-center gap-3 mt-2">
        <span class="text-xs px-2 py-0.5 rounded-full border capitalize" :class="statusBadge(wr.status)">
          {{ wr.status }}
        </span>
        <span class="text-sm text-gray-500">
          {{ wr.author }} wants to weave
          <span class="text-gray-300">{{ wr.source_stream }}</span>
          into
          <span class="text-gray-300">{{ wr.target_stream }}</span>
        </span>
      </div>
    </div>

    <!-- Description -->
    <div v-if="wr.description" class="p-4 bg-gray-900 border border-gray-800 rounded-lg mb-6 text-sm text-gray-300 whitespace-pre-wrap">{{ wr.description }}</div>

    <!-- Diff -->
    <div class="mb-6">
      <h3 class="text-sm font-medium text-gray-400 mb-3">Changes</h3>
      <DiffViewer
        :owner="owner"
        :loom="loomName"
        :source-stream="wr.source_stream"
        :target-stream="wr.target_stream"
      />
    </div>

    <!-- Reviews -->
    <div v-if="reviews.length" class="mb-6">
      <h3 class="text-sm font-medium text-gray-400 mb-3">Reviews</h3>
      <div v-for="r in reviews" :key="r.id" class="flex items-start gap-3 mb-3 p-3 bg-gray-900 border border-gray-800 rounded-lg">
        <div class="flex-1">
          <div class="flex items-center gap-2">
            <span class="text-sm font-medium text-white">{{ r.reviewer }}</span>
            <span class="text-xs capitalize" :class="reviewBadge(r.status)">{{ r.status.replace('_', ' ') }}</span>
            <span class="text-xs text-gray-600">{{ timeAgo(r.created_at) }}</span>
          </div>
          <p v-if="r.body" class="text-sm text-gray-400 mt-1">{{ r.body }}</p>
        </div>
      </div>
    </div>

    <!-- Comments -->
    <div class="mb-6">
      <h3 class="text-sm font-medium text-gray-400 mb-3">Discussion</h3>
      <div v-if="comments.length === 0" class="text-sm text-gray-600 mb-4">No comments yet</div>
      <div v-for="c in comments" :key="c.id" class="mb-3 p-3 bg-gray-900 border border-gray-800 rounded-lg">
        <div class="flex items-center gap-2 mb-1">
          <span class="text-sm font-medium text-white">{{ c.author }}</span>
          <span class="text-xs text-gray-600">{{ timeAgo(c.created_at) }}</span>
        </div>
        <p class="text-sm text-gray-300 whitespace-pre-wrap">{{ c.body }}</p>
      </div>

      <!-- Add comment -->
      <form v-if="auth.user" @submit.prevent="addComment" class="mt-4">
        <textarea
          v-model="newComment"
          rows="3"
          placeholder="Leave a comment..."
          class="w-full px-3 py-2 bg-gray-900 border border-gray-700 rounded-md text-white text-sm focus:border-indigo-500 focus:outline-none resize-y"
        />
        <div class="flex gap-2 mt-2">
          <button
            type="submit"
            :disabled="submittingComment || !newComment.trim()"
            class="px-3 py-1.5 bg-gray-700 hover:bg-gray-600 disabled:opacity-50 text-white text-sm rounded-md transition-colors"
          >
            Comment
          </button>
          <button
            type="button"
            @click="showReviewForm = !showReviewForm"
            class="px-3 py-1.5 text-gray-400 hover:text-gray-200 text-sm"
          >
            {{ showReviewForm ? 'Hide review' : 'Add review' }}
          </button>
        </div>
      </form>
    </div>

    <!-- Review form -->
    <div v-if="showReviewForm && auth.user" class="mb-6 p-4 bg-gray-900 border border-gray-800 rounded-lg">
      <h3 class="text-sm font-medium text-gray-400 mb-3">Submit Review</h3>
      <textarea
        v-model="reviewBody"
        rows="3"
        placeholder="Review comment..."
        class="w-full px-3 py-2 bg-gray-800 border border-gray-700 rounded-md text-white text-sm focus:border-indigo-500 focus:outline-none resize-y mb-3"
      />
      <div class="flex items-center gap-3">
        <select
          v-model="reviewStatus"
          class="px-3 py-1.5 bg-gray-800 border border-gray-700 rounded-md text-white text-sm focus:outline-none"
        >
          <option value="approved">Approve</option>
          <option value="changes_requested">Request changes</option>
          <option value="commented">Comment</option>
        </select>
        <button
          @click="submitReview"
          :disabled="submittingReview"
          class="px-4 py-1.5 bg-indigo-600 hover:bg-indigo-500 disabled:opacity-50 text-white text-sm rounded-md transition-colors"
        >
          Submit Review
        </button>
      </div>
    </div>

    <!-- Actions -->
    <div v-if="auth.user && wr.status === 'open'" class="flex gap-3 pt-4 border-t border-gray-800">
      <button
        @click="doWeave"
        class="px-4 py-2 bg-purple-600 hover:bg-purple-500 text-white text-sm rounded-md transition-colors"
      >
        Weave
      </button>
      <button
        v-if="auth.user.username === wr.author || auth.user.is_admin"
        @click="doClose"
        class="px-4 py-2 bg-gray-700 hover:bg-gray-600 text-white text-sm rounded-md transition-colors"
      >
        Close
      </button>
    </div>
  </div>
</template>
