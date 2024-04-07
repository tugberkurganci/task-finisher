package app_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"testing"
	"time"
)

func TestConcurrentTaskStatusCheck(t *testing.T) {
	// Goroutine'lerin sayısı
	numRequests := 100

	// Sonuçları toplamak için bir kanal oluştur
	results := make(chan bool, numRequests)

	for i := 0; i < numRequests; i++ {
		go func() {
			// Payload oluşturma
			payload := struct {
				Payload  string
				Deadline time.Time
			}{
				Payload:  "Test Content",
				Deadline: time.Now(),
			}

			// JSON formatına çevirme
			reqBody, err := json.Marshal(payload)
			if err != nil {
				t.Fatalf("JSON formatına dönüştürülürken bir hata oluştu: %v", err)
			}

			// Yeni bir istek oluştur
			req, err := http.NewRequest(http.MethodPost, "http://localhost:8080/tasks", bytes.NewBuffer(reqBody))
			if err != nil {
				t.Fatalf("İstek oluşturulurken bir hata oluştu: %v", err)
			}
			req.Header.Set("Content-Type", "application/json")

			// İstek gönderme
			client := &http.Client{}
			resp, err := client.Do(req)
			if err != nil {
				t.Errorf("İstek gönderilirken bir hata oluştu: %v", err)
				results <- false
				return
			}
			defer resp.Body.Close()

			// İstek sonucunu kontrol etme
			if resp.StatusCode != http.StatusCreated {
				t.Errorf("Beklenen durum kodu 201 değil, alınan: %d", resp.StatusCode)
				results <- false
				return
			}

			results <- true
		}()
	}

	// Gorutinlerin tamamlanmasını bekleyin ve sonuçları değerlendirin
	for i := 0; i < numRequests; i++ {
		if !<-results {
			fmt.Println("Bazı işlemler tamamlanamadı.")
			return
		}
	}

	fmt.Println("Tüm işlemler tamamlandı.")

	// Bir dakika bekleyin
	time.Sleep(1 * time.Minute)

	// Task'ın statusünü kontrol etme
	// Goroutine'lerin sayısı

	// Gorutinlerle eşzamanlı istek gönderme
	for i := 647; i < 747; i++ {
		go func(taskID int) {
			// HTTP isteği gönderme
			resp, err := http.Get(fmt.Sprintf("http://localhost:8080/tasks/%d", taskID))
			if err != nil {
				t.Errorf("İstek gönderilirken bir hata oluştu: %v", err)
				results <- false
				return
			}
			defer resp.Body.Close()

			// İstek sonucunu kontrol etme
			if resp.StatusCode != http.StatusOK {
				t.Errorf("Beklenen durum kodu 200 değil, alınan: %d", resp.StatusCode)
				results <- false
				return
			}

			// responseBody'yi oku ve çözümle
			body, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				t.Errorf("Yanıt gövdesini okurken bir hata oluştu: %v", err)
				results <- false
				return
			}

			var responseBody map[string]interface{}
			err = json.Unmarshal(body, &responseBody)
			if err != nil {
				t.Errorf("Yanıt çözümlenirken bir hata oluştu: %v", err)
				results <- false
				return
			}

			// "completed" alanını kontrol etme
			completed, ok := responseBody["Completed"].(bool)
			if !ok {
				t.Errorf("Yanıtta 'completed' alanı beklenen formatta değil")
				results <- false
				return
			}

			results <- completed
		}(i)
	}

	// Gorutinlerin tamamlanmasını bekleyin ve sonuçları değerlendirin
	for i := 0; i < numRequests; i++ {

		if !<-results {
			fmt.Println("Bazı işlemler tamamlanamadı.")
			fmt.Println(i)
			return
		}
	}

	fmt.Println("Tüm işlemler tamamlandı.")
}
