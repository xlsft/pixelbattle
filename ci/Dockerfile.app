FROM node:24 AS base
ENV PNPM_HOME="/pnpm"
ENV PATH="$PNPM_HOME:$PATH"
RUN corepack enable

FROM base AS build
WORKDIR /app

COPY package.json pnpm-lock.yaml ./
RUN pnpm install --frozen-lockfile

COPY . .

ENV NODE_OPTIONS="--max-old-space-size=4096"
RUN pnpm build

FROM base AS prod
WORKDIR /app
COPY --from=build /app/node_modules /app/node_modules
COPY --from=build /app/.output /app/.output

CMD ["node", ".output/server/index.mjs"]
