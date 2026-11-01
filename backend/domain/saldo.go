package domain

import (
	"errors"
	"math"
	"time"
)

// Saldo soma o livro de uma posicao por saldo: deposito e rendimento entram, retirada sai.
func Saldo(movs []Movimento) float64 {
	var s float64
	for _, m := range movs {
		switch m.Tipo {
		case Deposito, Rendimento:
			s += m.Valor
		case Retirada:
			s -= m.Valor
		}
	}
	return s
}

// AjustarSaldo cria o movimento que leva o saldo atual ao novo: rendimento se subiu, retirada se caiu.
// ok=false quando nao ha diferenca. Quem chama preenche ID e PosicaoID.
func AjustarSaldo(atual, novo float64, data time.Time) (Movimento, bool) {
	dif := round2(novo - atual)
	switch {
	case dif > 0:
		return Movimento{Tipo: Rendimento, Valor: dif, Data: data, Motivo: "ajuste de saldo"}, true
	case dif < 0:
		return Movimento{Tipo: Retirada, Valor: -dif, Data: data, Motivo: "ajuste de saldo"}, true
	}
	return Movimento{}, false
}

func round2(v float64) float64 { return math.Round(v*100) / 100 }

var ErrVendaMaiorQueSaldo = errors.New("venda maior que a quantidade em carteira")
