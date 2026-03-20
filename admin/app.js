// 全局状态
const state = {
    token: localStorage.getItem('admin_token'),
    users: [],
    models: [],
    selectedUserIds: new Set(),
    currentTab: 'dashboard',
    theme: localStorage.getItem('admin_theme') || 'light'
};

// API 基础 URL
const API_BASE_URL = getApiBaseUrl();

// API URL 检测
function getApiBaseUrl() {
    if (window.location.hostname === 'localhost' || window.location.hostname === '127.0.0.1') {
        return 'http://localhost:8080';
    }
    return '';
}

// 应用主题
function applyTheme(theme) {
    document.documentElement.setAttribute('data-theme', theme);
    localStorage.setItem('admin_theme', theme);
    state.theme = theme;
}

// 初始化主题
applyTheme(state.theme);

// 通用 fetch 函数
async function api(path, options = {}) {
    const token = state.token;
    const headers = {
        'Content-Type': 'application/json',
        ...(options.headers || {}),
    };

    if (token) {
        headers['Authorization'] = `Bearer ${token}`;
    }

    const url = API_BASE_URL + path;

    try {
        const response = await fetch(url, {
            ...options,
            headers,
        });

        const isJSON = (response.headers.get('content-type') || '').includes('application/json');
        const body = isJSON ? await response.json() : await response.text();

        if (response.status >= 200 && response.status < 300) {
            // Extract data from APIResponse structure if present
            if (body && body.data !== undefined) {
                return body.data;
            }
            return body;
        }

        // 处理 401 未授权
        if (response.status === 401) {
            logout();
            throw new Error('登录已过期，请重新登录');
        }

        // 检查是否是 APIResponse 错误格式
        if (body?.message) {
            throw new Error(body.message);
        }
        throw new Error(body?.error || body || `HTTP ${response.status}`);
    } catch (error) {
        console.error('API请求失败:', error);
        throw error;
    }
}

// 初始化
document.addEventListener('DOMContentLoaded', async () => {
    // 检查登录状态
    if (!state.token) {
        window.location.href = '/login.html';
        return;
    }

    try {
        // 验证 token 并获取管理员信息
        const me = await api('/api/admin/me');
        const userAvatar = document.querySelector('.user-avatar');
        const userName = document.querySelector('.user-name');
        if (userAvatar) userAvatar.textContent = (me.name || 'A')[0].toUpperCase();
        if (userName) userName.textContent = me.name || '管理员';
        console.log('管理员信息:', me);

        // 加载初始数据
        await loadDashboard();
        await loadUsers();
        await loadModels();
    } catch (error) {
        console.error('初始化失败:', error);
        if (error.message.includes('登录已过期')) {
            logout();
        }
    }
});

// 切换标签页
function switchTab(tabName) {
    state.currentTab = tabName;

    // 更新导航样式
    document.querySelectorAll('.nav-item').forEach(nav => {
        nav.classList.remove('active');
        if (nav.dataset.tab === tabName) {
            nav.classList.add('active');
        }
    });

    // 显示对应内容
    document.querySelectorAll('.tab-content').forEach(content => {
        content.classList.remove('active');
    });
    document.getElementById(tabName).classList.add('active');

    // 根据标签加载数据
    switch (tabName) {
        case 'dashboard':
            loadDashboard();
            break;
        case 'users':
            loadUsers();
            break;
        case 'models':
            loadModels();
            break;
    }
}

// 加载仪表板
async function loadDashboard() {
    try {
        console.log('开始加载仪表板数据...');
        // 并行加载数据
        const [usersData, modelsData] = await Promise.all([
            api('/api/admin/users'),
            api('/api/admin/models')
        ]);

        console.log('仪表板用户数据:', usersData);
        console.log('仪表板模型数据:', modelsData);

        // api() 函数已自动提取 data 字段，直接使用
        const users = usersData.users || [];
        const models = modelsData.models || [];

        console.log('提取的用户数组:', users);
        console.log('提取的模型数组:', models);
        console.log('用户数组长度:', users.length);
        console.log('模型数组长度:', models.length);

        const totalUsers = users.length;
        const activeUsers = users.filter(u => u.status === 'active').length;
        const inactiveUsers = totalUsers - activeUsers;

        console.log('统计数据 - 总用户:', totalUsers, '活跃用户:', activeUsers, '禁用用户:', inactiveUsers);

        // 更新统计数据
        document.getElementById('totalUsers').textContent = totalUsers;
        document.getElementById('activeUsers').textContent = activeUsers;
        document.getElementById('inactiveUsers').textContent = inactiveUsers;
        document.getElementById('totalModels').textContent = models.length;

        // 加载最近活动（模拟数据）
        const recentActivityEl = document.getElementById('recentActivity');

        if (users.length === 0 && models.length === 0) {
            recentActivityEl.innerHTML = `
                <div class="empty-state">
                    <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5">
                        <circle cx="12" cy="12" r="10"/>
                        <polyline points="12 6 12 12 16 14"/>
                    </svg>
                    <div class="empty-state-title">暂无活动记录</div>
                </div>
            `;
        } else {
            let activities = [];
            if (users.length > 0) {
                activities.push(`
                    <div class="activity-item" style="display: flex; justify-content: space-between; padding: 12px; border-bottom: 1px solid var(--border);">
                        <span>用户 ${users[0]?.name || 'Unknown'} 创建了新会话</span>
                        <span style="color: var(--text-secondary); font-size: 12px;">刚刚</span>
                    </div>
                `);
            }
            if (models.length > 0) {
                activities.push(`
                    <div class="activity-item" style="display: flex; justify-content: space-between; padding: 12px; border-bottom: 1px solid var(--border);">
                        <span>模型 ${models[0]?.name || 'Unknown'} 已更新</span>
                        <span style="color: var(--text-secondary); font-size: 12px;">5分钟前</span>
                    </div>
                `);
            }
            activities.push(`
                <div class="activity-item" style="display: flex; justify-content: space-between; padding: 12px;">
                    <span>系统正常运行中</span>
                    <span style="color: var(--text-secondary); font-size: 12px;">1小时前</span>
                </div>
            `);

            recentActivityEl.innerHTML = `<div style="display: flex; flex-direction: column;">${activities.join('')}</div>`;
        }

    } catch (error) {
        console.error('加载仪表板失败:', error);
        document.getElementById('recentActivity').innerHTML = `
            <div class="empty-state">
                <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5">
                    <circle cx="12" cy="12" r="10"/>
                    <line x1="12" y1="8" x2="12" y2="12"/>
                    <line x1="12" y1="16" x2="12.01" y2="16"/>
                </svg>
                <div class="empty-state-title">加载失败</div>
            </div>
        `;
    }
}

// 刷新仪表板
async function refreshDashboard() {
    const btn = event.target.closest('.btn');
    btn.disabled = true;
    btn.innerHTML = `
        <svg class="spinner" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" style="width: 16px; height: 16px;">
            <path d="M21 12a9 9 0 1 1-6.219-8.56"/>
        </svg>
        刷新中
    `;

    try {
        await loadDashboard();
    } finally {
        btn.disabled = false;
        btn.innerHTML = `
            <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                <path d="M23 4v6h-6"/>
                <path d="M20.49 15a9 9 0 1 1-2.12-9.36L23 10"/>
            </svg>
            刷新
        `;
    }
}

// 加载用户列表
async function loadUsers() {
    try {
        console.log('开始加载用户数据...');
        const userData = await api('/api/admin/users');
        console.log('用户API返回数据:', userData);

        // 检查数据结构
        if (!userData || typeof userData !== 'object') {
            console.error('用户数据结构错误:', userData);
            showAlert('user', 'error', '数据格式错误: 期望对象，但得到 ' + typeof userData);
            return;
        }

        // api() 函数已自动提取 data 字段，直接使用
        const users = userData.users || [];
        console.log('提取的用户数组:', users);
        console.log('用户数组长度:', users.length);

        // 检查用户数据
        if (!Array.isArray(users)) {
            console.error('用户数据不是数组:', users);
            showAlert('user', 'error', '数据格式错误: 期望数组，但得到 ' + typeof users);
            return;
        }

        if (users.length === 0) {
            console.log('没有用户数据');
            renderUsers([]);
            return;
        }

        // 验证第一个用户的数据结构
        if (users[0]) {
            console.log('第一个用户数据:', users[0]);
            console.log('用户ID:', users[0].id);
            console.log('用户名:', users[0].name);
            console.log('邮箱:', users[0].email);
            console.log('状态:', users[0].status);
            console.log('创建时间:', users[0].created_at);
            console.log('最后登录:', users[0].last_login_at);
            console.log('消息数量:', users[0].message_count);
            console.log('会话数量:', users[0].session_count);
        }

        state.users = users;
        renderUsers(users);
    } catch (error) {
        console.error('加载用户列表失败:', error);
        showAlert('user', 'error', '加载用户列表失败: ' + error.message);
    }
}

// 渲染用户列表
function renderUsers(users) {
    console.log('开始渲染用户列表，用户数量:', users.length);
    const tbody = document.getElementById('userTableBody');

    if (users.length === 0) {
        console.log('没有用户数据，显示空状态');
        tbody.innerHTML = `
            <tr>
                <td colspan="7">
                    <div class="empty-state">
                        <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5">
                            <path d="M17 21v-2a4 4 0 0 0-4-4H5a4 4 0 0 0-4 4v2"/>
                            <circle cx="9" cy="7" r="4"/>
                        </svg>
                        <div class="empty-state-title">暂无用户数据</div>
                    </div>
                </td>
            </tr>
        `;
        return;
    }

    console.log('渲染用户数据...');
    tbody.innerHTML = users.map(user => {
        console.log('渲染单个用户:', user);
        const statusClass = user.status === 'active' ? 'active' : 'inactive';
        const statusText = user.status === 'active' ? '激活' : user.status === 'inactive' ? '禁用' : '未知';
        const statusLabel = user.status === 'active' ? '禁用' : '启用';
        const statusBtnClass = user.status === 'active' ? 'danger' : 'success';

        const isSelected = state.selectedUserIds.has(user.id);
        const rowClass = isSelected ? 'selected' : '';

        return `
        <tr class="${rowClass}" data-user-id="${user.id}">
            <td>
                <input type="checkbox" class="user-checkbox"
                       value="${user.id}"
                       ${isSelected ? 'checked' : ''}
                       onchange="updateSelectionState()">
            </td>
            <td>${user.id || 'N/A'}</td>
            <td>${user.name || 'N/A'}</td>
            <td>${user.email || 'N/A'}</td>
            <td>
                <span class="status-badge ${statusClass}">
                    <span class="status-dot"></span>
                    ${statusText}
                </span>
            </td>
            <td>${user.created_at ? formatDate(user.created_at) : 'N/A'}</td>
            <td>${user.last_login_at ? formatDate(user.last_login_at) : '-'}</td>
            <td>
                <div class="action-buttons">
                    <button class="btn btn-secondary btn-icon" onclick="editUser('${user.id}')" title="编辑">
                        <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" width="14" height="14">
                            <path d="M11 4H4a2 2 0 0 0-2 2v14a2 2 0 0 0 2 2h14a2 2 0 0 0 2-2v-7"/>
                            <path d="M18.5 2.5a2.121 2.121 0 0 1 3 3L12 15l-4 1 1-4 9.5-9.5z"/>
                        </svg>
                    </button>
                    <button class="btn btn-${statusBtnClass} btn-icon" onclick="toggleUserStatus('${user.id}', '${user.status}')" title="${statusLabel}">
                        <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" width="14" height="14">
                            ${user.status === 'active'
                                ? '<rect x="3" y="3" width="18" height="18" rx="2" ry="2"/><line x1="9" y1="9" x2="15" y2="15"/><line x1="15" y1="9" x2="9" y2="15"/>'
                                : '<path d="M22 11.08V12a10 10 0 1 1-5.93-9.14"/><polyline points="22 4 12 14.01 9 11.01"/>'}
                        </svg>
                    </button>
                    <button class="btn btn-danger btn-icon" onclick="deleteUser('${user.id}')" title="删除">
                        <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" width="14" height="14">
                            <polyline points="3 6 5 6 21 6"/>
                            <path d="M19 6v14a2 2 0 0 1-2 2H7a2 2 0 0 1-2-2V6m3 0V4a2 2 0 0 1 2-2h4a2 2 0 0 1 2 2v2"/>
                        </svg>
                    </button>
                </div>
            </td>
        </tr>
    `;
}).join('');
}

// 加载模型列表
async function loadModels() {
    try {
        console.log('开始加载模型数据...');
        const modelData = await api('/api/admin/models');
        console.log('模型API返回数据:', modelData);

        // 检查数据结构
        if (!modelData || typeof modelData !== 'object') {
            console.error('模型数据结构错误:', modelData);
            showAlert('model', 'error', '数据格式错误: 期望对象，但得到 ' + typeof modelData);
            return;
        }

        // api() 函数已自动提取 data 字段，直接使用
        const models = modelData.models || [];
        console.log('提取的模型数组:', models);
        console.log('模型数组长度:', models.length);

        // 检查模型数据
        if (!Array.isArray(models)) {
            console.error('模型数据不是数组:', models);
            showAlert('model', 'error', '数据格式错误: 期望数组，但得到 ' + typeof models);
            return;
        }

        if (models.length === 0) {
            console.log('没有模型数据');
            renderModels([]);
            return;
        }

        // 验证第一个模型的数据结构
        if (models[0]) {
            console.log('第一个模型数据:', models[0]);
            console.log('模型ID:', models[0].id);
            console.log('模型名称:', models[0].name);
            console.log('模型ID:', models[0].model_id);
            console.log('状态:', models[0].status);
            console.log('描述:', models[0].description);
            console.log('创建时间:', models[0].created_at);
            console.log('消息数量:', models[0].message_count);
        }

        state.models = models;
        renderModels(models);
    } catch (error) {
        console.error('加载模型列表失败:', error);
        showAlert('model', 'error', '加载模型列表失败: ' + error.message);
    }
}

// 渲染模型列表
function renderModels(models) {
    console.log('开始渲染模型列表，模型数量:', models.length);
    const tbody = document.getElementById('modelTableBody');

    if (models.length === 0) {
        console.log('没有模型数据，显示空状态');
        tbody.innerHTML = `
            <tr>
                <td colspan="7">
                    <div class="empty-state">
                        <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5">
                            <rect x="3" y="3" width="18" height="18" rx="2" ry="2"/>
                            <line x1="9" y1="9" x2="15" y2="15"/>
                            <line x1="15" y1="9" x2="9" y2="15"/>
                        </svg>
                        <div class="empty-state-title">暂无模型数据</div>
                    </div>
                </td>
            </tr>
        `;
        return;
    }

    console.log('渲染模型数据...');
    try {
        tbody.innerHTML = models.map(model => {
            console.log('渲染单个模型:', model);

            // 安全处理可能为null或undefined的值
            const id = model.id || 'N/A';
            const name = model.name || 'N/A';
            const modelId = model.model_id || 'N/A';
            const status = model.status || 'unknown';
            const description = model.description || '-';
            const createdAt = model.created_at ? formatDate(model.created_at) : 'N/A';

            const statusClass = status === 'active' ? 'active' : 'inactive';
            const statusText = status === 'active' ? '激活' : status === 'inactive' ? '禁用' : '未知';
            const statusLabel = status === 'active' ? '禁用' : '启用';
            const statusBtnClass = status === 'active' ? 'danger' : 'success';

            return `
            <tr>
                <td>${id}</td>
                <td>${name}</td>
                <td>${modelId}</td>
                <td>
                    <span class="status-badge ${statusClass}">
                        <span class="status-dot"></span>
                        ${statusText}
                    </span>
                </td>
                <td>${description}</td>
                <td>${createdAt}</td>
                <td>
                    <div class="action-buttons">
                        <button class="btn btn-secondary btn-icon" onclick="editModel('${id}')" title="编辑">
                            <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" width="14" height="14">
                                <path d="M11 4H4a2 2 0 0 0-2 2v14a2 2 0 0 0 2 2h14a2 2 0 0 0 2-2v-7"/>
                                <path d="M18.5 2.5a2.121 2.121 0 0 1 3 3L12 15l-4 1 1-4 9.5-9.5z"/>
                            </svg>
                        </button>
                        <button class="btn btn-${statusBtnClass} btn-icon" onclick="toggleModelStatus('${id}', '${status}')" title="${statusLabel}">
                            <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" width="14" height="14">
                                ${status === 'active'
                                    ? '<rect x="3" y="3" width="18" height="18" rx="2" ry="2"/><line x1="9" y1="9" x2="15" y2="15"/><line x1="15" y1="9" x2="9" y2="15"/>'
                                    : '<path d="M22 11.08V12a10 10 0 1 1-5.93-9.14"/><polyline points="22 4 12 14.01 9 11.01"/>'}
                            </svg>
                        </button>
                        <button class="btn btn-danger btn-icon" onclick="deleteModel('${id}')" title="删除">
                            <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" width="14" height="14">
                                <polyline points="3 6 5 6 21 6"/>
                                <path d="M19 6v14a2 2 0 0 1-2 2H7a2 2 0 0 1-2-2V6m3 0V4a2 2 0 0 1 2-2h4a2 2 0 0 1 2 2v2"/>
                            </svg>
                        </button>
                    </div>
                </td>
            </tr>
        `;
        }).join('');
    } catch (error) {
        console.error('渲染模型列表失败:', error);
        showAlert('model', 'error', '渲染模型列表失败: ' + error.message);
        // 显示错误状态
        tbody.innerHTML = `
            <tr>
                <td colspan="7">
                    <div class="empty-state">
                        <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5">
                            <circle cx="12" cy="12" r="10"/>
                            <line x1="12" y1="8" x2="12" y2="12"/>
                            <line x1="12" y1="16" x2="12.01" y2="16"/>
                        </svg>
                        <div class="empty-state-title">加载数据时发生错误</div>
                    </div>
                </td>
            </tr>
        `;
    }
}

// 显示用户模态框
function showUserModal(userId = null) {
    const modal = document.getElementById('userModal');
    const form = document.getElementById('userForm');
    const title = document.getElementById('userModalTitle');
    const passwordGroup = document.getElementById('userPassword').closest('.form-group');

    form.reset();

    if (userId) {
        const user = state.users.find(u => u.id === userId);
        if (user) {
            title.textContent = '编辑用户';
            document.getElementById('userId').value = user.id;
            document.getElementById('userName').value = user.name;
            document.getElementById('userEmail').value = user.email;
            document.getElementById('userStatus').value = user.status;
            // 编辑时不显示密码输入框
            passwordGroup.style.display = 'none';
        }
    } else {
        title.textContent = '添加用户';
        document.getElementById('userId').value = '';
        passwordGroup.style.display = 'block';
    }

    modal.classList.add('active');
}

// 关闭用户模态框
function closeUserModal() {
    document.getElementById('userModal').classList.remove('active');
}

// 编辑用户
function editUser(userId) {
    showUserModal(userId);
}

// 保存用户
async function saveUser(event) {
    event.preventDefault();

    const userId = document.getElementById('userId').value;
    const userData = {
        name: document.getElementById('userName').value,
        email: document.getElementById('userEmail').value,
        status: document.getElementById('userStatus').value,
    };

    // 添加密码（仅创建时需要）
    if (!userId) {
        userData.password = document.getElementById('userPassword').value;
    }

    try {
        if (userId) {
            await api(`/api/admin/users/${userId}`, {
                method: 'PUT',
                body: JSON.stringify(userData)
            });
            showAlert('user', 'success', '用户更新成功');
        } else {
            await api('/api/admin/users', {
                method: 'POST',
                body: JSON.stringify(userData)
            });
            showAlert('user', 'success', '用户创建成功');
        }

        closeUserModal();
        loadUsers();
    } catch (error) {
        showAlert('user', 'error', error.message);
    }
}

// 切换用户状态
async function toggleUserStatus(userId, currentStatus) {
    const newStatus = currentStatus === 'active' ? 'inactive' : 'active';
    const action = newStatus === 'active' ? '启用' : '禁用';

    if (!confirm(`确定要${action}该用户吗？`)) {
        return;
    }

    try {
        await api(`/api/admin/users/${userId}`, {
            method: 'PUT',
            body: JSON.stringify({ status: newStatus })
        });

        showAlert('user', 'success', `用户${action}成功`);
        loadUsers();
    } catch (error) {
        showAlert('user', 'error', error.message);
    }
}

// 删除用户
async function deleteUser(userId) {
    if (!confirm('确定要删除该用户吗？此操作不可恢复！')) {
        return;
    }

    try {
        await api(`/api/admin/users/${userId}`, {
            method: 'DELETE'
        });

        showAlert('user', 'success', '用户删除成功');
        loadUsers();
    } catch (error) {
        showAlert('user', 'error', error.message);
    }
}

// 显示模型模态框
function showModelModal(modelId = null) {
    const modal = document.getElementById('modelModal');
    const form = document.getElementById('modelForm');
    const title = document.getElementById('modelModalTitle');

    form.reset();

    if (modelId) {
        const model = state.models.find(m => m.id === modelId);
        if (model) {
            title.textContent = '编辑模型';
            document.getElementById('modelId').value = model.id;
            document.getElementById('modelName').value = model.name;
            document.getElementById('modelIdInput').value = model.model_id;
            document.getElementById('modelStatus').value = model.status;
            document.getElementById('modelDescription').value = model.description || '';
        }
    } else {
        title.textContent = '添加模型';
        document.getElementById('modelId').value = '';
    }

    modal.classList.add('active');
}

// 关闭模型模态框
function closeModelModal() {
    document.getElementById('modelModal').classList.remove('active');
}

// 编辑模型
function editModel(modelId) {
    showModelModal(modelId);
}

// 保存模型
async function saveModel(event) {
    event.preventDefault();

    const modelId = document.getElementById('modelId').value;
    const modelData = {
        name: document.getElementById('modelName').value,
        model_id: document.getElementById('modelIdInput').value,
        status: document.getElementById('modelStatus').value,
        description: document.getElementById('modelDescription').value,
    };

    try {
        if (modelId) {
            await api(`/api/admin/models/${modelId}`, {
                method: 'PUT',
                body: JSON.stringify(modelData)
            });
            showAlert('model', 'success', '模型更新成功');
        } else {
            await api('/api/admin/models', {
                method: 'POST',
                body: JSON.stringify(modelData)
            });
            showAlert('model', 'success', '模型创建成功');
        }

        closeModelModal();
        loadModels();
    } catch (error) {
        showAlert('model', 'error', error.message);
    }
}

// 切换模型状态
async function toggleModelStatus(modelId, currentStatus) {
    const newStatus = currentStatus === 'active' ? 'inactive' : 'active';
    const action = newStatus === 'active' ? '启用' : '禁用';

    if (!confirm(`确定要${action}该模型吗？`)) {
        return;
    }

    try {
        await api(`/api/admin/models/${modelId}`, {
            method: 'PUT',
            body: JSON.stringify({ status: newStatus })
        });

        showAlert('model', 'success', `模型${action}成功`);
        loadModels();
    } catch (error) {
        showAlert('model', 'error', error.message);
    }
}

// 删除模型
async function deleteModel(modelId) {
    if (!confirm('确定要删除该模型吗？此操作不可恢复！')) {
        return;
    }

    try {
        await api(`/api/admin/models/${modelId}`, {
            method: 'DELETE'
        });

        showAlert('model', 'success', '模型删除成功');
        loadModels();
    } catch (error) {
        showAlert('model', 'error', error.message);
    }
}

// 登出
function logout() {
    localStorage.removeItem('admin_token');
    window.location.href = '/login.html';
}

// 显示提示信息
function showAlert(type, severity, message) {
    const container = document.getElementById(`${type}AlertContainer`);

    const alertEl = document.createElement('div');
    alertEl.className = `alert alert-${severity}`;
    alertEl.innerHTML = `
        ${severity === 'success'
            ? '<svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="M22 11.08V12a10 10 0 1 1-5.93-9.14"/><polyline points="22 4 12 14.01 9 11.01"/></svg>'
            : '<svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><circle cx="12" cy="12" r="10"/><line x1="12" y1="8" x2="12" y2="12"/><line x1="12" y1="16" x2="12.01" y2="16"/></svg>'
        }
        <span class="alert-message">${message}</span>
    `;

    container.innerHTML = '';
    container.appendChild(alertEl);

    // 自动消失
    setTimeout(() => {
        alertEl.style.opacity = '0';
        alertEl.style.transform = 'translateY(-10px)';
        setTimeout(() => alertEl.remove(), 200);
    }, 5000);
}

// 格式化日期
function formatDate(dateString) {
    if (!dateString) return '-';

    // 检查是否是时间戳（数字）
    if (typeof dateString === 'number') {
        const date = new Date(dateString);
        return date.toLocaleString('zh-CN', {
            year: 'numeric',
            month: '2-digit',
            day: '2-digit',
            hour: '2-digit',
            minute: '2-digit'
        });
    }

    // 检查是否是ISO格式字符串
    if (typeof dateString === 'string' && dateString.includes('T')) {
        const date = new Date(dateString);
        return date.toLocaleString('zh-CN', {
            year: 'numeric',
            month: '2-digit',
            day: '2-digit',
            hour: '2-digit',
            minute: '2-digit'
        });
    }

    // 默认处理
    const date = new Date(dateString);
    return date.toLocaleString('zh-CN', {
        year: 'numeric',
        month: '2-digit',
        day: '2-digit',
        hour: '2-digit',
        minute: '2-digit'
    });
}

// 点击模态框外部关闭
document.querySelectorAll('.modal-overlay').forEach(overlay => {
    overlay.addEventListener('click', (e) => {
        if (e.target === overlay) {
            overlay.classList.remove('active');
        }
    });
});

// ESC键关闭模态框
document.addEventListener('keydown', (e) => {
    if (e.key === 'Escape') {
        document.querySelectorAll('.modal-overlay.active').forEach(overlay => {
            overlay.classList.remove('active');
        });
    }
});
