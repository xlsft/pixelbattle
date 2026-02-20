import { useServer } from "~/composables/useServer"

export const useAuthStore = defineStore('Auth', () => {

    const token = ref<string>()
    const user = ref<UserModel>()

    const actions = {
        expired: () => { try {
            if (!token.value) return undefined
            const now = new Date()
            const exp = new Date(JSON.parse(atob(token.value.split('.')[1]!)).exp * 1000)
            if (!exp || exp < now) throw new Error('Expired')
            return false
        } catch { actions.logout(); return true }},
        getToken: () => {
            if (actions.expired()) return undefined
            else return token.value
        },
        update: async () => {
            if (!user.value?.id) return
            user.value = (await useServer<{ data: UserModel }>(`auth`)).data
            if (!user.value) actions.logout()
            return user.value
        },
        login: async (json: User) => {
            const response = await useServer<{ data: { user?: UserModel, token: string }}>('auth', { method: 'post', json, auth: false })
            if (response.data.user) user.value = response.data.user
            if (response.data.token) token.value = response.data.token
            return response
        },
        loginByInitData: async (json: { data: string }) => {
            const response = await useServer<{ data: { user?: UserModel, token: string }}>('auth/initdata', { method: 'post', json, auth: false })
            if (response.data.user) user.value = response.data.user
            if (response.data.token) token.value = response.data.token
            return response
        },
        logout: () => {
            token.value = undefined
            user.value = undefined
            if (import.meta.client) location.replace('/')
        },
    }
    
    return { ...actions, user, token }

}, { persist: true })