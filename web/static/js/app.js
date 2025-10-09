class S3Browser {
    constructor() {
        this.currentPath = '/';
        this.init();
    }

    init() {
        this.loadDirectory(this.currentPath);
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
        this.updateBreadcrumb();
        
        const fileList = document.getElementById('fileList');
        fileList.innerHTML = '';

        // Добавляем кнопку "Назад" если не в корне
        if (listing.path !== '/') {
            const parentPath = this.getParentPath(listing.path);
            const backButton = this.createBackButton(parentPath);
            fileList.appendChild(backButton);
        }

        // Сортируем: сначала папки, потом файлы
        const sortedFiles = listing.files.sort((a, b) => {
            if (a.is_directory && !b.is_directory) return -1;
            if (!a.is_directory && b.is_directory) return 1;
            return a.name.localeCompare(b.name);
        });

        sortedFiles.forEach(file => {
            const fileElement = this.createFileElement(file);
            fileList.appendChild(fileElement);
        });
    }

    createBackButton(parentPath) {
        const div = document.createElement('div');
        div.className = 'file-item';
        div.innerHTML = `
            <div class="file-icon folder">📁</div>
            <div class="file-info">
                <div class="file-name" onclick="browser.navigateTo('${parentPath}')">..</div>
            </div>
        `;
        return div;
    }

    createFileElement(file) {
        const div = document.createElement('div');
        div.className = 'file-item';
        
        const icon = file.is_directory ? '📁' : '📄';
        const size = file.is_directory ? '' : this.formatFileSize(file.size);
        
        div.innerHTML = `
            <div class="file-icon ${file.is_directory ? 'folder' : 'file'}">${icon}</div>
            <div class="file-info">
                <div class="file-name" onclick="browser.handleItemClick('${file.path}', ${file.is_directory})">
                    ${this.escapeHtml(file.name)}
                </div>
                ${size ? `<div class="file-size">${size}</div>` : ''}
            </div>
            <div class="file-actions">
                ${!file.is_directory ? `
                    <a href="/api/download?file=${encodeURIComponent(file.path)}" 
                       class="btn btn-download" download>Download</a>
                ` : ''}
            </div>
        `;
        
        return div;
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
        const sizes = ['B', 'KB', 'MB', 'GB'];
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
        const fileList = document.getElementById('fileList');
        fileList.innerHTML = '<div class="loading">Loading...</div>';
    }

    showError(message) {
        const fileList = document.getElementById('fileList');
        fileList.innerHTML = `<div class="loading" style="color: #e74c3c;">Error: ${message}</div>`;
    }
}

// Инициализация приложения
const browser = new S3Browser();