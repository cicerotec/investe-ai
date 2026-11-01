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

// Controle diz como uma categoria e acompanhada.
type Controle string

const (
	ControlePorSaldo Controle = "saldo" // livro de movimentos (reservas, renda fixa, dolar)
	ControlePorAtivo Controle = "ativo" // operacoes por ativo (acoes, FIIs, cripto)
)

// Categoria agrupa posicoes (acoes, FIIs, reserva de emergencia, cripto...).
type Categoria struct {
	ID         string
	CarteiraID string
	Nome       string
	Cor        string
	Ordem      int
	Controle   Controle
}

// Posicao liga uma instituicao a uma categoria. O saldo nao e guardado: vem dos movimentos.
type Posicao struct {
	ID            string
	CarteiraID    string
	InstituicaoID string
	CategoriaID   string
	Descricao     string
	Moeda         string // BRL ou USD
}

// TipoMovimento e a natureza de uma linha do livro de uma posicao por saldo.
type TipoMovimento string

const (
	Deposito   TipoMovimento = "deposito"
	Retirada   TipoMovimento = "retirada"
	Rendimento TipoMovimento = "rendimento"
)

// Movimento e uma linha do livro de uma posicao por saldo. Valor e sempre positivo; o sinal vem do tipo.
type Movimento struct {
	ID        string
	PosicaoID string
	Data      time.Time
	Tipo      TipoMovimento
	Valor     float64
	Motivo    string
}

// TipoAtivo classifica um ativo de bolsa ou cripto.
type TipoAtivo string

const (
	Acao   TipoAtivo = "acao"
	FII    TipoAtivo = "fii"
	Cripto TipoAtivo = "cripto"
)

// Ativo e um papel ou moeda acompanhado por operacoes.
type Ativo struct {
	ID         string
	CarteiraID string
	Ticker     string
	Tipo       TipoAtivo
}

// TipoOperacao e a natureza de uma operacao sobre um ativo.
type TipoOperacao string

const (
	Compra     TipoOperacao = "compra"
	Venda      TipoOperacao = "venda"
	Subscricao TipoOperacao = "subscricao"
)

// Operacao e uma compra, venda ou subscricao de um ativo. Impostos inclui corretagem e taxas.
type Operacao struct {
	ID            string
	AtivoID       string
	InstituicaoID string
	Data          time.Time
	Tipo          TipoOperacao
	Quantidade    float64
	Preco         float64
	Impostos      float64
}

// TipoEvento e uma mudanca que atinge o ativo inteiro, nao uma operacao da pessoa.
type TipoEvento string

const (
	TrocaTicker   TipoEvento = "troca_ticker"  // so muda o nome; quantidade e custo ficam
	Agrupamento   TipoEvento = "agrupamento"   // fator < 1 (ex.: 10 para 1 = 0,1)
	Desdobramento TipoEvento = "desdobramento" // fator > 1 (ex.: 1 para 2 = 2); tambem serve para 1 para 1,1
)

// EventoAtivo altera a quantidade de um ativo preservando o custo total.
type EventoAtivo struct {
	ID         string
	AtivoID    string
	Data       time.Time
	Tipo       TipoEvento
	Fator      float64 // multiplica a quantidade; ignorado na troca de ticker
	NovoTicker string  // so na troca de ticker
}
