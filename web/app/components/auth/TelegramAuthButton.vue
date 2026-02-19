<script setup lang="ts">
    const props = defineProps<{ id: number }>()

    const emits = defineEmits<{ data: [User] }>()
    const script = useScript({ src: 'https://telegram.org/js/telegram-widget.js?22', async: true, fetchpriority: 'high' }, { bundle: true, tagPriority: "critical" })

    const click = async () => {
        await script._loadPromise
        ;(window as any).Telegram.Login.widgetsOrigin = 'https://oauth.telegram.org'
        ;(window as any).Telegram.Login.auth({ bot_id: props.id }, (data: any) => {
            emits('data', data)
        })
    }
</script>

<template>
    <button mini @click="click">Войти через Telegram</button>
</template>