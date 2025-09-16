# URL Shortener Frontend

A modern, responsive frontend for the URL Shortener microservices application built with vanilla HTML, CSS, and JavaScript.

## Features

### 🔗 URL Management
- **Shorten URLs**: Convert long URLs into short, shareable links
- **Copy to Clipboard**: One-click copying of shortened URLs
- **URL Validation**: Client-side validation for URL format and length
- **Real-time Search**: Filter and search through your shortened URLs
- **Delete URLs**: Remove unwanted shortened URLs with confirmation

### 📊 Analytics & Statistics
- **Dashboard Overview**: View total URLs, clicks, and performance metrics
- **Individual URL Stats**: Track clicks for each shortened URL
- **Today's Activity**: See real-time activity for the current day
- **Most Popular**: Identify your most clicked URLs

### 🔔 Notifications
- **Slack Integration**: Send notifications to configured Slack channels
- **Real-time Feedback**: Toast notifications for all user actions
- **Error Handling**: Comprehensive error messages and retry options

### 🎨 User Experience
- **Responsive Design**: Works perfectly on desktop, tablet, and mobile
- **Modern UI**: Clean, professional interface with smooth animations
- **Dark Theme Support**: Beautiful gradient design
- **Keyboard Shortcuts**: Quick access with keyboard shortcuts
- **Auto-refresh**: Automatic data updates every 30 seconds

## Quick Start

### Prerequisites
Make sure the backend services are running. You can start them using:

```bash
# From the project root directory
make dev

# Or using Docker Compose directly
docker-compose up --build
```

### Running the Frontend

1. **Simple Local Server** (Recommended):
   ```bash
   cd frontend
   python3 -m http.server 3000
   ```
   Then open http://localhost:3000

2. **Using Node.js**:
   ```bash
   cd frontend
   npx serve -p 3000
   ```

3. **Using PHP**:
   ```bash
   cd frontend
   php -S localhost:3000
   ```

4. **Direct File Access**:
   Simply open `index.html` in your web browser (note: some features may not work due to CORS restrictions)

## API Integration

The frontend integrates with the following backend services:

- **Link Service** (port 8001): URL creation and deletion
- **Redirect Service** (port 8002): URL redirection and click tracking
- **Stats Service** (port 8003): Analytics and statistics
- **Notification Service** (port 8004): Slack notifications

All requests go through the Nginx gateway at `http://localhost:8080`.

### API Endpoints Used

```javascript
const API_ENDPOINTS = {
    generate: 'http://localhost:8080/api/generate',      // PUT - Create short URL
    redirect: 'http://localhost:8080/r',                 // GET - Redirect to original
    stats: 'http://localhost:8080/api/stats',            // GET - Get statistics
    delete: 'http://localhost:8080/api/delete',          // DELETE - Delete URL
    notifications: 'http://localhost:8080/api/notifications' // POST - Send notification
};
```

## Configuration

### Backend URL Configuration
If your backend is running on a different host or port, update the API base URL in `script.js`:

```javascript
const API_BASE_URL = 'http://your-backend-host:8080';
```

### CORS Configuration
For production deployment, ensure your backend services are configured to accept requests from your frontend domain.

## Features Overview

### URL Shortening Form
- Input validation (minimum 15 characters, valid URL format)
- Loading states with spinner animations
- Success feedback with shortened URL display
- Copy to clipboard functionality

### URL Management Dashboard
- Real-time search and filtering
- Individual URL statistics
- Bulk operations support
- Delete confirmation modals

### Statistics Dashboard
- Total URLs created
- Total clicks across all URLs
- Today's click count
- Most popular URL identification

### Notification System
- Toast notifications for user feedback
- Error handling with retry options
- Success/error/info message types
- Auto-dismiss after 5 seconds

## Keyboard Shortcuts

- `Ctrl/Cmd + K`: Focus search input
- `Escape`: Close modals and notifications

## Browser Support

- Chrome 60+
- Firefox 55+
- Safari 12+
- Edge 79+

## Technology Stack

- **HTML5**: Semantic markup and accessibility
- **CSS3**: Modern styling with flexbox and grid
- **Vanilla JavaScript**: No framework dependencies
- **Font Awesome**: Icons and visual elements
- **CSS Custom Properties**: Theme consistency
- **Fetch API**: Modern HTTP requests
- **Clipboard API**: Copy functionality

## File Structure

```
frontend/
├── index.html          # Main HTML file
├── styles.css          # Stylesheet with responsive design
├── script.js           # JavaScript application logic
└── README.md          # This file
```

## Development

### Adding New Features

1. **New API Endpoint**: Update `API_ENDPOINTS` in `script.js`
2. **New UI Component**: Add HTML structure and corresponding CSS
3. **New Functionality**: Implement in `script.js` with proper error handling

### Styling Guidelines

- Use CSS custom properties for theme consistency
- Follow BEM methodology for CSS class naming
- Ensure responsive design for all components
- Use semantic HTML elements

### JavaScript Guidelines

- Use modern ES6+ features
- Implement proper error handling
- Use async/await for API calls
- Follow consistent naming conventions

## Deployment

### Static File Hosting
The frontend consists of static files and can be deployed to any web server:

- **Nginx**: Place files in web root directory
- **Apache**: Upload to htdocs or public_html
- **CDN**: Use services like Cloudflare, AWS CloudFront, etc.

### Docker Deployment
Create a simple Dockerfile for containerized deployment:

```dockerfile
FROM nginx:alpine
COPY . /usr/share/nginx/html
EXPOSE 80
```

### Environment-Specific Configuration
For different environments, update the `API_BASE_URL` in `script.js`:

```javascript
// Development
const API_BASE_URL = 'http://localhost:8080';

// Production
const API_BASE_URL = 'https://your-production-api.com';
```

## Troubleshooting

### Common Issues

1. **CORS Errors**:
   - Ensure backend services allow requests from your frontend domain
   - Use proper CORS headers in your backend configuration

2. **API Connection Failed**:
   - Verify backend services are running
   - Check the API_BASE_URL configuration
   - Ensure all required services are healthy

3. **Features Not Working**:
   - Check browser console for JavaScript errors
   - Verify browser compatibility
   - Ensure JavaScript is enabled

### Debug Mode
Open browser developer tools to see:
- Network requests and responses
- JavaScript console errors
- Performance metrics

## Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Test thoroughly
5. Submit a pull request

## License

This project is licensed under the MIT License - see the LICENSE file for details.
