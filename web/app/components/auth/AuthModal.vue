<script setup lang="ts">
    import TelegramAuthButton from './TelegramAuthButton.vue';
    import { useAuthStore } from '~/store/auth.store';

    const auth = useAuthStore()
</script>

<template>
    <TelegramAuthButton :id="7964362622" @data="async (data) => await auth.login(data)"/>
</template>
