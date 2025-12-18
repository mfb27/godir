/**
 * 格式化文件大小
 * @param {number} bytes - 字节数
 * @returns {string} 格式化后的文件大小字符串
 */
function formatFileSize(bytes) {
    if (bytes === 0) return '0 B';
    const k = 1024;
    const sizes = ['B', 'KB', 'MB', 'GB'];
    const i = Math.floor(Math.log(bytes) / Math.log(k));
    return Math.round(bytes / Math.pow(k, i) * 100) / 100 + ' ' + sizes[i];
}

/**
 * 更新上传进度
 * @param {string} itemId - 进度项ID
 * @param {number} percent - 进度百分比
 * @param {string} status - 状态文本
 */
function updateProgress(itemId, percent, status) {
    const item = document.getElementById(itemId);
    if (!item) return;

    const percentEl = item.querySelector('.progress-percent');
    const fillEl = item.querySelector('.progress-fill');

    if (percentEl) {
        percentEl.textContent = status || `${percent}%`;
    }
    if (fillEl) {
        fillEl.style.width = `${percent}%`;
    }
}

/**
 * 检查用户登录状态
 * @returns {Object|null} 用户信息对象或null
 */
function checkLoginStatus() {
    const token = localStorage.getItem('token');
    if (!token) {
        window.location.href = 'login.html';
        return null;
    }

    try {
        const payload = JSON.parse(atob(token.split('.')[1]));
        const userId = payload.userId;
        const username = payload.username || '用户';

        return { userId, username, token };
    } catch (e) {
        console.error('解析token失败:', e);
        localStorage.removeItem('token');
        window.location.href = 'login.html';
        return null;
    }
}

/**
 * 检查用户信息并在页面上显示用户头像和相关信息
 * @param {object} userInfo - 用户信息对象
 */
function displayUserInfo(userInfo) {
    if (!userInfo) return;
    
    const { userId, username } = userInfo;
    
    const avatarContainer = document.getElementById('userAvatarContainer');
    if (!avatarContainer) return;
    
    const avatar = document.getElementById('userAvatar');
    const dropdownAvatar = document.getElementById('dropdownAvatar');
    const dropdownName = document.getElementById('dropdownName');
    const dropdownId = document.getElementById('dropdownId');
    
    // 获取用户名首字母作为头像
    const avatarText = username.charAt(0).toUpperCase();
    
    if (avatar) avatar.textContent = avatarText;
    if (dropdownAvatar) dropdownAvatar.textContent = avatarText;
    if (dropdownName) dropdownName.textContent = username;
    if (dropdownId) dropdownId.textContent = userId;
    
    avatarContainer.style.display = 'block';
}

// /**
//  * 加载用户详细信息并显示
//  * @param {string} token - 用户认证token
//  */
// async function loadUserDetailsForDisplay(token) {
//     try {
//         const response = await fetch(`${API_BASE_URL}/user/profile`, {
//             headers: {
//                 'Authorization': `Bearer ${token}`
//             }
//         });
        
//         if (!response.ok) {
//             throw new Error(`HTTP error! status: ${response.status}`);
//         }
        
//         const data = await response.json();
//         if (data.code === 0 && data.data) {
//             displayUserDetails(data.data);
//         }
//     } catch (error) {
//         console.error('加载用户信息失败:', error);
//         // 如果加载详细信息失败，回退到基本显示
//         const userInfo = checkLoginStatus();
//         if (userInfo) {
//             displayUserInfo(userInfo);
//         }
//     }
// }

/**
 * 显示消息提示
 * @param {string} message - 提示消息
 * @param {string} type - 消息类型: success, error, info
 */
function showMessage(message, type = 'info') {
    // 如果已经存在提示框，则移除
    const existingToast = document.getElementById('toast-message');
    if (existingToast) {
        existingToast.remove();
    }

    // 创建提示框元素
    const toast = document.createElement('div');
    toast.id = 'toast-message';
    toast.style.cssText = `
        position: fixed;
        top: 20px;
        right: 20px;
        padding: 16px 24px;
        border-radius: 8px;
        color: white;
        font-weight: 500;
        z-index: 9999;
        box-shadow: 0 4px 12px rgba(0, 0, 0, 0.15);
        transform: translateX(100%);
        transition: transform 0.3s ease;
        ${type === 'success' ? 'background: #52c41a;' : ''}
        ${type === 'error' ? 'background: #ff4d4f;' : ''}
        ${type === 'info' ? 'background: #1890ff;' : ''}
    `;

    toast.textContent = message;
    document.body.appendChild(toast);

    // 动画显示
    setTimeout(() => {
        toast.style.transform = 'translateX(0)';
    }, 100);

    // 3秒后自动移除
    setTimeout(() => {
        toast.style.transform = 'translateX(100%)';
        setTimeout(() => {
            if (toast.parentNode) {
                toast.parentNode.removeChild(toast);
            }
        }, 300);
    }, 3000);
}

// 导出函数（如果使用模块系统）
if (typeof module !== 'undefined' && typeof exports !== 'undefined') {
    module.exports = {
        formatFileSize,
        updateProgress,
        checkLoginStatus,
        showMessage
    };
}


// 页面加载时检查登录状态并获取用户信息
const userInfo = checkLoginStatus();
if (!userInfo) {
    window.location.href = 'login.html';
}

const { userId, username, token } = userInfo;

// 全局变量
let uploadQueue = [];
let isUploading = false;

// DOM元素
const fileList = document.getElementById('fileList');
const emptyState = document.getElementById('emptyState');
const selectAllCheckbox = document.getElementById('selectAllCheckbox');
const batchDeleteBtn = document.getElementById('batchDeleteBtn');
const uploadProgress = document.getElementById('uploadProgress');

// 处理文件选择
async function handleFileSelect(event) {
    const files = Array.from(event.target.files);
    if (files.length === 0) return;

    const userInfo = checkLoginStatus();
    if (!userInfo) return;

    // 显示上传进度
    const progressContainer = document.getElementById('uploadProgress');
    const progressList = document.getElementById('progressList');
    
    // 添加防御性检查，确保元素存在再操作
    if (progressContainer && progressList) {
        progressContainer.style.display = 'block';
        progressList.innerHTML = '';

        // 为每个文件创建进度项
        const progressItems = {};
        files.forEach(file => {
            const itemId = `progress-${Date.now()}-${Math.random()}`;
            progressItems[file.name] = itemId;
            const progressItem = document.createElement('div');
            progressItem.className = 'progress-item';
            progressItem.id = itemId;
            progressItem.innerHTML = `
                <div class="progress-header">
                    <div class="progress-name">${file.name}</div>
                    <div class="progress-percent">0%</div>
                </div>
                <div class="progress-bar">
                    <div class="progress-fill" style="width: 0%"></div>
                </div>
            `;
            progressList.appendChild(progressItem);
        });

        // 上传每个文件
        for (const file of files) {
            try {
                await uploadFile(file, userInfo, progressItems[file.name]);
            } catch (error) {
                console.error(`上传文件 ${file.name} 失败:`, error);
                updateProgress(progressItems[file.name], 0, '上传失败');
            }
        }

        // 上传完成后刷新文件列表
        setTimeout(() => {
            loadFileList();
            progressContainer.style.display = 'none';
        }, 1000);
    } else {
        // 如果找不到进度条元素，仍然尝试上传文件
        console.warn('Progress container or list not found in DOM');
        for (const file of files) {
            try {
                await uploadFile(file, userInfo, null);
            } catch (error) {
                console.error(`上传文件 ${file.name} 失败:`, error);
                showMessage(`上传文件 ${file.name} 失败: ${error.message}`, 'error');
            }
        }
        
        // 上传完成后刷新文件列表
        setTimeout(() => {
            loadFileList();
        }, 1000);
    }

    // 清空文件选择
    event.target.value = '';
}

// 上传文件
async function uploadFile(file, userInfo, progressItemId) {
    try {
        // 1. 获取临时token
        if (progressItemId) {
            updateProgress(progressItemId, 0, '获取上传凭证...');
        }
        const tokenResponse = await fetch(`${API_BASE_URL}/material/upload-token`, {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
                'Authorization': `Bearer ${userInfo.token}`
            },
            body: JSON.stringify({
                fileName: file.name,
                fileSize: file.size,
                contentType: file.type || 'application/octet-stream'
            })
        });

        const tokenData = await tokenResponse.json();
        if (handleApiResult(tokenData)) {
            throw new Error('token expired');
        }
        if (tokenData.code !== 0) {
            throw new Error(tokenData.msg || '获取上传凭证失败');
        }

        const { accessKeyId, secretAccessKey, sessionToken, bucket, key, endpoint } = tokenData.data;

        // 2. 配置AWS S3
        if (progressItemId) {
            updateProgress(progressItemId, 5, '初始化上传...');
        }
        AWS.config.update({
            accessKeyId: accessKeyId,
            secretAccessKey: secretAccessKey,
            region: 'us-east-1'
        });

        const s3 = new AWS.S3({
            endpoint: endpoint,
            s3ForcePathStyle: true,
            signatureVersion: 'v4',
            s3DisableBodySigning: true
        });

        // 3. 上传到MinIO
        if (progressItemId) {
            updateProgress(progressItemId, 10, '上传中...');
        }
        await new Promise((resolve, reject) => {
            const params = {
                Bucket: bucket,
                Key: key,
                Body: file,
                ContentType: file.type || 'application/octet-stream'
            };

            s3.upload(params)
                .on('httpUploadProgress', (evt) => {
                    if (progressItemId) {
                        const percent = Math.round((evt.loaded / evt.total) * 80) + 10;
                        updateProgress(progressItemId, percent, '上传中...');
                    }
                })
                .send((err, data) => {
                    if (err) {
                        reject(err);
                    } else {
                        if (progressItemId) {
                            updateProgress(progressItemId, 95, '上传成功');
                        }
                        resolve(data);
                    }
                });
        });

        // 4. 保存文件信息
        if (progressItemId) {
            updateProgress(progressItemId, 95, '保存文件信息...');
        }
        const saveResponse = await fetch(`${API_BASE_URL}/material/save`, {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
                'Authorization': `Bearer ${userInfo.token}`
            },
            body: JSON.stringify({
                fileName: file.name,
                fileSize: file.size,
                contentType: file.type || 'application/octet-stream',
                bucket: bucket,
                key: key,
                url: `${endpoint}/${bucket}/${key}`
            })
        });

        const saveData = await saveResponse.json();
        if (handleApiResult(saveData)) {
            throw new Error('token expired');
        }
        if (saveData.code !== 0) {
            throw new Error(saveData.msg || '保存文件信息失败');
        }

        if (progressItemId) {
            updateProgress(progressItemId, 100, '完成');
        }
    } catch (error) {
        console.error('上传失败:', error);
        if (progressItemId) {
            updateProgress(progressItemId, 0, `失败: ${error.message}`);
        } else {
            showMessage(`上传失败: ${error.message}`, 'error');
        }
        throw error;
    }
}

// 加载文件列表
async function loadFileList() {
    const userInfo = checkLoginStatus();
    if (!userInfo) return;

    const { token } = userInfo;

    try {
        const response = await fetch(`${API_BASE_URL}/material/list`, {
            headers: {
                'Authorization': `Bearer ${token}`
            }
        });

        if (!response.ok) {
            throw new Error(`HTTP error! status: ${response.status}`);
        }

        const data = await response.json();
        if (handleApiResult(data)) return;
        const fileList = document.getElementById('fileList');
        const emptyState = document.getElementById('emptyState');

        if (data.code === 0 && data.data && data.data.materials && data.data.materials.length > 0) {
            fileList.style.display = 'grid';
            emptyState.style.display = 'none';
            fileList.innerHTML = data.data.materials.map(file => {
                // 判断是否为图片或视频类型文件（这些文件有缩略图）
                const isMedia = file.contentType && (file.contentType.startsWith('image/') || file.contentType.startsWith('video/'));
                // 判断是否为视频类型，用于显示右上角三角标识
                const isVideo = file.contentType && file.contentType.startsWith('video/');
                // 如果是媒体文件，则使用封面预览URL作为封面，否则使用文件预览URL
                // 兼容不同的字段名
                const coverSrc = isMedia ? 
                    (file.coverPreviewUrl || file.CoverPreviewURL || file.previewUrl || file.PreviewURL) : 
                    (file.previewUrl || file.PreviewURL);
                
                // 兼容不同的下载链接字段名
                const downloadUrl = file.downloadUrl || file.DownloadURL || file.presignedUrl || file.PresignedURL || '';
                
                // 兼容不同的预览链接字段名
                const previewUrl = file.previewUrl || file.PreviewURL || file.signedGetUrl || file.SignedGetURL || '';
                
                // 把交互移动到封面和右上角菜单中，鼠标悬停时显示左上复选框
                const safeName = (file.fileName || file.FileName || '').replace(/'/g, "\\'");
                return `
                <div class="file-card">
                    <div class="file-cover-wrapper">
                        <input type="checkbox" class="file-checkbox" data-id="${file.id || file.ID}" onchange="updateBatchDeleteButton()">
                        ${coverSrc ? 
                            `<div class="file-cover" style="background-image: url('${coverSrc}')" onclick="previewFile('${previewUrl}', '${file.contentType || file.ContentType}', '${safeName}')"></div>` :
                            `<div class="file-cover" style="background-color:#f5f5f5;display:flex;align-items:center;justify-content:center;" onclick="previewFile('${previewUrl}', '${file.contentType || file.ContentType}', '${safeName}')">📄</div>`
                        }
                        ${isVideo ? `
                            <div class="video-corner" onclick="previewFile('${previewUrl}', '${file.contentType || file.ContentType}', '${safeName}')">
                                <svg viewBox="0 0 24 24" xmlns="http://www.w3.org/2000/svg"><path d="M8 5v14l11-7z"/></svg>
                            </div>
                        ` : ''}
                        <button class="file-menu-btn" onclick="toggleFileMenu(this); event.stopPropagation();">⋯</button>
                        <div class="file-menu-dropdown">
                            <button onclick="previewFile('${previewUrl}', '${file.contentType || file.ContentType}', '${safeName}'); this.closest('.file-menu-dropdown').classList.remove('visible')">预览</button>
                            <button onclick="downloadFile('${downloadUrl}', '${safeName}'); this.closest('.file-menu-dropdown').classList.remove('visible')">下载</button>
                            <button onclick="openRenameModal(${file.id || file.ID}, '${safeName}'); this.closest('.file-menu-dropdown').classList.remove('visible')">重命名</button>
                            <button onclick="openPublishModal(${file.id || file.ID}); this.closest('.file-menu-dropdown').classList.remove('visible')">发布</button>
                            <button onclick="deleteFile(${file.id || file.ID}); this.closest('.file-menu-dropdown').classList.remove('visible')">删除</button>
                        </div>
                    </div>
                    <div class="file-name" title="${file.fileName || file.FileName}">${file.fileName || file.FileName}</div>
                    <div class="file-size">${formatFileSize(file.fileSize || file.FileSize)}</div>
                    <div class="file-date">${file.createdAt || file.CreatedAt}</div>
                </div>
            `}).join('');
        } else {
            fileList.style.display = 'none';
            emptyState.style.display = 'block';
            emptyState.innerHTML = `
                <div class="empty-state-content">
                    <svg width="64" height="64" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5">
                        <path d="M21 15v4a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2v-4"></path>
                        <polyline points="7,10 12,15 17,10"></polyline>
                        <line x1="12" y1="15" x2="12" y2="3"></line>
                    </svg>
                    <h3>暂无素材文件</h3>
                    <p>点击上方上传按钮或拖拽文件到页面，开始上传您的第一个素材文件</p>
                </div>
            `;
        }
    } catch (error) {
        console.error('加载文件列表失败:', error);
        showMessage('加载文件列表失败: ' + error.message, 'error');
    }
}

// 当用户点击搜索按钮或回车时调用
async function onMaterialSearch() {
    const q = (document.getElementById('materialSearchInput') || {}).value || '';
    await searchFiles(q);
}

// 搜索素材并渲染结果（如果 q 为空则加载全部列表）
async function searchFiles(q) {
    const userInfo = checkLoginStatus();
    if (!userInfo) return;

    const { token } = userInfo;

    if (!q || q.trim() === '') {
        // 为空则加载全部
        return await loadFileList();
    }

    try {
        const url = new URL(`${API_BASE_URL}/material/search`, window.location.origin);
        url.searchParams.set('q', q);

        const response = await fetch(url.toString(), {
            headers: {
                'Authorization': `Bearer ${token}`
            }
        });

        if (!response.ok) {
            throw new Error(`HTTP error! status: ${response.status}`);
        }

        const data = await response.json();
        if (handleApiResult(data)) return;

        const fileList = document.getElementById('fileList');
        const emptyState = document.getElementById('emptyState');

        if (data.code === 0 && data.data && data.data.materials && data.data.materials.length > 0) {
            fileList.style.display = 'grid';
            emptyState.style.display = 'none';
            fileList.innerHTML = data.data.materials.map(file => {
                const isMedia = file.contentType && (file.contentType.startsWith('image/') || file.contentType.startsWith('video/'));
                const isVideo = file.contentType && file.contentType.startsWith('video/');
                const coverSrc = isMedia ? (file.coverPreviewUrl || file.coverUrl || file.previewUrl) : (file.previewUrl || '');
                const downloadUrl = file.downloadUrl || file.downloadURL || '';
                const previewUrl = file.previewUrl || file.PreviewURL || file.previewUrl || '';
                const safeName = (file.fileName || file.FileName || '').replace(/'/g, "\\'");
                return `
                <div class="file-card">
                    <div class="file-cover-wrapper">
                        <input type="checkbox" class="file-checkbox" data-id="${file.id || file.ID}" onchange="updateBatchDeleteButton()">
                        ${coverSrc ? `<div class="file-cover" style="background-image: url('${coverSrc}')" onclick="previewFile('${previewUrl}', '${file.contentType || file.ContentType}', '${safeName}')"></div>` : `<div class="file-cover" style="background-color:#f5f5f5;display:flex;align-items:center;justify-content:center;" onclick="previewFile('${previewUrl}', '${file.contentType || file.ContentType}', '${safeName}')">📄</div>`}
                        ${isVideo ? `
                            <div class="video-corner" onclick="previewFile('${previewUrl}', '${file.contentType || file.ContentType}', '${safeName}')">
                                <svg viewBox="0 0 24 24" xmlns="http://www.w3.org/2000/svg"><path d="M8 5v14l11-7z"/></svg>
                            </div>
                        ` : ''}
                        <button class="file-menu-btn" onclick="toggleFileMenu(this); event.stopPropagation();">⋯</button>
                        <div class="file-menu-dropdown">
                            <button onclick="previewFile('${previewUrl}', '${file.contentType || file.ContentType}', '${safeName}'); this.closest('.file-menu-dropdown').classList.remove('visible')">预览</button>
                            <button onclick="downloadFile('${downloadUrl}', '${safeName}'); this.closest('.file-menu-dropdown').classList.remove('visible')">下载</button>
                            <button onclick="openRenameModal(${file.id || file.ID}, '${safeName}'); this.closest('.file-menu-dropdown').classList.remove('visible')">重命名</button>
                            <button onclick="openPublishModal(${file.id || file.ID}); this.closest('.file-menu-dropdown').classList.remove('visible')">发布</button>
                            <button onclick="deleteFile(${file.id || file.ID}); this.closest('.file-menu-dropdown').classList.remove('visible')">删除</button>
                        </div>
                    </div>
                    <div class="file-name" title="${file.fileName || file.FileName}">${file.fileName || file.FileName}</div>
                    <div class="file-size">${formatFileSize(file.fileSize || file.FileSize)}</div>
                    <div class="file-date">${file.createdAt || file.CreatedAt}</div>
                </div>
            `}).join('');
        } else {
            fileList.style.display = 'none';
            emptyState.style.display = 'block';
            emptyState.innerHTML = `
                <div class="empty-state-content">
                    <h3>未找到匹配的素材</h3>
                    <p>请尝试其他关键词或清空搜索查看全部文件</p>
                </div>
            `;
        }
    } catch (error) {
        console.error('搜索素材失败:', error);
        showMessage('搜索素材失败: ' + error.message, 'error');
    }
}

// 根据文件类型预览文件
// - 图片: 在弹窗中显示
// - 视频: 在弹窗中播放
// - 文本/JSON/XML: 在新窗口打开
// - 其他类型: 提示用户下载
function previewFile(url, contentType, fileName) {
    // Handle case where URL is not available
    if (!url) {
        alert('该文件不支持预览');
        return;
    }
    
    if (contentType && contentType.startsWith('image/')) {
        // 图片文件在弹窗中显示
        openImagePreview(url, fileName);
    } else if (contentType && contentType.startsWith('video/')) {
        // 视频文件在弹窗中播放
        openVideoPreview(url, fileName);
    } else if (contentType && (contentType.startsWith('text/') || contentType.includes('json') || contentType.includes('xml'))) {
        // 文本文件在新窗口打开
        window.open(url, '_blank');
    } else {
        // 其他文件类型提示下载
        if (confirm('该文件类型不支持在线预览，是否要下载该文件？')) {
            const link = document.createElement('a');
            link.href = url;
            link.download = fileName;
            link.style.display = 'none';
            document.body.appendChild(link);
            link.click();
            document.body.removeChild(link);
        }
    }
}

// 打开图片预览弹窗
function openImagePreview(url, fileName) {
    // 创建图片预览弹窗
    const modal = document.createElement('div');
    modal.style.cssText = `
        position: fixed;
        top: 0;
        left: 0;
        width: 100%;
        height: 100%;
        background: rgba(0, 0, 0, 0.9);
        display: flex;
        justify-content: center;
        align-items: center;
        z-index: 2000;
    `;
    
    modal.innerHTML = `
        <div style="
            position: relative;
            max-width: 90%;
            max-height: 90%;
        ">
            <div style="
                display: flex;
                justify-content: space-between;
                align-items: center;
                margin-bottom: 15px;
                padding: 0 20px;
            ">
                <h3 style="
                    margin: 0;
                    color: #fff;
                    font-size: 18px;
                ">${fileName}</h3>
                <button onclick="this.closest('div').parentElement.parentElement.remove()" style="
                    background: #ff4d4f;
                    color: white;
                    border: none;
                    border-radius: 50%;
                    width: 30px;
                    height: 30px;
                    cursor: pointer;
                    font-size: 18px;
                    display: flex;
                    align-items: center;
                    justify-content: center;
                ">×</button>
            </div>
            <div style="
                display: flex;
                justify-content: center;
                align-items: center;
                max-height: 80vh;
            ">
                <img src="${url}" style="
                    max-width: 100%;
                    max-height: 70vh;
                    object-fit: contain;
                " alt="${fileName}">
            </div>
        </div>
    `;
    
    document.body.appendChild(modal);
    
    // 点击遮罩层关闭弹窗
    modal.addEventListener('click', (e) => {
        if (e.target === modal) {
            modal.remove();
        }
    });
    
    // 按ESC键关闭弹窗
    const closeOnEscape = (e) => {
        if (e.key === 'Escape') {
            modal.remove();
            document.removeEventListener('keydown', closeOnEscape);
        }
    };
    document.addEventListener('keydown', closeOnEscape);
}

// 打开视频预览弹窗
function openVideoPreview(url, fileName) {
    // 创建视频预览弹窗
    const modal = document.createElement('div');
    modal.style.cssText = `
        position: fixed;
        top: 0;
        left: 0;
        width: 100%;
        height: 100%;
        background: rgba(0, 0, 0, 0.9);
        display: flex;
        justify-content: center;
        align-items: center;
        z-index: 2000;
    `;
    
    modal.innerHTML = `
        <div style="
            position: relative;
            max-width: 90%;
            max-height: 90%;
        ">
            <div style="
                display: flex;
                justify-content: space-between;
                align-items: center;
                margin-bottom: 15px;
                padding: 0 20px;
            ">
                <h3 style="
                    margin: 0;
                    color: #fff;
                    font-size: 18px;
                ">${fileName}</h3>
                <button onclick="this.closest('div').parentElement.parentElement.remove()" style="
                    background: #ff4d4f;
                    color: white;
                    border: none;
                    border-radius: 50%;
                    width: 30px;
                    height: 30px;
                    cursor: pointer;
                    font-size: 18px;
                    display: flex;
                    align-items: center;
                    justify-content: center;
                ">×</button>
            </div>
            <div style="
                display: flex;
                justify-content: center;
                align-items: center;
                max-height: 80vh;
            ">
                <video controls autoplay style="
                    max-width: 100%;
                    max-height: 70vh;
                ">
                    <source src="${url}" type="video/mp4">
                    <source src="${url}" type="video/webm">
                    <source src="${url}" type="video/ogg">
                    您的浏览器不支持视频播放。
                </video>
            </div>
        </div>
    `;
    
    document.body.appendChild(modal);
    
    // 点击遮罩层关闭弹窗
    modal.addEventListener('click', (e) => {
        if (e.target === modal) {
            modal.remove();
        }
    });
    
    // 按ESC键关闭弹窗
    const closeOnEscape = (e) => {
        if (e.key === 'Escape') {
            modal.remove();
            document.removeEventListener('keydown', closeOnEscape);
        }
    };
    document.addEventListener('keydown', closeOnEscape);
}

// 下载文件
function downloadFile(url, fileName) {
    if (!url) {
        showMessage('无法下载此文件', 'error');
        return;
    }
    
    const link = document.createElement('a');
    link.href = url;
    link.download = fileName;
    link.style.display = 'none';
    document.body.appendChild(link);
    link.click();
    document.body.removeChild(link);
}

// 删除单个文件
async function deleteFile(fileId) {
    if (!confirm('确定要删除这个文件吗？')) {
        return;
    }

    const userInfo = checkLoginStatus();
    if (!userInfo) return;

    const { token } = userInfo;

    try {
        // 使用批量删除接口，传入单个文件ID
        const response = await fetch(`${API_BASE_URL}/material/delete`, {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
                'Authorization': `Bearer ${token}`
            },
            body: JSON.stringify({
                ids: [fileId]
            })
        });

        // 检查响应状态
        if (!response.ok) {
            throw new Error(`HTTP error! status: ${response.status}`);
        }

        const responseText = await response.text();
        if (!responseText) {
            throw new Error('服务器返回空响应');
        }

        let data;
        try {
            data = JSON.parse(responseText);
        } catch (parseError) {
            throw new Error('服务器响应不是有效的JSON格式');
        }

        if (data.code === 0) {
            showMessage('文件删除成功', 'success');
            // 重新加载文件列表
            await loadFileList();
        } else {
            throw new Error(data.message || '删除失败');
        }
    } catch (error) {
        console.error('删除文件失败:', error);
        showMessage('删除文件失败: ' + error.message, 'error');
    }
}

// 获取选中的文件ID
function getSelectedFileIds() {
    const checkboxes = document.querySelectorAll('.file-checkbox:checked');
    return Array.from(checkboxes).map(cb => parseInt(cb.dataset.id));
}

// 更新批量删除按钮状态
function updateBatchDeleteButton() {
    const selectedCount = getSelectedFileIds().length;
    const batchDeleteBtn = document.getElementById('batchDeleteBtn');
    
    if (selectedCount > 0) {
        batchDeleteBtn.style.display = 'flex';
        batchDeleteBtn.innerHTML = `
            <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                <path d="M3 6h18"></path>
                <path d="M19 6v14a2 2 0 0 1-2 2H7a2 2 0 0 1-2-2V6m3 0V4a2 2 0 0 1 2-2h4a2 2 0 0 1 2 2v2"></path>
            </svg>
            批量删除 (${selectedCount})
        `;
    } else {
        batchDeleteBtn.style.display = 'none';
    }
}

// 全选/取消全选
function toggleSelectAll() {
    const selectAllCheckbox = document.getElementById('selectAllCheckbox');
    const checkboxes = document.querySelectorAll('.file-checkbox');
    
    checkboxes.forEach(checkbox => {
        checkbox.checked = selectAllCheckbox.checked;
    });
    
    updateBatchDeleteButton();
}

// 批量删除文件
async function batchDeleteFiles() {
    const selectedIds = getSelectedFileIds();
    
    if (selectedIds.length === 0) {
        showMessage('请先选择要删除的文件', 'info');
        return;
    }

    if (!confirm(`确定要删除选中的 ${selectedIds.length} 个文件吗？`)) {
        return;
    }

    const userInfo = checkLoginStatus();
    if (!userInfo) return;

    const { token } = userInfo;

    try {
        // 使用批量删除接口
        const response = await fetch(`${API_BASE_URL}/material/delete`, {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
                'Authorization': `Bearer ${token}`
            },
            body: JSON.stringify({
                ids: selectedIds
            })
        });

        if (!response.ok) {
            throw new Error(`HTTP error! status: ${response.status}`);
        }

        const responseText = await response.text();
        if (!responseText) {
            throw new Error('服务器返回空响应');
        }

        let data;
        try {
            data = JSON.parse(responseText);
        } catch (parseError) {
            throw new Error('服务器响应不是有效的JSON格式');
        }

        if (data.code === 0) {
            showMessage(`成功删除 ${selectedIds.length} 个文件`, 'success');
            
            // 重置全选按钮
            document.getElementById('selectAllCheckbox').checked = false;
            updateBatchDeleteButton();
            
            // 重新加载文件列表
            await loadFileList();
        } else {
            throw new Error(data.message || '删除失败');
        }
    } catch (error) {
        console.error('批量删除文件失败:', error);
        showMessage('批量删除文件失败: ' + error.message, 'error');
    }
}

// 打开发布弹窗
function openPublishModal(materialId) {
    // 移除已存在的弹窗
    const existingModal = document.getElementById('publishModal');
    if (existingModal) {
        existingModal.remove();
    }
    
    // 创建弹窗
    const modal = document.createElement('div');
    modal.id = 'publishModal';
    modal.className = 'modal';
    modal.innerHTML = `
        <div class="modal-content">
            <div class="modal-header">
                <h3>发布素材</h3>
                <span class="close">&times;</span>
            </div>
            <div class="modal-body">
                <form id="publishForm">
                    <input type="hidden" id="publishMaterialId" value="${materialId}">
                    <div class="form-group">
                        <label for="publishDescription">描述（可选）:</label>
                        <textarea id="publishDescription" rows="4" placeholder="请输入描述内容"></textarea>
                    </div>
                    <div class="form-actions">
                        <button type="button" class="btn btn-secondary" id="cancelPublish">取消</button>
                        <button type="submit" class="btn btn-primary">确定</button>
                    </div>
                </form>
            </div>
        </div>
    `;
    
    document.body.appendChild(modal);
    
    // 添加事件监听器
    const closeModal = () => {
        modal.remove();
    };
    
    modal.querySelector('.close').addEventListener('click', closeModal);
    modal.querySelector('#cancelPublish').addEventListener('click', closeModal);
    
    // 点击模态框外部关闭
    modal.addEventListener('click', (e) => {
        if (e.target === modal) {
            closeModal();
        }
    });
    
    // 表单提交事件
    modal.querySelector('#publishForm').addEventListener('submit', async (e) => {
        e.preventDefault();
        
        const materialId = document.getElementById('publishMaterialId').value;
        const description = document.getElementById('publishDescription').value;
        
        try {
            const userInfo = checkLoginStatus();
            if (!userInfo) {
                showMessage('请先登录', 'error');
                return;
            }
            
            const response = await fetch(`${API_BASE_URL}/material/publish`, {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                    'Authorization': `Bearer ${userInfo.token}`
                },
                body: JSON.stringify({
                    materialId: parseInt(materialId),
                    description: description
                })
            });
            
            if (!response.ok) {
                throw new Error(`HTTP error! status: ${response.status}`);
            }
            
            const result = await response.json();
            
            if (result.code !== 0) {
                throw new Error(result.msg || '发布失败');
            }
            
            showMessage('发布成功', 'success');
            closeModal();
        } catch (error) {
            console.error('发布失败:', error);
            showMessage('发布失败: ' + error.message, 'error');
        }
    });
}

// 打开重命名对话框
function openRenameModal(materialId, currentFileName) {
    // 移除已存在的弹窗
    const existingModal = document.getElementById('renameModal');
    if (existingModal) {
        existingModal.remove();
    }
    
    // 创建弹窗
    const modal = document.createElement('div');
    modal.id = 'renameModal';
    modal.className = 'modal';
    modal.innerHTML = `
        <div class="modal-content">
            <div class="modal-header">
                <h3>重命名素材</h3>
                <span class="close">&times;</span>
            </div>
            <div class="modal-body">
                <form id="renameForm">
                    <input type="hidden" id="renameMaterialId" value="${materialId}">
                    <div class="form-group">
                        <label for="newFileName">新文件名:</label>
                        <input type="text" id="newFileName" value="${currentFileName}" placeholder="请输入新文件名" required>
                    </div>
                    <div class="form-actions">
                        <button type="button" class="btn btn-secondary" id="cancelRename">取消</button>
                        <button type="submit" class="btn btn-primary">确定</button>
                    </div>
                </form>
            </div>
        </div>
    `;
    
    document.body.appendChild(modal);
    
    // 添加事件监听器
    const closeModal = () => {
        modal.remove();
    };
    
    modal.querySelector('.close').addEventListener('click', closeModal);
    modal.querySelector('#cancelRename').addEventListener('click', closeModal);
    
    // 点击模态框外部关闭
    modal.addEventListener('click', (e) => {
        if (e.target === modal) {
            closeModal();
        }
    });
    
    // 文件名输入框自动选中
    const fileNameInput = modal.querySelector('#newFileName');
    fileNameInput.focus();
    fileNameInput.select();
    
    // 表单提交事件
    modal.querySelector('#renameForm').addEventListener('submit', async (e) => {
        e.preventDefault();
        
        const materialId = document.getElementById('renameMaterialId').value;
        const newName = document.getElementById('newFileName').value.trim();
        
        if (!newName) {
            showMessage('文件名不能为空', 'error');
            return;
        }
        
        try {
            const userInfo = checkLoginStatus();
            if (!userInfo) {
                showMessage('请先登录', 'error');
                return;
            }
            
            const response = await fetch(`${API_BASE_URL}/material/update-name`, {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                    'Authorization': `Bearer ${userInfo.token}`
                },
                body: JSON.stringify({
                    materialId: parseInt(materialId),
                    newName: newName
                })
            });
            
            if (!response.ok) {
                throw new Error(`HTTP error! status: ${response.status}`);
            }
            
            const result = await response.json();
            
            if (result.code !== 0) {
                throw new Error(result.msg || '重命名失败');
            }
            
            showMessage('重命名成功', 'success');
            closeModal();
            // 重新加载文件列表
            await loadFileList();
        } catch (error) {
            console.error('重命名失败:', error);
            showMessage('重命名失败: ' + error.message, 'error');
        }
    });
}


// 退出登录
async function handleLogout() {
    const token = localStorage.getItem('token');
    if (token) {
        try {
            await fetch(`${API_BASE_URL}/auth/logout`, {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                    'Authorization': `Bearer ${token}`
                }
            });
        } catch (error) {
            console.error('退出登录请求失败:', error);
        }
    }
    localStorage.removeItem('token');
    window.location.href = 'login.html';
}

// 初始化页面
// document.addEventListener('DOMContentLoaded', async function() {
//     // 检查登录状态
//     const userInfo = checkLoginStatus();
//     if (!userInfo) {
//         window.location.href = 'login.html';
//         return;
//     }

//     // 显示用户信息
//     displayUserInfo(userInfo);

//     // 加载用户详细信息并显示真实头像
//     // loadUserDetailsForDisplay(userInfo.token);
    
//     // 加载文件列表
//     await loadFileList();
// });

// 清空素材搜索并恢复全部列表
function clearMaterialSearch() {
    const input = document.getElementById('materialSearchInput');
    if (!input) return;
    if ((input.value || '').trim() === '') {
        // 如果已经为空，直接加载全部列表
        loadFileList();
        input.focus();
        return;
    }
    input.value = '';
    // 触发搜索逻辑（空会回退到全部列表）
    onMaterialSearch();
    input.focus();
}

// Ensure clear button hides when input is cleared programmatically
function updateMaterialClearBtnVisibility() {
    const input = document.getElementById('materialSearchInput');
    const clearBtn = document.getElementById('materialClearBtn');
    if (!input || !clearBtn) return;
    clearBtn.style.display = (input.value || '').trim() ? 'inline-block' : 'none';
}

// 绑定回车触发搜索与清除按钮显示控制
(function() {
    const input = document.getElementById('materialSearchInput');
    const clearBtn = document.getElementById('materialClearBtn');
    if (!input) return;

    input.addEventListener('keydown', function(e) {
        if (e.key === 'Enter') {
            e.preventDefault();
            onMaterialSearch();
        }
    });

    if (clearBtn) {
        // 根据输入内容显示/隐藏清除按钮（初始态）
        clearBtn.style.display = (input.value || '').trim() ? 'inline-block' : 'none';
        input.addEventListener('input', function() {
            clearBtn.style.display = (input.value || '').trim() ? 'inline-block' : 'none';
        });
        // 当通过代码清空时也更新清除按钮显示
        // 覆盖原 clearMaterialSearch 简单行为：在清空后调用此方法
        const origClear = window.clearMaterialSearch;
        window.clearMaterialSearch = function() {
            if (origClear) origClear();
            updateMaterialClearBtnVisibility();
        }
    }
})();

// 切换右上角菜单显示
function toggleFileMenu(btn) {
    const wrapper = btn.closest('.file-cover-wrapper');
    const dropdown = wrapper ? wrapper.querySelector('.file-menu-dropdown') : null;

    // 关闭所有不属于当前 wrapper 的打开菜单
    document.querySelectorAll('.file-menu-dropdown.visible').forEach(d => {
        if (d !== dropdown) d.classList.remove('visible');
    });

    if (!dropdown) return;
    dropdown.classList.toggle('visible');
}

// 点击任意位置关闭打开的菜单（除非点击在菜单内部）
document.addEventListener('click', function(e) {
    const openMenus = document.querySelectorAll('.file-menu-dropdown.visible');
    if (!openMenus || openMenus.length === 0) return;

    // 如果点击在任一 .file-menu-dropdown 或 .file-menu-btn 上，则不关闭
    let node = e.target;
    while (node) {
        if (node.classList && (node.classList.contains('file-menu-dropdown') || node.classList.contains('file-menu-btn'))) {
            return;
        }
        node = node.parentElement;
    }

    openMenus.forEach(d => d.classList.remove('visible'));
});