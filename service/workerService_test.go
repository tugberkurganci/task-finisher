package service

import (
	"fmt"
	"konzek-mid/mocks/repository"
	"konzek-mid/models"

	"testing"
	"time"

	"github.com/golang/mock/gomock"
)

// MockTaskRepository oluşturmak için yardımcı fonksiyon
func setup(t *testing.T) (*gomock.Controller, *repository.MockTaskRepository) {
	ctrl := gomock.NewController(t)
	mockRepo := repository.NewMockTaskRepository(ctrl)
	return ctrl, mockRepo
}

func TestWorker_ProcessWithRetry(t *testing.T) {
	// Test için hazırlıkları yap
	ctrl, mockRepo := setup(t)
	defer ctrl.Finish()

	// Worker oluştur
	taskQueue := make(chan models.Task, 1)
	complete := make(chan bool, 1)
	worker := NewWorker(1, taskQueue, complete, mockRepo)

	// Görev verisi oluştur
	task := models.Task{
		ID:         1,
		Payload:    "Test payload",
		MaxRetries: 3,
	}

	// Mock repository'den beklenen çağrıları ayarla
	mockRepo.EXPECT().MarkTaskCompleted(task).Return(nil).Times(1)

	// Worker'ı başlat
	worker.Start()

	// Worker'a görevi gönder
	taskQueue <- task

	// Worker'ın işlemesi için biraz bekleyelim
	time.Sleep(1 * time.Second)

	// Assertler
	select {
	case <-complete:
		// Görev başarıyla tamamlandı
	case <-time.After(5 * time.Second):
		t.Error("Görev işlenemedi")
	}
}

// waiting fail
func TestWorker_ProcessWithRetry_MaxRetriesExceeded(t *testing.T) {
	// Test için hazırlıkları yap
	ctrl, mockRepo := setup(t)
	defer ctrl.Finish()

	// Worker oluştur
	taskQueue := make(chan models.Task, 1)
	complete := make(chan bool, 1)
	worker := NewWorker(1, taskQueue, complete, mockRepo)

	// Görev verisi oluştur
	task := models.Task{
		ID:         1,
		Payload:    "Test payload",
		MaxRetries: 3,
	}

	// Mock repository'den beklenen çağrıları ayarla

	mockRepo.EXPECT().MarkTaskCompleted(task).Return(fmt.Errorf("Task failed")).Times(4)

	// Worker'ı başlat
	worker.Start()

	// Worker'a görevi gönder
	taskQueue <- task

	// Worker'ın işlemesi için biraz bekleyelim
	time.Sleep(1 * time.Second)

	// Assertler
	select {
	case <-complete:
		t.Error("Görevin işlenmesi beklenenden fazla sürdü.")
	default:
		// Görev tamamlanmadı
	}
}

func TestWorker_ProcessWithRetry_Success(t *testing.T) {
	// Test için hazırlıkları yap
	ctrl, mockRepo := setup(t)
	defer ctrl.Finish()

	// Worker oluştur
	taskQueue := make(chan models.Task, 1)
	complete := make(chan bool, 1)
	worker := NewWorker(1, taskQueue, complete, mockRepo)

	// Görev verisi oluştur
	task := models.Task{
		ID:         1,
		Payload:    "Test payload",
		MaxRetries: 3,
	}

	// Mock repository'den beklenen çağrıları ayarla
	mockRepo.EXPECT().MarkTaskCompleted(task).Return(nil)

	// Worker'ı başlat
	worker.Start()

	// Worker'a görevi gönder
	taskQueue <- task

	// Worker'ın işlemesi için biraz bekleyelim
	time.Sleep(1 * time.Second)

	// Assertler
	select {
	case <-complete:
		// Görev tamamlandı, başarılı test
	default:
		t.Error("Görevin işlenmesi beklenenden uzun sürdü veya tamamlanmadı.")
	}
}
func TestWorker_CompleteTask_Success(t *testing.T) {
	// Test için hazırlıkları yap
	ctrl, mockRepo := setup(t)
	defer ctrl.Finish()

	// Mock repository'den beklenen çağrıları ayarla
	task := models.Task{ID: 1}
	mockRepo.EXPECT().MarkTaskCompleted(task).Return(nil)

	// Worker oluştur
	worker := NewWorker(1, make(chan models.Task), make(chan bool), mockRepo)

	// Görevi tamamla
	err := worker.CompleteTask(task)

	// Hata kontrolü
	if err != nil {
		t.Errorf("Beklenmedik bir hata oluştu: %v", err)
	}
}

func TestWorker_CompleteTask_Error(t *testing.T) {
	// Test için hazırlıkları yap
	ctrl, mockRepo := setup(t)
	defer ctrl.Finish()

	// Mock repository'den beklenen çağrıları ayarla
	task := models.Task{ID: 1}
	mockRepo.EXPECT().MarkTaskCompleted(task).Return(fmt.Errorf("Task completion failed"))

	// Worker oluştur
	worker := NewWorker(1, make(chan models.Task), make(chan bool), mockRepo)

	// Görevi tamamla
	err := worker.(*WorkerImpl).CompleteTask(task)

	// Hata kontrolü
	if err == nil {
		t.Error("Hata bekleniyordu ancak oluşmadı.")
	}
}
