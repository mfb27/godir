const API_BASE_URL = 'http://192.168.31.67:8081'; // 根据实际端口修改

// 切换标签页
function switchTab(tab) {
    const loginTab = document.querySelector('.tab-container .tab:first-child');
    const registerTab = document.querySelector('.tab-container .tab:last-child');
    const loginForm = document.getElementById('loginForm');
    const registerForm = document.getElementById('registerForm');
    const authTabs = document.getElementById('authTabs');

    if (tab === 'login') {
        loginTab.classList.add('active');
        registerTab.classList.remove('active');
        loginForm.classList.add('active');
        registerForm.classList.remove('active');
    } else {
        loginTab.classList.remove('active');
        registerTab.classList.add('active');
        loginForm.classList.remove('active');
        registerForm.classList.add('active');
    }
}

// 检查是否已登录
function checkLoginStatus() {
    const token = localStorage.getItem('token');
    if (token) {
        try {
            // 解析token获取用户信息
            const payload = JSON.parse(atob(token.split('.')[1]));
            const userId = payload.userId;
            const username = payload.username || '用户';
            
            // 显示用户信息
            showUserInfo(userId, username);
            
            // 如果已经登录，则直接跳转到主页（可选）
            // window.location.href = 'index.html';
        } catch (e) {
            console.error('解析token失败:', e);
            localStorage.removeItem('token');
        }
        return;
    }
    // 未登录时显示标签页
    document.getElementById('authTabs').style.display = 'block';
}

// 显示消息
function showMessage(text, type = 'error') {
    const messageEl = document.getElementById('message');
    messageEl.textContent = text;
    messageEl.className = `message ${type}`;
    messageEl.style.display = 'block';
    setTimeout(() => {
        messageEl.style.display = 'none';
    }, 3000);
}

// 显示用户信息
function showUserInfo(userId, userName) {
    document.getElementById('authTabs').style.display = 'none';
    document.getElementById('loginForm').style.display = 'none';
    document.getElementById('registerForm').style.display = 'none';
    document.getElementById('userInfo').style.display = 'block';
    document.getElementById('userId').textContent = userId;
    document.getElementById('userName').textContent = userName;
}

// 隐藏用户信息
function hideUserInfo() {
    document.getElementById('authTabs').style.display = 'block';
    document.getElementById('loginForm').style.display = 'block';
    document.getElementById('registerForm').style.display = 'none';
    document.getElementById('userInfo').style.display = 'none';
    
    // 默认激活登录表单
    switchTab('login');
}

// 处理登录
async function handleLogin() {
    const username = document.getElementById('username').value;
    const password = document.getElementById('password').value;

    if (!username || !password) {
        showMessage('请输入用户名和密码', 'error');
        return;
    }

    const btn = event.target;
    btn.disabled = true;
    btn.textContent = '登录中...';

    try {
        const response = await fetch(`${API_BASE_URL}/auth/login`, {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify({
                username: username,
                password: password
            })
        });

        const data = await response.json();

        if (data.code === 0) {
            // 保存token
            localStorage.setItem('token', data.data.token);
            showMessage('登录成功！', 'success');
            // 延迟跳转，让用户看到成功消息
            setTimeout(() => {
                window.location.href = 'index.html';
            }, 1500);
        } else {
            showMessage(data.msg || '登录失败', 'error');
        }
    } catch (error) {
        showMessage('网络错误: ' + error.message, 'error');
    } finally {
        btn.disabled = false;
        btn.textContent = '登录';
    }
}

// 处理退出
async function handleLogout() {
    const token = localStorage.getItem('token');
    
    if (!token) {
        localStorage.removeItem('token');
        hideUserInfo();
        showMessage('已退出登录', 'success');
        return;
    }

    const btn = event.target;
    btn.disabled = true;
    btn.textContent = '退出中...';

    try {
        const response = await fetch(`${API_BASE_URL}/auth/logout`, {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
                'Authorization': `Bearer ${token}`
            }
        });

        const data = await response.json();

        // 无论服务端返回什么，都清除本地token
        localStorage.removeItem('token');
        hideUserInfo();
        document.getElementById('username').value = '';
        document.getElementById('password').value = '';
        showMessage(data.msg || '退出成功', 'success');
    } catch (error) {
        // 即使请求失败，也清除本地token
        localStorage.removeItem('token');
        hideUserInfo();
        showMessage('已退出登录', 'success');
    } finally {
        btn.disabled = false;
        btn.textContent = '退出登录';
    }
}

// 处理注册
async function handleRegister() {
    const username = document.getElementById('regUsername').value;
    const password = document.getElementById('regPassword').value;

    if (!username || !password) {
        showMessage('请填写完整信息', 'error');
        return;
    }

    if (password.length < 6) {
        showMessage('密码长度至少6位', 'error');
        return;
    }

    const btn = event.target;
    btn.disabled = true;
    btn.textContent = '注册中...';

    try {
        const response = await fetch(`${API_BASE_URL}/auth/register`, {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify({
                username: username,
                password: password
            })
        });

        const data = await response.json();

        if (data.code === 0) {
            showMessage('注册成功！请登录', 'success');
            // 清空注册表单
            document.getElementById('regUsername').value = '';
            document.getElementById('regPassword').value = '';
            // 切换到登录页面并自动填充用户名
            setTimeout(() => {
                switchTab('login');
                document.getElementById('username').value = username;
                document.getElementById('password').focus();
            }, 1500);
        } else {
            showMessage(data.msg || '注册失败', 'error');
        }
    } catch (error) {
        showMessage('网络错误: ' + error.message, 'error');
    } finally {
        btn.disabled = false;
        btn.textContent = '注册';
    }
}

// 支持回车键登录
document.getElementById('password').addEventListener('keypress', function(e) {
    if (e.key === 'Enter') {
        handleLogin();
    }
});

// 支持回车键注册
document.getElementById('regPassword').addEventListener('keypress', function(e) {
    if (e.key === 'Enter') {
        handleRegister();
    }
});

// 页面加载时检查登录状态
document.addEventListener('DOMContentLoaded', function() {
    checkLoginStatus();
});