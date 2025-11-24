class S3Browser {
    constructor() {
        this.currentPath = '/';
        this.files = [];
        this.searchQuery = '';
        this.sortColumn = 'name';
        this.sortDirection = 'asc';
        this.init();
    }

    init() {
        this.cacheDom();
        this.setupControls();
        this.loadVersion();
        this.updateCurrentYear();
        this.registerServiceWorker();
        this.loadDirectory(this.currentPath);
    }

    cacheDom() {
        this.fileListElement = document.getElementById('fileList');
        this.searchInput = document.getElementById('fileSearch');
        this.clearSearchButton = document.getElementById('clearSearch');
        this.sortSelect = document.getElementById('sortSelect');
    }

    setupControls() {
        if (this.searchInput) {
            this.searchInput.addEventListener('input', (event) => {
                this.searchQuery = event.target.value;
                this.renderFileTable();
            });
        }

        if (this.clearSearchButton) {
            this.clearSearchButton.addEventListener('click', () => {
                if (!this.searchInput) return;
                this.searchQuery = '';
                this.searchInput.value = '';
                this.renderFileTable();
            });
        }

        if (this.sortSelect) {
            this.sortSelect.value = `${this.sortColumn}-${this.sortDirection}`;
            this.sortSelect.addEventListener('change', (event) => {
                const [column, direction] = event.target.value.split('-');
                this.sortColumn = column;
                this.sortDirection = direction;
                this.renderFileTable();
            });
        }
    }

    async loadDirectory(path) {
        try {
            this.showLoading();
            
            const response = await fetch(`/api/list?path=${encodeURIComponent(path)}`);
            const data = await response.json();
            
            if (!response.ok) {
                throw new Error(data.error);
            }
            
            this.renderDirectory(data);
        } catch (error) {
            this.showError(error.message);
        }
    }

    renderDirectory(listing) {
        this.currentPath = listing.path;
        this.files = Array.isArray(listing.files) ? listing.files : [];
        this.updateBreadcrumb();
        this.renderFileTable();
    }

    renderFileTable() {
        if (!this.fileListElement) return;

        this.fileListElement.innerHTML = '';

        const table = document.createElement('div');
        table.className = 'file-table';
        table.appendChild(this.createTableHeader());

        if (this.currentPath !== '/') {
            table.appendChild(this.createParentRow(this.getParentPath(this.currentPath)));
        }

        const processedFiles = this.getProcessedFiles();

        if (processedFiles.length === 0) {
            this.fileListElement.appendChild(table);
            this.fileListElement.appendChild(this.createEmptyState());
            this.updateFooterStats(0);
            return;
        }

        processedFiles.forEach(file => {
            table.appendChild(this.createFileRow(file));
        });

        this.fileListElement.appendChild(table);
        this.updateFooterStats(processedFiles.length);
    }

    createTableHeader() {
        const header = document.createElement('div');
        header.className = 'file-row file-row-head';
        header.innerHTML = `
            <div class="file-cell cell-name">Имя</div>
            <div class="file-cell cell-size">Размер</div>
            <div class="file-cell cell-date">Изменён</div>
            <div class="file-cell cell-actions">Действия</div>
        `;
        return header;
    }

    createEmptyState() {
        const empty = document.createElement('div');
        empty.className = 'empty-state';
        if (this.searchQuery) {
            empty.textContent = 'Ничего не найдено. Попробуйте изменить запрос.';
        } else {
            empty.textContent = 'Папка пуста.';
        }
        return empty;
    }

    createParentRow(parentPath) {
        const row = document.createElement('div');
        row.className = 'file-row file-row-parent';
        row.innerHTML = `
            <div class="file-cell cell-name">
                <div class="file-icon folder">📁</div>
                <button class="file-name link-button" onclick="browser.navigateTo(${JSON.stringify(parentPath)})">..</button>
            </div>
            <div class="file-cell cell-size">—</div>
            <div class="file-cell cell-date">—</div>
            <div class="file-cell cell-actions">
                <button class="btn btn-ghost" onclick="browser.navigateTo(${JSON.stringify(parentPath)})">Назад</button>
            </div>
        `;
        return row;
    }

    createFileRow(file) {
        const row = document.createElement('div');
        row.className = 'file-row';

        const icon = this.getIconForItem(file);
        const size = file.is_directory ? '—' : this.formatFileSize(file.size);
        const modified = file.last_modified ? this.formatDate(file.last_modified) : '—';

        row.innerHTML = `
            <div class="file-cell cell-name">
                <div class="file-icon ${file.is_directory ? 'folder' : 'file'}">${icon}</div>
                <button class="file-name link-button" onclick="browser.handleItemClick(${JSON.stringify(file.path)}, ${file.is_directory})">
                    ${this.escapeHtml(file.name)}
                </button>
            </div>
            <div class="file-cell cell-size">${size}</div>
            <div class="file-cell cell-date">${modified}</div>
            <div class="file-cell cell-actions">
                ${!file.is_directory ? `
                    <a href="/api/download?file=${encodeURIComponent(file.path)}" 
                       class="btn btn-download" download>Скачать</a>
                ` : ''}
            </div>
        `;

        return row;
    }

    getProcessedFiles() {
        if (!Array.isArray(this.files)) return [];

        let files = [...this.files];

        if (this.searchQuery.trim()) {
            const query = this.searchQuery.trim().toLowerCase();
            files = files.filter(file => (file.name || '').toLowerCase().includes(query));
        }

        return files.sort((a, b) => this.sortFiles(a, b));
    }

    sortFiles(a, b) {
        if (a.is_directory !== b.is_directory) {
            return a.is_directory ? -1 : 1;
        }

        const direction = this.sortDirection === 'asc' ? 1 : -1;

        switch (this.sortColumn) {
            case 'size': {
                const sizeA = a.is_directory ? -1 : (a.size || 0);
                const sizeB = b.is_directory ? -1 : (b.size || 0);
                if (sizeA === sizeB) break;
                return direction * (sizeA - sizeB);
            }
            case 'date': {
                const dateA = a.last_modified ? new Date(a.last_modified).getTime() : 0;
                const dateB = b.last_modified ? new Date(b.last_modified).getTime() : 0;
                if (dateA === dateB) break;
                return direction * (dateA - dateB);
            }
            default:
                return direction * this.collator().compare(a.name, b.name);
        }

        return direction * this.collator().compare(a.name, b.name);
    }

    collator() {
        if (!this._collator) {
            this._collator = new Intl.Collator(undefined, { sensitivity: 'base', numeric: true });
        }
        return this._collator;
    }

    getIconForItem(file) {
        if (file.is_directory) return '📁';
        const mime = (file.mime_type || '').toLowerCase();
        const name = (file.name || '').toLowerCase();
        if (mime.startsWith('image/')) return '🖼️';
        if (mime.startsWith('video/')) return '🎬';
        if (mime.startsWith('audio/')) return '🎵';
        if (mime === 'application/pdf') return '📕';
        if (mime.includes('zip') || mime.includes('tar') || mime.includes('gzip') || mime.includes('7z')) return '🗜️';
        if (mime.includes('excel') || name.endsWith('.xlsx') || name.endsWith('.xls') || name.endsWith('.csv')) return '📊';
        if (mime.includes('word') || name.endsWith('.doc') || name.endsWith('.docx')) return '📝';
        if (mime.includes('powerpoint') || name.endsWith('.ppt') || name.endsWith('.pptx')) return '📈';
        if (mime.startsWith('text/') || name.endsWith('.txt') || name.endsWith('.md')) return '📄';
        if (name.endsWith('.json') || name.endsWith('.yaml') || name.endsWith('.yml') || name.endsWith('.xml')) return '🧾';
        if (name.endsWith('.js') || name.endsWith('.ts') || name.endsWith('.go') || name.endsWith('.py') || name.endsWith('.java') || name.endsWith('.rb') || name.endsWith('.php') || name.endsWith('.cpp') || name.endsWith('.c') || name.endsWith('.cs') || name.endsWith('.sh')) return '💻';
        return '📦';
    }

    handleItemClick(path, isDirectory) {
        if (isDirectory) {
            this.navigateTo(path);
        } else {
            window.open(`/api/download?file=${encodeURIComponent(path)}`, '_blank');
        }
    }

    navigateTo(path) {
        this.loadDirectory(path);
    }

    formatDate(dateString) {
        const date = new Date(dateString);
        if (Number.isNaN(date.getTime())) {
            return '—';
        }

        return new Intl.DateTimeFormat('ru-RU', {
            year: 'numeric',
            month: 'short',
            day: '2-digit',
            hour: '2-digit',
            minute: '2-digit'
        }).format(date);
    }

    updateBreadcrumb() {
        const breadcrumb = document.getElementById('breadcrumb');
        const parts = this.currentPath.split('/').filter(p => p);
        
        let html = '<a href="javascript:void(0)" onclick="browser.navigateTo(\'/\')">Home</a>';
        
        let currentPath = '';
        parts.forEach(part => {
            currentPath += '/' + part;
            html += ` / <a href="javascript:void(0)" onclick="browser.navigateTo('${currentPath}')">${this.escapeHtml(part)}</a>`;
        });
        
        breadcrumb.innerHTML = html;
    }

    getParentPath(path) {
        if (path === '/') return '/';
        const parts = path.split('/').filter(p => p);
        parts.pop();
        return parts.length ? '/' + parts.join('/') : '/';
    }

    formatFileSize(bytes) {
        if (bytes === 0) return '0 B';
        const k = 1024;
        const sizes = ['B', 'KB', 'MB', 'GB', 'TB'];
        const i = Math.floor(Math.log(bytes) / Math.log(k));
        return parseFloat((bytes / Math.pow(k, i)).toFixed(2)) + ' ' + sizes[i];
    }

    escapeHtml(unsafe) {
        return unsafe
            .replace(/&/g, "&amp;")
            .replace(/</g, "&lt;")
            .replace(/>/g, "&gt;")
            .replace(/"/g, "&quot;")
            .replace(/'/g, "&#039;");
    }

    showLoading() {
        if (this.fileListElement) {
            this.fileListElement.innerHTML = '<div class="loading">Loading...</div>';
        }
    }

    showError(message) {
        if (this.fileListElement) {
            this.fileListElement.innerHTML = `<div class="loading" style="color: #e74c3c;">Error: ${message}</div>`;
        }
        this.updateConnectionStatus(false);
    }

    async loadVersion() {
        try {
            const response = await fetch('/api/version');
            const data = await response.json();
            if (response.ok) {
                // Используем версию из API, если она есть, иначе показываем dev
                const version = data.version || 'dev';
                document.getElementById('appVersion').textContent = version.startsWith('v') ? version : `v${version}`;
                
                // Дополнительная информация в консоли для разработчиков
                if (data.git_commit) {
                    console.log(`S3 Smart Browser ${version} (${data.git_commit.substring(0, 8)})`);
                }
            }
        } catch (error) {
            console.warn('Failed to load version:', error);
            // Fallback на dev версию
            document.getElementById('appVersion').textContent = 'vdev';
        }
    }

    updateFooterStats(fileCount) {
        document.getElementById('fileCount').textContent = fileCount;
        this.updateConnectionStatus(true);
    }

    updateConnectionStatus(connected) {
        const statusElement = document.getElementById('connectionStatus');
        if (connected) {
            statusElement.textContent = 'Connected';
            statusElement.className = 'stat-value status-connected';
        } else {
            statusElement.textContent = 'Disconnected';
            statusElement.className = 'stat-value status-disconnected';
        }
    }

    updateCurrentYear() {
        const yearElement = document.getElementById('currentYear');
        if (yearElement) {
            yearElement.textContent = new Date().getFullYear();
        }
    }

    async registerServiceWorker() {
        if ('serviceWorker' in navigator) {
            try {
                const registration = await navigator.serviceWorker.register('/sw.js');
                console.log('Service Worker registered successfully:', registration);
                
                // Проверяем обновления
                registration.addEventListener('updatefound', () => {
                    const newWorker = registration.installing;
                    newWorker.addEventListener('statechange', () => {
                        if (newWorker.state === 'installed' && navigator.serviceWorker.controller) {
                            // Новый Service Worker установлен, показываем уведомление
                            this.showUpdateNotification();
                        }
                    });
                });
                
            } catch (error) {
                console.error('Service Worker registration failed:', error);
            }
        }
    }

    showUpdateNotification() {
        // Создаем уведомление об обновлении
        const notification = document.createElement('div');
        notification.className = 'update-notification';
        notification.innerHTML = `
            <div class="update-content">
                <span>🔄 New version available!</span>
                <button onclick="window.location.reload()" class="update-btn">Update</button>
                <button onclick="this.parentElement.parentElement.remove()" class="close-btn">×</button>
            </div>
        `;
        
        document.body.appendChild(notification);
        
        // Автоматически скрываем через 10 секунд
        setTimeout(() => {
            if (notification.parentElement) {
                notification.remove();
            }
        }, 10000);
    }
}

// Инициализация приложения
const browser = new S3Browser();