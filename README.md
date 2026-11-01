# Investe AI

App informativo para ver, na palma da mão, onde estão os investimentos: quanto há em cada categoria e em quais bancos ou corretoras.
Substitui a planilha do Google Sheets, que ficou difícil de manter (muitas abas).

Segue o mesmo padrão do projeto `financas-ai`, com stack própria e separada na AWS.

## Princípios
- **Informativo.** Não calcula imposto nem recomenda nada: mostra posições e totais.
- **Manutenção mínima.** Atualizar um valor é um toque, sem abas nem fórmulas para manter.
- **Celular primeiro.** Instalável na tela inicial (PWA), usado pelo navegador também.
- **Dados sensíveis.** Valores reais nunca entram no git, em exemplos ou em logs.

## Categorias
Ações, FIIs, reserva de emergência, reserva de oportunidade, cripto, dólar e renda fixa. Cada carteira usa as que fizerem sentido (os filhos, por exemplo, têm a própria carteira com renda fixa).
A lista é configurável, para criar ou renomear categorias sem mexer no código.

## Modelo de dados
| Entidade | Campos principais |
|---|---|
| **Carteira** | nome, moeda base |
| **Acesso** | pessoa (login), carteira, papel (dono, coautor ou leitor) |
| **Instituição** | carteira, nome, tipo (banco, corretora, exchange) |
| **Categoria** | carteira, nome, cor, ordem, tipo de controle (saldo ou ativo), meta opcional |
| **Posição** | carteira, instituição, categoria, descrição, moeda, saldo (calculado) |
| **Movimento** | posição, data, tipo (depósito, retirada, rendimento), valor, motivo (texto livre) |
| **Ativo** | carteira, ticker, tipo (ação, FII, cripto) |
| **Operação** | ativo, instituição, data, tipo (compra, venda, subscrição), quantidade, preço, impostos |
| **Evento do ativo** | ativo, data, tipo (troca de ticker, agrupamento, desdobramento), proporção |
| **Histórico** | valor total da carteira e por categoria ao longo do tempo (para a evolução) |

Há dois jeitos de controlar uma categoria:
- **Por saldo** (reservas, renda fixa, dólar): a posição é um livro de movimentos, e o saldo é a soma. O motivo de cada movimento fica registrado.
  "Ajustar saldo" cria sozinho o movimento de rendimento (ou a diferença), para dar pouco trabalho.
- **Por ativo** (ações, FIIs, cripto): a posição vem das operações. O servidor calcula quantidade, preço médio, lucro e % de alocação,
  e aplica os eventos do ativo (por exemplo, `WIZC3` virou `WIZS3`; agrupamento na proporção de 1 para 1,1).

As reservas de emergência e de oportunidade são categorias da carteira, com os bancos dentro. Isso permite as duas visões:
- por categoria: quanto há em reserva de emergência e em quais bancos;
- por instituição: tudo o que está em cada banco ou corretora.

Totais, médias, lucro, alocação e evolução são sempre calculados no servidor, nunca digitados.

## Carteiras e acesso
A unidade do app é a **carteira**: o conjunto de investimentos de uma pessoa (ou família). Cada carteira é um espaço isolado, com as próprias
instituições, categorias e posições. Quem usa o app pode ter acesso a várias carteiras, com um papel em cada uma.
É o mesmo mecanismo do `financas-ai` (Cognito + vínculo da pessoa com o espaço + permissões por papel em `auth/authz.go`), só que aqui uma pessoa
pode estar em mais de um espaço.

| Papel na carteira | O que pode |
|---|---|
| **dono** | tudo: cadastros, posições, valores, convidar e remover pessoas, apagar a carteira |
| **coautor** | ver e atualizar posições e cadastros; não gerencia pessoas nem apaga a carteira |
| **leitor** | somente leitura: resumo, posições e evolução |

Exemplos de uso:
- Você tem a sua carteira (dono) e uma carteira para cada filho, da qual também é dono. O filho pode ser convidado como leitor, quando tiver login.
- Um amigo tem a carteira dele (dono) e convida você como coautor. Ele vê tudo o que você faz e pode remover seu acesso.
- O app seleciona a carteira ativa; o resumo e as visões por categoria e instituição valem só para ela.

Regras:
- O isolamento é **no servidor**: toda leitura e escrita é feita dentro de uma carteira e só passa se a pessoa tiver papel nela.
- Carteiras não se misturam: não há total consolidado entre carteiras de donos diferentes.
- Uma carteira sempre tem pelo menos um dono.
- O que não está liberado explicitamente é negado, com testes ao lado do código.

## Telas
1. **Resumo:** patrimônio total, total e % de alocação por categoria, evolução e progresso das metas.
2. **Categoria:** posições da categoria agrupadas por instituição.
3. **Instituição:** posições da instituição agrupadas por categoria.
4. **Atualizar:** ajustar saldo, lançar movimento (com motivo) e registrar operação de compra ou venda.
5. **Cadastros:** instituições, categorias e posições.
6. **Carteiras e acessos:** trocar de carteira, criar uma nova e convidar ou remover pessoas (só o dono).

## Stack (mesmo padrão do financas-ai)
- **Front:** HTML, CSS e JavaScript puros, sem build.
- **Backend:** AWS Lambda em Go com Function URL.
- **Dados:** DynamoDB.
- **Login:** Cognito (Hosted UI + PKCE).
- **Hospedagem:** S3 privado + CloudFront, região `sa-east-1`, infra em AWS SAM.

Roadmap e decisões em [`BACKLOG.md`](BACKLOG.md).
