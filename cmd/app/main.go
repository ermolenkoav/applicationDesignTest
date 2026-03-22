// Ниже реализован сервис бронирования номеров в отеле. В предметной области
// выделены два понятия: Order — заказ, который включает в себя даты бронирования
// и контакты пользователя, и RoomAvailability — количество свободных номеров на
// конкретный день.
//
// Задание:
// - провести рефакторинг кода с выделением слоев и абстракций
// - применить best-practices там где это имеет смысл
// - исправить имеющиеся в реализации логические и технические ошибки и неточности
package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"applicationDesignTest/internal/api/rest"
	"applicationDesignTest/internal/logg"
	"applicationDesignTest/internal/repo"
	"applicationDesignTest/internal/service"
)

func main() {
	logg.Info("up and running!")

	r := repo.New()
	svc := service.NewBookingService(r)
	srv := rest.NewServer(svc)

	go func() {
		if err := srv.ListenAndServe(); err != nil {
			logg.Fatal("listen and serve: %v", err)
		}
	}()

	done := make(chan os.Signal, 1)
	signal.Notify(done, syscall.SIGINT, syscall.SIGTERM)
	<-done

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		logg.Fatal("api shutdown: %v", err)
	}

	logg.Info("graceful shutdown!")
}
