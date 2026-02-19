import tailwindcss from "@tailwindcss/vite";

export default defineNuxtConfig({
    compatibilityDate: '2025-07-15',
    devtools: { enabled: true },
    css: ['~/assets/css/main.css'],  
    modules: ['@pinia/nuxt', 'pinia-plugin-persistedstate/nuxt', '@nuxt/scripts'],
    ssr: false,
    vite: { 
        plugins: [ tailwindcss() as any ], 
        server: { allowedHosts: ['.tuna.am'] },
    },
    app: {
        head: {
            title: "xlsft`s pixelbattle",
            meta: [
                { name: 'description', content: 'xlsft`s pixelbattle — это масштабная онлайн-битва пикселей, где каждый игрок влияет на общее полотно. Захватывай территорию, объединяйся с другими и оставь свой след в истории.' },
                { name: 'keywords', content: 'pixelbattle,pixel battle,xlsft,пиксель батл,пиксельная битва,онлайн игра,пиксель арт,multiplayer,canvas game,browser game,real-time game' },
                { name: 'robots', content: 'index, follow' },
                { name: 'author', content: 'xlsft' },
                { name: 'publisher', content: 'xlsft' },
                { name: 'application-name', content: "xlsft`s pixelbattle" },
                { name: 'theme-color', content: '#000000' },
                { name: 'color-scheme', content: 'dark' },
                { name: 'referrer', content: 'strict-origin-when-cross-origin' },
                { name: 'format-detection', content: 'telephone=no' },

                // Apple
                { name: 'apple-mobile-web-app-title', content: "xlsft`s pixelbattle" },
                { name: 'apple-mobile-web-app-capable', content: 'yes' },
                { name: 'apple-mobile-web-app-status-bar-style', content: 'black-translucent' },

                // Open Graph
                { property: 'og:title', content: 'xlsft`s pixelbattle — участвуй в битве пикселей' },
                { property: 'og:description', content: 'xlsft`s pixelbattle. Один холст и бесконечная битва за каждый пиксель. Присоединяйся!' },
                { property: 'og:type', content: 'website' },
                { property: 'og:url', content: 'https://pixelbattle.xlsft.com' },
                { property: 'og:site_name', content: 'xlsft`s pixelbattle' },
                { property: 'og:locale', content: 'ru_RU' },
                { property: 'og:image', content: 'https://pixelbattle.xlsft.com/og_image.png' },
                { property: 'og:image:width', content: '1200' },
                { property: 'og:image:height', content: '627' },
                { property: 'og:image:alt', content: 'xlsft`s pixelbattle — участвуй в битве пикселей' },

                // Twitter
                { name: 'twitter:card', content: 'summary_large_image' },
                { name: 'twitter:title', content: 'xlsft`s pixelbattle — глобальная битва пикселей' },
                { name: 'twitter:description', content: 'xlsft`s pixelbattle. Один холст и бесконечная битва за каждый пиксель. Присоединяйся!' },
                { name: 'twitter:image', content: 'https://pixelbattle.xlsft.com/og_image.png' },
                { name: 'twitter:site', content: '@xlsft' },
                { name: 'twitter:creator', content: '@xlsft' },
            ]
        }
    },

})