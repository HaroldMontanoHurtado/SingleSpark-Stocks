import { defineStore } from 'pinia'
import axios from 'axios'

export const useStocksStore = defineStore('stocks', {
    state: () => ({
        items: [] as any[],
        loading: false,
        error: null as string | null
    }),
    actions: {
        async fetch() {
        this.loading = true
        this.error = null
        try {
            const res = await axios.get('/api/stocks') // proxy via vite
            this.items = res.data
        } catch (e: any) {
            this.error = e?.message ?? 'fetch error'
        } finally {
            this.loading = false
        }
        }
    }
})
