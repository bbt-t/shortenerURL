package handler

import (
	"context"
	"fmt"
	"github.com/bbt-t/shortenerURL/internal/adapter/storage"
	"github.com/bbt-t/shortenerURL/pkg"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/gofrs/uuid"
)

type DeleteHandler struct {
	repo           storage.DatabaseRepository
	delURLsCh      chan Job
	expireDuration time.Duration
}

type Job struct {
	URL    string
	UserID uuid.UUID
}

func NewDeleteHandler(repo storage.DatabaseRepository, delURLsChBufSize int) *DeleteHandler {
	h := &DeleteHandler{
		repo:           repo,
		delURLsCh:      make(chan Job, delURLsChBufSize),
		expireDuration: 100 * time.Millisecond,
	}
	h.Loop()
	return h
}

func (d *DeleteHandler) Loop() {
	go func() {
		jobs := make([]Job, 0)
		ticker := time.NewTicker(d.expireDuration)
		for {
			select {
			case job, ok := <-d.delURLsCh:
				if !ok {
					return
				}
				jobs = append(jobs, job)
			case <-ticker.C:
				d.DeleteURLs(jobs)
				jobs = make([]Job, 0)
			}
		}
	}()
}

func (d *DeleteHandler) DeleteURLs(jobsToDelete []Job) {
	if len(jobsToDelete) == 0 {
		return
	}
	jobsByUser := make(map[uuid.UUID][]string)
	for _, job := range jobsToDelete {
		jobsByUser[job.UserID] = append(jobsByUser[job.UserID], job.URL)
	}
	for userID, userJobs := range jobsByUser {
		go func(userID uuid.UUID, userJobs []string) {
			err := d.repo.DelURLArray(context.Background(), userID, userJobs)
			if err != nil {
				log.Printf("Error while deleting urls: %v\n", err)
			}
		}(userID, userJobs)
	}
}

func (d *DeleteHandler) Handle(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	payload, errBody := io.ReadAll(r.Body)
	if errBody != nil {
		http.Error(
			w,
			fmt.Sprintf("Incorrent request body: %s", payload),
			http.StatusBadRequest,
		)
		return
	}
	userID, _ := uuid.FromString(fmt.Sprintf("%v", r.Context().Value("user_id")))

	ctx := context.Background()

	ids := pkg.ConvertStrToSlice(string(payload))
	go d.repo.DelURLArray(ctx, userID, ids)

	w.WriteHeader(http.StatusAccepted)
}

//func (s ShortenerHandler) deleteURL(w http.ResponseWriter, r *http.Request) {
//	defer r.Body.Close()
//
//	payload, errBody := io.ReadAll(r.Body)
//	if errBody != nil {
//		http.Error(
//			w,
//			fmt.Sprintf("Incorrent request body: %s", payload),
//			http.StatusBadRequest,
//		)
//		return
//	}
//	userID, _ := uuid.FromString(fmt.Sprintf("%v", r.Context().Value("user_id")))
//
//	ctx := context.Background()
//	go s.s.DelURLArray(ctx, userID, payload)
//
//	w.WriteHeader(http.StatusAccepted)
//}
