package router

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"prueba4/models"
	"sync"
)

const (
	totalUsers = 10000 
	userPerRequest = 500 
	api = "https://randomuser.me/api?results=500"
)

func FetchUsers(ctx context.Context)([]models.User, error){
	
	req, err := http.NewRequestWithContext(ctx, "GET", api, nil)
	if err != nil {
		return nil, err 
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err 
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err 
	}

	var result models.ApiResponse 
	err = json.Unmarshal(body, &result)

	fmt.Println(result)

	if err != nil {
		return nil, err
	}
	if len(result.Results) == 0 {
		return nil, fmt.Errorf("error: respuesta vacía")
	}
	return result.Results, nil 
}

func Worker( ctx context.Context, jobs chan models.Job, results chan []models.User, wg *sync.WaitGroup){
	defer wg.Done()
	for job := range jobs {
		
		var users []models.User 
		var err error 

		for intento := 1; intento <= 3; intento++{
			users, err = FetchUsers(ctx)
			if err == nil && len(users) > 0{
				fmt.Printf("Worker completo trabajo %d", job.ID)
				results <- users 
				break 
			}
			fmt.Printf("Worker, reintentando trabajo %d", job.ID)
		}
		if err != nil || len(users)==0 {
			fmt.Println(err)
		}
	}
	
}