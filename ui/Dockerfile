# Use a lightweight web server image
FROM nginx:alpine

# Copy the HTML and JavaScript files to the web server's default directory
COPY index.html /usr/share/nginx/html/

# Expose port 80 for the web server
EXPOSE 80

# Start the web server when the container starts
CMD ["nginx", "-g", "daemon off;"]
