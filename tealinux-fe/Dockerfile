# ============================================
# Stage 1: Build
# ============================================
FROM node:20-alpine AS builder

WORKDIR /app

# Copy package files first for better layer caching
COPY package.json package-lock.json ./

# Install dependencies
RUN npm ci

# Copy source code
COPY . .

# Build the Astro project
RUN npm run build:prod

# ============================================
# Stage 2: Production
# ============================================
FROM node:20-alpine AS runtime

WORKDIR /app

# Copy built output and production dependencies
COPY --from=builder /app/dist ./dist
COPY --from=builder /app/node_modules ./node_modules
COPY --from=builder /app/package.json ./

# Set environment
ENV NODE_ENV=production
ENV HOST=0.0.0.0
ENV PORT=4321

EXPOSE 4321

# Start the Astro SSR server
CMD ["node", "./dist/server/entry.mjs"]
