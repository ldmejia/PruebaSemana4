package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	router "prueba4/Router"
	"prueba4/models"
	"sync"
	"time"
)

const (
	totalUsers = 10000 
	userPerRequest = 500 
	workers = 5
	totalJobs = totalUsers/userPerRequest
)

func HandleUsers(w http.ResponseWriter, r *http.Request){
	ctx, cancel := context.WithTimeout(r.Context(), 45*time.Second)
	defer cancel()

	jobs := make(chan models.Job, totalJobs)
	results := make(chan []models.User, totalJobs)

	var wg sync.WaitGroup

	for i:=1; i<= workers; i++{
		wg.Add(1)
		go router.Worker(ctx, jobs, results, &wg)	
	}

	go func() {
		for i := 0; i < totalJobs; i++{
			select {
			case <- ctx.Done():
			case jobs <- models.Job{ID: i}:
			}
		}
		close(jobs)
	}()

	go func(){
		wg.Wait()
		close(results)
	}()

	var allUsers []models.User 

	Loop:
	for {
		select{
	case <-ctx.Done():
		http.Error(w, "Tiempo excedido", http.StatusRequestTimeout)
		return
	case resultados, ok := <- results: 
		if !ok {
			break Loop
		}
			allUsers = append(allUsers, resultados...)
			if len(allUsers) >= totalUsers {
				break Loop
			}
		}
	}

	if len(allUsers) < totalUsers{
		http.Error(w, fmt.Sprintf("No hay 10,000 usuarios, solo %d", len(allUsers)), http.StatusServiceUnavailable)
		return 
	}

allUsers = allUsers[:totalUsers]
w.Header().Set("Content-Type", "application/json")

json.NewEncoder(w).Encode(Responder(allUsers))
}

func Responder(usuarios []models.User)(map[string]interface{}){
	
	respuesta := map[string]interface{}{
		"info": map[string]interface{}{
			"type": "success", 
			"totalUsers": len(usuarios), 
		}, 
		"results": usuarios, 
	}
	return respuesta
}