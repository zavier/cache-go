package server

import (
	"../cache"
	"bufio"
	"fmt"
	"io"
	"log"
	"net"
	"strconv"
	"strings"
)

type Server struct {
	cache.Cache
}

func (s *Server) Listen() {
	l, e := net.Listen("tcp", ":12346")
	if e != nil {
		panic(e)
	}
	for {
		c, e := l.Accept()
		if e != nil {
			panic(e)
		}
		go s.process(c)
	}
}

func New(c cache.Cache) *Server {
	return &Server{c}
}

func (s *Server) process(conn net.Conn) {
	log.Println("processing...")
	defer conn.Close()
	r := bufio.NewReader(conn)
	for {
		op, e := r.ReadByte()
		if e != nil {
			if e != io.EOF {
				log.Println("close connection due to error", e)
			}
			return
		}
		if op == 'S' {
			e = s.set(conn, r)
		} else if op == 'G' {
			e = s.get(conn, r)
		} else if op == 'D' {
			e = s.del(conn, r)
		} else {
			log.Println("close connection due to invalid operation:", op)
			return
		}
		if e != nil {
			log.Println("close connection due to error", e)
			return
		}
	}
}

func (s *Server) get(conn net.Conn, r *bufio.Reader) error {
	k, e := s.readKey(r)
	if e != nil {
		return e
	}
	v, e := s.Get(k)
	return sendResponse(v, e, conn)
}

func (s *Server) set(conn net.Conn, r *bufio.Reader) error {
	log.Println("seting")
	k, v, e := s.readKeyAndValue(r)
	log.Println("k: ", k, "v: ", v)
	if e != nil {
		return e
	}
	return sendResponse(nil, s.Set(k, v), conn)
}

func (s *Server) del(conn net.Conn, r *bufio.Reader) error {
	k, e := s.readKey(r)
	if e != nil {
		return e
	}
	return sendResponse(nil, s.Del(k), conn)
}

func (s *Server) readKey(r *bufio.Reader) (string, error) {
	klen, e := readLen(r)
	if e != nil {
		return "", e
	}

	k := make([]byte, klen)
	_, e = io.ReadFull(r, k)
	if e != nil {
		return "", e
	}
	return string(k), nil
}

func (s *Server) readKeyAndValue(r *bufio.Reader) (string, []byte, error) {
	klen, e := readLen(r)
	if e != nil {
		return "", nil, e
	}
	log.Println("klen", klen)
	vlen, e := readLen(r)
	log.Println("vlen", vlen)
	k := make([]byte, klen)
	_, e = io.ReadFull(r, k)
	if e != nil {
		return "", nil, e
	}

	v := make([]byte, vlen)
	_, e = io.ReadFull(r, v)
	if e != nil {
		return "", nil, e
	}

	return string(k), v, nil
}

// S<klen><SP><vlen><SP><key><value>
func readLen(r *bufio.Reader) (int, error) {
	// 读取开始到空格的字符（包括空格）
	tmp, e := r.ReadString(' ')
	if e != nil {
		return 0, e
	}

	// 去掉空格后转为数值
	l, e := strconv.Atoi(strings.TrimSpace(tmp))
	if e != nil {
		return 0, e
	}
	return l, nil
}

func sendResponse(value []byte, err error, conn net.Conn) error {
	if err != nil {
		errString := err.Error()
		// 组装错误信息
		tmp := fmt.Sprintf("-%d ", len(errString)) + errString
		_, e := conn.Write([]byte(tmp))
		return e
	}
	vlen := fmt.Sprintf("%d ", len(value))
	_, e := conn.Write(append([]byte(vlen), value...))
	return e
}
