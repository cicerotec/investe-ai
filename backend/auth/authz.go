// Autorizacao por papel dentro de uma carteira. Tudo aqui e puro (sem AWS), para ser testado em authz_test.go.
// Regra de ouro: o que nao esta liberado explicitamente e negado.
package auth

const (
	PapelDono    = "dono"
	PapelCoautor = "coautor"
	PapelLeitor  = "leitor"
)

// Acao e uma coisa que alguem pode tentar fazer; cada rota da API vira uma Acao.
type Acao string

const (
	AcaoLerCarteira      Acao = "carteira.ler"
	AcaoLerPosicoes      Acao = "posicoes.ler"
	AcaoLerHistorico     Acao = "historico.ler"
	AcaoEscreverPosicoes Acao = "posicoes.escrever" // criar, atualizar valor, excluir
	AcaoEscreverCadastro Acao = "cadastro.escrever" // instituicoes e categorias
	AcaoGerenciarAcessos Acao = "acessos.gerenciar" // convidar e remover pessoas
	AcaoApagarCarteira   Acao = "carteira.apagar"
)

// Permitido diz se o papel pode fazer a acao. Papel ou acao desconhecidos: false.
func Permitido(papel string, a Acao) bool {
	switch papel {
	case PapelDono:
		switch a {
		case AcaoLerCarteira, AcaoLerPosicoes, AcaoLerHistorico,
			AcaoEscreverPosicoes, AcaoEscreverCadastro, AcaoGerenciarAcessos, AcaoApagarCarteira:
			return true
		}
	case PapelCoautor:
		switch a {
		case AcaoLerCarteira, AcaoLerPosicoes, AcaoLerHistorico,
			AcaoEscreverPosicoes, AcaoEscreverCadastro:
			return true
		}
	case PapelLeitor:
		switch a {
		case AcaoLerCarteira, AcaoLerPosicoes, AcaoLerHistorico:
			return true
		}
	}
	return false
}
