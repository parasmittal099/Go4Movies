# Use the official Node.js image
FROM node:20-alpine

# Set the working directory
WORKDIR /app

# Copy package files first to leverage Docker caching
COPY package.json package-lock.json ./

# Install dependencies
RUN npm install

# Copy the rest of the application code
COPY . .

# Next.js default dev port
EXPOSE 3000

# Run install on container start, then launch dev server
CMD ["sh", "-c", "npm install && npm run dev"]
