<template>
    <div class="p-4">
        <div class="flex items-center justify-between mb-4">
        <h2 class="text-xl font-semibold">Stocks</h2>
        <button @click="reload" class="px-3 py-1 rounded bg-blue-600 text-white">Refrescar</button>
        </div>

        <div v-if="store.loading">Cargando...</div>
        <div v-if="store.error" class="text-red-600">{{ store.error }}</div>

        <table v-if="!store.loading && store.items.length" class="min-w-full bg-white">
        <thead>
            <tr>
            <th class="px-4 py-2 text-left">Ticker</th>
            <th class="px-4 py-2 text-left">Company</th>
            <th class="px-4 py-2 text-left">Brokerage</th>
            <th class="px-4 py-2 text-left">Rating From</th>
            <th class="px-4 py-2 text-left">Rating To</th>
            <th class="px-4 py-2 text-left">Target From</th>
            <th class="px-4 py-2 text-left">Target To</th>
            </tr>
        </thead>
        <tbody>
            <tr v-for="(s, idx) in store.items" :key="idx" class="hover:bg-gray-50">
            <td class="border px-4 py-2">{{ getField(s, ['Ticker','ticker']) }}</td>
            <td class="border px-4 py-2">{{ getField(s, ['Company','company']) }}</td>
            <td class="border px-4 py-2">{{ getField(s, ['Brokerage','brokerage']) }}</td>
            <td class="border px-4 py-2">{{ getField(s, ['Rating From','rating_from']) }}</td>
            <td class="border px-4 py-2">{{ getField(s, ['Rating To','rating_to']) }}</td>
            <td class="border px-4 py-2">{{ getField(s, ['Target From','target_from']) }}</td>
            <td class="border px-4 py-2">{{ getField(s, ['Target To','target_to']) }}</td>
            </tr>
        </tbody>
        </table>

        <div v-if="!store.loading && !store.items.length">No hay acciones aún. Ejecuta ingestión en backend o presiona Refrescar.</div>
    </div>
</template>

<script lang="ts">
import { defineComponent } from 'vue'
import { useStocksStore } from '@/stores/stocks'

export default defineComponent({
    name: 'StocksTable',
    setup() {
        const store = useStocksStore()
        const reload = async () => {
        await store.fetch()
        }
        return { store, reload, getField }
    }
})

function getField(item: any, keys: string[]) {
    for (const k of keys) {
        if (item[k] !== undefined && item[k] !== null && item[k] !== '') return item[k]
    }
    // fallback: try raw_json parsed
    if (item.raw_json) {
        try {
        const r = typeof item.raw_json === 'string' ? JSON.parse(item.raw_json) : item.raw_json
        for (const k of ['Ticker','ticker','Company','company']) {
            if (r[k]) return r[k]
        }
        } catch {}
    }
    return '-'
}
</script>

<style scoped>
table { border-collapse: collapse; width: 100%; }
th, td { text-align: left; padding: 8px; }
</style>
