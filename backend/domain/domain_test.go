package domain

import (
	"errors"
	"math"
	"testing"
	"time"
)

func dia(d int) time.Time { return time.Date(2024, 1, d, 0, 0, 0, 0, time.UTC) }

func perto(t *testing.T, nome string, got, want float64) {
	t.Helper()
	if math.Abs(got-want) > 1e-6 {
		t.Errorf("%s: got %v, want %v", nome, got, want)
	}
}

func TestSaldoSomaEntradasESubtraiRetiradas(t *testing.T) {
	movs := []Movimento{
		{Tipo: Deposito, Valor: 1000},
		{Tipo: Rendimento, Valor: 50},
		{Tipo: Retirada, Valor: 300},
	}
	perto(t, "saldo", Saldo(movs), 750)
}

func TestAjustarSaldo(t *testing.T) {
	m, ok := AjustarSaldo(750, 760.5, dia(1))
	if !ok || m.Tipo != Rendimento {
		t.Fatalf("subiu: esperava rendimento, got %+v ok=%v", m, ok)
	}
	perto(t, "valor", m.Valor, 10.5)
	m, ok = AjustarSaldo(750, 700, dia(1))
	if !ok || m.Tipo != Retirada {
		t.Fatalf("caiu: esperava retirada, got %+v ok=%v", m, ok)
	}
	perto(t, "valor", m.Valor, 50)
	if _, ok := AjustarSaldo(750, 750, dia(1)); ok {
		t.Error("sem diferenca nao deve criar movimento")
	}
}

func TestPrecoMedioDeDuasCompras(t *testing.T) {
	r, err := ResumirAtivo("PETR4", []Operacao{
		{Data: dia(1), Tipo: Compra, Quantidade: 10, Preco: 10},
		{Data: dia(2), Tipo: Compra, Quantidade: 10, Preco: 20},
	}, nil)
	if err != nil {
		t.Fatal(err)
	}
	perto(t, "quantidade", r.Quantidade, 20)
	perto(t, "preco medio", r.PrecoMedio, 15)
}

func TestVendaTiraCustoPeloPrecoMedioEApuraLucro(t *testing.T) {
	r, err := ResumirAtivo("PETR4", []Operacao{
		{Data: dia(1), Tipo: Compra, Quantidade: 10, Preco: 10},
		{Data: dia(2), Tipo: Compra, Quantidade: 10, Preco: 20},
		{Data: dia(3), Tipo: Venda, Quantidade: 5, Preco: 30, Impostos: 1},
	}, nil)
	if err != nil {
		t.Fatal(err)
	}
	perto(t, "quantidade", r.Quantidade, 15)
	perto(t, "preco medio continua", r.PrecoMedio, 15)
	perto(t, "lucro", r.LucroRealizado, 5*(30-15)-1)
}

func TestVendaMaiorQueACarteiraDaErro(t *testing.T) {
	_, err := ResumirAtivo("X", []Operacao{
		{Data: dia(1), Tipo: Compra, Quantidade: 1, Preco: 10},
		{Data: dia(2), Tipo: Venda, Quantidade: 2, Preco: 10},
	}, nil)
	if !errors.Is(err, ErrVendaMaiorQueSaldo) {
		t.Fatalf("got %v", err)
	}
}

func TestDesdobramentoPreservaCusto(t *testing.T) {
	// 10 acoes a 11, depois 1 para 1,1: 11 acoes, custo total igual (110), medio 10.
	r, err := ResumirAtivo("AXIA3", []Operacao{
		{Data: dia(1), Tipo: Compra, Quantidade: 10, Preco: 11},
	}, []EventoAtivo{{Data: dia(2), Tipo: Desdobramento, Fator: 1.1}})
	if err != nil {
		t.Fatal(err)
	}
	perto(t, "quantidade", r.Quantidade, 11)
	perto(t, "custo", r.CustoTotal, 110)
	perto(t, "preco medio", r.PrecoMedio, 10)
}

func TestTrocaDeTickerMantemPosicao(t *testing.T) {
	r, err := ResumirAtivo("WIZC3", []Operacao{
		{Data: dia(1), Tipo: Compra, Quantidade: 10, Preco: 5},
	}, []EventoAtivo{{Data: dia(2), Tipo: TrocaTicker, NovoTicker: "WIZS3"}})
	if err != nil {
		t.Fatal(err)
	}
	if r.Ticker != "WIZS3" {
		t.Errorf("ticker: %s", r.Ticker)
	}
	perto(t, "quantidade", r.Quantidade, 10)
}

func TestEventoNoMesmoDiaVaiAntesDaOperacao(t *testing.T) {
	// Agrupamento 10 para 1 e venda de 1 papel no mesmo dia: a venda ja usa a quantidade nova.
	r, err := ResumirAtivo("X", []Operacao{
		{Data: dia(1), Tipo: Compra, Quantidade: 10, Preco: 1},
		{Data: dia(2), Tipo: Venda, Quantidade: 1, Preco: 10},
	}, []EventoAtivo{{Data: dia(2), Tipo: Agrupamento, Fator: 0.1}})
	if err != nil {
		t.Fatal(err)
	}
	perto(t, "quantidade", r.Quantidade, 0)
}

func TestFatorInvalidoDaErro(t *testing.T) {
	_, err := ResumirAtivo("X", nil, []EventoAtivo{{Data: dia(1), Tipo: Agrupamento, Fator: 0}})
	if err == nil {
		t.Error("fator zero deve dar erro")
	}
}
