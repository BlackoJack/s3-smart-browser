// S3 Smart Browser Service Worker
const CACHE_NAME = 's3-smart-browser-v1.0.8';
const STATIC_CACHE = 's3-static-v1';
const DYNAMIC_CACHE = 's3-dynamic-v1';

// Файлы для кэширования
const STATIC_FILES = [
    '/',
    '/static/css/style.css',
    '/static/js/app.js',
    '/static/images/logo.svg',
    '/static/images/favicon.svg',
    '/static/manifest.json'
];

// Установка Service Worker
self.addEventListener('install', (event) => {
    console.log('Service Worker: Installing...');
    
    event.waitUntil(
        caches.open(STATIC_CACHE)
            .then((cache) => {
                console.log('Service Worker: Caching static files');
                return cache.addAll(STATIC_FILES);
            })
            .then(() => {
                console.log('Service Worker: Installation complete');
                return self.skipWaiting();
            })
            .catch((error) => {
                console.error('Service Worker: Installation failed', error);
            })
    );
});

// Активация Service Worker
self.addEventListener('activate', (event) => {
    console.log('Service Worker: Activating...');
    
    event.waitUntil(
        caches.keys()
            .then((cacheNames) => {
                return Promise.all(
                    cacheNames.map((cacheName) => {
                        if (cacheName !== STATIC_CACHE && cacheName !== DYNAMIC_CACHE) {
                            console.log('Service Worker: Deleting old cache', cacheName);
                            return caches.delete(cacheName);
                        }
                    })
                );
            })
            .then(() => {
                console.log('Service Worker: Activation complete');
                return self.clients.claim();
            })
    );
});

// Перехват запросов
self.addEventListener('fetch', (event) => {
    const { request } = event;
    const url = new URL(request.url);
    
    // Пропускаем API запросы и внешние ресурсы
    if (url.pathname.startsWith('/api/') || 
        url.hostname !== location.hostname) {
        return;
    }
    
    event.respondWith(
        caches.match(request)
            .then((cachedResponse) => {
                // Возвращаем кэшированную версию если есть
                if (cachedResponse) {
                    console.log('Service Worker: Serving from cache', request.url);
                    return cachedResponse;
                }
                
                // Иначе загружаем из сети
                return fetch(request)
                    .then((response) => {
                        // Проверяем валидность ответа
                        if (!response || response.status !== 200 || response.type !== 'basic') {
                            return response;
                        }
                        
                        // Кэшируем динамические ресурсы
                        if (request.method === 'GET') {
                            const responseClone = response.clone();
                            caches.open(DYNAMIC_CACHE)
                                .then((cache) => {
                                    cache.put(request, responseClone);
                                });
                        }
                        
                        return response;
                    })
                    .catch((error) => {
                        console.error('Service Worker: Fetch failed', error);
                        
                        // Возвращаем fallback для навигационных запросов
                        if (request.mode === 'navigate') {
                            return caches.match('/');
                        }
                        
                        throw error;
                    });
            })
    );
});

// Обработка push уведомлений (для будущего использования)
self.addEventListener('push', (event) => {
    console.log('Service Worker: Push notification received');
    
    const options = {
        body: event.data ? event.data.text() : 'S3 Smart Browser notification',
        icon: '/static/images/icon-192x192.png',
        badge: '/static/images/favicon.svg',
        vibrate: [100, 50, 100],
        data: {
            dateOfArrival: Date.now(),
            primaryKey: 1
        },
        actions: [
            {
                action: 'explore',
                title: 'Open S3 Browser',
                icon: '/static/images/favicon.svg'
            },
            {
                action: 'close',
                title: 'Close',
                icon: '/static/images/favicon.svg'
            }
        ]
    };
    
    event.waitUntil(
        self.registration.showNotification('S3 Smart Browser', options)
    );
});

// Обработка кликов по уведомлениям
self.addEventListener('notificationclick', (event) => {
    console.log('Service Worker: Notification clicked');
    
    event.notification.close();
    
    if (event.action === 'explore') {
        event.waitUntil(
            clients.openWindow('/')
        );
    }
});

// Периодическая синхронизация (для будущего использования)
self.addEventListener('sync', (event) => {
    console.log('Service Worker: Background sync', event.tag);
    
    if (event.tag === 'background-sync') {
        event.waitUntil(
            // Здесь можно добавить логику синхронизации
            Promise.resolve()
        );
    }
});

// Обработка ошибок
self.addEventListener('error', (event) => {
    console.error('Service Worker: Error occurred', event.error);
});

self.addEventListener('unhandledrejection', (event) => {
    console.error('Service Worker: Unhandled promise rejection', event.reason);
});
