// 获取用户名首字母作为头像
function getAvatarText(username) {
    if (!username) return 'U';
    return username.charAt(0).toUpperCase();
}

/**
 * 检查用户登录状态（合并版）
 * @returns {Object|null} 用户信息对象或null
 */
function checkLoginStatusMerged() {
    const token = localStorage.getItem('token');
    if (!token) {
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
        return null;
    }
}

/**
 * 加载用户详细信息并显示
 * @param {string} token - 用户认证token
 */
async function loadUserDetailsForDisplay(token) {
    try {
        console.log('加载用户详细信息并显示');
        const response = await fetch(`${API_BASE_URL}/user/profile`, {
            headers: {
                'Authorization': `Bearer ${token}`
            }
        });
        
        if (!response.ok) {
            throw new Error(`HTTP error! status: ${response.status}`);
        }
        
        const data = await response.json();
        if (handleApiResult(data)) return;
        if (data.code === 0 && data.data) {
            displayUserDetails(data.data);
        }
    } catch (error) {
        console.error('加载用户信息失败:', error);
        // 如果加载详细信息失败，回退到基本显示
        const userInfo = checkLoginStatusMerged();
        if (userInfo) {
            displayUserInfo(userInfo);
        }
    }
}

/**
 * 显示个人中心页面
 */
function showProfilePage() {
    // 隐藏其他页面，显示个人中心页面
    document.getElementById('homePage').classList.add('hidden');
    document.getElementById('materialPage').classList.add('hidden');
    document.getElementById('profilePage').classList.remove('hidden');
    
    // 更新左侧导航栏活动状态
    document.querySelectorAll('.sidebar-menu-item').forEach(item => {
        item.classList.remove('active');
    });
    
    // 找到个人中心链接并设置为活动状态
    const profileLinks = document.querySelectorAll('a[href="#"]');
    profileLinks.forEach(link => {
        if (link.textContent.includes('个人中心')) {
            link.classList.add('active');
        }
    });
    
    // 加载用户信息
    loadUserProfile();
}

/**
 * 显示主页
 */
function showHomePage() {
    console.log('点击主页，调用showHomePage函数');
    
    // 隐藏其他页面，显示主页
    const materialPage = document.getElementById('materialPage');
    const profilePage = document.getElementById('profilePage');
    if (materialPage) {
        materialPage.classList.add('hidden');
    }
    if (profilePage) {
        profilePage.classList.add('hidden');
    }
    const homePage = document.getElementById('homePage');
    if (homePage) {
        homePage.classList.remove('hidden');
    }

    // 更新导航按钮状态
    const toggleToHome = document.getElementById('toggleToHome');
    const toggleToMaterial = document.getElementById('toggleToMaterial');
    if (toggleToHome) toggleToHome.style.display = 'none';
    if (toggleToMaterial) toggleToMaterial.style.display = 'inline-block';

    // 更新左侧导航栏活动状态
    document.querySelectorAll('.sidebar-menu-item').forEach(item => {
        item.classList.remove('active');
    });
    const activeLink = document.querySelector('.sidebar-menu-item[href="#"]');
    if (activeLink) activeLink.classList.add('active');

    // 检查登录状态
    const userInfo = checkLoginStatusMerged();
    const mainContainer = document.getElementById('mainContainer');
    const unauthorizedContainer = document.getElementById('unauthorizedContainer');
    
    if (!userInfo) {
        console.log('用户未登录');
        if (mainContainer) mainContainer.style.display = 'none';
        if (unauthorizedContainer) unauthorizedContainer.style.display = 'block';
    } else {
        console.log('用户已登录，准备加载发布列表');
        if (mainContainer) mainContainer.style.display = 'block';
        if (unauthorizedContainer) unauthorizedContainer.style.display = 'none';

        // 加载用户详细信息并显示真实头像
        loadUserDetailsForDisplay(userInfo.token);

        // 加载发布的素材列表
        loadPublishedMaterials();
    }
}

/**
 * 加载用户个人信息
 */
async function loadUserProfile() {
    try {
        const userInfo = checkLoginStatusMerged();
        if (!userInfo) {
            showMessage('请先登录', 'error');
            return;
        }
        
        const response = await fetch(`${API_BASE_URL}/user/profile`, {
            headers: {
                'Authorization': `Bearer ${userInfo.token}`
            }
        });
        
        if (!response.ok) {
            throw new Error(`HTTP error! status: ${response.status}`);
        }
        
        const data = await response.json();
        if (data.code === 0 && data.data) {
            // 填充表单数据
            document.getElementById('nickname').value = data.data.nickname || '';
            document.getElementById('gender').value = data.data.gender || 0;
            
            // 显示头像
            const avatarImg = document.getElementById('avatarImage');
            const avatarPlaceholder = document.getElementById('avatarPlaceholder');
            if (data.data.avatar) {
                avatarImg.src = data.data.avatar;
                avatarImg.style.display = 'block';
                avatarPlaceholder.style.display = 'none';
            } else {
                avatarImg.style.display = 'none';
                avatarPlaceholder.style.display = 'block';
            }
        } else if (data.code === 10000001) {
            // 跳转到登录页面
            window.location.href = 'login.html';
        }else{
            throw new Error(data.msg || '获取用户信息失败');
        }
    } catch (error) {
        console.error('加载用户信息失败:', error);
        showMessage('加载用户信息失败: ' + error.message, 'error');
    }
}

// 添加头像预览功能
document.addEventListener('DOMContentLoaded', function() {
    const avatarInput = document.getElementById('avatar');
    const avatarPreview = document.getElementById('avatarPreview');
    const avatarImage = document.getElementById('avatarImage');
    const avatarPlaceholder = document.getElementById('avatarPlaceholder');
    
    if (avatarInput) {
        avatarInput.addEventListener('change', function(e) {
            const file = e.target.files[0];
            if (file) {
                const reader = new FileReader();
                reader.onload = function(e) {
                    avatarImage.src = e.target.result;
                    avatarImage.style.display = 'block';
                    avatarPlaceholder.style.display = 'none';
                };
                reader.readAsDataURL(file);
            }
        });
        
        avatarPreview.addEventListener('click', function() {
            avatarInput.click();
        });
    }
    
    // 表单提交事件
    const profileForm = document.getElementById('profileForm');
    if (profileForm) {
        profileForm.addEventListener('submit', async function(e) {
            e.preventDefault();
            await saveUserProfile();
        });
    }
});

/**
 * 保存用户个人信息
 */
async function saveUserProfile() {
    try {
        const userInfo = checkLoginStatusMerged();
        if (!userInfo) {
            showMessage('请先登录', 'error');
            return;
        }
        
        const nickname = document.getElementById('nickname').value;
        const gender = parseInt(document.getElementById('gender').value);
        const avatarFile = document.getElementById('avatar').files[0];
        
        // 如果选择了头像文件，先上传头像
        let avatarUrl = null;
        if (avatarFile) {
            avatarUrl = await uploadAvatar(avatarFile, userInfo.token);
        }
        
        // 准备用户数据
        const userData = {
            nickname: nickname,
            gender: gender
        };
        
        // 如果上传了新头像，也包含在用户数据中
        if (avatarUrl) {
            userData.avatar = avatarUrl;
        }
        
        // 更新用户信息
        const response = await fetch(`${API_BASE_URL}/user/profile`, {
            method: 'PUT',
            headers: {
                'Content-Type': 'application/json',
                'Authorization': `Bearer ${userInfo.token}`
            },
            body: JSON.stringify(userData)
        });
        
        if (!response.ok) {
            throw new Error(`HTTP error! status: ${response.status}`);
        }
        
        const data = await response.json();
        if (data.code === 0) {
            showMessage('个人信息保存成功', 'success');
        } else {
            throw new Error(data.msg || '保存失败');
        }
    } catch (error) {
        console.error('保存用户信息失败:', error);
        showMessage('保存用户信息失败: ' + error.message, 'error');
    }
}

/**
 * 上传头像文件
 */
async function uploadAvatar(file, token) {
    try {
        const formData = new FormData();
        formData.append('file', file);
        
        const response = await fetch(`${API_BASE_URL}/user/avatar`, {
            method: 'POST',
            headers: {
                'Authorization': `Bearer ${token}`
            },
            body: formData
        });
        
        if (!response.ok) {
            throw new Error(`HTTP error! status: ${response.status}`);
        }
        
        const data = await response.json();
        if (data.code === 0 && data.data) {
            return data.data.url;
        } else {
            throw new Error(data.msg || '头像上传失败');
        }
    } catch (error) {
        console.error('头像上传失败:', error);
        throw error;
    }
}

// 处理退出登录
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

    // 清除token并跳转到登录页
    localStorage.removeItem('token');
    window.location.href = 'login.html';
}

// 设置用户头像点击事件来切换下拉菜单
function setupUserDropdown() {
    console.log('设置用户头像点击事件');
    const avatarContainer = document.getElementById('userAvatarContainer');
    if (!avatarContainer) return;

    const dropdownMenu = avatarContainer.querySelector('.user-dropdown');
    if (!dropdownMenu) return;

    // 统一使用容器的 .show 类来控制显示（与 CSS 配合）
    let hoverTimeout = null;
    const clearHide = () => {
        if (hoverTimeout) {
            clearTimeout(hoverTimeout);
            hoverTimeout = null;
        }
    };
    const show = () => {
        console.log('show');
        clearHide();
        avatarContainer.classList.add('show');
    };
    const hide = () => {
        clearHide();
        avatarContainer.classList.remove('show');
    };
    const scheduleHide = (delay = 300) => {
        clearHide();
        hoverTimeout = setTimeout(() => {
            avatarContainer.classList.remove('show');
            hoverTimeout = null;
        }, delay);
    };

    const isTouch = ('ontouchstart' in window) || (navigator.maxTouchPoints && navigator.maxTouchPoints > 0);
    if (isTouch) {
        // 触屏：点击切换显示状态
        avatarContainer.addEventListener('click', function(e) {
            e.stopPropagation();
            avatarContainer.classList.toggle('show');
        });

        // 点击页面其他地方时隐藏下拉菜单
        document.addEventListener('click', function() {
            if (avatarContainer.classList.contains('show')) {
                avatarContainer.classList.remove('show');
            }
        });

        // 防止点击下拉菜单内部时关闭
        dropdownMenu.addEventListener('click', function(e) {
            e.stopPropagation();
        });
    } else {
        console.log('非触屏设备，设置鼠标悬停事件');
        // 非触屏：通过鼠标进入/移出控制显示，并加入延迟以避免闪烁
        // 使容器可聚焦以支持键盘操作
        try { avatarContainer.setAttribute('tabindex', '0'); } catch (e) {}

        avatarContainer.addEventListener('mouseenter', function() {
            console.log('鼠标进入头像容器');
            show();
        });
        avatarContainer.addEventListener('mouseleave', function() {
            // 鼠标离开后延迟隐藏，给用户移动到下拉菜单的时间
            console.log('鼠标离开头像容器');
            scheduleHide(300);
        });

        // 当鼠标进入下拉菜单时取消隐藏，离开时延迟隐藏，防止在头像和菜单之间移动时闪烁
        dropdownMenu.addEventListener('mouseenter', function() {
            console.log('鼠标进入下拉菜单');
            clearHide();
            show();
        });
        dropdownMenu.addEventListener('mouseleave', function() {
            console.log('鼠标离开下拉菜单');
            scheduleHide(300);
        });

        // 点击页面任意其他地方立即隐藏下拉
        document.addEventListener('click', function(e) {
            if (!avatarContainer.contains(e.target)) {
                hide();
            }
        });

        // 键盘支持：回车或空格切换显示
        avatarContainer.addEventListener('keydown', function(e) {
            if (e.key === 'Enter' || e.key === ' ' || e.key === 'Spacebar') {
                e.preventDefault();
                if (avatarContainer.classList.contains('show')) hide(); else show();
            }
        });
    }
}

/**
 * 加载发布的素材列表
 */
async function loadPublishedMaterials() {
    console.log('开始加载发布列表');
    const publishedListContainer = document.getElementById('publishedMaterialsList');
    if (!publishedListContainer) {
        console.log('找不到publishedMaterialsList容器');
        return;
    }

    // 显示加载状态
    publishedListContainer.innerHTML = '<div class="loading-placeholder">加载中...</div>';

    try {
        console.log('发送请求到/public/published接口');
        const token = localStorage.getItem('token');
        const headers = {};
        if (token) headers['Authorization'] = `Bearer ${token}`;
        const response = await fetch(`${API_BASE_URL}/public/published`, { headers });
        console.log('收到响应:', response);
        
        if (!response.ok) {
            throw new Error(`HTTP error! status: ${response.status}`);
        }

        const data = await response.json();
        console.log('解析后的数据:', data);

        if (data.code === 0 && data.data && data.data.list) {
            renderPublishedMaterials(data.data.list);
        } else {
            publishedListContainer.innerHTML = '<div class="empty-state"><div class="empty-text">暂无发布内容</div></div>';
        }
    } catch (error) {
        console.error('加载发布列表失败:', error);
        publishedListContainer.innerHTML = '<div class="empty-state"><div class="empty-text">加载失败，请稍后重试</div></div>';
    }
}

/**
 * 渲染发布的素材列表
 * @param {Array} publishedMaterials - 发布的素材列表
 */
function renderPublishedMaterials(publishedMaterials) {
    const container = document.getElementById('publishedMaterialsList');
    if (!container) return;

    if (!publishedMaterials || publishedMaterials.length === 0) {
        container.innerHTML = '<div class="empty-state"><div class="empty-text">暂无发布内容</div></div>';
        return;
    }

    container.innerHTML = publishedMaterials.map(pm => {
        // 获取用户名（优先使用昵称，然后是用户名）
        const userName = (pm.user && (pm.user.nickname || pm.user.username)) || '未知用户';
        const avatarText = userName ? userName.charAt(0).toUpperCase() : 'U';
        
        // 判断是否为图片或视频类型文件
        const isImage = pm.material && pm.material.contentType && pm.material.contentType.startsWith('image/');
        const isVideo = pm.material && pm.material.contentType && pm.material.contentType.startsWith('video/');
        const isMedia = isImage || isVideo;
        
        // 封面图片URL - 优先使用预签名封面URL，然后是预签名URL，最后是普通URL
        let coverSrc = null;
        if (isMedia) {
            coverSrc = pm.material.coverPreviewUrl || 
                      pm.material.CoverPreviewURL || 
                      pm.material.coverUrl || 
                      pm.material.CoverURL ||
                      pm.material.previewUrl || 
                      pm.material.PreviewURL ||
                      pm.material.url || 
                      pm.material.URL;
        }
        
        // 原图URL用于预览
        const previewSrc = pm.material.previewUrl || 
                          pm.material.PreviewURL || 
                          pm.material.url || 
                          pm.material.URL;
        
        // 头像HTML - 如果有头像URL则显示真实头像，否则显示文字头像
        let avatarHtml = '';
        if (pm.user && pm.user.avatar) {
            // 显示真实头像
            avatarHtml = `<img src="${pm.user.avatar}" alt="${userName}" class="user-avatar-img" style="width: 32px; height: 32px; border-radius: 50%; object-fit: cover;">`;
        } else {
            // 显示文字头像
            avatarHtml = `<div class="user-avatar-small">${avatarText}</div>`;
        }
        

        // 绑定点赞按钮事件
        setTimeout(() => {
            document.querySelectorAll('.like-btn').forEach(btn => {
                btn.onclick = async function(e) {
                    e.stopPropagation();
                    const publishId = this.getAttribute('data-publish-id');
                    const liked = this.style.color === 'rgb(231, 76, 60)' || this.style.color === '#e74c3c';
                    const token = localStorage.getItem('token');
                    if (!token) {
                        showToast('请先登录后点赞', 'error');
                        return;
                    }
                    try {
                        const url = liked ? `${API_BASE_URL}/material/published/unlike` : `${API_BASE_URL}/material/published/like`;
                        console.log('点赞请求', {url, publishId, liked, token});
                        const resp = await fetch(url, {
                            method: 'POST',
                            headers: {
                                'Content-Type': 'application/json',
                                'Authorization': `Bearer ${token}`
                            },
                            body: JSON.stringify({ publishId: Number(publishId) })
                        });
                        const data = await resp.json();
                        console.log('点赞响应', data);
                        if (data.code === 0 && data.data) {
                            // 更新按钮颜色和计数
                            this.querySelector('.like-count').textContent = data.data.likesCount;
                            this.style.color = data.data.liked ? '#e74c3c' : '#888';
                        } else {
                            showToast(data.message || '操作失败', 'error');
                        }
                    } catch (err) {
                        showToast('网络错误', 'error');
                    }
                };
            });
        }, 100);

        return `
            <div class="published-material-card">
                <div class="user-info-header">
                    ${avatarHtml}
                    <div class="user-name">${escapeHtml(userName)}</div>
                </div>
                ${pm.description ? `<div class="publish-description-content">${escapeHtml(pm.description)}</div>` : ''}
                <div class="material-preview">
                    ${coverSrc ? 
                        `<div class="material-cover" style="background-image: url('${coverSrc}')" onclick="previewMedia('${previewSrc}', '${pm.material.contentType}')">
                            ${isVideo ? `<div class="video-corner" onclick="previewMedia('${previewSrc}', '${pm.material.contentType}')"><svg viewBox="0 0 24 24" xmlns="http://www.w3.org/2000/svg"><path d="M8 5v14l11-7z"/></svg></div>` : ''}
                        </div>` :
                        (isMedia ?
                            `<div class="material-cover" style="background-image: url('${pm.material.url || pm.material.URL}')" onclick="previewMedia('${previewSrc}', '${pm.material.contentType}')">
                                ${isVideo ? `<div class="video-corner" onclick="previewMedia('${previewSrc}', '${pm.material.contentType}')"><svg viewBox="0 0 24 24" xmlns="http://www.w3.org/2000/svg"><path d="M8 5v14l11-7z"/></svg></div>` : ''}
                            </div>` :
                            `<div class="material-icon">${getFileTypeIcon(pm.material.contentType)}</div>`)
                    }
                </div>
                <div class="material-meta" style="margin-top:8px;">
                    <button class="like-btn" data-publish-id="${pm.id}" style="background:none;border:none;cursor:pointer;color:${pm.liked ? '#e74c3c' : '#888'};font-size:18px;vertical-align:middle;">
                        <span class="like-icon">&#x2764;</span>
                        <span class="like-count">${pm.likesCount || 0}</span>
                    </button>
                </div>
            </div>
        `;
        
    }).join('');
}

// HTML转义函数，防止XSS攻击
function escapeHtml(text) {
    if (!text) return '';
    const map = {
        '&': '&amp;',
        '<': '&lt;',
        '>': '&gt;',
        '"': '&quot;',
        "'": '&#039;'
    };
    
    return text.toString().replace(/[&<>"']/g, function(m) { return map[m]; });
}

// 获取文件类型图标
function getFileTypeIcon(contentType) {
    if (!contentType) return '📄';
    
    if (contentType.startsWith('image/')) {
        return '🖼️';
    } else if (contentType.startsWith('video/')) {
        return '🎬';
    } else if (contentType.startsWith('audio/')) {
        return '🎵';
    } else if (contentType.includes('pdf')) {
        return '📑';
    } else {
        return '📄';
    }
}

/**
 * 预览媒体文件（图片或视频）
 * @param {string} mediaUrl - 媒体文件URL
 * @param {string} contentType - 文件类型
 */
function previewMedia(mediaUrl, contentType) {
    // 创建遮罩层
    const overlay = document.createElement('div');
    overlay.style.position = 'fixed';
    overlay.style.top = '0';
    overlay.style.left = '0';
    overlay.style.width = '100%';
    overlay.style.height = '100%';
    overlay.style.backgroundColor = 'rgba(0, 0, 0, 0.9)';
    overlay.style.display = 'flex';
    overlay.style.justifyContent = 'center';
    overlay.style.alignItems = 'center';
    overlay.style.zIndex = '9999';
    overlay.style.cursor = 'pointer';
    
    // 点击遮罩层或按ESC键关闭预览
    overlay.onclick = function() {
        // 暂停视频播放
        if (mediaElement && mediaElement.tagName === 'VIDEO') {
            mediaElement.pause();
        }
        document.body.removeChild(overlay);
    };
    
    // 监听键盘事件，按ESC键关闭预览
    const closeOnEsc = function(event) {
        if (event.key === 'Escape') {
            // 暂停视频播放
            if (mediaElement && mediaElement.tagName === 'VIDEO') {
                mediaElement.pause();
            }
            document.body.removeChild(overlay);
            document.removeEventListener('keydown', closeOnEsc);
        }
    };
    document.addEventListener('keydown', closeOnEsc);

    // 创建媒体元素（图片或视频）
    let mediaElement;
    if (contentType.startsWith('image/')) {
        mediaElement = document.createElement('img');
        mediaElement.src = mediaUrl;
        mediaElement.style.maxWidth = '90%';
        mediaElement.style.maxHeight = '90%';
        mediaElement.style.objectFit = 'contain';
        mediaElement.style.borderRadius = '4px';
    } else if (contentType.startsWith('video/')) {
        mediaElement = document.createElement('video');
        mediaElement.src = mediaUrl;
        mediaElement.controls = true;
        mediaElement.autoplay = true;
        mediaElement.style.maxWidth = '90%';
        mediaElement.style.maxHeight = '90%';
        mediaElement.style.objectFit = 'contain';
        mediaElement.style.borderRadius = '4px';
    } else {
        // 对于其他类型，默认显示图片
        mediaElement = document.createElement('img');
        mediaElement.src = mediaUrl;
        mediaElement.style.maxWidth = '90%';
        mediaElement.style.maxHeight = '90%';
        mediaElement.style.objectFit = 'contain';
        mediaElement.style.borderRadius = '4px';
    }
    
    mediaElement.onclick = function(e) {
        e.stopPropagation();
    };

    // 添加加载动画
    const loading = document.createElement('div');
    loading.textContent = '加载中...';
    loading.style.color = 'white';
    loading.style.fontSize = '18px';
    
    // 媒体文件加载完成后移除加载动画
    if (mediaElement.tagName === 'IMG') {
        mediaElement.onload = function() {
            if (overlay.contains(loading)) {
                overlay.removeChild(loading);
            }
            overlay.appendChild(mediaElement);
        };
        
        mediaElement.onerror = function() {
            if (overlay.contains(loading)) {
                loading.textContent = '媒体加载失败';
            }
        };
    } else if (mediaElement.tagName === 'VIDEO') {
        mediaElement.onloadeddata = function() {
            if (overlay.contains(loading)) {
                overlay.removeChild(loading);
            }
            overlay.appendChild(mediaElement);
        };
        
        mediaElement.onerror = function() {
            if (overlay.contains(loading)) {
                loading.textContent = '视频加载失败';
            }
        };
    }

    // 添加元素到遮罩层并显示
    overlay.appendChild(loading);
    document.body.appendChild(overlay);
}

// 页面加载时初始化
document.addEventListener('DOMContentLoaded', function() {
    showHomePage();
    
    // 启用用户下拉菜单交互逻辑
    setupUserDropdown();
});