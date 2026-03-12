import { describe, it, expect, vi, beforeEach } from 'vitest'
import { setActivePinia, createPinia } from 'pinia'

// Mock localStorage before importing the store (since the store reads localStorage at module level)
const storageMap = new Map<string, string>()
const localStorageMock = {
  getItem: vi.fn((key: string) => storageMap.get(key) ?? null),
  setItem: vi.fn((key: string, value: string) => { storageMap.set(key, value) }),
  removeItem: vi.fn((key: string) => { storageMap.delete(key) }),
  clear: vi.fn(() => { storageMap.clear() }),
  get length() { return storageMap.size },
  key: vi.fn((index: number) => [...storageMap.keys()][index] ?? null),
}
Object.defineProperty(globalThis, 'localStorage', { value: localStorageMock, writable: true })

// Mock the router
vi.mock('@/router', () => ({
  router: { push: vi.fn() },
}))

// Mock the users API
vi.mock('@/api/users', () => ({
  login: vi.fn(),
  register: vi.fn(),
  getCurrentUser: vi.fn(),
}))

import { useAuthStore } from '@/stores/auth'
import { router } from '@/router'
import * as usersApi from '@/api/users'

const mockRouter = vi.mocked(router)
const mockLogin = vi.mocked(usersApi.login)
const mockRegister = vi.mocked(usersApi.register)
const mockGetCurrentUser = vi.mocked(usersApi.getCurrentUser)

beforeEach(() => {
  storageMap.clear()
  vi.clearAllMocks()
  setActivePinia(createPinia())
})

describe('useAuthStore', () => {
  describe('initial state', () => {
    it('has null user initially', () => {
      const store = useAuthStore()
      expect(store.user).toBeNull()
    })

    it('reads token from localStorage on init', () => {
      storageMap.set('token', 'saved-token')
      const store = useAuthStore()
      expect(store.token).toBe('saved-token')
    })

    it('has null token when localStorage is empty', () => {
      const store = useAuthStore()
      expect(store.token).toBeNull()
    })

    it('isAuthenticated is false initially', () => {
      const store = useAuthStore()
      expect(store.isAuthenticated).toBe(false)
    })

    it('isAdmin is false initially', () => {
      const store = useAuthStore()
      expect(store.isAdmin).toBe(false)
    })
  })

  describe('login', () => {
    it('sets user and token from API response', async () => {
      const fakeUser = {
        id: 'u1',
        username: 'alice',
        email: 'alice@example.com',
        display_name: 'Alice',
        bio: '',
        is_admin: false,
        token: 'jwt-token-123',
      }
      mockLogin.mockResolvedValue(fakeUser)

      const store = useAuthStore()
      await store.login('alice', 'password123')

      expect(mockLogin).toHaveBeenCalledWith({ login: 'alice', password: 'password123' })
      expect(store.token).toBe('jwt-token-123')
      expect(store.user).toEqual(fakeUser)
    })

    it('saves token to localStorage', async () => {
      mockLogin.mockResolvedValue({ token: 'new-token' })

      const store = useAuthStore()
      await store.login('alice', 'pass')

      expect(localStorageMock.setItem).toHaveBeenCalledWith('token', 'new-token')
      expect(storageMap.get('token')).toBe('new-token')
    })

    it('navigates to home after login', async () => {
      mockLogin.mockResolvedValue({ token: 't' })

      const store = useAuthStore()
      await store.login('alice', 'pass')

      expect(mockRouter.push).toHaveBeenCalledWith('/')
    })

    it('isAuthenticated becomes true after login', async () => {
      mockLogin.mockResolvedValue({
        id: 'u1',
        username: 'alice',
        email: 'a@b.com',
        display_name: 'Alice',
        bio: '',
        is_admin: false,
        token: 't',
      })

      const store = useAuthStore()
      await store.login('alice', 'pass')

      expect(store.isAuthenticated).toBe(true)
    })
  })

  describe('register', () => {
    it('calls the register API and sets state', async () => {
      const fakeUser = {
        id: 'u2',
        username: 'bob',
        email: 'bob@example.com',
        display_name: 'Bob',
        bio: '',
        is_admin: false,
        token: 'reg-token',
      }
      mockRegister.mockResolvedValue(fakeUser)

      const store = useAuthStore()
      await store.register('bob', 'bob@example.com', 'password')

      expect(mockRegister).toHaveBeenCalledWith({
        username: 'bob',
        email: 'bob@example.com',
        password: 'password',
      })
      expect(store.token).toBe('reg-token')
      expect(store.user).toEqual(fakeUser)
      expect(storageMap.get('token')).toBe('reg-token')
      expect(mockRouter.push).toHaveBeenCalledWith('/')
    })
  })

  describe('logout', () => {
    it('clears user, token, and localStorage', async () => {
      mockLogin.mockResolvedValue({
        id: 'u1',
        username: 'alice',
        email: 'a@b.com',
        display_name: 'Alice',
        bio: '',
        is_admin: true,
        token: 't',
      })

      const store = useAuthStore()
      await store.login('alice', 'pass')

      // Verify logged in state first
      expect(store.isAuthenticated).toBe(true)
      expect(store.isAdmin).toBe(true)

      store.logout()

      expect(store.user).toBeNull()
      expect(store.token).toBeNull()
      expect(localStorageMock.removeItem).toHaveBeenCalledWith('token')
    })

    it('navigates to login page', () => {
      const store = useAuthStore()
      store.logout()
      expect(mockRouter.push).toHaveBeenCalledWith('/login')
    })
  })

  describe('fetchUser', () => {
    it('sets user from getCurrentUser API', async () => {
      const fakeUser = {
        id: 'u1',
        username: 'alice',
        email: 'a@b.com',
        display_name: 'Alice',
        bio: 'Hello',
        is_admin: false,
      }
      mockGetCurrentUser.mockResolvedValue(fakeUser)

      const store = useAuthStore()
      await store.fetchUser()

      expect(store.user).toEqual(fakeUser)
      expect(store.isAuthenticated).toBe(true)
    })

    it('calls logout on API error', async () => {
      mockGetCurrentUser.mockRejectedValue(new Error('Unauthorized'))

      const store = useAuthStore()
      await store.fetchUser()

      expect(store.user).toBeNull()
      expect(store.token).toBeNull()
      expect(mockRouter.push).toHaveBeenCalledWith('/login')
    })
  })

  describe('isAdmin computed', () => {
    it('returns true when user is admin', async () => {
      mockLogin.mockResolvedValue({
        id: 'u1',
        username: 'admin',
        email: 'admin@example.com',
        display_name: 'Admin',
        bio: '',
        is_admin: true,
        token: 't',
      })

      const store = useAuthStore()
      await store.login('admin', 'pass')

      expect(store.isAdmin).toBe(true)
    })

    it('returns false when user is not admin', async () => {
      mockLogin.mockResolvedValue({
        id: 'u1',
        username: 'user',
        email: 'user@example.com',
        display_name: 'User',
        bio: '',
        is_admin: false,
        token: 't',
      })

      const store = useAuthStore()
      await store.login('user', 'pass')

      expect(store.isAdmin).toBe(false)
    })
  })
})
