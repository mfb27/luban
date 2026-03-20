# 批量模型管理功能设计文档

## 项目概述
为Luban管理后台的模型管理模块添加批量操作功能，支持批量删除、禁用和激活模型。

## 设计目标
- 实现原子性批量操作
- 提供用户友好的界面交互
- 保持数据一致性
- 与现有系统架构兼容

## 技术方案

### 后端实现

#### 1. 数据库模型更新
```go
// internal/model/models.go
type Model struct {
    ID          string    `gorm:"type:varchar(64);primaryKey" json:"id"`
    Name        string    `gorm:"type:varchar(128);not null" json:"name"`
    ModelID     string    `gorm:"type:varchar(128);uniqueIndex" json:"model_id"`
    Status      string    `gorm:"type:varchar(16);default:'active'" json:"status"`
    Description string    `gorm:"type:text" json:"description"`
    CreatedAt   time.Time `json:"created_at"`
    UpdatedAt   time.Time `json:"updated_at"`
    DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"` // 软删除支持
}
```

#### 2. 新增API端点
```go
// internal/handler/admin.go
func (a *AdminApp) adminBatchDeleteModels(c *gin.Context) {
    var req struct {
        IDs []string `json:"ids" binding:"required"`
    }

    if err := c.ShouldBindJSON(&req); err != nil {
        response.NewResponseHelper(c).Error(response.CodeInvalidParam, "invalid request")
        return
    }

    // 开启事务
    tx := a.db.Begin()
    defer func() {
        if r := recover(); r != nil {
            tx.Rollback()
            panic(r)
        }
    }()

    var successCount, errorCount int
    var errorMessages []string

    for _, id := range req.IDs {
        var model model.Model
        if err := tx.First(&model, "id = ?", id).Error; err != nil {
            errorCount++
            errorMessages = append(errorMessages, fmt.Sprintf("模型 %s 不存在", id))
            continue
        }

        // 检查是否有关联的消息
        var messageCount int64
        tx.Model(&model.Message{}).Where("model_id = ?", id).Count(&messageCount)
        if messageCount > 0 {
            errorCount++
            errorMessages = append(errorMessages, fmt.Sprintf("模型 %s 有关联消息，无法删除", id))
            continue
        }

        if err := tx.Delete(&model).Error; err != nil {
            errorCount++
            errorMessages = append(errorMessages, fmt.Sprintf("删除模型 %s 失败: %v", id, err))
            continue
        }

        successCount++
    }

    if errorCount > 0 {
        tx.Rollback()
        response.NewResponseHelper(c).Error(response.CodeInvalidParam, fmt.Sprintf("批量删除失败: %d个模型处理失败，%d个成功", errorCount, successCount))
        return
    }

    if err := tx.Commit().Error; err != nil {
        response.NewResponseHelper(c).Error(response.CodeDatabaseError, "事务提交失败")
        return
    }

    response.NewResponseHelper(c).Success(gin.H{
        "count": successCount,
    })
}

func (a *AdminApp) adminBatchActivateModels(c *gin.Context) {
    var req struct {
        IDs []string `json:"ids" binding:"required"`
    }

    if err := c.ShouldBindJSON(&req); err != nil {
        response.NewResponseHelper(c).Error(response.CodeInvalidParam, "invalid request")
        return
    }

    // 开启事务
    tx := a.db.Begin()
    defer tx.Rollback()

    var successCount, errorCount int
    var errorMessages []string

    for _, id := range req.IDs {
        var model model.Model
        if err := tx.First(&model, "id = ?", id).Error; err != nil {
            errorCount++
            errorMessages = append(errorMessages, fmt.Sprintf("模型 %s 不存在", id))
            continue
        }

        if err := tx.Model(&model).Update("status", "active").Error; err != nil {
            errorCount++
            errorMessages = append(errorMessages, fmt.Sprintf("激活模型 %s 失败: %v", id, err))
            continue
        }

        successCount++
    }

    if errorCount > 0 {
        response.NewResponseHelper(c).Error(response.CodeInvalidParam, fmt.Sprintf("批量激活失败: %d个模型处理失败，%d个成功", errorCount, successCount))
        return
    }

    if err := tx.Commit().Error; err != nil {
        response.NewResponseHelper(c).Error(response.CodeDatabaseError, "事务提交失败")
        return
    }

    response.NewResponseHelper(c).Success(gin.H{
        "count": successCount,
    })
}

func (a *AdminApp) adminBatchDeactivateModels(c *gin.Context) {
    var req struct {
        IDs []string `json:"ids" binding:"required"`
    }

    if err := c.ShouldBindJSON(&req); err != nil {
        response.NewResponseHelper(c).Error(response.CodeInvalidParam, "invalid request")
        return
    }

    // 开启事务
    tx := a.db.Begin()
    defer tx.Rollback()

    var successCount, errorCount int
    var errorMessages []string

    for _, id := range req.IDs {
        var model model.Model
        if err := tx.First(&model, "id = ?", id).Error; err != nil {
            errorCount++
            errorMessages = append(errorMessages, fmt.Sprintf("模型 %s 不存在", id))
            continue
        }

        if err := tx.Model(&model).Update("status", "inactive").Error; err != nil {
            errorCount++
            errorMessages = append(errorMessages, fmt.Sprintf("禁用模型 %s 失败: %v", id, err))
            continue
        }

        successCount++
    }

    if errorCount > 0 {
        response.NewResponseHelper(c).Error(response.CodeInvalidParam, fmt.Sprintf("批量禁用失败: %d个模型处理失败，%d个成功", errorCount, successCount))
        return
    }

    if err := tx.Commit().Error; err != nil {
        response.NewResponseHelper(c).Error(response.CodeDatabaseError, "事务提交失败")
        return
    }

    response.NewResponseHelper(c).Success(gin.H{
        "count": successCount,
    })
}
```

#### 3. 路由注册
```go
// internal/handler/admin.go
func (a *App) registerAdminRoutes() {
    // ... 现有路由 ...

    // 模型管理 - 批量操作
    adminGroup.POST("/models/batch/delete", adminApp.adminBatchDeleteModels)
    adminGroup.POST("/models/batch/activate", adminApp.adminBatchActivateModels)
    adminGroup.POST("/models/batch/deactivate", adminApp.adminBatchDeactivateModels)
}
```

### 前端实现

#### 1. 表格结构更新
```html
<!-- admin/index.html -->
<table class="table" id="modelTable">
    <thead>
        <tr>
            <th width="40">
                <input type="checkbox" id="selectAll" class="form-checkbox">
            </th>
            <th>ID</th>
            <th>模型名称</th>
            <th>模型ID</th>
            <th>状态</th>
            <th>描述</th>
            <th>创建时间</th>
            <th>操作</th>
        </tr>
    </thead>
    <tbody id="modelTableBody">
        <!-- 模型数据 -->
    </tbody>
</table>
```

#### 2. 浮动操作栏
```html
<!-- admin/index.html -->
<div id="floatingActionBar" class="floating-action-bar" style="display: none;">
    <div class="action-bar-content">
        <span class="selected-count">已选择 <span id="selectedCount">0</span> 个模型</span>
        <div class="action-buttons">
            <button class="btn btn-danger" id="batchDeleteBtn">
                <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" width="16" height="16">
                    <polyline points="3 6 5 6 21 6"/>
                    <path d="M19 6v14a2 2 0 0 1-2 2H7a2 2 0 0 1-2-2V6m3 0V4a2 2 0 0 1 2-2h4a2 2 0 0 1 2 2v2"/>
                </svg>
                删除
            </button>
            <button class="btn btn-success" id="batchActivateBtn">
                <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" width="16" height="16">
                    <rect x="3" y="3" width="18" height="18" rx="2" ry="2"/>
                    <line x1="9" y1="9" x2="15" y2="15"/>
                    <line x1="15" y1="9" x2="9" y2="15"/>
                </svg>
                激活
            </button>
            <button class="btn btn-warning" id="batchDeactivateBtn">
                <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" width="16" height="16">
                    <path d="M22 11.08V12a10 10 0 1 1-5.93-9.14"/>
                    <polyline points="22 4 12 14.01 9 11.01"/>
                </svg>
                禁用
            </button>
            <button class="btn btn-secondary" id="clearSelectionBtn">
                <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" width="16" height="16">
                    <path d="M6 18L18 6M6 6l12 12"/>
                </svg>
                清除
            </button>
        </div>
    </div>
</div>
```

#### 3. JavaScript 实现
```javascript
// admin/app.js
// 全局状态
const state = {
    token: localStorage.getItem('admin_token'),
    users: [],
    models: [],
    currentTab: 'dashboard',
    theme: localStorage.getItem('admin_theme') || 'light',
    selectedModels: new Set(), // 存储选中的模型ID
};

// 初始化模型表格事件
function initModelTableEvents() {
    const modelTable = document.getElementById('modelTable');

    // 全选/取消全选
    document.getElementById('selectAll').addEventListener('change', function() {
        const isChecked = this.checked;
        const checkboxes = modelTable.querySelectorAll('tbody input[type="checkbox"]');

        checkboxes.forEach(checkbox => {
            checkbox.checked = isChecked;
            const modelId = checkbox.dataset.modelId;
            if (isChecked) {
                state.selectedModels.add(modelId);
            } else {
                state.selectedModels.delete(modelId);
            }
        });

        updateFloatingActionBar();
    });

    // 单个选择
    modelTable.addEventListener('change', function(e) {
        if (e.target.type === 'checkbox' && e.target.dataset.modelId) {
            const modelId = e.target.dataset.modelId;
            if (e.target.checked) {
                state.selectedModels.add(modelId);
            } else {
                state.selectedModels.delete(modelId);
            }
            updateFloatingActionBar();
        }
    });

    // 批量操作按钮
    document.getElementById('batchDeleteBtn').addEventListener('click', () => {
        handleBatchOperation('delete');
    });

    document.getElementById('batchActivateBtn').addEventListener('click', () => {
        handleBatchOperation('activate');
    });

    document.getElementById('batchDeactivateBtn').addEventListener('click', () => {
        handleBatchOperation('deactivate');
    });

    document.getElementById('clearSelectionBtn').addEventListener('click', () => {
        clearSelection();
    });
}

// 更新浮动操作栏
function updateFloatingActionBar() {
    const actionBar = document.getElementById('floatingActionBar');
    const selectedCount = document.getElementById('selectedCount');

    if (state.selectedModels.size > 0) {
        selectedCount.textContent = state.selectedModels.size;
        actionBar.style.display = 'block';
    } else {
        actionBar.style.display = 'none';
    }
}

// 清除选择
function clearSelection() {
    state.selectedModels.clear();
    document.getElementById('selectAll').checked = false;
    const checkboxes = document.querySelectorAll('tbody input[type="checkbox"]');
    checkboxes.forEach(checkbox => {
        checkbox.checked = false;
    });
    updateFloatingActionBar();
}

// 处理批量操作
async function handleBatchOperation(operation) {
    if (state.selectedModels.size === 0) return;

    const actionText = {
        delete: '删除',
        activate: '激活',
        deactivate: '禁用'
    };

    const actionVerb = {
        delete: '删除',
        activate: '激活',
        deactivate: '禁用'
    };

    if (!confirm(`确定要${actionVerb[operation]}${state.selectedModels.size}个模型吗？`)) {
        return;
    }

    try {
        showLoadingOverlay();

        const response = await api(`/api/admin/models/batch/${operation}`, {
            method: 'POST',
            body: JSON.stringify({ ids: Array.from(state.selectedModels) })
        });

        showAlert('model', 'success', `${actionText[operation]}成功`);
        clearSelection();
        loadModels();
    } catch (error) {
        showAlert('model', 'error', error.message);
    } finally {
        hideLoadingOverlay();
    }
}

// 显示加载遮罩
function showLoadingOverlay() {
    const overlay = document.createElement('div');
    overlay.className = 'loading-overlay';
    overlay.innerHTML = `
        <div class="loading-content">
            <div class="spinner"></div>
            <span>操作中...</span>
        </div>
    `;
    document.body.appendChild(overlay);
}

// 隐藏加载遮罩
function hideLoadingOverlay() {
    const overlay = document.querySelector('.loading-overlay');
    if (overlay) {
        overlay.remove();
    }
}

// 在模型渲染时添加选择框
function renderModels(models) {
    // ... 现有渲染逻辑 ...

    tbody.innerHTML = models.map(model => {
        // ... 现有渲染逻辑 ...

        return `
        <tr>
            <td>
                <input type="checkbox" data-model-id="${model.id}" class="form-checkbox">
            </td>
            <td>${id}</td>
            <!-- 其他列 -->
        </tr>
        `;
    }).join('');

    // 重新绑定事件
    initModelTableEvents();
}
```

#### 4. CSS 样式
```css
/* admin/index.html */
.floating-action-bar {
    position: fixed;
    bottom: 20px;
    left: 50%;
    transform: translateX(-50%);
    background: var(--bg-elevated);
    border: 1px solid var(--border);
    border-radius: var(--radius-lg);
    padding: var(--space-3) var(--space-4);
    box-shadow: var(--shadow-lg);
    display: flex;
    align-items: center;
    gap: var(--space-4);
    z-index: 100;
    transition: all var(--transition-base);
}

.action-bar-content {
    display: flex;
    align-items: center;
    gap: var(--space-4);
}

.selected-count {
    font-weight: 500;
    color: var(--text-primary);
}

.action-buttons {
    display: flex;
    gap: var(--space-2);
}

.loading-overlay {
    position: fixed;
    inset: 0;
    background: rgba(0, 0, 0, 0.5);
    backdrop-filter: blur(4px);
    z-index: 200;
    display: flex;
    align-items: center;
    justify-content: center;
}

.loading-content {
    background: var(--bg-elevated);
    border-radius: var(--radius-xl);
    padding: var(--space-5);
    display: flex;
    align-items: center;
    gap: var(--space-3);
    box-shadow: var(--shadow-xl);
}
```

## 安全考虑

1. **权限控制**：所有批量操作需要管理员权限
2. **事务处理**：确保原子性操作，失败时回滚
3. **数据验证**：检查模型存在性和关联数据
4. **防误操作**：确认对话框防止意外删除

## 测试计划

1. **单元测试**：
   - 批量删除API端点
   - 批量激活/禁用API端点
   - 事务回滚测试

2. **集成测试**：
   - 前端选择功能
   - 浮动操作栏交互
   - 批量操作流程

3. **边界测试**：
   - 空选择集
   - 大量模型选择
   - 网络错误处理

## 部署计划

1. 数据库迁移（软删除支持）
2. 后端API实现
3. 前端UI更新
4. 测试验证
5. 生产部署

## 依赖项

- GORM 1.25+（支持软删除）
- Gin 1.9+（路由和中间件）
- 现有前端框架（HTML/CSS/JS）

## 风险评估

- **数据一致性**：通过事务确保
- **用户体验**：浮动操作栏提供清晰反馈
- **性能**：批量操作可能影响性能，需要优化
- **兼容性**：与现有系统完全兼容

## 成功标准

1. 批量操作功能正常工作
2. 原子性保证
3. 用户界面友好
4. 错误处理完善
5. 性能可接受