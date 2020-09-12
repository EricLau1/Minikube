# Criando CronJobs no Linux

Um cronjob executa uma tarefa no sistema operacional 
sempre que atingir determinado horário ou data 
definida na criação do mesmo.

## Comandos

```bash
    # Verificar se o usuario pode criar um cronjob
    crontab -e

    # Comando para mudar o editor do cronjob
    select-editor

    # Visualizar cron jobs criados
    crontab -l

    # Verificar se o daemon do Cron está rodando
    pgrep cron

    # Ver os logs do cron
    sudo cat /var/log/syslog | grep cron
```

- O arquivo deve conter o seguinte conteúdo:

```vim
# Edit this file to introduce tasks to be run by cron.
# 
# Each task to run has to be defined through a single line
# indicating with different fields when the task will be run
# and what command to run for the task
# 
# To define the time you can provide concrete values for
# minute (m), hour (h), day of month (dom), month (mon),
# and day of week (dow) or use '*' in these fields (for 'any').
# 
# Notice that tasks will be started based on the cron's system
# daemon's notion of time and timezones.
# 
# Output of the crontab jobs (including errors) is sent through
# email to the user the crontab file belongs to (unless redirected).
# 
# For example, you can run a backup of all your user accounts
# at 5 a.m every week with:
# 0 5 * * 1 tar -zcf /var/backups/home.tgz /home/
# 
# For more information see the manual pages of crontab(5) and cron(8)
# 
# m h  dom mon dow   command
 ```

 - Tradução:

Edite esse arquivo para introduzir tarefas a serem executadas pelo cron.

Cada tarefa a ser executada deve ser definida através de uma única linha
indicando com campos diferentes quando a tarefa será executada
e qual comando executar para a tarefa.

Para definir o tempo, você pode fornecer valores concretos para
minuto (m), hora (h), dia do mês (dom), mês (mon),
e dia da semana (dow) ou usar '__*__' nesses campos (para 'qualquer').

Observe que as tarefas serão iniciadas com base na noção de horário e fuso horário do daemon do sistema do cron.

A saída dos trabalhos do crontab (incluindo erros) é enviada por email ao usuário ao qual o arquivo crontab pertence (a menos que seja redirecionado).

Por exemplo, você pode executar um backup de todas as suas contas de usuário
às 5 da manhã todas as semanas com:
0 5 * * 1 tar -zcf /var/backups/home.tgz /home/

Para mais informações, consulte as páginas de manual do crontab(5) e cron(8)

Exemplo de cronjob que pode ser colocado no arquivo:

```vim
30 10 * * * echo "Cron job executado: $(date)" >> cronjob-example.txt
```

> Sempre que adicionar uma linha no arquivo, faça uma quebra de linha para evitar bugs.

A tarefa acima será executada todas as vezes que o horário for 10:30 de qualquer dia do mês, em qualquer mês, e, em qualquer dia da semana (seguindo o exemplo descrito no próprio arquivo: `m h dom mon dow command`), ou seja todos os dias as 10:30 o comando `echo "Cron job executado: $(date)" >> cronjob-example.txt` será executado.

- A sigla "__dom__" é para `day of month` (dia do mês).
- A sigla "__mon__" é para `month` (mês).
- A sigla "__dow__" pe para `day of week` (dia da semana).


## Criando um cronjobs no Kubernetes

1. Criar o arquivo `cronjob.yaml`

Este CronJob imprime o horário atual e uma string uma vez por minuto:

```yaml
# cronjob.yaml
apiVersion: batch/v1beta1
kind: CronJob
metadata:
  name: hello
spec:
  schedule: "*/1 * * * *"
  jobTemplate:
    spec:
      template:
        spec:
          containers:
          - name: hello
            image: busybox
            args:
            - /bin/sh
            - -c
            - date; echo "Hello, World!"
          restartPolicy: OnFailure
```


2. Rodar o comando:

```bash
    kubectl apply -f cronjob.yaml
```