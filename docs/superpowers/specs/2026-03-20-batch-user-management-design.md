# Batch User Management Design Spec

## Overview

Add batch management functionality to the user management module in the Luban admin panel. This feature allows administrators to select multiple users and perform batch operations: delete, disable, and activate.

## Requirements

- Support batch operations on multiple users simultaneously
- Maximum of 50 users per batch operation
- Validate all user IDs exist before performing operations
- Confirmation dialog before executing batch operations
- Clear UI feedback for selection state and operation results
- Loading state during batch operations to prevent duplicate requests

## Backend Implementation

### New Data Models

Add to `internal/model/admin.go`:

```go
// BatchUserStatusRequest 批量更新用户状态请求
type BatchUserStatusRequest struct {
    UserIDs []string `json:"user_ids" binding:"required,min=1,max=50"`
    Status  string   `json:"status" binding:"required,oneof=active inactive"`
}

// BatchDeleteRequest 批量删除用户请求
type BatchDeleteRequest struct {
    UserIDs []string `json:"user_ids" binding:"required,min=1,max=50"`
}
```

### New API Endpoints

#### 1. Batch Update User Status

**Endpoint:** `PUT /api/admin/users/batch/status`

**Request Body:**
```json
{
  "user_ids": ["uuid1", "uuid2", "uuid3"],
  "status": "active"
}
```

**Success Response:**
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "updated_count": 3
  },
  "requestId": "uuid"
}
```

**Error Responses:**
- `10000001` - Invalid parameters (empty user_ids, invalid status, or >50 IDs)
- `10000004` - User not found (one or more user IDs don't exist)
- `20000003` - Database error

#### 2. Batch Delete Users

**Endpoint:** `DELETE /api/admin/users/batch`

**Request Body:**
```json
{
  "user_ids": ["uuid1", "uuid2", "uuid3"]
}
```

**Success Response:**
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "deleted_count": 3
  },
  "requestId": "uuid"
}
```

**Error Responses:**
- `10000001` - Invalid parameters (empty user_ids or >50 IDs)
- `10000004` - User not found (one or more user IDs don't exist)
- `20000003` - Database error

### Handler Implementation

Add to `internal/handler/admin.go`:

```go
// adminBatchUpdateUserStatus 批量更新用户状态
func (a *AdminApp) adminBatchUpdateUserStatus(c *gin.Context) {
    var req model.BatchUserStatusRequest

    if err := c.ShouldBindJSON(&req); err != nil {
        response.NewResponseHelper(c).Error(response.CodeInvalidParam, err.Error())
        return
    }

    // 去重用户ID
    uniqueUserIDs := make(map[string]bool)
    for _, id := range req.UserIDs {
        uniqueUserIDs[id] = true
    }
    req.UserIDs = make([]string, 0, len(uniqueUserIDs))
    for id := range uniqueUserIDs {
        req.UserIDs = append(req.UserIDs, id)
    }

    // 验证所有用户存在
    var count int64
    if err := a.db.Model(&model.User{}).Where("id IN ?", req.UserIDs).Count(&count).Error; err != nil {
        response.NewResponseHelper(c).Error(response.CodeDatabaseError, "failed to verify users")
        return
    }
    if int(count) != len(req.UserIDs) {
        response.NewResponseHelper(c).Error(response.CodeNotFound, "one or more users not found")
        return
    }

    // 批量更新
    result := a.db.Model(&model.User{}).
        Where("id IN ?", req.UserIDs).
        Update("status", req.Status)

    if result.Error != nil {
        response.NewResponseHelper(c).Error(response.CodeDatabaseError, "failed to update user status")
        return
    }

    response.NewResponseHelper(c).Success(gin.H{
        "updated_count": result.RowsAffected,
    })
}

// adminBatchDeleteUsers 批量删除用户
func (a *AdminApp) adminBatchDeleteUsers(c *gin.Context) {
    var req model.BatchDeleteRequest

    if err := c.ShouldBindJSON(&req); err != nil {
        response.NewResponseHelper(c).Error(response.CodeInvalidParam, err.Error())
        return
    }

    // 去重用户ID
    uniqueUserIDs := make(map[string]bool)
    for _, id := range req.UserIDs {
        uniqueUserIDs[id] = true
    }
    req.UserIDs = make([]string, 0, len(uniqueUserIDs))
    for id := range uniqueUserIDs {
        req.UserIDs = append(req.UserIDs, id)
    }

    // 验证所有用户存在
    var count int64
    if err := a.db.Model(&model.User{}).Where("id IN ?", req.UserIDs).Count(&count).Error; err != nil {
        response.NewResponseHelper(c).Error(response.CodeDatabaseError, "failed to verify users")
        return
    }
    if int(count) != len(req.UserIDs) {
        response.NewResponseHelper(c).Error(response.CodeNotFound, "one or more users not found")
        return
    }

    // 批量删除（物理删除，与现有 adminDeleteUser 保持一致）
    result := a.db.Where("id IN ?", req.UserIDs).Delete(&model.User{})

    if result.Error != nil {
        response.NewResponseHelper(c).Error(response.CodeDatabaseError, "failed to delete users")
        return
    }

    response.NewResponseHelper(c).Success(gin.H{
        "deleted_count": result.RowsAffected,
    })
}
```

### Route Registration

Add to `registerAdminRoutes()` in `internal/handler/admin.go`:

```go
// 用户管理
adminGroup.GET("/users", adminApp.adminGetUsers)
adminGroup.POST("/users", adminApp.adminCreateUser)
adminGroup.PUT("/users/:id", adminApp.adminUpdateUser)
adminGroup.PATCH("/users/:id/status", adminApp.adminToggleUserStatus)
adminGroup.DELETE("/users/:id", adminApp.adminDeleteUser)
adminGroup.PUT("/users/batch/status", adminApp.adminBatchUpdateUserStatus) // NEW
adminGroup.DELETE("/users/batch", adminApp.adminBatchDeleteUsers)         // NEW
```

## Frontend Implementation

### UI Changes in `admin/index.html`

#### 1. Table Header - Add Checkbox Column

```html
<thead>
    <tr>
        <th style="width: 40px;">
            <input type="checkbox" id="selectAllUsers" onclick="toggleSelectAll()">
        </th>
        <th>ID</th>
        <th>用户名</th>
        <th>邮箱</th>
        <th>状态</th>
        <th>创建时间</th>
        <th>最后登录</th>
        <th>操作</th>
    </tr>
</thead>
```

#### 2. Table Body - Add Row Checkbox

Update `renderUsers()` function to include checkbox:
```javascript
return `
<tr data-user-id="${user.id}">
    <td>
        <input type="checkbox" class="user-checkbox" value="${user.id}" onchange="updateSelectionState()">
    </td>
    <td>${user.id || 'N/A'}</td>
    ...
</tr>
`;
```

#### 3. Selection Toolbar - Add New UI Element

Add after top bar in users tab:
```html
<div id="selectionToolbar" class="selection-toolbar" style="display: none;">
    <div class="selection-info">
        <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" width="20" height="20">
            <path d="M9 11l3 3L22 4"/>
            <path d="M21 12v7a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2V5a2 2 0 0 1 2-2h11"/>
        </svg>
        <span id="selectedCount">已选择 0 个用户</span>
    </div>
    <div class="selection-actions">
        <button class="btn btn-success" onclick="batchUpdateStatus('active')">
            <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                <path d="M22 11.08V12a10 10 0 1 1-5.93-9.14"/>
                <polyline points="22 4 12 14.01 9 11.01"/>
            </svg>
            激活
        </button>
        <button class="btn btn-warning" onclick="batchUpdateStatus('inactive')">
            <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                <rect x="3" y="3" width="18" height="18" rx="2" ry="2"/>
                <line x1="9" y1="9" x2="15" y2="15"/>
                <line x1="15" y1="9" x2="9" y2="15"/>
            </svg>
            禁用
        </button>
        <button class="btn btn-danger" onclick="batchDeleteUsers()">
            <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                <polyline points="3 6 5 6 21 6"/>
                <path d="M19 6v14a2 2 0 0 1-2 2H7a2 2 0 0 1-2-2V6m3 0V4a2 2 0 0 1 2-2h4a2 2 0 0 1 2 2v2"/>
            </svg>
            删除
        </button>
        <button class="btn btn-ghost" onclick="clearSelection()">
            取消选择
        </button>
    </div>
</div>
```

#### 4. Selection Toolbar Styles

Add to CSS:
```css
/* Selection Toolbar */
.selection-toolbar {
    display: flex;
    align-items: center;
    justify-content: space-between;
    background: var(--primary);
    color: white;
    padding: var(--space-3) var(--space-4);
    border-radius: var(--radius-md);
    margin-bottom: var(--space-4);
    box-shadow: var(--shadow-md);
    animation: slideDown var(--transition-base);
}

@keyframes slideDown {
    from {
        opacity: 0;
        transform: translateY(-10px);
    }
    to {
        opacity: 1;
        transform: translateY(0);
    }
}

.selection-info {
    display: flex;
    align-items: center;
    gap: var(--space-2);
    font-weight: 500;
}

.selection-actions {
    display: flex;
    gap: var(--space-2);
}

.selection-actions .btn {
    background: rgba(255, 255, 255, 0.15);
    color: white;
    border-color: rgba(255, 255, 255, 0.3);
}

.selection-actions .btn:hover {
    background: rgba(255, 255, 255, 0.25);
}

.selection-actions .btn-danger {
    background: rgba(220, 38, 38, 0.8);
    border-color: rgba(220, 38, 38, 1);
}

.selection-actions .btn-danger:hover {
    background: rgba(220, 38, 38, 1);
}

.btn-warning {
    background: var(--warning-bg);
    color: var(--warning-text);
    border-color: var(--warning-border);
}

.btn-warning:hover:not(:disabled) {
    background: #fcd34d;
    border-color: #fbbf24;
}

/* Table Checkboxes */
.user-checkbox {
    width: 18px;
    height: 18px;
    cursor: pointer;
    accent-color: var(--primary);
}

tr.selected {
    background: var(--primary-light);
}
```

### JavaScript Changes in `admin/app.js`

#### 1. Add to Global State

```javascript
const state = {
    token: localStorage.getItem('admin_token'),
    users: [],
    models: [],
    selectedUserIds: new Set(), // NEW
    currentTab: 'dashboard',
    theme: localStorage.getItem('admin_theme') || 'light'
};
```

#### 2. Selection Management Functions

```javascript
// 更新选择状态
function updateSelectionState() {
    const checkboxes = document.querySelectorAll('.user-checkbox:checked');
    // 自动去重（Set会自动去重）
    state.selectedUserIds = new Set(Array.from(checkboxes).map(cb => cb.value));

    const toolbar = document.getElementById('selectionToolbar');
    const countSpan = document.getElementById('selectedCount');

    if (state.selectedUserIds.size > 0) {
        toolbar.style.display = 'flex';
        countSpan.textContent = `已选择 ${state.selectedUserIds.size} 个用户`;
    } else {
        toolbar.style.display = 'none';
    }

    // 更新"全选"复选框状态
    const selectAll = document.getElementById('selectAllUsers');
    const allCheckboxes = document.querySelectorAll('.user-checkbox');
    if (allCheckboxes.length > 0) {
        selectAll.checked = Array.from(allCheckboxes).every(cb => cb.checked);
        selectAll.indeterminate = Array.from(allCheckboxes).some(cb => cb.checked) &&
                                !Array.from(allCheckboxes).every(cb => cb.checked);
    }
}

// 全选/取消全选
function toggleSelectAll() {
    const selectAll = document.getElementById('selectAllUsers');
    const checkboxes = document.querySelectorAll('.user-checkbox');
    checkboxes.forEach(cb => cb.checked = selectAll.checked);
    updateSelectionState();
}

// 清除选择
function clearSelection() {
    document.querySelectorAll('.user-checkbox').forEach(cb => cb.checked = false);
    document.getElementById('selectAllUsers').checked = false;
    updateSelectionState();
}

// 批量更新用户状态
async function batchUpdateStatus(status) {
    const userIds = Array.from(state.selectedUserIds);
    const action = status === 'active' ? '激活' : '禁用';

    if (userIds.length === 0) {
        showAlert('user', 'error', '请先选择用户');
        return;
    }

    if (!confirm(`确定要${action}选中的 ${userIds.length} 个用户吗？`)) {
        return;
    }

    // 禁用按钮防止重复点击
    const buttons = document.querySelectorAll('.selection-actions .btn');
    buttons.forEach(btn => {
        btn.disabled = true;
        btn.classList.add('btn-loading');
    });

    try {
        const result = await api('/api/admin/users/batch/status', {
            method: 'PUT',
            body: JSON.stringify({
                user_ids: userIds,
                status: status
            })
        });

        showAlert('user', 'success', `成功${action} ${result.updated_count} 个用户`);
        clearSelection();
        loadUsers();
    } catch (error) {
        showAlert('user', 'error', error.message);
    } finally {
        // 恢复按钮状态
        buttons.forEach(btn => {
            btn.disabled = false;
            btn.classList.remove('btn-loading');
        });
    }
}

// 批量删除用户
async function batchDeleteUsers() {
    const userIds = Array.from(state.selectedUserIds);

    if (userIds.length === 0) {
        showAlert('user', 'error', '请先选择用户');
        return;
    }

    if (!confirm(`确定要删除选中的 ${userIds.length} 个用户吗？此操作不可恢复！`)) {
        return;
    }

    // 禁用按钮防止重复点击
    const buttons = document.querySelectorAll('.selection-actions .btn');
    buttons.forEach(btn => {
        btn.disabled = true;
        btn.classList.add('btn-loading');
    });

    try {
        const result = await api('/api/admin/users/batch', {
            method: 'DELETE',
            body: JSON.stringify({
                user_ids: userIds
            })
        });

        showAlert('user', 'success', `成功删除 ${result.deleted_count} 个用户`);
        clearSelection();
        loadUsers();
    } catch (error) {
        showAlert('user', 'error', error.message);
    } finally {
        // 恢复按钮状态
        buttons.forEach(btn => {
            btn.disabled = false;
            btn.classList.remove('btn-loading');
        });
    }
}
```

#### 3. Update renderUsers() Function

Modify to include checkbox in each row and handle selected state:
```javascript
function renderUsers(users) {
    // ... existing empty state check ...

    tbody.innerHTML = users.map(user => {
        const statusClass = user.status === 'active' ? 'active' : 'inactive';
        const statusText = user.status === 'active' ? '激活' : user.status === 'inactive' ? '禁用' : '未知';
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
                    <!-- existing action buttons -->
                </div>
            </td>
        </tr>
    `;
    }).join('');
}
```

## API Documentation Updates

Add to `API.md`:

### 23. 批量更新用户状态

**PUT** `/api/admin/users/batch/status`

**请求头：**
```
Authorization: Bearer {admin_token}
```

**请求体：**
```json
{
  "user_ids": ["uuid1", "uuid2", ...],
  "status": "active"
}
```

**参数限制：**
- `user_ids`: 1-50个用户ID数组
- `status`: "active" 或 "inactive"

**响应示例：**
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "updated_count": 5
  },
  "requestId": "uuid"
}
```

### 24. 批量删除用户

**DELETE** `/api/admin/users/batch`

**请求头：**
```
Authorization: Bearer {admin_token}
```

**请求体：**
```json
{
  "user_ids": ["uuid1", "uuid2", ...]
}
```

**参数限制：**
- `user_ids`: 1-50个用户ID数组

**响应示例：**
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "deleted_count": 5
  },
  "requestId": "uuid"
}
```

## Testing Checklist

- [ ] Batch activate 5 users
- [ ] Batch disable 5 users
- [ ] Batch delete 5 users
- [ ] Select all checkbox works correctly
- [ ] Partial selection displays correct count
- [ ] Confirmation dialog appears before operations
- [ ] Success message shows correct count
- [ ] Selection is cleared after successful operation
- [ ] Error handling for invalid user IDs
- [ ] Error handling for exceeding 50 user limit
- [ ] Empty selection prevents batch operations
- [ ] Transaction rollback on partial failures

## Files to Modify

1. `internal/model/admin.go` - Add batch request models
2. `internal/handler/admin.go` - Add batch handlers and routes
3. `admin/index.html` - Add checkbox column and selection toolbar
4. `admin/app.js` - Add selection management and batch operation functions
5. `API.md` - Document new endpoints
