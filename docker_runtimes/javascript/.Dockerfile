# Use an official node runtime as a parent image
FROM node:18-alpine3.17

# Set the working directory to /app
WORKDIR /app

# Copy the current directory contents into the container at /app
COPY . /app

# Run app.py when the container launches
CMD ["node", "app.js"]