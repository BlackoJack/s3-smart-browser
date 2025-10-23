# S3 Smart Browser API Examples

## List Directory Contents

### Get root directory
```bash
curl "http://localhost:8080/api/list"
```

### Get specific directory
```bash
curl "http://localhost:8080/api/list?path=/documents"
```

### Response format
```json
{
  "path": "/documents",
  "files": [
    {
      "name": "folder1",
      "path": "/documents/folder1",
      "size": 0,
      "is_directory": true
    },
    {
      "name": "file.pdf",
      "path": "/documents/file.pdf",
      "size": 1024000,
      "is_directory": false
    }
  ]
}
```

## Download File

### Download file directly
```bash
curl -O "http://localhost:8080/api/download?file=/documents/file.pdf"
```

### Download with custom filename
```bash
curl -o "my-file.pdf" "http://localhost:8080/api/download?file=/documents/file.pdf"
```

### Get download URL (for JavaScript)
```javascript
const downloadUrl = `/api/download?file=${encodeURIComponent(filePath)}`;
window.open(downloadUrl, '_blank');
```

## Error Handling

### Error response format
```json
{
  "error": "File not found"
}
```

### Common HTTP status codes
- `200` - Success
- `400` - Bad Request (missing file parameter)
- `404` - File not found
- `500` - Internal Server Error

## JavaScript Integration

### Fetch directory listing
```javascript
async function loadDirectory(path) {
  try {
    const response = await fetch(`/api/list?path=${encodeURIComponent(path)}`);
    const data = await response.json();
    
    if (!response.ok) {
      throw new Error(data.error);
    }
    
    return data;
  } catch (error) {
    console.error('Error loading directory:', error);
    throw error;
  }
}
```

### Download file programmatically
```javascript
function downloadFile(filePath) {
  const downloadUrl = `/api/download?file=${encodeURIComponent(filePath)}`;
  const link = document.createElement('a');
  link.href = downloadUrl;
  link.download = filePath.split('/').pop();
  document.body.appendChild(link);
  link.click();
  document.body.removeChild(link);
}
```
