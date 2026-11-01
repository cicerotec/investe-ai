package domain

import (
	"fmt"
	"sort"
)

// ResumoAtivo e a posicao de um ativo calculada das operacoes e eventos.
type ResumoAtivo struct {
	Ticker         string
	Quantidade     float64
	CustoTotal     float64 // custo do que ainda esta em carteira
	PrecoMedio     float64
	LucroRealizado float64 // das vendas, ja descontados os impostos da venda
}

// ResumirAtivo aplica operacoes e eventos em ordem de data e devolve a posicao.
// No mesmo dia, o evento vale antes das operacoes (a negociacao ja usa a nova quantidade).
// Compra e subscricao somam quantidade e custo (preco x quantidade + impostos); a venda tira custo pelo preco medio.
func ResumirAtivo(ticker string, ops []Operacao, evs []EventoAtivo) (ResumoAtivo, error) {
	type passo struct {
		op *Operacao
		ev *EventoAtivo
	}
	var passos []passo
	for i := range ops {
		passos = append(passos, passo{op: &ops[i]})
	}
	for i := range evs {
		passos = append(passos, passo{ev: &evs[i]})
	}
	quando := func(p passo) int64 {
		if p.op != nil {
			return p.op.Data.UnixNano()
		}
		return p.ev.Data.UnixNano()
	}
	sort.SliceStable(passos, func(i, j int) bool {
		di, dj := quando(passos[i]), quando(passos[j])
		if di != dj {
			return di < dj
		}
		return passos[i].ev != nil && passos[j].ev == nil // evento antes da operacao
	})

	r := ResumoAtivo{Ticker: ticker}
	for _, p := range passos {
		switch {
		case p.ev != nil:
			switch p.ev.Tipo {
			case TrocaTicker:
				if p.ev.NovoTicker != "" {
					r.Ticker = p.ev.NovoTicker
				}
			case Agrupamento, Desdobramento:
				if p.ev.Fator <= 0 {
					return ResumoAtivo{}, fmt.Errorf("evento %s com fator invalido", p.ev.Tipo)
				}
				r.Quantidade *= p.ev.Fator
			}
		case p.op.Tipo == Compra || p.op.Tipo == Subscricao:
			r.Quantidade += p.op.Quantidade
			r.CustoTotal += p.op.Quantidade*p.op.Preco + p.op.Impostos
		case p.op.Tipo == Venda:
			if p.op.Quantidade > r.Quantidade+1e-9 {
				return ResumoAtivo{}, ErrVendaMaiorQueSaldo
			}
			pm := 0.0
			if r.Quantidade > 0 {
				pm = r.CustoTotal / r.Quantidade
			}
			r.LucroRealizado += p.op.Quantidade*(p.op.Preco-pm) - p.op.Impostos
			r.CustoTotal -= p.op.Quantidade * pm
			r.Quantidade -= p.op.Quantidade
		}
	}
	if r.Quantidade > 1e-9 {
		r.PrecoMedio = r.CustoTotal / r.Quantidade
	} else {
		r.Quantidade, r.CustoTotal = 0, 0
	}
	return r, nil
}
