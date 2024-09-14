package utils

import (
	"container/heap"
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"sync"
)

var (
	WORKERS    = 1  //количество рабочих
	WORKERSCAP = 50 //размер очереди каждого рабочего
)

type Worker struct {
	files   chan *multipart.FileHeader // канал для файлов
	pending int                        // количество оставшихся изображений
	index   int                        // позиция в куче
	wg      *sync.WaitGroup
}

func generator(out chan *multipart.FileHeader, files []*multipart.FileHeader) {
	for _, fl := range files {
		out <- fl
	}
	out <- &multipart.FileHeader{Size: 0}
}

func (w *Worker) work(done chan *Worker, out chan interface{}) {
	for {
		file := <-w.files
		w.wg.Add(1)
		serve(file, out)
		w.wg.Done()
		done <- w
	}
}

type Pool []*Worker

func (p Pool) Less(i, j int) bool { return p[i].pending < p[j].pending }
func (p Pool) Len() int           { return len(p) }
func (p Pool) Swap(i, j int) {
	if i >= 0 && i < len(p) && j >= 0 && j < len(p) {
		p[i], p[j] = p[j], p[i]
		p[i].index, p[j].index = i, j
	}
}

func (p *Pool) Push(x interface{}) {
	n := len(*p)
	worker := x.(*Worker)
	worker.index = n
	*p = append(*p, worker)
}

func (p *Pool) Pop() interface{} {
	old := *p
	n := len(old)
	item := old[n-1]
	item.index = -1
	*p = old[:n-1]
	return item
}

type Balancer struct {
	pool     Pool                       // куча рабочих
	done     chan *Worker               // канал уведомления о завершении для рабочих
	requests chan *multipart.FileHeader // канал для получения новых заданий
	flowCtrl chan bool                  // канал для PMFC (Poor Man's Flow Control)
	queue    int                        // количество незавершенных заданий, которые были переданы рабочим
	wg       *sync.WaitGroup
}

func (b *Balancer) init(in chan *multipart.FileHeader, out chan interface{}) {
	b.requests = make(chan *multipart.FileHeader)
	b.flowCtrl = make(chan bool)
	b.done = make(chan *Worker)
	b.wg = new(sync.WaitGroup)

	// Запускаем Flow contol
	go func() {
		for {
			b.requests <- <-in // получаем задание и пересылаем его
			<-b.flowCtrl       // ждём получения подтверждения
		}
	}()

	// Инициализируем кучу и создаём рабочих
	heap.Init(&b.pool)
	for i := 0; i < WORKERS; i++ {
		w := &Worker{
			files:   make(chan *multipart.FileHeader, WORKERSCAP),
			pending: 0,
			index:   0,
			wg:      b.wg,
		}
		go w.work(b.done, out) // запустили рабочего
		heap.Push(&b.pool, w)  // и отправили его в кучу
	}
}

// quit - канал уведомлений от главного цикла
func (b *Balancer) balance(quit chan bool) {
	lastJobs := false
	for {
		select {
		case <-quit: // пршило указание на остановку работы
			b.wg.Wait()  // ждём всех рабочих
			quit <- true // отправляем сигнал что закончили
		case file := <-b.requests: // Получено новое задание от Flow controller
			if file.Size > 0 { // если полученный файл не nil
				b.dispatch(file) // dispatch отправляет файл рабочим
			} else {
				lastJobs = true // поднимаем флаг завершения
			}
		case w := <-b.done: // пришло уведомление что рабочий окончил работу
			b.completed(w)
			if lastJobs {
				if w.pending == 0 { // если у рабочего кончились задания, то удаляем его из кучи
					heap.Remove(&b.pool, w.index)
				}
				if len(b.pool) == 0 { // если куча стала пуста
					quit <- true
				}
			}
		}
	}
}

// функция отправки задания
func (b *Balancer) dispatch(file *multipart.FileHeader) {
	w := heap.Pop(&b.pool).(*Worker) // берем из кучи самого незагруженного работника
	w.files <- file                  // отправляем ему задания
	w.pending++                      // добавляем ему "весу"
	heap.Push(&b.pool, w)            // и отправляем обратно его в кучу
	if b.queue++; b.queue < WORKERS*WORKERSCAP {
		b.flowCtrl <- true
	}
}

func (b *Balancer) completed(w *Worker) {
	w.pending--
	heap.Remove(&b.pool, w.index)
	heap.Push(&b.pool, w)
	if b.queue--; b.queue == WORKERS*WORKERSCAP-1 {
		b.flowCtrl <- true
	}
}

func serve(file *multipart.FileHeader, out chan interface{}) {
	src, err := file.Open()
	if err != nil {
		fmt.Printf("Err1 - %v", err)
		out <- err
		return
	}
	defer src.Close()
	dst, err := os.Create(file.Filename)
	if err != nil {
		out <- err
		return
	}
	defer dst.Close()
	if _, err = io.Copy(dst, src); err != nil {
		out <- err
		return
	}
	out <- file.Filename
}

func ServeFiles(files []*multipart.FileHeader) ([]string, error) {
	filesChan := make(chan *multipart.FileHeader)
	quit := make(chan bool)
	balancer := new(Balancer)
	outChan := make(chan interface{})
	balancer.init(filesChan, outChan)

	go balancer.balance(quit)
	go generator(filesChan, files)

	returnValues := make([]string, 0)
	for {
		select {
		case photoUrl := <-outChan:
			if strUrl, ok := photoUrl.(string); ok {
				returnValues = append(returnValues, strUrl)
			} else {
				<-quit
			}
		case <-quit:
			return returnValues, nil // Fix
		}
	}
}
