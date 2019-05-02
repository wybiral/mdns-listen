package packet

import (
	"net"

	"github.com/miekg/dns"
)

type Packet struct {
	From      string        `json:"from"`
	ID        uint16        `json:"id"`
	Questions []*Question   `json:"questions,omitempty"`
	Answers   []interface{} `json:"answers,omitempty"`
}

type Question struct {
	Name  string `json:"name"`
	Type  string `json:"type"`
	Class string `json:"class"`
}

type A struct {
	Name string `json:"name"`
	Type string `json:"type"`
	A    string `json:"a"`
}

type AAAA struct {
	Name string `json:"name"`
	Type string `json:"type"`
	AAAA string `json:"aaaa"`
}

type SRV struct {
	Name     string `json:"name"`
	Type     string `json:"type"`
	Priority uint16 `json:"priority"`
	Weight   uint16 `json:"weight"`
	Port     uint16 `json:"port"`
	Target   string `json:"target"`
}

// Parse one packet
func Parse(packet []byte, from net.Addr) (*Packet, error) {
	m := &dns.Msg{}
	err := m.Unpack(packet)
	if err != nil {
		return nil, err
	}
	p := &Packet{
		From: from.String(),
		ID:   m.Id,
	}
	for _, x := range m.Question {
		p.Questions = append(p.Questions, parseQuestion(x))
	}
	for _, x := range m.Answer {
		h := x.Header()
		switch v := x.(type) {
		case *dns.A:
			p.Answers = append(p.Answers, parseA(h, v))
		case *dns.AAAA:
			p.Answers = append(p.Answers, parseAAAA(h, v))
		case *dns.SRV:
			p.Answers = append(p.Answers, parseSRV(h, v))
		}
	}
	return p, nil
}

func parseQuestion(x dns.Question) *Question {
	return &Question{
		Name:  x.Name,
		Type:  dns.TypeToString[x.Qtype],
		Class: dns.ClassToString[x.Qclass],
	}
}

func parseA(h *dns.RR_Header, v *dns.A) *A {
	return &A{
		Name: h.Name,
		Type: dns.TypeToString[h.Rrtype],
		A:    v.A.String(),
	}
}

func parseAAAA(h *dns.RR_Header, v *dns.AAAA) *AAAA {
	return &AAAA{
		Name: h.Name,
		Type: dns.TypeToString[h.Rrtype],
		AAAA: v.AAAA.String(),
	}
}

func parseSRV(h *dns.RR_Header, v *dns.SRV) *SRV {
	return &SRV{
		Name:     h.Name,
		Type:     dns.TypeToString[h.Rrtype],
		Priority: v.Priority,
		Weight:   v.Weight,
		Port:     v.Port,
		Target:   v.Target,
	}
}
