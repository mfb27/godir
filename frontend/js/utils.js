const API_BASE_URL = 'http://192.168.31.67:8080';

// 简单的全局消息提示
function showToast(msg, type = 'info') {
    let toast = document.createElement('div');
    toast.className = 'toast toast-' + type;
    toast.textContent = msg;
    document.body.appendChild(toast);
    setTimeout(() => {
        toast.style.opacity = '0';
        setTimeout(() => document.body.removeChild(toast), 500);
    }, 2000);
}

/**
 * 检查用户登录状态
 * @returns {object|null} 用户信息对象或null（未登录）
 */
function checkLoginStatus() {
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
 * 检查用户登录状态并重定向到登录页面（如果未登录）
 * @returns {object|null} 用户信息对象或null（未登录并重定向）
 */
function checkLoginAndRedirect() {
    const userInfo = checkLoginStatus();
    if (!userInfo) {
        window.location.href = 'login.html';
        return null;
    }
    return userInfo;
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

/**
 * 安全地获取用户信息，包括完整的错误处理
 * @returns {Promise<object|null>} 用户信息或null
 */
async function getUserInfoSafe() {
    const userInfoStr = localStorage.getItem('userInfo');
    if (!userInfoStr) {
        return null;
    }

    try {
        const userInfo = JSON.parse(userInfoStr);
        if (!userInfo.token) {
            return null;
        }
        return userInfo;
    } catch (e) {
        console.error('用户信息解析失败:', e);
        return null;
    }
}

/**
 * 加载用户详细信息，包括头像
 * @param {string} token - 用户认证token
 * @returns {Promise<object|null>} 用户详细信息或null
 */
async function loadUserDetails(token) {
    try {
        const response = await fetch(`${API_BASE_URL}/user/profile`, {
            headers: {
                'Authorization': `Bearer ${token}`
            }
        });
        
        if (!response.ok) {
            throw new Error(`HTTP error! status: ${response.status}`);
        }
        
        const data = await response.json();
        if (data.code === 0 && data.data) {
            return data.data;
        }
    } catch (error) {
        console.error('加载用户信息失败:', error);
    }
    return null;
}

/**
 * 显示用户详细信息，包括真实头像
 * @param {object} userDetails - 用户详细信息
 */
function displayUserDetails(userDetails) {
    if (!userDetails) return;
    
    const avatarContainer = document.getElementById('userAvatarContainer');
    if (!avatarContainer) return;
    
    const avatar = document.getElementById('userAvatar');
    const dropdownAvatar = document.getElementById('dropdownAvatar');
    const dropdownName = document.getElementById('dropdownName');
    const dropdownId = document.getElementById('dropdownId');
    
    // 显示用户昵称，如果没有则显示用户名
    const displayName = userDetails.nickname || userDetails.username;
    
    // 显示用户ID
    if (dropdownId) dropdownId.textContent = userDetails.id;
    
    // 显示用户名
    if (dropdownName) dropdownName.textContent = displayName;
    
    // 如果有头像URL，则显示真实头像
    if (userDetails.avatar && avatar) {
        // 清除文本内容
        avatar.textContent = '';
        dropdownAvatar.textContent = '';
        
        // 创建图片元素
        const avatarImg = document.createElement('img');
        avatarImg.src = userDetails.avatar;
        avatarImg.alt = '用户头像';
        avatarImg.style.cssText = `
            width: 100%;
            height: 100%;
            object-fit: cover;
            border-radius: 50%;
        `;
        
        // 创建下拉菜单中的头像图片元素
        const dropdownAvatarImg = avatarImg.cloneNode(true);
        
        // 添加到容器中
        avatar.appendChild(avatarImg);
        dropdownAvatar.appendChild(dropdownAvatarImg);
    } else {
        // 如果没有头像URL，则显示文字头像
        const avatarText = displayName.charAt(0).toUpperCase();
        if (avatar) avatar.textContent = avatarText;
        if (dropdownAvatar) dropdownAvatar.textContent = avatarText;
    }
    
    avatarContainer.style.display = 'block';

}

/**
 * 全局统一处理接口返回值
 * 如果发现 token 过期（code=10000001），清理本地登录状态并跳转到登录页
 * @param {object} data 接口返回的 JSON 对象
 * @returns {boolean} 如果已处理（例如跳转），返回 true，调用方应停止后续处理
 */
function handleApiResult(data) {
    if (!data || typeof data.code === 'undefined') return false;
    if (data.code === 10000001) {
        // token expired
        try {
            localStorage.removeItem('token');
            localStorage.removeItem('userInfo');
        } catch (e) {}
        // 如果当前不是在 login 页面，则跳转
        if (!window.location.href.endsWith('login.html')) {
            window.location.href = 'login.html';
        }
        return true;
    }
    return false;
}