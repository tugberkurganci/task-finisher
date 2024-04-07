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

func TestEndToEndTaskStatusCheckWıthOutInterval(t *testing.T) {
	// Başlama zamanını kaydet

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

	// HTTP isteği gönderme
	req, err := http.NewRequest(http.MethodPost, "http://localhost:8080/tasks", bytes.NewBuffer(reqBody))
	if err != nil {
		t.Fatalf("İstek oluşturulurken bir hata oluştu: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")

	// İstek gönderme
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		t.Fatalf("İstek gönderilirken bir hata oluştu: %v", err)
	}
	defer resp.Body.Close()

	// İstek sonucunu kontrol etme
	if resp.StatusCode != http.StatusCreated {
		t.Errorf("Beklenen durum kodu 201 değil, alınan: %d", resp.StatusCode)
	}

	// İstekin tamamlanması için biraz bekleme süresi ekleyebilirsiniz
	time.Sleep(20 * time.Second)

	// Task'ın statusünü kontrol etme
	req2, err := http.NewRequest(http.MethodGet, "http://localhost:8080/tasks/343", nil)
	if err != nil {
		t.Fatalf("İstek oluşturulurken bir hata oluştu: %v", err)
	}

	resp2, err := client.Do(req2)
	if err != nil {
		t.Fatalf("İstek gönderilirken bir hata oluştu: %v", err)
	}
	defer resp2.Body.Close()

	// İstek sonucunu kontrol etme
	if resp2.StatusCode != http.StatusOK {
		t.Errorf("Beklenen durum kodu 200 değil, alınan: %d", resp2.StatusCode)
	}

	// responseBody'yi oku ve çözümle
	body, err := ioutil.ReadAll(resp2.Body)
	if err != nil {
		t.Fatalf("Yanıt gövdesini okurken bir hata oluştu: %v", err)
	}
	fmt.Println("Yanıt gövdesi:", string(body))

	var responseBody map[string]interface{}
	err = json.Unmarshal(body, &responseBody)
	if err != nil {
		t.Fatalf("Yanıt çözümlenirken bir hata oluştu: %v", err)
	}

	// "completed" alanını kontrol etme
	completed, ok := responseBody["Completed"].(bool)
	if !ok {
		t.Errorf("Yanıtta 'completed' alanı beklenen formatta değil")
	}

	// "completed" değerini kontrol etme
	if completed {
		fmt.Println("completed")
	} else {
		fmt.Println("not completed")
	}

}
func TestEndToEndTaskStatusCheckWithInterval(t *testing.T) {
	// Başlama zamanını kaydet

	// Payload oluşturma
	payload := struct {
		Payload  string
		Deadline time.Time
		Interval string
	}{
		Payload:  "Test Content",
		Deadline: time.Now(), // 2 saat sonrasına deadline belirleyin
		Interval: "1h",
	}

	// JSON formatına çevirme
	reqBody, err := json.Marshal(payload)
	if err != nil {
		t.Fatalf("JSON formatına dönüştürülürken bir hata oluştu: %v", err)
	}

	// HTTP isteği gönderme
	req, err := http.NewRequest(http.MethodPost, "http://localhost:8080/tasks", bytes.NewBuffer(reqBody))
	if err != nil {
		t.Fatalf("İstek oluşturulurken bir hata oluştu: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")

	// İstek gönderme
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		t.Fatalf("İstek gönderilirken bir hata oluştu: %v", err)
	}
	defer resp.Body.Close()

	// İstek sonucunu kontrol etme
	if resp.StatusCode != http.StatusCreated {
		t.Errorf("Beklenen durum kodu 201 değil, alınan: %d", resp.StatusCode)
	}

	// İstekin tamamlanması için biraz bekleme süresi ekleyebilirsiniz
	time.Sleep(20 * time.Second)

	// Task'ın statusünü kontrol etme
	req2, err := http.NewRequest(http.MethodGet, "http://localhost:8080/tasks/344", nil)
	if err != nil {
		t.Fatalf("İstek oluşturulurken bir hata oluştu: %v", err)
	}

	resp2, err := client.Do(req2)
	if err != nil {
		t.Fatalf("İstek gönderilirken bir hata oluştu: %v", err)
	}
	defer resp2.Body.Close()

	// İstek sonucunu kontrol etme
	if resp2.StatusCode != http.StatusOK {
		t.Errorf("Beklenen durum kodu 200 değil, alınan: %d", resp2.StatusCode)
	}

	// responseBody'yi oku ve çözümle
	body, err := ioutil.ReadAll(resp2.Body)
	if err != nil {
		t.Fatalf("Yanıt gövdesini okurken bir hata oluştu: %v", err)
	}
	fmt.Println("Yanıt gövdesi:", string(body))

	var responseBody map[string]interface{}
	err = json.Unmarshal(body, &responseBody)
	if err != nil {
		t.Fatalf("Yanıt çözümlenirken bir hata oluştu: %v", err)
	}

	// "completed" alanını kontrol etme
	completed, ok := responseBody["Completed"].(bool)
	if !ok {
		t.Errorf("Yanıtta 'completed' alanı beklenen formatta değil")
	}

	// "completed" değerini kontrol etme
	if completed {
		t.Errorf("completed olmamalı: %v", completed)
	} else {
		fmt.Println("not completed")
	}

	// Deadline'ı kontrol etme
	deadline, err := time.Parse(time.RFC3339Nano, responseBody["Deadline"].(string))
	if err != nil {
		t.Fatalf("Yanıtta 'deadline' alanı beklenen formatta değil: %v", err)
	}
	now := time.Now()
	if deadline.Before(now) {
		t.Errorf("Deadline zaten geçmiş: %v", deadline)
	} else {
		fmt.Println(deadline)
	}

}
