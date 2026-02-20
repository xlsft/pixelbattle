import { useUUID } from "@xlsft/nuxt";
import ky, { type HTTPError, type Options as _Options } from "ky"
import { useAuthStore } from "~/store/auth.store";

export const useServer = async <ResponseData = unknown, Response = ResponseData & ResponseError>(
    endpoint: string,
    options?: RequestOptions
): Promise<Response> => {

    const config = useRuntimeConfig().public
    const auth = typeof options?.auth === 'string' ? options.auth : options?.auth !== false ? await useAuthStore().getToken() : undefined 
    const id = useUUID()
    const headers = {
        'Accept': 'application/json',
        ...options?.auth !== false && auth ? { 'Authorization': `Bearer ${auth}` } : {},
        ...options?.headers ? options.headers : {},
    }
    
    let response: Response

    try {
        response = await ky(endpoint, {
            ...options,
            prefixUrl: `${config.baseUrl}/api`,
            headers,
            timeout: false,
            retry: 0,
            hooks: {
                beforeRequest: [(request) => Object.entries(headers).forEach(([key, value]) => request.headers.set(key, value as string))],
                beforeError: [
                    async (error: HTTPError) => {
                        response = (await error.response.json()) as Response;
                        error.name = 'ResponseError';
                        return error;
                    },
                ],
                ...options?.hooks ? options.hooks : {},
            },
        }).json<Response>(); 
        if (import.meta.client && (response as any).redirect) {
            const router = useRouter()
            router.push((response as any).redirect)
        }
        return response;
    } catch (error) {
        console.error(error);
        return response!;
    }
};