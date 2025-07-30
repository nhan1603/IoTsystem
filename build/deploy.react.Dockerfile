FROM --platform=linux/amd64 node:22.6-alpine AS build-stage
# Set the working directory in the container
WORKDIR /app
# Copy package.json and package-lock.json
COPY package*.json ./
# Install dependencies
RUN npm install
# Copy the rest of the application code
COPY . .
# Build the app
RUN npm run build

# Stage 2
# Use Nginx to serve the static files
FROM nginx:alpine
COPY nginx.conf /etc/nginx/conf.d/default.conf
# Copy the build output to replace the default nginx contents.
COPY --from=build-stage /app/dist /usr/share/nginx/html
EXPOSE 5339
# Start Nginx server
CMD ["nginx", "-g", "daemon off;"]
