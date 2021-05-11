// 模拟键盘输入Demo
package main

import (
"fmt"
"log"
"sync"
"time"
)

type ProcManager struct {
	// process manager
	pid int
	mp map[int]*Process
	mutex sync.Mutex
	k *Keyboard
	lastid int
}

func NewProcManager(k *Keyboard) *ProcManager {
	return &ProcManager{k: k, lastid: -1, mp: make(map[int]*Process, 0)}
}

const (
	RUNNING = iota
	DEAD
)

type Process struct {
	status int
	id int
	buf []byte
	end bool
	sigint bool
	readable bool
	mutex sync.Mutex
	keyboard *Keyboard
}

func (pm *ProcManager) NewProcess() *Process {
	p := &Process{status: RUNNING, id: pm.pid}
	pm.pid++
	pm.mp[p.id] = p
	go p.Run()
	return p
}

func (pm *ProcManager) Check() {
	for _, v := range pm.mp {
		v.mutex.Lock()
		id := v.id
		status := v.status
		v.mutex.Unlock()
		log.Println(id, status)
	}
}

func (pm *ProcManager) RemoveDead() {
	var rm []int
	for k, v := range pm.mp {
		v.mutex.Lock()
		if v.status == DEAD {
			rm = append(rm, k)
		}
		v.mutex.Unlock()
	}
	for _, item := range rm {
		delete(pm.mp, item)
	}
}

func (pm *ProcManager) Swit(id int) {
	process, ok := pm.mp[id]
	if !ok {
		log.Printf("No %d process exists", id)
		return
	}
	process.SetKeyboard(pm.k)
}

func (p *Process) Run()  {
	log.Printf("Process %v running\n", p.id)
	for true {
		time.Sleep(time.Millisecond * 300)
		p.Read()
		p.mutex.Lock()
		if p.readable {
			tmp := p.buf
			p.buf = make([]byte, 0)
			tmpend := p.end
			p.end = false
			tmpint := p.sigint
			p.sigint = false
			p.readable = false
			p.mutex.Unlock()
			// handle the data
			if tmpint {
				p.mutex.Lock()
				p.status = DEAD
				p.mutex.Unlock()
				log.Printf("Process %v Interrupted\n", p.id)
				return
			}
			log.Print(tmp)
			log.Print(string(tmp))
			if tmpend {
				log.Printf("Process %v Input End\n", p.id)
			}
			continue
		} else {
			p.mutex.Unlock()
		}
	}
}

func (p *Process) SetKeyboard(k *Keyboard)  {
	p.mutex.Lock()
	defer p.mutex.Unlock()
	p.keyboard = k
}

func (p *Process) Read()  {
	p.mutex.Lock()
	defer p.mutex.Unlock()
	if p.keyboard == nil {
		return
	}
	p.keyboard.mutex.Lock()
	buf := p.keyboard.buf
	p.keyboard.buf = make([]byte, 0)
	p.keyboard.mutex.Unlock()

	for _, c := range buf {
		if c == '\n' {
			p.buf = append(p.buf, c)
			p.readable = true
			p.end = false
			p.sigint = false
			return
		} else if c == 'D' {
			p. readable = true
			p. end = true
			p.sigint = false
			return
		} else if c == 'C' {
			p.readable = true
			p.end = false
			p.sigint = true
			return
		} else {
			p.buf = append(p.buf, c)
		}
	}
}

type Keyboard struct {
	buf []byte
	mutex sync.Mutex
}

func (k *Keyboard) input(){
	var c byte
	var err error
	k.mutex.Lock()
	defer k.mutex.Unlock()
	for true {
		_, err = fmt.Scanf("%c", &c)
		if err != nil {
			panic(err)
		}
		k.buf = append(k.buf, c)
		if c == '\n' || c == 'D' || c == 'C' {
			return
		}
	}
}

func NewKeyboard() *Keyboard {
	return &Keyboard{buf : make([]byte, 0)}
}

func main() {
	k := NewKeyboard()
	pm := NewProcManager(k)
	for true{
		var cmd int
		log.Println("0. Create process")
		log.Println("1. Keyboard Input")
		log.Println("2. Keyboard Switch")
		log.Println("3. Check process")
		log.Println("4. Recover process")
		fmt.Scan(&cmd)
		if cmd == 0 {
			pm.NewProcess()
		} else if cmd == 1 {
			k.input()
		} else if cmd == 2 {
			log.Println("Input pid")
			var id int
			fmt.Scan(&id)
			pm.Swit(id)
		} else if cmd == 3 {
			pm.Check()
		} else if cmd == 4 {
			pm.RemoveDead()
		}
	}
}

