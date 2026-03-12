import { ofetch } from 'ofetch'

export const api = ofetch.create({
  baseURL: '/api/v1',
  onRequest({ options }) {
    const token = localStorage.getItem('token')
    if (token) {
      options.headers = new Headers(options.headers)
      options.headers.set('Authorization', `Bearer ${token}`)
    }
  },
})
