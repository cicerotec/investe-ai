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
| **Instituição** | nome, tipo (banco, corretora, exchange) |
| **Categoria** | nome, cor, ordem |
| **Posição** | carteira, instituição, categoria, descrição do ativo, moeda (BRL ou USD), valor atual, data da última atualização |
| **Carteira** | nome, moeda base |
| **Acesso** | pessoa (login), carteira, papel (dono, coautor ou leitor) |
| **Histórico** | posição, valor, data (um registro a cada atualização, para ver a evolução) |

Instituições, categorias e posições pertencem a uma carteira. Uma posição liga uma instituição a uma categoria. Isso permite as duas visões:
- por categoria: quanto há em reserva de emergência e em quais bancos;
- por instituição: tudo o que está em cada banco ou corretora.

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
1. **Resumo:** patrimônio total e total por categoria.
2. **Categoria:** posições da categoria agrupadas por instituição.
3. **Instituição:** posições da instituição agrupadas por categoria.
4. **Atualizar valores:** lista das posições com edição rápida do valor.
5. **Cadastros:** instituições, categorias e posições.
6. **Carteiras e acessos:** trocar de carteira, criar uma nova e convidar ou remover pessoas (só o dono).

## Stack (mesmo padrão do financas-ai)
- **Front:** HTML, CSS e JavaScript puros, sem build.
- **Backend:** AWS Lambda em Go com Function URL.
- **Dados:** DynamoDB.
- **Login:** Cognito (Hosted UI + PKCE).
- **Hospedagem:** S3 privado + CloudFront, região `sa-east-1`, infra em AWS SAM.

Roadmap e decisões em [`BACKLOG.md`](BACKLOG.md).
