FROM node:20-slim AS base
ENV PNPM_HOME="/pnpm"
ENV PATH="$PNPM_HOME:$PATH"
# Set the HOME environment variable to a writable directory
ENV HOME="/app"
RUN corepack enable

# Update and install ca-certificates. This package contains a set of trusted CA (Certificate Authority) certificates.
RUN apt-get update && apt-get -yq install ca-certificates

# Next.js uses sharp for optimizing images.
RUN npm install -g --arch=x64 --platform=linux --libc=glibc sharp

FROM base AS builder
# Check https://github.com/nodejs/docker-node/tree/b4117f9333da4138b03a546ec926ef50a31506c3#nodealpine to understand why libc6-compat might be needed.
# Set working directory
WORKDIR /app
RUN pnpm install turbo --global
COPY . .
RUN npx turbo prune web --docker

# Add lockfile and package.json's of isolated subworkspace
FROM base AS installer
WORKDIR /app

# First install the dependencies (as they change less often)
COPY .gitignore .gitignore
COPY --from=builder /app/out/json/ .
COPY --from=builder /app/out/pnpm-lock.yaml ./pnpm-lock.yaml
RUN pnpm install

# Build the project
COPY --from=builder /app/out/full/ ./
COPY turbo.json turbo.json

RUN npx turbo run build --filter=web

FROM base AS runner
WORKDIR /app

# For image optimization
ENV NEXT_SHARP_PATH=/usr/local/lib/node_modules/sharp
COPY --from=base --chown=nextjs:nodejs /usr/local/lib/node_modules/sharp /usr/local/lib/node_modules/sharp

# Don't run production as root
RUN addgroup --system --gid 1001 nodejs
RUN adduser --system --uid 1001 nextjs

# Ensure the non-root user has permissions to the HOME directory
RUN mkdir -p $HOME/.cache && chown -R nextjs:nodejs $HOME

USER nextjs

COPY --from=installer /app/apps/web/next.config.js .
COPY --from=installer /app/apps/web/package.json .

# Automatically leverage output traces to reduce image size
# https://nextjs.org/docs/advanced-features/output-file-tracing
COPY --from=installer --chown=nextjs:nodejs /app/apps/web/.next/standalone ./
COPY --from=installer --chown=nextjs:nodejs /app/apps/web/.next/static ./apps/web/.next/static
# COPY --from=installer --chown=nextjs:nodejs /app/apps/web/public ./apps/web/public

# Expose the port Next.js will run on
EXPOSE 3000

# Start the Next.js application
CMD node apps/web/server.js