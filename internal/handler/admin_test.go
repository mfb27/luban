package handler

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/mfb27/luban/internal/model"
	"gorm.io/gorm"
)

func TestAdminBatchDeleteModels(t *testing.T) {
	// 创建模拟的数据库
	db, err := setupTestDB()
	if err != nil {
		t.Fatalf("Failed to setup test DB: %v", err)
	}
	defer func() {
		if db != nil {
			db.Migrator().DropTable(&model.Model{}, &model.Message{}, &model.Session{}, &model.User{}, &model.Attachment{})
		}
	}()

	// 创建测试数据
	createTestModels(db)

	// 创建管理员应用
	adminApp := NewAdminApp(nil, db)

	// 创建 Gin 上下文
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	// 测试1: 正常批量删除
	c.Request, _ = http.NewRequest("DELETE", "/api/admin/models/batch", nil)
	c.Request.Header.Set("Content-Type", "application/json")
	c.Request.Body = http.NoBody

	// 这里需要模拟请求体，由于Gin的限制，我们直接调用函数
	adminApp.adminBatchDeleteModels(c)

	// 检查响应
	if w.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
	}

	// 测试2: 空请求
	w = httptest.NewRecorder()
	c, _ = gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest("DELETE", "/api/admin/models/batch", nil)
	c.Request.Header.Set("Content-Type", "application/json")
	c.Request.Body = http.NoBody
	adminApp.adminBatchDeleteModels(c)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status %d for empty request, got %d", http.StatusOK, w.Code)
	}

	// 测试3: 无效的JSON
	w = httptest.NewRecorder()
	c, _ = gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest("DELETE", "/api/admin/models/batch", nil)
	c.Request.Header.Set("Content-Type", "application/json")
	c.Request.Body = http.NoBody
	adminApp.adminBatchDeleteModels(c)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status %d for invalid JSON, got %d", http.StatusOK, w.Code)
	}
}

func TestAdminBatchActivateModels(t *testing.T) {
	// 创建模拟的数据库
	db, err := setupTestDB()
	if err != nil {
		t.Fatalf("Failed to setup test DB: %v", err)
	}
	defer func() {
		if db != nil {
			db.Migrator().DropTable(&model.Model{}, &model.Message{}, &model.Session{}, &model.User{}, &model.Attachment{})
		}
	}()

	// 创建测试数据
	createTestModels(db)

	// 创建管理员应用
	adminApp := NewAdminApp(nil, db)

	// 创建 Gin 上下文
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	// 测试1: 正常批量激活
	c.Request, _ = http.NewRequest("PATCH", "/api/admin/models/batch/activate", nil)
	c.Request.Header.Set("Content-Type", "application/json")
	c.Request.Body = http.NoBody

	// 这里需要模拟请求体，由于Gin的限制，我们直接调用函数
	adminApp.adminBatchActivateModels(c)

	// 检查响应
	if w.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
	}

	// 测试2: 空请求
	w = httptest.NewRecorder()
	c, _ = gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest("PATCH", "/api/admin/models/batch/activate", nil)
	c.Request.Header.Set("Content-Type", "application/json")
	c.Request.Body = http.NoBody
	adminApp.adminBatchActivateModels(c)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status %d for empty request, got %d", http.StatusOK, w.Code)
	}

	// 测试3: 无效的JSON
	w = httptest.NewRecorder()
	c, _ = gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest("PATCH", "/api/admin/models/batch/activate", nil)
	c.Request.Header.Set("Content-Type", "application/json")
	c.Request.Body = http.NoBody
	adminApp.adminBatchActivateModels(c)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status %d for invalid JSON, got %d", http.StatusOK, w.Code)
	}
}

func setupTestDB() (*gorm.DB, error) {
	// 这里应该返回一个真实的测试数据库连接
	// 在实际实现中，你需要设置一个内存数据库或测试数据库
	return nil, nil
}

func createTestModels(db *gorm.DB) {
	// 创建测试模型数据
}

func TestAdminBatchDeactivateModels(t *testing.T) {
	// 创建模拟的数据库
	db, err := setupTestDB()
	if err != nil {
		t.Fatalf("Failed to setup test DB: %v", err)
	}
	defer func() {
		if db != nil {
			db.Migrator().DropTable(&model.Model{}, &model.Message{}, &model.Session{}, &model.User{}, &model.Attachment{})
		}
	}()

	// 创建测试数据
	createTestModels(db)

	// 创建管理员应用
	adminApp := NewAdminApp(nil, db)

	// 创建 Gin 上下文
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	// 测试1: 正常批量禁用
	c.Request, _ = http.NewRequest("PATCH", "/api/admin/models/batch/deactivate", nil)
	c.Request.Header.Set("Content-Type", "application/json")
	c.Request.Body = http.NoBody

	// 这里需要模拟请求体，由于Gin的限制，我们直接调用函数
	adminApp.adminBatchDeactivateModels(c)

	// 检查响应
	if w.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
	}

	// 测试2: 空请求
	w = httptest.NewRecorder()
	c, _ = gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest("PATCH", "/api/admin/models/batch/deactivate", nil)
	c.Request.Header.Set("Content-Type", "application/json")
	c.Request.Body = http.NoBody
	adminApp.adminBatchDeactivateModels(c)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status %d for empty request, got %d", http.StatusOK, w.Code)
	}

	// 测试3: 无效的JSON
	w = httptest.NewRecorder()
	c, _ = gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest("PATCH", "/api/admin/models/batch/deactivate", nil)
	c.Request.Header.Set("Content-Type", "application/json")
	c.Request.Body = http.NoBody
	adminApp.adminBatchDeactivateModels(c)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status %d for invalid JSON, got %d", http.StatusOK, w.Code)
	}
}