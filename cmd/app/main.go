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

	"applicationDesignTest/internal/logg"
	"applicationDesignTest/internal/repo"
	"applicationDesignTest/internal/rest"
	"applicationDesignTest/internal/service"
)

func main() {
	ctx := context.Background()
	logg.Info("up and running!")

	repo := repo.NewPersistent()

	service := service.NewBookingService(repo)

	api := rest.NewServe(service)
	go func() {
		if err := api.ListenAndServe(); err != nil {
			logg.Fatal("listen and serve: ", err)
		}
	}()

	done := make(chan os.Signal, 1)
	signal.Notify(done, syscall.SIGINT, syscall.SIGTERM)
	<-done

	if err := api.Shutdown(ctx); err != nil {
		logg.Fatal("api shutdown: ", err)
	}

	logg.Info("graceful shutdown!")
}
