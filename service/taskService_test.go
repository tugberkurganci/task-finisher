package service

import (
	"fmt"
	"konzek-mid/mocks/repository"
	"konzek-mid/models"

	"reflect"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

var mockRepo *repository.MockTaskRepository
var service TaskService

func setupv1(t *testing.T) func() {
	ctrl := gomock.NewController(t)
	mockRepo = repository.NewMockTaskRepository(ctrl)
	service = NewTaskService(mockRepo)

	return func() {
		service = nil
		ctrl.Finish()
	}
}
func TestEnqueueTask_Success(t *testing.T) {
	// Test için hazırlıkları yap
	td := setupv1(t)
	defer td()

	// Mock repository'den beklenen değerlerin ayarlanması
	mockRepo.EXPECT().InsertTask(gomock.Any()).Return(models.Task{}, nil)

	// Servis fonksiyonunun çağrılması
	result, err := service.EnqueueTask(models.Task{})
	fmt.Println(result)
	// Hata kontrolü
	if err != nil {
		t.Error(err)
	}

	assert.Empty(t, result)
}
func TestGetTaskStatus_Success(t *testing.T) {
	// Test için hazırlıkları yap
	defer setupv1(t)()

	// Mock repository'den beklenen değerlerin ayarlanması
	expectedTask := models.Task{ID: 1, Payload: "Test Content"}
	mockRepo.EXPECT().GetTaskByID(1).Return(expectedTask, nil)

	// Servis fonksiyonunun çağrılması
	task, err := service.GetTaskStatus(1)

	// Hata kontrolü
	if err != nil {
		t.Errorf("Beklenmeyen hata: %v", err)
	}

	// Doğrulama
	assert.Equal(t, expectedTask, task, "Dönen görev beklenen görevle eşleşmiyor")
}
func TestScheduleTasks_Success(t *testing.T) {
	// Test için hazırlıkları yap
	defer setupv1(t)()

	// TaskServiceImpl tipinde bir değişken oluştur
	taskService := service.(*TaskServiceImpl)

	// Görevlerin takip edileceği slice
	var scheduledTasks []models.Task

	// Mock repository'den beklenen değerlerin ayarlanması
	mockRepo.EXPECT().GetPastScheduledTasks().Return([]models.Task{
		{ID: 1, Payload: "Task 1", Priority: 2},
		{ID: 2, Payload: "Task 2", Priority: 1},
		// Daha fazla görev eklenebilir
	}, nil).AnyTimes() // Her çağrıda aynı görevleri dönmesini sağlıyoruz

	// Servis fonksiyonunun çağrılması
	taskService.ScheduleTasks()

	// Test için beklemeyi başlat
	go func() {
		time.Sleep(15 * time.Second) // Testi bekleme süresi kadar uzatalım
		close(taskService.TaskQueue) // TaskQueue kanalını kapat
	}()

	// TaskQueue kanalından görevleri al
	for task := range taskService.TaskQueue {
		scheduledTasks = append(scheduledTasks, task)
	}

	// Hata kontrolü ve doğrulama
	// ScheduledTasks slice'ında beklenen görev sırası olmalıdır.
	expectedTasks := []models.Task{
		{ID: 2, Payload: "Task 2", Priority: 1},
		{ID: 1, Payload: "Task 1", Priority: 2},
		// Daha fazla görev varsa, sıraya eklenmelidir.
	}
	if !reflect.DeepEqual(expectedTasks, scheduledTasks) {
		t.Errorf("Beklenen görev sırası alınan görev sırasıyla eşleşmiyor")
	}
}

func TestStartWorkers(t *testing.T) {
	// Test için hazırlıkları yap
	defer setupv1(t)()

	// TaskServiceImpl tipinde bir değişken oluştur
	taskService := service.(*TaskServiceImpl)
	// Servis fonksiyonunun çağrılması
	taskService.StartWorkers()

	// Hata kontrolü
	// Beklenen davranışlar test edilebilir
	time.Sleep(5 * time.Second) // Testi bekleme süresi kadar uzatalım

	// İşçi sayısını kontrol et
	expectedWorkerCount := 3
	actualWorkerCount := len(taskService.Workers)
	if actualWorkerCount != expectedWorkerCount {
		t.Errorf("Beklenen işçi sayısıyla eşleşmeyen çalışan işçi sayısı: %d", actualWorkerCount)
	}

}
