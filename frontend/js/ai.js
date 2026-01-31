let cozeWebSDK = null;
let currentAppId = null;
let aiAppsList = [];

async function loadAiApps() {
    try {
        const response = await fetch(`${API_BASE_URL}/ai/apps`);
        const data = await response.json();
        if (data.code === 0 && data.data) {
            aiAppsList = data.data.apps || [];
            renderAiApps();
        } else {
            console.error('加载AI应用列表失败:', data.msg || '未知错误');
        }
    } catch (error) {
        console.error('加载AI应用列表失败:', error);
    }
}

function renderAiApps() {
    const grid = document.getElementById('aiAppGrid');
    if (!grid) return;
    
    grid.innerHTML = '';
    
    if (aiAppsList.length === 0) {
        grid.innerHTML = '<div style="text-align: center; color: #999; padding: 40px;">暂无AI应用</div>';
        return;
    }
    
    aiAppsList.forEach(app => {
        const card = document.createElement('div');
        card.className = 'ai-app-card';
        card.onclick = () => openAiApp(app);
        
        const icon = app.icon || '🤖';
        const iconEl = document.createElement('div');
        iconEl.className = 'ai-app-icon';
        iconEl.textContent = icon;
        
        const nameEl = document.createElement('div');
        nameEl.className = 'ai-app-name';
        nameEl.textContent = app.name;
        
        const descEl = document.createElement('div');
        descEl.className = 'ai-app-desc';
        descEl.textContent = app.desc || '';
        
        card.appendChild(iconEl);
        card.appendChild(nameEl);
        card.appendChild(descEl);
        grid.appendChild(card);
    });
}

function initAiPage() {
    showAiAppList();
    loadAiApps();
}

function showAiAppList() {
    document.getElementById('aiAppList').classList.remove('hidden');
    document.getElementById('aiChatPage').classList.add('hidden');
    if (cozeWebSDK) {
        cozeWebSDK = null;
    }
}

function openAiApp(app) {
    if (!app || !app.appId) {
        console.error('应用数据无效:', app);
        return;
    }
    
    currentAppId = app.appId;
    document.getElementById('aiAppList').classList.add('hidden');
    document.getElementById('aiChatPage').classList.remove('hidden');
    document.getElementById('aiChatTitle').textContent = app.name;
    
    initAiChat(app);
}

function backToAiAppList() {
    showAiAppList();
}

function initAiChat(app) {
    if (cozeWebSDK) {
        return;
    }
    
    if (typeof CozeWebSDK === 'undefined') {
        console.error('CozeWebSDK未加载');
        return;
    }
    
    const userInfo = checkLoginStatus();
    const nickname = userInfo ? userInfo.username : '用户';
    const userId = userInfo ? userInfo.userId : '';
    
    const container = document.getElementById('aiApp');
    container.innerHTML = '';
    
    cozeWebSDK = new CozeWebSDK.AppWebSDK({
        token: 'pat_aYiVuZJtV1rOBp0MpQUIr6r4nciF3N6j7aDWJ7f10WrFS6gtWZfYFxWjDGiIlXre',
        appId: app.appId,
        container: '#aiApp',
        userInfo: { 
            id: userId, 
            url: '', 
            nickname: nickname, 
        }, 
        ui: {
            className: 'coze-app-sdk'
        }
    });
}
