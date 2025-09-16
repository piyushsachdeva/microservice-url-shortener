const http = require('http');
const httpProxy = require('http-proxy');

// Create a proxy server
const proxy = httpProxy.createProxyServer({});

// Add CORS headers to all responses
const addCorsHeaders = (res) => {
  res.setHeader('Access-Control-Allow-Origin', '*');
  res.setHeader('Access-Control-Allow-Methods', 'GET, POST, PUT, DELETE, OPTIONS');
  res.setHeader('Access-Control-Allow-Headers', 'DNT,User-Agent,X-Requested-With,If-Modified-Since,Cache-Control,Content-Type,Range,Authorization');
};

// Create the server
const server = http.createServer((req, res) => {
  // Add CORS headers
  addCorsHeaders(res);
  
  // Handle OPTIONS requests for CORS preflight
  if (req.method === 'OPTIONS') {
    res.writeHead(204);
    res.end();
    return;
  }
  
  // Route requests based on path
  if (req.url.startsWith('/api/generate') || req.url.startsWith('/api/delete')) {
    // Route to link service on port 8001
    req.url = req.url.replace('/api', '');
    proxy.web(req, res, { target: 'http://localhost:8001' });
  } else if (req.url.startsWith('/api/stats')) {
    // Route to stats service on port 8003
    req.url = req.url.replace('/api', '');
    proxy.web(req, res, { target: 'http://localhost:8003' });
  } else if (req.url.match(/^\/[a-zA-Z0-9]{8}$/)) {
    // Route to redirect service on port 8002
    req.url = '/redirect' + req.url;
    proxy.web(req, res, { target: 'http://localhost:8002' });
  } else if (req.url === '/health') {
    res.writeHead(200, { 'Content-Type': 'text/plain' });
    res.end('OK');
  } else {
    res.writeHead(404);
    res.end('Not Found');
  }
});

// Handle proxy errors
proxy.on('error', (err, req, res) => {
  console.error('Proxy error:', err);
  res.writeHead(500);
  res.end('Internal Server Error');
});

const PORT = 8080;
server.listen(PORT, () => {
  console.log(`Gateway server running on http://localhost:${PORT}`);
});
