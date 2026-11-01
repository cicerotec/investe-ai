package auth

import "testing"

var todasAcoes = []Acao{
	AcaoLerCarteira, AcaoLerPosicoes, AcaoLerHistorico,
	AcaoEscreverPosicoes, AcaoEscreverCadastro, AcaoGerenciarAcessos, AcaoApagarCarteira,
}

func TestPermitidoPorPapel(t *testing.T) {
	esperado := map[string]map[Acao]bool{
		PapelDono: {
			AcaoLerCarteira: true, AcaoLerPosicoes: true, AcaoLerHistorico: true,
			AcaoEscreverPosicoes: true, AcaoEscreverCadastro: true, AcaoGerenciarAcessos: true, AcaoApagarCarteira: true,
		},
		PapelCoautor: {
			AcaoLerCarteira: true, AcaoLerPosicoes: true, AcaoLerHistorico: true,
			AcaoEscreverPosicoes: true, AcaoEscreverCadastro: true,
		},
		PapelLeitor: {AcaoLerCarteira: true, AcaoLerPosicoes: true, AcaoLerHistorico: true},
	}
	for papel, acoes := range esperado {
		for _, a := range todasAcoes {
			if got := Permitido(papel, a); got != acoes[a] {
				t.Errorf("papel %s, acao %s: got %v, want %v", papel, a, got, acoes[a])
			}
		}
	}
}

func TestCoautorNaoGerenciaPessoasNemApaga(t *testing.T) {
	for _, a := range []Acao{AcaoGerenciarAcessos, AcaoApagarCarteira} {
		if Permitido(PapelCoautor, a) {
			t.Errorf("coautor nao pode %s", a)
		}
	}
}

func TestPapelOuAcaoDesconhecidoENegado(t *testing.T) {
	if Permitido("", AcaoLerPosicoes) || Permitido("admin", AcaoLerPosicoes) {
		t.Error("papel desconhecido deve ser negado")
	}
	if Permitido(PapelDono, Acao("qualquer.coisa")) {
		t.Error("acao desconhecida deve ser negada")
	}
}
