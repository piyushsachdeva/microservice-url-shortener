// URL Shortener Frontend Script v2.3 - Proxy version
// API Configuration
console.log('Script starting to load...');
const API_BASE_URL = ''; // Use relative URLs for API calls
const DISPLAY_BASE_URL = window.location.origin; // Use full domain for display URLs

// Wait for backend to be ready before initializing
async function waitForBackend() {
    const maxRetries = 10;
    const retryDelay = 1000; // 1 second
    
    for (let i = 0; i < maxRetries; i++) {
        try {
            const response = await fetch(`${API_BASE_URL}/api/stats/health`);
            if (response.ok) {
                console.log('Backend is ready!');
                return true;
            }
        } catch (error) {
            console.log(`Waiting for backend... attempt ${i + 1}/${maxRetries}`);
        }
        await new Promise(resolve => setTimeout(resolve, retryDelay));
    }
    console.warn('Backend may not be ready, proceeding anyway...');
    return false;
}

// API Endpoints
const API_ENDPOINTS = {
    generate: `${API_BASE_URL}/api/generate`,
    redirect: `${DISPLAY_BASE_URL}/r`, // Full URL for display
    links: `${API_BASE_URL}/api/links`,
    stats: `${API_BASE_URL}/api/stats`,
    delete: `${API_BASE_URL}/api/delete`
};

// Global state
let urlData = [];
let filteredData = [];

// DOM Elements
console.log('Loading DOM elements...');
const elements = {
    shortenForm: document.getElementById('shortenForm'),
    longUrlInput: document.getElementById('longUrl'),
    refreshButton: document.getElementById('refreshBtn'),
    searchInput: document.getElementById('searchInput'),
    loadingSpinner: document.getElementById('loadingSpinner'),
    urlList: document.getElementById('urlList'),
    toast: document.getElementById('toast'),
    toastMessage: document.getElementById('toastMessage'),
    confirmModal: document.getElementById('confirmModal'),
    confirmMessage: document.getElementById('confirmMessage'),
    confirmOk: document.getElementById('confirmOk'),
    confirmCancel: document.getElementById('confirmCancel'),
    shortenBtn: document.getElementById('shortenBtn'),
    shortUrlInput: document.getElementById('shortUrl'),
    originalUrl: document.getElementById('originalUrl'),
    createdAt: document.getElementById('createdAt'),
    shortId: document.getElementById('shortId'),
    resultSection: document.getElementById('resultSection'),
    copyBtn: document.getElementById('copyBtn'),
    toastClose: document.getElementById('toastClose'),
    // Statistics elements
    totalUrls: document.getElementById('totalUrls'),
    totalClicks: document.getElementById('totalClicks'),
    todayClicks: document.getElementById('todayClicks'),
    mostPopular: document.getElementById('mostPopular')
};

console.log('Elements loaded:', Object.keys(elements).filter(key => elements[key] === null));

// Toast notification system
function showToast(message, type = 'info') {
    elements.toastMessage.textContent = message;
    elements.toast.className = `toast ${type} show`;
    
    setTimeout(() => {
        hideToast();
    }, 5000);
}

function hideToast() {
    elements.toast.classList.remove('show');
}

// Modal system
function showModal(message, onConfirm) {
    console.log('showModal called with message:', message);
    elements.confirmMessage.textContent = message;
    elements.confirmModal.classList.remove('hidden');
    elements.confirmModal.classList.add('show');
    console.log('Modal classes after adding show:', elements.confirmModal.className);
    
    elements.confirmOk.onclick = () => {
        console.log('Confirm OK clicked');
        hideModal();
        onConfirm();
    };
}

function hideModal() {
    elements.confirmModal.classList.remove('show');
    elements.confirmModal.classList.add('hidden');
}

// Loading state management
function setLoading(isLoading, element = elements.shortenBtn) {
    if (isLoading) {
        element.disabled = true;
        const originalHTML = element.innerHTML;
        element.dataset.originalHTML = originalHTML;
        element.innerHTML = '<i class="fas fa-spinner fa-spin"></i> Loading...';
    } else {
        element.disabled = false;
        element.innerHTML = element.dataset.originalHTML || element.innerHTML;
    }
}

// URL validation
function isValidURL(string) {
    try {
        const url = new URL(string);
        return url.protocol === 'http:' || url.protocol === 'https:';
    } catch (_) {
        return false;
    }
}

// Format date
function formatDate(dateString) {
    const date = new Date(dateString);
    return date.toLocaleString();
}

// Format relative time
function formatRelativeTime(dateString) {
    const date = new Date(dateString);
    const now = new Date();
    const diff = now - date;
    
    const seconds = Math.floor(diff / 1000);
    const minutes = Math.floor(seconds / 60);
    const hours = Math.floor(minutes / 60);
    const days = Math.floor(hours / 24);
    
    if (days > 0) return `${days} day${days > 1 ? 's' : ''} ago`;
    if (hours > 0) return `${hours} hour${hours > 1 ? 's' : ''} ago`;
    if (minutes > 0) return `${minutes} minute${minutes > 1 ? 's' : ''} ago`;
    return 'Just now';
}

// API functions
async function apiRequest(url, options = {}) {
    try {
        console.log('Making API request to:', url);
        console.log('Full URL resolved to:', new URL(url, window.location.href).href);
        
        const response = await fetch(url, {
            headers: {
                'Content-Type': 'application/json',
                ...options.headers
            },
            ...options
        });
        
        console.log('Response status:', response.status);
        console.log('Response URL:', response.url);
        
        if (!response.ok) {
            throw new Error(`HTTP error! status: ${response.status}`);
        }
        
        const text = await response.text();
        return text ? JSON.parse(text) : null;
    } catch (error) {
        console.error('API request failed:', error);
        throw error;
    }
}

// Shorten URL function
async function shortenUrl(longUrl) {
    return await apiRequest(API_ENDPOINTS.generate, {
        method: 'PUT',
        body: JSON.stringify({ long: longUrl })
    });
}

// Get statistics function
async function getStats() {
    console.log('getStats called');
    console.log('API_BASE_URL:', API_BASE_URL);
    console.log('API_ENDPOINTS.stats:', API_ENDPOINTS.stats);
    return await apiRequest(API_ENDPOINTS.stats);
}

// Delete URL function
async function deleteUrl(id) {
    return await apiRequest(API_ENDPOINTS.delete, {
        method: 'DELETE',
        body: JSON.stringify({ id: id })
    });
}

// URL shortening form handler
if (elements.shortenForm) {
    elements.shortenForm.addEventListener('submit', async (e) => {
        e.preventDefault();
        
        const longUrl = elements.longUrlInput.value.trim();
        
        if (!longUrl) {
            showToast('Please enter a URL', 'error');
            return;
        }
    
    if (longUrl.length < 15) {
        showToast('URL must be at least 15 characters long', 'error');
        return;
    }
    
    if (!isValidURL(longUrl)) {
        showToast('Please enter a valid URL (must start with http:// or https://)', 'error');
        return;
    }
    
    setLoading(true);
    
    try {
        const result = await shortenUrl(longUrl);
        
        // Populate result section
        const shortUrl = `${API_ENDPOINTS.redirect}/${result.id}`;
        elements.shortUrlInput.value = shortUrl;
        elements.originalUrl.textContent = result.original_url;
        elements.createdAt.textContent = formatDate(result.created_at);
        elements.shortId.textContent = result.id;
        
        // Show result section
        elements.resultSection.classList.remove('hidden');
        
        // Clear form
        elements.longUrlInput.value = '';
        
        // Refresh URL list
        await loadUrlList();
        
        showToast('URL shortened successfully!', 'success');
    } catch (error) {
        console.error('Error shortening URL:', error);
        showToast('Failed to shorten URL. Please try again.', 'error');
    } finally {
        setLoading(false);
    }
    });
} else {
    console.error('shortenForm element not found');
}

// Copy to clipboard function
if (elements.copyBtn) {
    elements.copyBtn.addEventListener('click', async () => {
        try {
        await navigator.clipboard.writeText(elements.shortUrlInput.value);
        showToast('URL copied to clipboard!', 'success');
        
        // Visual feedback
        const originalText = elements.copyBtn.innerHTML;
        elements.copyBtn.innerHTML = '<i class="fas fa-check"></i> Copied!';
        setTimeout(() => {
            elements.copyBtn.innerHTML = originalText;
        }, 2000);
    } catch (error) {
        // Fallback for older browsers
        elements.shortUrlInput.select();
        document.execCommand('copy');
        showToast('URL copied to clipboard!', 'success');
    }
    });
} else {
    console.error('copyBtn element not found');
}

// Load URL list function
async function loadUrlList() {
    elements.loadingSpinner.classList.remove('hidden');
    
    try {
        // Try to get basic URL list from link service first
        urlData = await apiRequest(API_ENDPOINTS.links);
        
        // Try to enrich with stats data if stats service is available
        try {
            const statsData = await getStats();
            // If stats are available, use them (they might have additional fields)
            if (statsData && Array.isArray(statsData)) {
                urlData = statsData;
            }
        } catch (statsError) {
            console.log('Stats service unavailable, using basic link data:', statsError);
            // Continue with basic link data, don't fail
        }
        
        filteredData = [...urlData];
        updateStatistics();
        renderUrlList();
    } catch (error) {
        console.error('Error loading URL list:', error);
        elements.urlList.innerHTML = `
            <div class="error-message">
                <i class="fas fa-exclamation-triangle"></i>
                <p>Failed to load URLs. Please check if the backend is running.</p>
                <button onclick="loadUrlList()" class="btn btn-primary">
                    <i class="fas fa-retry"></i> Retry
                </button>
            </div>
        `;
        showToast('Failed to load URL list', 'error');
    } finally {
        elements.loadingSpinner.classList.add('hidden');
    }
}

// Render URL list function
function renderUrlList() {
    if (filteredData.length === 0) {
        elements.urlList.innerHTML = `
            <div class="empty-state">
                <i class="fas fa-link" style="font-size: 3rem; color: #a0aec0; margin-bottom: 1rem;"></i>
                <p style="color: #718096; text-align: center;">No URLs found. Create your first short URL above!</p>
            </div>
        `;
        return;
    }
    
    elements.urlList.innerHTML = filteredData.map(url => {
        const clickCount = url.stats ? url.stats.length : 0;
        const shortUrl = `${API_ENDPOINTS.redirect}/${url.id}`;
        
        return `
            <div class="url-item" data-id="${url.id}">
                <div class="url-item-header">
                    <div class="url-item-info">
                        <h4>
                            <i class="fas fa-link"></i> 
                            <a href="${shortUrl}" target="_blank" style="color: #667eea; text-decoration: none;">
                                ${shortUrl}
                            </a>
                        </h4>
                        <p><strong>Original:</strong> ${url.original_url}</p>
                        <p><strong>Created:</strong> ${formatRelativeTime(url.created_at)}</p>
                    </div>
                    <div class="url-item-actions">
                        <button onclick="copyUrl('${shortUrl}')" class="btn btn-copy" title="Copy URL">
                            <i class="fas fa-copy"></i>
                        </button>
                        <button onclick="confirmDelete('${url.id}')" class="btn btn-danger" title="Delete URL">
                            <i class="fas fa-trash"></i>
                        </button>
                    </div>
                </div>
                <div class="url-stats">
                    <div class="stats-row">
                        <div class="stat-item">
                            <h5>${clickCount}</h5>
                            <p>Total Clicks</p>
                        </div>
                        <div class="stat-item">
                            <h5>${getTodayClicks(url.stats)}</h5>
                            <p>Today's Clicks</p>
                        </div>
                        <div class="stat-item">
                            <h5>${getUniqueClicks(url.stats)}</h5>
                            <p>Unique Clicks</p>
                        </div>
                    </div>
                </div>
            </div>
        `;
    }).join('');
}

// Update statistics function
function updateStatistics() {
    const totalUrls = urlData.length;
    const totalClicks = urlData.reduce((sum, url) => sum + (url.stats ? url.stats.length : 0), 0);
    const todayClicks = urlData.reduce((sum, url) => sum + getTodayClicks(url.stats), 0);
    
    // Find most popular URL
    let mostPopular = '-';
    if (urlData.length > 0) {
        const popularUrl = urlData.reduce((max, url) => {
            const clicks = url.stats ? url.stats.length : 0;
            const maxClicks = max.stats ? max.stats.length : 0;
            return clicks > maxClicks ? url : max;
        });
        mostPopular = popularUrl.id;
    }
    
    elements.totalUrls.textContent = totalUrls;
    elements.totalClicks.textContent = totalClicks;
    elements.todayClicks.textContent = todayClicks;
    elements.mostPopular.textContent = mostPopular;
}

// Helper functions for statistics
function getTodayClicks(stats) {
    if (!stats) return 0;
    const today = new Date().toDateString();
    return stats.filter(stat => new Date(stat.created_at).toDateString() === today).length;
}

function getUniqueClicks(stats) {
    // For simplicity, assuming each click is unique
    // In a real implementation, you'd track by IP or user ID
    return stats ? stats.length : 0;
}

// Search functionality
if (elements.searchInput) {
    elements.searchInput.addEventListener('input', (e) => {
        const searchTerm = e.target.value.toLowerCase();
        
        if (searchTerm === '') {
        filteredData = [...urlData];
    } else {
        filteredData = urlData.filter(url => 
            url.original_url.toLowerCase().includes(searchTerm) ||
            url.id.toLowerCase().includes(searchTerm)
        );
    }
    
    renderUrlList();
    });
} else {
    console.error('searchInput element not found');
}

// Refresh button handler
if (elements.refreshButton) {
    elements.refreshButton.addEventListener('click', () => {
        loadUrlList();
    });
} else {
    console.error('refreshButton element not found');
}

// Copy URL function
async function copyUrl(url) {
    try {
        await navigator.clipboard.writeText(url);
        showToast('URL copied to clipboard!', 'success');
    } catch (error) {
        showToast('Failed to copy URL', 'error');
    }
}

// Confirm delete function
function confirmDelete(id) {
    console.log('confirmDelete called with id:', id);
    const url = urlData.find(u => u.id === id);
    console.log('Found URL:', url);
    const shortUrl = `${API_ENDPOINTS.redirect}/${id}`;
    
    showModal(
        `Are you sure you want to delete this URL?\n\n${shortUrl}\n\nThis action cannot be undone.`,
        () => deleteUrlHandler(id)
    );
}

// Delete URL handler
async function deleteUrlHandler(id) {
    console.log('deleteUrlHandler called with id:', id);
    try {
        console.log('Calling deleteUrl API...');
        await deleteUrl(id);
        console.log('Delete API call successful');
        showToast('URL deleted successfully!', 'success');
        await loadUrlList();
    } catch (error) {
        console.error('Error deleting URL:', error);
        showToast('Failed to delete URL', 'error');
    }
}

// Event listeners for modal and toast (with null checks)
if (elements.toastClose) {
    elements.toastClose.addEventListener('click', hideToast);
} else {
    console.error('toastClose element not found');
}

if (elements.confirmCancel) {
    elements.confirmCancel.addEventListener('click', hideModal);
} else {
    console.error('confirmCancel element not found');
}

if (elements.confirmModal) {
    // Close modal when clicking outside
    elements.confirmModal.addEventListener('click', (e) => {
        if (e.target === elements.confirmModal) {
            hideModal();
        }
    });
} else {
    console.error('confirmModal element not found');
}

// Keyboard shortcuts
document.addEventListener('keydown', (e) => {
    // ESC to close modal or toast
    if (e.key === 'Escape') {
        hideModal();
        hideToast();
    }
    
    // Ctrl/Cmd + K to focus search
    if ((e.ctrlKey || e.metaKey) && e.key === 'k') {
        e.preventDefault();
        elements.searchInput.focus();
    }
});

// Auto-refresh every 30 seconds
setInterval(() => {
    if (document.visibilityState === 'visible') {
        loadUrlList();
    }
}, 30000);

// Page visibility change handler
document.addEventListener('visibilitychange', () => {
    if (!document.hidden) {
        loadUrlList();
    }
});

// Initialize the application
document.addEventListener('DOMContentLoaded', async () => {
    console.log('URL Shortener Frontend initialized');
    
    // Wait for backend to be ready
    await waitForBackend();
    
    // Load initial data
    loadUrlList();
    
    // Focus on URL input
    elements.longUrlInput.focus();
    
    // Check if backend is running
    setTimeout(async () => {
        try {
            await fetch(`${API_BASE_URL}/health`);
            console.log('Backend connection established');
        } catch (error) {
            showToast('Backend server is not running. Please start the backend first.', 'error');
        }
    }, 1000);
});

// Export functions for global access
window.copyUrl = copyUrl;
window.confirmDelete = confirmDelete;
window.loadUrlList = loadUrlList;
