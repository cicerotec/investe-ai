// Tipos e regras puras do dominio (sem AWS, sem HTTP).
package domain

import "time"

// Carteira e o espaco isolado de investimentos de uma pessoa ou familia.
type Carteira struct {
	ID        string
	Nome      string
	MoedaBase string // "BRL"
	CriadaEm  time.Time
}

// Instituicao e um banco, corretora ou exchange de uma carteira.
type Instituicao struct {
	ID         string
	CarteiraID string
	Nome       string
	Tipo       string // banco, corretora, exchange
}

// Categoria agrupa posicoes (acoes, FIIs, reserva de emergencia, cripto...).
type Categoria struct {
	ID         string
	CarteiraID string
	Nome       string
	Cor        string
	Ordem      int
}

// Posicao liga uma instituicao a uma categoria com um valor atual.
// Quantidade e opcional: so faz sentido para ativos (cripto, acoes).
type Posicao struct {
	ID            string
	CarteiraID    string
	InstituicaoID string
	CategoriaID   string
	Descricao     string
	Moeda         string // BRL ou USD
	Valor         float64
	Quantidade    *float64
	AtualizadaEm  time.Time
}
