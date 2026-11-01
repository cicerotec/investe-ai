# Backlog

## Decisões
- Stack própria na AWS, separada do `financas-ai`, com o mesmo padrão de arquitetura.
- Valores atualizados manualmente no início; sem cotações automáticas.
- Sem importação da planilha antiga por enquanto. Se vier a ser feita, a partir de um CSV só com a estrutura.
- Valores reais nunca entram no repositório.
- Só o dono convida ou remove pessoas; o coautor não convida.
- O leitor também vê o histórico e a evolução das posições da carteira.
- Posição = uma linha do que se tem (instituição + categoria + valor). Vender ativos ou retirar/devolver da reserva é só atualizar o valor da posição. Passar dinheiro de uma carteira para outra (ex.: cripto para a carteira de um filho) já funciona ajustando as duas posições à mão.
- Começa só com o saldo atualizado, sem registrar cada movimentação.

## Fase 1: base
- [ ] Estrutura do repositório (`web/`, `backend/`, `scripts/`, `template.yaml`) a partir do padrão do `financas-ai`.
- [ ] Infra SAM: tabela DynamoDB, Cognito, Lambda, S3 e CloudFront.
- [ ] Login com Cognito e vínculo pessoa ↔ carteira com papel (dono, coautor, leitor), uma pessoa em várias carteiras.
- [ ] Criar carteira, seletor de carteira ativa e convite por e-mail (como o `seed.sh` do `financas-ai`).
- [ ] Isolamento por carteira no servidor, com testes de que ninguém lê ou grava em carteira sem papel nela.
- [ ] CRUD de instituições, categorias e posições.

## Fase 2: visualização
- [ ] Resumo: patrimônio total e total por categoria.
- [ ] Visão por categoria, agrupada por instituição.
- [ ] Visão por instituição, agrupada por categoria.
- [ ] Tela de atualização rápida de valores.
- [ ] PWA instalável no celular.

## Fase 3: histórico
- [ ] Registro do valor a cada atualização.
- [ ] Gráfico de evolução do patrimônio e por categoria.
- [ ] Aviso de posições desatualizadas há muito tempo.

## Ideias futuras
- Cotações automáticas de ações, FIIs, cripto e dólar.
- Conversão de dólar para real na visão consolidada.
- Transferência entre carteiras: um só lançamento que diminui uma posição e aumenta a outra, com os dois lados ligados. Só entre carteiras em que a pessoa é dono ou coautor. Dois casos: em reais (vende a cripto, passa o BRL e ele vira renda fixa na carteira do filho) e do próprio ativo (a cripto em si muda de carteira).
- Registrar movimentações (aporte, resgate, venda) além do saldo, para ver o que entrou e saiu de cada posição.
- Metas por categoria (por exemplo, valor-alvo da reserva de emergência).

## Perguntas em aberto
- Dólar: guardar o valor em USD e converter pela cotação, ou guardar já em BRL?
